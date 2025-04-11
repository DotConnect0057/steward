package run

import (
	"sync"
	"fmt"
	"os"
	// "time"

	"text/tabwriter"
	"steward/utils"
	"steward/pkg/exec"
	"steward/pkg/pkgman"
	"steward/pkg/common"
)

// var log = logrus.New()
var logger = utils.SetupLogging(false)

type ValidationMode int

type TaskStatus struct {
    HostName       string
	Status         string
    Packages       string
    Templates      string
	Procedures     string
	TotalTasks     int
	PkgTasks       int
	TemplateTasks  int
	ProcedureTasks int
}

const (
    ExactMatch ValidationMode = iota // Exact string match
    LazyMatch                        // Partial or regex-based match
)

func DisplayProgress(totalTasks int, completedTasks int,  tasks []TaskStatus, mu *sync.Mutex) {
    writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

    // Clear the screen
    fmt.Print("\033[H\033[2J")

    // Print the header
    // fmt.Fprintln(writer, "TASK\tSTATUS\tPROGRESS")
    fmt.Fprintf(writer, "Total Task: %d\tCompleted Task: %d\t\n", totalTasks, completedTasks)
    fmt.Fprintf(writer, "----------------------------------------------------\n")

    // Print each task's status
    mu.Lock()
    for _, task := range tasks {
        fmt.Fprintf(writer, "%s\tPackages: %s\tTemplates: %s\tProcedures: %s\t%s\n", task.HostName, task.Packages, task.Templates, task.Procedures, task.Status)
    }
    mu.Unlock()

    writer.Flush()
}

