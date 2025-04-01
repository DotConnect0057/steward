package main

import (
	// "bytes"
	// "fmt"
	// "io/ioutil"

	// "github.com/sirupsen/logrus"
    // "go.uber.org/zap"
    "fmt"
	"steward/utils"
	"steward/pkg/ssh"
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
    client, err := ssh.SetupSSHClient("192.168.100.14", port, "admin", "admin", "")
    if err != nil {
        logger.Fatalf("Error setting up SSH client: %v", err)
    }
    logger.Info("SSH connection successfully established")

    err = ssh.RunRemoteCommand(client, "ls -l")
    if err != nil {
        logger.Fatalf("Error executing command: %v", err)
    }
    logger.Info("Step 1: Command executed successfully")

    err = ssh.RunRemoteCommandWithValidation(client, "echo 'Hello'", "Hello", 0)
    if err != nil {
        logger.Fatalf("Error executing command with validation: %v", err)
    }
    logger.Info("Step 2: Command executed successfully with validation")

    fmt.Println("Done !")
}