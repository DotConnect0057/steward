package ssh

import (
    "testing"
)

func TestSetupSSHClient(t *testing.T) {
    host := "192.168.100.14"
    port := "22"
	user := "admin"
    password := "admin"
    keyPath := "" // No SSH key used in this test

    client, err := SetupSSHClient(host, port, user, password, keyPath)
    if err != nil {
        t.Fatalf("Failed to set up SSH client: %v", err)
    }
    defer client.Close()

    t.Log("SSH client setup successfully")
}

func TestRunRemoteCommand(t *testing.T) {
    host := "192.168.100.14"
    port := "22"
    user := "admin"
    password := "admin"
    keyPath := "" // No SSH key used in this test

    client, err := SetupSSHClient(host, port, user, password, keyPath)
    if err != nil {
        t.Fatalf("Failed to set up SSH client: %v", err)
    }
    defer client.Close()

    command := "ls -l"
    if err := RunRemoteCommand(client, command); err != nil {
        t.Fatalf("Failed to run remote command: %v", err)
    }

    t.Log("Remote command executed successfully")
}

func TestRunRemoteCommandWithValidation(t *testing.T) {
    host := "192.168.100.14"
    port := "22"
    user := "admin"
    password := "admin"
    keyPath := "" // No SSH key used in this test

    client, err := SetupSSHClient(host, port, user, password, keyPath)
    if err != nil {
        t.Fatalf("Failed to set up SSH client: %v", err)
    }
    defer client.Close()

    command := "echo Hello"
    expectedOutput := "Hello\n"

    // Test ExactMatch
    if err := RunRemoteCommandWithValidation(client, command, expectedOutput, ExactMatch); err != nil {
        t.Fatalf("ExactMatch validation failed: %v", err)
    }
    t.Log("ExactMatch validation succeeded")

    // Test LazyMatch
    expectedOutput = "Hello"
    if err := RunRemoteCommandWithValidation(client, command, expectedOutput, LazyMatch); err != nil {
        t.Fatalf("LazyMatch validation failed: %v", err)
    }
    t.Log("LazyMatch validation succeeded")
}