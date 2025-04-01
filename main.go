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
)

// var log = logrus.New()
var logger = utils.SetupLogging(false)

type ValidationMode int

const (
    ExactMatch ValidationMode = iota // Exact string match
    LazyMatch                        // Partial or regex-based match
)

func main() {
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