func ApplyConfig(config *common.Config) error {
    // Generate common configuration files
    for _, template := range config.Common.Templates {
        err := common.GenerateConfig(template.TemplateFile, template.OutputFile, template.Data)
        if err != nil {
            logger.Fatalf("Error generating config: %v", err)
            return err
        }
        logger.Infof("Generated config file: %s", template.OutputFile)
    }

    // Use a WaitGroup to manage parallel execution
    var wg sync.WaitGroup
    var mu sync.Mutex // To protect shared resources like logging
    var firstError error

    // Generate host-specific configuration files and install packages in parallel
    for _, host := range config.Hosts {
        wg.Add(1)
        go func(host common.Host) {
            defer wg.Done()

            // Generate host-specific configuration files
            for _, template := range host.Templates {
                err := common.GenerateConfig(template.TemplateFile, template.OutputFile, template.Data)
                if err != nil {
                    mu.Lock()
                    logger.Errorf("Error generating config for host %s: %v", host.Host, err)
                    if firstError == nil {
                        firstError = err
                    }
                    mu.Unlock()
                    return
                }
                mu.Lock()
                logger.Infof("Generated config file for host %s: %s", host.Host, template.OutputFile)
                mu.Unlock()
            }

            // SSH client configuration
            sshClient, err := exec.SetupSSHClient(host.Host, "22", host.User, host.Password, "")
            if err != nil {
                mu.Lock()
                logger.Errorf("Error setting up SSH client for host %s: %v", host.Host, err)
                if firstError == nil {
                    firstError = err
                }
                mu.Unlock()
                return
            }
            defer sshClient.Close()
            mu.Lock()
            logger.Infof("SSH client set up for host: %s", host.Host)
            mu.Unlock()

            // Execute commands on the remote host
            aptman := pkgman.NewAptManager(sshClient)
            mu.Lock()
            logger.Infof("Installing packages for host: %s", host.Host)
            mu.Unlock()

            // Update Repositories
            err = aptman.UpdateRepo(host.Password)
            if err != nil {
                mu.Lock()
                logger.Errorf("Error updating repositories on host %s: %v", host.Host, err)
                if firstError == nil {
                    firstError = err
                }
                mu.Unlock()
                return
            }
            mu.Lock()
            logger.Infof("Updated repositories on host %s", host.Host)
            mu.Unlock()

            // Install common standard packages
            // Iterate over and install each common standard package
            for _, pkg := range config.Common.Packages.Standard {
                err := aptman.InstallPackage(host.Password, pkg)
                if err != nil {
                    mu.Lock()
                    logger.Errorf("Error installing package %s on host %s: %v", pkg, host.Host, err)
                    if firstError == nil {
                        firstError = err
                    }
                    mu.Unlock()
                    return
                }
                mu.Lock()
                logger.Infof("Installed package %s on host %s", pkg, host.Host)
                mu.Unlock()
            }

            // Install host-specific standard packages
            for _, pkg := range host.Packages.Standard {
                err := aptman.InstallPackage(host.Password, pkg)
                if err != nil {
                    mu.Lock()
                    logger.Errorf("Error installing package %s on host %s: %v", pkg, host.Host, err)
                    if firstError == nil {
                        firstError = err
                    }
                    mu.Unlock()
                    return
                }
                mu.Lock()
                logger.Infof("Installed package %s on host %s", pkg, host.Host)
                mu.Unlock()
            }

            // Install common third-party packages
            for _, pkg := range config.Common.Packages.ThirdParty {
                // Install the GPG key
                err = aptman.InstallGPGKey(host.Password, pkg.Name, pkg.GPGKeyURL)
                if err != nil {
                    mu.Lock()
                    logger.Errorf("Error installing GPG key %s on host %s: %v", pkg.Name, host.Host, err)
                    if firstError == nil {
                        firstError = err
                    }
                    mu.Unlock()
                    return
                }
                mu.Lock()
                logger.Infof("Installed GPG key %s on host %s", pkg.Name, host.Host)
                mu.Unlock()

                // Add the repository
                err := aptman.AddRepository(host.Password, pkg.Name, pkg.Repo)
                if err != nil {
                    mu.Lock()
                    logger.Errorf("Error adding repository %s on host %s: %v", pkg.Name, host.Host, err)
                    if firstError == nil {
                        firstError = err
                    }
                    mu.Unlock()
                    return
                }

                // Install the package
                mu.Lock()
                logger.Infof("Install package %s on host %s", pkg.Name, host.Host)
                mu.Unlock()
                for _, dep := range pkg.Packages {
                    err = aptman.InstallPackage(host.Password, dep)
                    if err != nil {
                        mu.Lock()
                        logger.Errorf("Error installing third-party package %s on host %s: %v", dep, host.Host, err)
                        if firstError == nil {
                            firstError = err
                        }
                        mu.Unlock()
                        return
                    }
                    mu.Lock()
                    logger.Infof("Installed third-party package %s on host %s", dep, host.Host)
                    mu.Unlock()
                }
            }
        }(host)
    }

    // Wait for all goroutines to finish
    wg.Wait()

    // Return the first error encountered, if any
    return firstError
}

