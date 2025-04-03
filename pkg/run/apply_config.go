package run

import (
	"steward/utils"
	"steward/pkg/exec"
	"steward/pkg/pkgman"
	"steward/pkg/common"
)

// var log = logrus.New()
var logger = utils.SetupLogging(true)

type ValidationMode int

const (
    ExactMatch ValidationMode = iota // Exact string match
    LazyMatch                        // Partial or regex-based match
)

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

    // Generate host-specific configuration files
    for _, host := range config.Hosts {
        for _, template := range host.Templates {
            err := common.GenerateConfig(template.TemplateFile, template.OutputFile, template.Data)
            if err != nil {
                logger.Fatalf("Error generating config for host %s: %v", host.Host, err)
                return err
            }
            logger.Infof("Generated config file for host %s: %s", host.Host, template.OutputFile)
        }
    }

	// Install packages for each host
	for _, host := range config.Hosts {
		// SSH client configuration
		sshClient, err := exec.SetupSSHClient(host.Host, "22", host.User, host.Password, "")
		if err != nil {
			logger.Errorf("Error setting up SSH client for host %s: %v", host.Host, err)
			return err
		}
		defer sshClient.Close()
		logger.Infof("SSH client set up for host: %s", host.Host)

		// Execute commands on the remote host
		aptman := pkgman.NewAptManager(sshClient)
		logger.Infof("Installing packages for host: %s", host.Host)

		// Update Repositories
		err = aptman.UpdateRepo()
		if err != nil {
			logger.Errorf("Error updating repositories on host %s: %v", host.Host, err)
			return err
		}
		logger.Infof("Updated repositories on host %s", host.Host)

		// Install common standard packages
		for _, pkg := range config.Common.Packages.Standard {
			err := aptman.InstallPackage(pkg)
			if err != nil {
				logger.Errorf("Error installing package %s on host %s: %v", pkg, host.Host, err)
				return err
			}
			logger.Infof("Installed package %s on host %s", pkg, host.Host)
		}
		// Install host specific standard packages
		for _, pkg := range host.Packages.Standard {
			err := aptman.InstallPackage(pkg)
			if err != nil {
				logger.Errorf("Error installing package %s on host %s: %v", pkg, host.Host, err)
				return err
			}
			logger.Infof("Installed package %s on host %s", pkg, host.Host)
		}

		// Install common third-party packages
		for _, pkg := range config.Common.Packages.ThirdParty {
			// Install the GPG key
			err = aptman.InstallGPGKey(pkg.Name, pkg.GPGKeyURL)
			if err != nil {
				logger.Errorf("Error installing GPG key %s on host %s: %v", pkg.Name, host.Host, err)
				return err
			}
			logger.Infof("Installed GPG key %s on host %s", pkg.Name, host.Host)

			// Add the repository
			err := aptman.AddRepository(pkg.Name, pkg.Repo)
			if err != nil {
				logger.Errorf("Error adding repository %s on host %s: %v", pkg.Name, host.Host, err)
				return err
			}

			// Install the package
			logger.Infof("Install package %s on host %s", pkg.Name, host.Host)
			for _, dep := range pkg.Packages {
				err = aptman.InstallPackage(dep)
				if err != nil {
					logger.Errorf("Error installing third-party package %s on host %s: %v", dep, host.Host, err)
					return err
				}
				logger.Infof("Installed third-party package %s on host %s", dep, host.Host)
			}
		}

		// Install host specific third-party packages
		for _, pkg := range host.Packages.ThirdPartyPackages {
			// Install the GPG key
			err = aptman.InstallGPGKey(pkg.Name, pkg.GPGKeyURL)
			if err != nil {
				logger.Errorf("Error installing GPG key %s on host %s: %v", pkg.Name, host.Host, err)
				return err
			}
			logger.Infof("Installed GPG key %s on host %s", pkg.Name, host.Host)

			// Add the repository
			err := aptman.AddRepository(pkg.Name, pkg.Repo)
			if err != nil {
				logger.Errorf("Error adding repository %s on host %s: %v", pkg.Name, host.Host, err)
				return err
			}

			// Install the package
			logger.Infof("Install package %s on host %s", pkg.Name, host.Host)
			for _, dep := range pkg.Packages {
				err = aptman.InstallPackage(dep)
				if err != nil {
					logger.Errorf("Error installing third-party package %s on host %s: %v", dep, host.Host, err)
					return err
				}
				logger.Infof("Installed third-party package %s on host %s", dep, host.Host)
			}
		}

	}

    return nil
}