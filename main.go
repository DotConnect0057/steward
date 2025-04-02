package main

import (
	// "bytes"
	// "fmt"
	// "io/ioutil"

	// "github.com/sirupsen/logrus"
    // "go.uber.org/zap"
    "fmt"
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

func main() {

	// Load configuration
	config, err := common.LoadConfig("./config/steward/config.yaml")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}

	// Generate common configuration files
	for _, template := range config.Common.Templates {
		err := common.GenerateConfig(template.TemplateFile, template.OutputFile, template.Data)
		if err != nil {
			logger.Fatalf("Error generating config: %v", err)
		}
		logger.Infof("Generated config file: %s", template.OutputFile)
	}
	// Geenrate host-specific configuration files
	for _, host := range config.Hosts {
		for _, template := range host.Templates {
			err := common.GenerateConfig(template.TemplateFile, template.OutputFile, template.Data)
			if err != nil {
				logger.Fatalf("Error generating config: %v", err)
			}
			logger.Infof("Generated config file for host %s: %s", host.Host, template.OutputFile)
		}
	}

    port := "22"
    client, err := exec.SetupSSHClient("192.168.100.14", port, "admin", "admin", "")
    if err != nil {
        logger.Fatalf("Error setting up SSH client: %v", err)
    }
    logger.Info("SSH connection successfully established")

    err = exec.RunRemoteCommand(client, "ls -l")
    if err != nil {
        logger.Fatalf("Error executing command: %v", err)
    }
    logger.Info("Step 1: Command executed successfully")

	// add gpg key
	aptManager := pkgman.NewAptManager(client)
	repo := "https://pkgs.k8s.io/core:/stable:/v1.32/deb/"
	keyUrl := "https://pkgs.k8s.io/core:/stable:/v1.32/deb/Release.key"

	if err := aptManager.InstallGPGKey("k8s", keyUrl); err != nil {
		logger.Fatalf("Error installing GPG key: %v", err)
	}
	if err := aptManager.AddRepository("k8s", repo); err != nil {
		logger.Fatalf("Error adding repository: %v", err)
	}
	if err := aptManager.InstallPackage("kubelet"); err != nil {
		logger.Fatalf("Error installing package: %v", err)
	}
	logger.Info("Step 2: GPG key and repository added successfully")
	
    // err = ssh.RunRemoteCommandWithValidation(client, "echo 'Hello'", "Hello", 0)
    // if err != nil {
    //     logger.Fatalf("Error executing command with validation: %v", err)
    // }
    // logger.Info("Step 2: Command executed successfully with validation")

    fmt.Println("Done !")
}