package exec

import (
	"os"
	"bytes"
	"fmt"
    "strings"

	// "github.com/sirupsen/logrus"
	"steward/utils"
	"golang.org/x/crypto/ssh"
)

type ValidationMode int

var logger = utils.SetupLogging(false)

const (
    ExactMatch ValidationMode = iota // Exact string match
    LazyMatch                        // Partial or substring match
)

// SetupSSHClient sets up an SSH client. Password is optional.
// If password is not provided, use SSH key authentication.
func SetupSSHClient(host string, port string, user string, password string, keyPath string) (*ssh.Client, error) {
    var authMethods []ssh.AuthMethod

    // Add password authentication if provided
    if password != "" {
        authMethods = append(authMethods, ssh.Password(password))
    }

    // Add public key authentication if keyPath is provided
    if keyPath != "" {
        privateKey, err := os.ReadFile(keyPath)
        if err != nil {
            return nil, fmt.Errorf("failed to load SSH key: %v", err)
        }
        // Create the Signer for this private key.
        signer, err := ssh.ParsePrivateKey(privateKey)
        if err != nil {
            logger.Fatalf("Unable to parse private key: %v", err)
        }
        authMethods = append(authMethods, ssh.PublicKeys(signer))
    } else if password == "" {
        // If no password and no keyPath, return an error
        return nil, fmt.Errorf("either password or SSH key must be provided")
    }

    // Configure the SSH client
    sshConfig := &ssh.ClientConfig{
        User:            user,
        Auth:            authMethods,
        HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Skip host key verification (not recommended for production)
    }

    // Connect to the SSH server
    client, err := ssh.Dial("tcp", host+":"+port, sshConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to SSH server: %v", err)
    }

    return client, nil
}

// RunRemoteCommand executes a remote command on the SSH server
func RunRemoteCommand(client *ssh.Client, command string) error {
    session, err := client.NewSession()
    if err != nil {
        return err
    }
    defer session.Close()

    var stdoutBuf, stderrBuf bytes.Buffer
    session.Stdout = &stdoutBuf
    session.Stderr = &stderrBuf

    if err := session.Run(command); err != nil {
        logger.Errorf("Command execution error on host %s: %v\nStdout: %s\nStderr: %s",
            client.RemoteAddr(), err, stdoutBuf.String(), stderrBuf.String())
        return fmt.Errorf("command execution error: %v", err)
    }

    logger.Infof("Command executed successfully on host %s: %s", client.RemoteAddr(), command)
    return nil
}

// RunRemoteCommandWithOutput executes a command on the remote server and returns its output
func RunRemoteCommandWithOutput(client *ssh.Client, command string) (string, error) {
    session, err := client.NewSession()
    if err != nil {
        return "", fmt.Errorf("failed to create SSH session: %w", err)
    }
    defer session.Close()

    var stdoutBuf, stderrBuf bytes.Buffer
    session.Stdout = &stdoutBuf
    session.Stderr = &stderrBuf

    if err := session.Run(command); err != nil {
        return "", fmt.Errorf("command execution failed: %w\nstderr: %s", err, stderrBuf.String())
    }

    return stdoutBuf.String(), nil
}

// RunRemoteCommandWithValidation executes a remote command and validates its output, return error if not valid.
func RunRemoteCommandWithValidation(client *ssh.Client, command string, expectedOutput string, mode ValidationMode) error {
    session, err := client.NewSession()
    if err != nil {
        return err
    }
    defer session.Close()

    var stdoutBuf, stderrBuf bytes.Buffer
    session.Stdout = &stdoutBuf
    session.Stderr = &stderrBuf

    // Run the command
    if err := session.Run(command); err != nil {
        logger.Errorf("Command execution error on host %s: %v\nStdout: %s\nStderr: %s",
            client.RemoteAddr(), err, stdoutBuf.String(), stderrBuf.String())
        return fmt.Errorf("command execution error: %v", err)
    }

    // Capture the output
    output := stdoutBuf.String()

    // Log the output
    logger.Infof("Command executed successfully on host %s: %s\nOutput: %s", client.RemoteAddr(), command, output)

    // Validate the output based on the mode
    switch mode {
    case ExactMatch:
        if output != expectedOutput {
            logger.Errorf("Output validation failed (Exact Match) on host %s: expected '%s', got '%s'",
                client.RemoteAddr(), expectedOutput, output)
            return fmt.Errorf("output validation failed: expected '%s', got '%s'", expectedOutput, output)
        }
    case LazyMatch:
        // Use strings.Contains for substring matching
        if !strings.Contains(output, expectedOutput) {
            logger.Errorf("Output validation failed (Lazy Match) on host %s: expected substring '%s', got '%s'",
                client.RemoteAddr(), expectedOutput, output)
            return fmt.Errorf("output validation failed: expected substring '%s', got '%s'", expectedOutput, output)
        }
    default:
        return fmt.Errorf("unknown validation mode")
    }

    logger.Infof("Output validation succeeded on host %s: %s", client.RemoteAddr(), command)
    return nil
}