func ApplyConfigWithProgress(config *common.Config) {
    // Initialize task statuses
    var tasks []TaskStatus
    for _, host := range config.Hosts {

		// Count the number of third-party packages for common each host
		thirdPartyPkgCount := 0
		for _, pkg := range config.Common.Packages.ThirdParty {
			thirdPartyPkgCount += len(pkg.Packages)
		}		

		pkgTasks := len(host.Packages.Standard) + len(config.Common.Packages.Standard) + thirdPartyPkgCount
		templateTasks := len(config.Common.Templates) + len(host.Templates)
		procedureTasks := len(config.Common.CustomProcedures)

		tasks = append(tasks, TaskStatus{
            HostName: host.Host,
            Status:   "Pending",
			Packages: fmt.Sprintf("0/%d", pkgTasks),
			Templates: fmt.Sprintf("0/%d", templateTasks),
			Procedures: fmt.Sprintf("0/%d", procedureTasks),
			TotalTasks: pkgTasks + templateTasks + procedureTasks,
			PkgTasks: pkgTasks,
			TemplateTasks: templateTasks,
			ProcedureTasks: procedureTasks,
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

    // Simulate task execution for each host
    for i, host := range config.Hosts {
        wg.Add(1)
        go func(taskIndex int, host common.Host) {
            defer wg.Done()

			completedPkgTasks := 0
			completedTemplateTasks := 0
			completedProcedureTasks := 0
			// pkgTasks := len(host.Packages.Standard) + len(config.Common.Packages.Standard)
			logger.Infof("Starting task for host: %s", host.Host)

            // SSH client configuration
            sshClient, err := exec.SetupSSHClient(host.Host, host.Port, host.User, host.Password, "")
            if err != nil {
                mu.Lock()
                logger.Errorf("Error setting up SSH client for host %s: %v", host.Host, err)
                mu.Unlock()
                return
            }
            defer sshClient.Close()
            mu.Lock()
            logger.Infof("SSH client set up for host: %s", host.Host)
            mu.Unlock()

            // Execute commands on the remote host
            aptman := pkgman.NewAptManager(sshClient)
            mu.Lock()
            logger.Infof("Installing packages for host: %s", host.Host)
            mu.Unlock()
			tasks[taskIndex].Status = "In Progress"
			DisplayProgress(totalAllHostsTasks, completedTotalTasks, tasks, &mu)

			// Update Repositories to avoid errors of missing links
			err = aptman.UpdateRepo(host.Password)
			if err != nil {
				mu.Lock()
				logger.Errorf("Error updating repositories on host %s: %v", host.Host, err)
				mu.Unlock()
				return
			}

            // Install common standard packages
            for _, pkg := range config.Common.Packages.Standard {
                err := aptman.InstallPackage(host.Password, pkg)
                if err != nil {
                    mu.Lock()
					tasks[taskIndex].Status = "Error"
                    logger.Errorf("Error installing package %s on host %s: %v", pkg, host.Host, err)
                    mu.Unlock()
                    return
                }
                mu.Lock()
                logger.Infof("Installed package %s on host %s", pkg, host.Host)
                mu.Unlock()
				completedPkgTasks++
				completedTotalTasks++
				tasks[taskIndex].Packages = fmt.Sprintf("%d/%d", completedPkgTasks, tasks[taskIndex].PkgTasks)
				tasks[taskIndex].Status = "In Progress"
				DisplayProgress(totalAllHostsTasks, completedTotalTasks, tasks, &mu)
            }

			// Install host-specific standard packages
			for _, pkg := range host.Packages.Standard {
				err := aptman.InstallPackage(host.Password, pkg)
				if err != nil {
					mu.Lock()
					tasks[taskIndex].Status = "Error"
					logger.Errorf("Error installing package %s on host %s: %v", pkg, host.Host, err)
					mu.Unlock()
					return
				}
				mu.Lock()
				logger.Infof("Installed package %s on host %s", pkg, host.Host)
				mu.Unlock()
				completedPkgTasks++
				completedTotalTasks++
				tasks[taskIndex].Packages = fmt.Sprintf("%d/%d", completedPkgTasks, tasks[taskIndex].PkgTasks)
				tasks[taskIndex].Status = "In Progress"
				DisplayProgress(totalAllHostsTasks, completedTotalTasks, tasks, &mu)
			}

			// Install common third-party packages
			for _, pkg := range config.Common.Packages.ThirdParty {
				// Install the GPG key skip if empty
				if pkg.GPGKeyURL == "" {
					mu.Lock()
					logger.Infof("GPG key URL is empty for package %s on host %s, skipping installation", pkg.Name, host.Host)
					mu.Unlock()
				}
				// Install the GPG key
				err = aptman.InstallGPGKey(host.Password, pkg.Name, pkg.GPGKeyURL)
				if err != nil {
					mu.Lock()
					tasks[taskIndex].Status = "Error"
					logger.Errorf("Error installing GPG key %s on host %s: %v", pkg.Name, host.Host, err)
					mu.Unlock()
					return
				}
				mu.Lock()
				logger.Infof("Installed GPG key %s on host %s", pkg.Name, host.Host)
				mu.Unlock()
				// Add the repository
				err = aptman.AddRepository(host.Password, pkg.Name, pkg.Repo)
				if err != nil {
					mu.Lock()
					tasks[taskIndex].Status = "Error"
					logger.Errorf("Error adding repository %s on host %s: %v", pkg.Name, host.Host, err)
					mu.Unlock()
					return
				}
				mu.Lock()
				logger.Infof("Added repository %s on host %s", pkg.Name, host.Host)
				mu.Unlock()
				// Install the package
				mu.Lock()
				logger.Infof("Install package %s on host %s", pkg.Name, host.Host)
				mu.Unlock()
				for _, dep := range pkg.Packages {
					err = aptman.InstallPackage(host.Password, dep)
					if err != nil {
						mu.Lock()
						tasks[taskIndex].Status = "Error"
						logger.Errorf("Error installing third-party package %s on host %s: %v", dep, host.Host, err)
						mu.Unlock()
						return
					}
					mu.Lock()
					logger.Infof("Installed third-party package %s on host %s", dep, host.Host)
					mu.Unlock()
					completedPkgTasks++
					completedTotalTasks++
					tasks[taskIndex].Packages = fmt.Sprintf("%d/%d", completedPkgTasks, tasks[taskIndex].PkgTasks)
					tasks[taskIndex].Status = "In Progress"
					DisplayProgress(totalAllHostsTasks, completedTotalTasks, tasks, &mu)
				}
			}

			// Generate and Transfer common configuration files
			for _, template := range config.Common.Templates {
				err := common.GenerateConfig(template.TemplateFile, template.OutputFile, template.Data)
				if err != nil {
					mu.Lock()
					tasks[taskIndex].Status = "Error"
					logger.Errorf("Error generating config for host %s, template: %s, %v", template.TemplateFile, host.Host, err)
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
					tasks[taskIndex].Status = "Error"
					logger.Errorf("Error transferring file %s to host %s: %v", template.OutputFile, host.Host, err)
					mu.Unlock()
					return
				}
				// mu.Lock()
				// logger.Infof("Generated config file for host %s: %s", host.Host, template.OutputFile)
				// mu.Unlock()
				completedTemplateTasks++
				completedTotalTasks++
				tasks[taskIndex].Templates = fmt.Sprintf("%d/%d", completedTemplateTasks, tasks[taskIndex].TemplateTasks)
				tasks[taskIndex].Status = "In Progress"
				DisplayProgress(totalAllHostsTasks, completedTotalTasks, tasks, &mu)
			}

			// Generate and Transfer host-specific configuration files
			for _, template := range host.Templates {
				err := common.GenerateConfig(template.TemplateFile, template.OutputFile, template.Data)
				if err != nil {
					mu.Lock()
					tasks[taskIndex].Status = "Error"
					logger.Errorf("Error generating config for host %s: %v", host.Host, err)
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
					tasks[taskIndex].Status = "Error"
					logger.Errorf("Error transferring file %s to host %s: %v", template.OutputFile, host.Host, err)
					mu.Unlock()
					return
				}
				// mu.Lock()
				// logger.Infof("Generated config file for host %s: %s", host.Host, template.OutputFile)
				// mu.Unlock()
				completedTemplateTasks++
				completedTotalTasks++
				tasks[taskIndex].Templates = fmt.Sprintf("%d/%d", completedTemplateTasks, tasks[taskIndex].TemplateTasks)
				tasks[taskIndex].Status = "In Progress"
				DisplayProgress(totalAllHostsTasks, completedTotalTasks, tasks, &mu)	
			}

			// Run custom procedures
			for _, procedure := range config.Common.CustomProcedures {
				if procedure.Sudo {
					err = exec.RunRemoteCommandWithSudoValidation(sshClient, fmt.Sprintf("sudo %s", procedure.Command), procedure.ExpectedOutput, exec.LazyMatch, host.Password)
				} else {
					err = exec.RunRemoteCommandWithSudoValidation(sshClient, procedure.Command, procedure.ExpectedOutput, exec.LazyMatch, host.Password)
				}
				if err != nil {
					mu.Lock()
					tasks[taskIndex].Status = "Error"
					logger.Errorf("Error running custom procedure %s on host %s: %v", procedure.Command, host.Host, err)
					mu.Unlock()
					return
				}
				mu.Lock()
				logger.Infof("Run custom procedure %s on host %s", procedure.Command, host.Host)
				mu.Unlock()
				completedProcedureTasks++
				completedTotalTasks++
				tasks[taskIndex].Procedures = fmt.Sprintf("%d/%d", completedProcedureTasks, tasks[taskIndex].ProcedureTasks)
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

}

// func MockApplyConfigWithProgress(config *common.Config) {
//     // Initialize task statuses
//     var tasks []TaskStatus
//     for _, host := range config.Hosts {
// 		// totalTasks := len(host.Packages.Standard) + len(config.Common.Packages.Standard) + len(config.Common.Templates) + len(host.Templates)
// 		pkgTasks := len(host.Packages.Standard) + len(config.Common.Packages.Standard)
// 		templateTasks := len(config.Common.Templates) + len(host.Templates)

// 		tasks = append(tasks, TaskStatus{
//             HostName: host.Host,
//             Status:   "Pending",
// 			Packages: fmt.Sprintf("0/%d", pkgTasks),
// 			Templates: fmt.Sprintf("0/%d", templateTasks),
//         })
//     }

//     var mu sync.Mutex
//     var wg sync.WaitGroup

//     // Display initial progress
//     DisplayProgress(tasks, &mu)

//     // Simulate task execution for each host
//     for i, host := range config.Hosts {
//         wg.Add(1)
//         go func(taskIndex int, host common.Host) {
//             defer wg.Done()

//             totalTasks := len(host.Packages.Standard) + len(config.Common.Packages.Standard) + len(config.Common.Templates) + len(host.Templates)
// 			pkgTasks := len(host.Packages.Standard) + len(config.Common.Packages.Standard)
// 			templateTasks := len(config.Common.Templates) + len(host.Templates)
//             completedTasks := 0
// 			completedPkgTasks := 0
// 			completedTemplateTasks := 0

//             // Simulate task progress
//             for completedTasks < totalTasks {
//                 // Simulate work with different speeds for each host
//                 time.Sleep(time.Duration(700+taskIndex*100) * time.Millisecond)

//                 mu.Lock()
//                 completedTasks++
// 				if completedPkgTasks < pkgTasks {
// 					completedPkgTasks++
//                 if completedTemplateTasks < templateTasks {
//                     completedTemplateTasks++
//                     completedTasks++
//                 }
// 				}
//                 tasks[taskIndex].Status = "In Progress"
// 				tasks[taskIndex].Packages = fmt.Sprintf("%d/%d", completedPkgTasks, pkgTasks)
//                 tasks[taskIndex].Templates = fmt.Sprintf("%d/%d", completedTemplateTasks, templateTasks)
//                 if completedTasks == totalTasks {
//                     tasks[taskIndex].Status = "Completed"
//                 }
//                 mu.Unlock()

//                 // Update the progress display
//                 DisplayProgress(tasks, &mu)
//             }
//         }(i, host)
//     }

//     // Wait for all tasks to complete
//     wg.Wait()

//     // Final display
//     DisplayProgress(tasks, &mu)
// }

func GetConfigCounts(config *common.Config) (int) {
    // Count the number of common templates
    numCommonPackages := len(config.Common.Packages.Standard)
    fmt.Printf("Number of common packages: %d\n", numCommonPackages)

    // Count the number of hosts
    numHosts := len(config.Hosts)
    fmt.Printf("Number of hosts: %d\n", numHosts)

    // Count the number of packages for each host
    for i, host := range config.Hosts {
        numHostTemplates := len(host.Templates)
        numHostPackages := len(host.Packages.Standard)
        fmt.Printf("Host %d (%s): %d templates, %d packages\n", i+1, host.Host, numHostTemplates, numHostPackages)
    }

	return numCommonPackages
}