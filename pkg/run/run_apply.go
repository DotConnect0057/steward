package run

import (
	"fmt"
	"os"
	"sync"
	"text/tabwriter"

	"steward/pkg/common"
	"steward/pkg/exec"
	"steward/pkg/pkgman"
	"steward/utils"
)

var logger = utils.SetupLogging(false)

type TaskStatus struct {
	HostName      string
	Status        string
	Application   string
	Configuration string
	Command       string
	TotalTasks    int
	AppTasks      int
	ConfigTasks   int
	CommandTasks  int
}

func DisplayProgress(totalTasks int, completedTasks int, tasks []TaskStatus, mu *sync.Mutex) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Clear the screen
	fmt.Print("\033[H\033[2J")

	// Print the header
	fmt.Fprintf(writer, "Total Task: %d\tCompleted Task: %d\t\n", totalTasks, completedTasks)
	fmt.Fprintf(writer, "\n")
	fmt.Fprintf(writer, "HOST\tPACKAGES\tCONFIGURATION\tCOMMANDS\tSTATUS\n")

	// Print each task's status
	// mu.Lock()
	for _, task := range tasks {
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n", task.HostName, task.Application, task.Configuration, task.Command, task.Status)
	}
	// mu.Unlock()

	writer.Flush()
}

func ApplyConfigWithProgress(config *common.Config) *common.Config {
	// Initialize task statuses
	var tasks []TaskStatus
	for _, host := range config.Hosts {

		// Count the number of core application and external application
		commonExternalApps := len(config.Common.Application.External)
		hostExternalApps := len(host.Application.External)

		// for _, app := range config.Common.Application.External {
		// 	commonExternalApps += len(app)
		// }

		// for _, app := range host.Application.External {
		// 	hostExternalApps += len(app.Packages)
		// }

		appTasks := len(host.Application.Core) + len(config.Common.Application.Core) + commonExternalApps + hostExternalApps
		configTasks := len(host.Configuration) + len(config.Common.Configuration)
		commandTasks := len(host.Commands) + len(config.Common.Commands)

		tasks = append(tasks, TaskStatus{
			HostName:      host.Host,
			Status:        "Pending",
			Application:   fmt.Sprintf("0/%d", appTasks),
			Configuration: fmt.Sprintf("0/%d", configTasks),
			Command:       fmt.Sprintf("0/%d", commandTasks),
			TotalTasks:    appTasks + configTasks + commandTasks,
			AppTasks:      appTasks,
			ConfigTasks:   configTasks,
			CommandTasks:  commandTasks,
		})
	}

	var mu sync.Mutex
	var wg sync.WaitGroup

	totalAllHostsTasks := 0
	for _, task := range tasks {
		totalAllHostsTasks += task.TotalTasks
	}
	completedTotalTasks := 0

	// Display initial progress
	DisplayProgress(totalAllHostsTasks, completedTotalTasks, tasks, &mu)

	// Process each host
	for i, host := range config.Hosts {
		wg.Add(1)
		go func(taskIndex int, host common.Host) {
			defer wg.Done()

			completedAppTasks := 0
			completedConfigTasks := 0
			completedCommandTasks := 0

			logger.Infof("Starting tasks for host: %s", host.Host)

			// SSH client configuration
			sshClient, err := exec.SetupSSHClient(host.Host, host.Port, host.User, host.Password, host.SSHKey)
			if err != nil {
				mu.Lock()
				logger.Errorf("Error setting up SSH client for host %s: %v", host.Host, err)
				tasks[taskIndex].Status = "Error"
				mu.Unlock()
				return
			}
			defer sshClient.Close()

			// Initialize apt package manager
			aptman := pkgman.NewAptManager(sshClient)
			tasks[taskIndex].Status = "In Progress"
			DisplayProgress(totalAllHostsTasks, completedTotalTasks, tasks, &mu)

			err = aptman.UpdateRepo(host.Password)
			if err != nil {
				mu.Lock()
				logger.Errorf("Error updating apt repository on host %s: %v", host.Host, err)
				tasks[taskIndex].Status = "Error"
				mu.Unlock()
				return
			}
			mu.Lock()
			logger.Infof("Updated apt repository on host %s", host.Host)
			mu.Unlock()

			// Install common core packages
			for pkgIndex, pkg := range config.Common.Application.Core {
				fullPkgName := ""
				if pkg.Version != "" {
					fullPkgName = fmt.Sprintf("%s=%s", pkg.Name, pkg.Version)
				} else {
					fullPkgName = pkg.Name
				}

				err := aptman.InstallPackage(host.Password, fullPkgName)
				if err != nil {
					mu.Lock()
					logger.Errorf("Error installing package %s on host %s: %v", pkg, host.Host, err)
					tasks[taskIndex].Status = "Error"
					mu.Unlock()
					return
				}
				mu.Lock()
				appVersion, err := aptman.FetchInstalledVersion(pkg.Name)
				if err != nil {
					logger.Errorf("Error fetching installed version of package %s on host %s: %v", pkg, host.Host, err)
					tasks[taskIndex].Status = "Error"
					return
				}

				// replace app name in config with app name and version
				config.Common.Application.Core[pkgIndex].Version = appVersion
				logger.Infof("Installed package %s on host %s", pkg, host.Host)
				mu.Unlock()

				// update progress
				completedAppTasks++
				completedTotalTasks++
				tasks[taskIndex].Application = fmt.Sprintf("%d/%d", completedAppTasks, tasks[taskIndex].AppTasks)
				tasks[taskIndex].Status = "In Progress"
				DisplayProgress(totalAllHostsTasks, completedTotalTasks, tasks, &mu)
			}

			// Install host specific core packages
			for pkgIndex, pkg := range host.Application.Core {
				fullPkgName := ""
				if pkg.Version != "" {
					fullPkgName = fmt.Sprintf("%s=%s", pkg.Name, pkg.Version)
				} else {
					fullPkgName = pkg.Name
				}
				err := aptman.InstallPackage(host.Password, fullPkgName)
				if err != nil {
					mu.Lock()
					logger.Errorf("Error installing package %s on host %s: %v", pkg, host.Host, err)
					tasks[taskIndex].Status = "Error"
					mu.Unlock()
					return
				}
				mu.Lock()
				appVersion, err := aptman.FetchInstalledVersion(pkg.Name)
				if err != nil {
					logger.Errorf("Error fetching installed version of package %s on host %s: %v", pkg, host.Host, err)
					tasks[taskIndex].Status = "Error"
					return
				}
				// replace app name in config with app name and version
				config.Hosts[i].Application.Core[pkgIndex].Version = appVersion
				// config.Common.Application.Core[pkg] = appVersion
				logger.Infof("Installed package %s on host %s", pkg, host.Host)
				mu.Unlock()

				// update progress
				completedAppTasks++
				completedTotalTasks++
				tasks[taskIndex].Application = fmt.Sprintf("%d/%d", completedAppTasks, tasks[taskIndex].AppTasks)
				tasks[taskIndex].Status = "In Progress"
				DisplayProgress(totalAllHostsTasks, completedTotalTasks, tasks, &mu)
			}

			// Install common external packages
			for pkgIndex, pkg := range config.Common.Application.External {

				// Install GPG key skip if empty
				if pkg.GPGKeyURL != "" {
					err := aptman.InstallGPGKey(host.Password, pkg.Name, pkg.GPGKeyURL)
					if err != nil {
						mu.Lock()
						logger.Errorf("Error installing GPG key %s on host %s: %v", pkg.Name, host.Host, err)
						tasks[taskIndex].Status = "Error"
						mu.Unlock()
						return
					}
					mu.Lock()
					logger.Infof("Installed GPG key %s on host %s", pkg.Name, host.Host)
					mu.Unlock()
				}

				// Install repo skip if empty
				if pkg.Repo != "" {
					err := aptman.AddRepository(host.Password, pkg.Name, pkg.Repo)
					if err != nil {
						mu.Lock()
						logger.Errorf("Error adding repo %s on host %s: %v", pkg.Name, host.Host, err)
						tasks[taskIndex].Status = "Error"
						mu.Unlock()
						return
					}
					mu.Lock()
					logger.Infof("Added repo %s on host %s", pkg.Name, host.Host)
					mu.Unlock()
				}

				// Install packages
				fullPkgName := ""
				if pkg.Version != "" {
					fullPkgName = fmt.Sprintf("%s=%s", pkg.Name, pkg.Version)
				} else {
					fullPkgName = pkg.Name
				}

				err := aptman.InstallPackage(host.Password, fullPkgName)
				if err != nil {
					mu.Lock()
					logger.Errorf("Error installing package %s on host %s: %v", pkg, host.Host, err)
					tasks[taskIndex].Status = "Error"
					mu.Unlock()
					return
				}
				mu.Lock()
				appVersion, err := aptman.FetchInstalledVersion(pkg.Name)
				if err != nil {
					logger.Errorf("Error fetching installed version of package %s on host %s: %v", pkg, host.Host, err)
					tasks[taskIndex].Status = "Error"
					return
				}

				// replace app name in config with app name and version
				config.Common.Application.Core[pkgIndex].Version = appVersion
				logger.Infof("Installed package %s on host %s", pkg, host.Host)

				// update progress
				completedAppTasks++
				completedTotalTasks++
				tasks[taskIndex].Application = fmt.Sprintf("%d/%d", completedAppTasks, tasks[taskIndex].AppTasks)
				tasks[taskIndex].Status = "In Progress"
				DisplayProgress(totalAllHostsTasks, completedTotalTasks, tasks, &mu)
				mu.Unlock()

				// for pkgIndex, pkg := range app {
				// 	fullPkgName := ""
				// 	if ver != "" {
				// 		fullPkgName = fmt.Sprintf("%s=%s", pkg, ver)
				// 	} else {
				// 		fullPkgName = pkg
				// 	}
				// 	err := aptman.InstallPackage(host.Password, fullPkgName)
				// 	if err != nil {
				// 		mu.Lock()
				// 		logger.Errorf("Error installing package %s on host %s: %v", pkg, host.Host, err)
				// 		tasks[taskIndex].Status = "Error"
				// 		mu.Unlock()
				// 		return
				// 	}
				// 	mu.Lock()
				// 	appVersion, err := aptman.FetchInstalledVersion(pkg)
				// 	if err != nil {
				// 		logger.Errorf("Error fetching installed version of package %s on host %s: %v", pkg, host.Host, err)
				// 		tasks[taskIndex].Status = "Error"
				// 		return
				// 	}
				// 	// replace app name in config with app name and version
				// 	config.Common.Application.External[0].Packages[pkg] = appVersion
				// 	logger.Infof("Installed package %s on host %s", pkg, host.Host)
				// 	mu.Unlock()

			}
			// Install host-specific external packages
			for pkgIndex, pkg := range host.Application.External {
				// Install GPG key skip if empty
				if pkg.GPGKeyURL != "" {
					err := aptman.InstallGPGKey(host.Password, pkg.Name, pkg.GPGKeyURL)
					if err != nil {
						mu.Lock()
						logger.Errorf("Error installing GPG key %s on host %s: %v", pkg.Name, host.Host, err)
						tasks[taskIndex].Status = "Error"
						mu.Unlock()
						return
					}
					mu.Lock()
					logger.Infof("Installed GPG key %s on host %s", pkg.Name, host.Host)
					mu.Unlock()
				}
				// Install repo skip if empty
				if pkg.Repo != "" {
					err := aptman.AddRepository(host.Password, pkg.Name, pkg.Repo)
					if err != nil {
						mu.Lock()
						logger.Errorf("Error adding repo %s on host %s: %v", pkg.Name, host.Host, err)
						tasks[taskIndex].Status = "Error"
						mu.Unlock()
						return
					}
					mu.Lock()
					logger.Infof("Added repo %s on host %s", pkg.Name, host.Host)
					mu.Unlock()
				}
				// Install packages
				fullPkgName := ""
				if pkg.Version != "" {
					fullPkgName = fmt.Sprintf("%s=%s", pkg.Name, pkg.Version)
				} else {
					fullPkgName = pkg.Name
				}

				err := aptman.InstallPackage(host.Password, fullPkgName)
				if err != nil {
					mu.Lock()
					logger.Errorf("Error installing package %s on host %s: %v", pkg, host.Host, err)
					tasks[taskIndex].Status = "Error"
					mu.Unlock()
					return
				}
				mu.Lock()
				appVersion, err := aptman.FetchInstalledVersion(pkg.Name)
				if err != nil {
					logger.Errorf("Error fetching installed version of package %s on host %s: %v", pkg, host.Host, err)
					tasks[taskIndex].Status = "Error"
					return
				}

				// replace app name in config with app name and version
				config.Hosts[i].Application.External[pkgIndex].Version = appVersion
				logger.Infof("Installed package %s on host %s", pkg, host.Host)

				// update progress
				completedAppTasks++
				completedTotalTasks++
				tasks[taskIndex].Application = fmt.Sprintf("%d/%d", completedAppTasks, tasks[taskIndex].AppTasks)
				tasks[taskIndex].Status = "In Progress"
				DisplayProgress(totalAllHostsTasks, completedTotalTasks, tasks, &mu)
				mu.Unlock()

			}

			// Generate and transfer common configuration templates
			for _, template := range config.Common.Configuration {
				err := common.GenerateConfig(template.TemplateFile, template.OutputFile, template.Data)
				if err != nil {
					mu.Lock()
					logger.Errorf("Error generating config for template %s on host %s: %v", template.Name, host.Host, err)
					tasks[taskIndex].Status = "Error"
					mu.Unlock()
					return
				}
				if template.Sudo {
					err = exec.TransferFileWithRoot(sshClient, template.OutputFile, template.RemoteFile, host.Password)
				} else {
					err = exec.TransferFile(sshClient, template.OutputFile, template.RemoteFile)
				}
				if err != nil {
					mu.Lock()
					logger.Errorf("Error transferring file %s to host %s: %v", template.OutputFile, host.Host, err)
					tasks[taskIndex].Status = "Error"
					mu.Unlock()
					return
				}
				mu.Lock()
				logger.Infof("Transferred file %s to host %s", template.OutputFile, host.Host)
				mu.Unlock()

				completedConfigTasks++
				completedTotalTasks++
				tasks[taskIndex].Configuration = fmt.Sprintf("%d/%d", completedConfigTasks, tasks[taskIndex].ConfigTasks)
				tasks[taskIndex].Status = "In Progress"
				DisplayProgress(totalAllHostsTasks, completedTotalTasks, tasks, &mu)
			}

			// Generate and transfer host-specific configuration templates
			for _, template := range host.Configuration {
				err := common.GenerateConfig(template.TemplateFile, template.OutputFile, template.Data)
				if err != nil {
					mu.Lock()
					logger.Errorf("Error generating config for template %s on host %s: %v", template.Name, host.Host, err)
					tasks[taskIndex].Status = "Error"
					mu.Unlock()
					return
				}
				if template.Sudo {
					err = exec.TransferFileWithRoot(sshClient, template.OutputFile, template.RemoteFile, host.Password)
				} else {
					err = exec.TransferFile(sshClient, template.OutputFile, template.RemoteFile)
				}
				if err != nil {
					mu.Lock()
					logger.Errorf("Error transferring file %s to host %s: %v", template.OutputFile, host.Host, err)
					tasks[taskIndex].Status = "Error"
					mu.Unlock()
					return
				}
				mu.Lock()
				logger.Infof("Transferred file %s to host %s", template.OutputFile, host.Host)
				mu.Unlock()

				// update progress
				completedConfigTasks++
				completedTotalTasks++
				tasks[taskIndex].Configuration = fmt.Sprintf("%d/%d", completedConfigTasks, tasks[taskIndex].ConfigTasks)
				tasks[taskIndex].Status = "In Progress"
				DisplayProgress(totalAllHostsTasks, completedTotalTasks, tasks, &mu)
			}

			// Execute common commands
			for _, command := range config.Common.Commands {
				if command.Sudo {
					err = exec.RunRemoteCommandWithSudoValidation(sshClient, fmt.Sprintf("sudo %s", command.Command), command.ExpectedOutput, exec.LazyMatch, host.Password)
				} else {
					err = exec.RunRemoteCommandWithValidation(sshClient, command.Command, command.ExpectedOutput, exec.LazyMatch)
				}
				if err != nil {
					mu.Lock()
					logger.Errorf("Error executing command %s on host %s: %v", command.Name, host.Host, err)
					tasks[taskIndex].Status = "Error"
					mu.Unlock()
					return
				}
				mu.Lock()
				logger.Infof("Executed command %s on host %s", command.Name, host.Host)
				mu.Unlock()

				completedCommandTasks++
				completedTotalTasks++
				tasks[taskIndex].Command = fmt.Sprintf("%d/%d", completedCommandTasks, tasks[taskIndex].CommandTasks)
				tasks[taskIndex].Status = "In Progress"
				DisplayProgress(totalAllHostsTasks, completedTotalTasks, tasks, &mu)
			}

			// Execute host-specific commands
			for _, command := range host.Commands {
				if command.Sudo {
					err = exec.RunRemoteCommandWithSudoValidation(sshClient, fmt.Sprintf("sudo %s", command.Command), command.ExpectedOutput, exec.LazyMatch, host.Password)
				} else {
					err = exec.RunRemoteCommandWithValidation(sshClient, command.Command, command.ExpectedOutput, exec.LazyMatch)
				}
				if err != nil {
					mu.Lock()
					logger.Errorf("Error executing command %s on host %s: %v", command.Name, host.Host, err)
					tasks[taskIndex].Status = "Error"
					mu.Unlock()
					return
				}
				mu.Lock()
				logger.Infof("Executed command %s on host %s", command.Name, host.Host)
				mu.Unlock()

				// update progress
				completedCommandTasks++
				completedTotalTasks++
				tasks[taskIndex].Command = fmt.Sprintf("%d/%d", completedCommandTasks, tasks[taskIndex].CommandTasks)
				tasks[taskIndex].Status = "In Progress"
				DisplayProgress(totalAllHostsTasks, completedTotalTasks, tasks, &mu)
			}

			tasks[taskIndex].Status = "Completed"
			DisplayProgress(totalAllHostsTasks, completedTotalTasks, tasks, &mu)
		}(i, host)
	}

	// Wait for all tasks to complete
	wg.Wait()

	// Final display
	DisplayProgress(totalAllHostsTasks, completedTotalTasks, tasks, &mu)
	return config
}
