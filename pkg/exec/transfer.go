package exec

import (
	"os"
	"fmt"
    "path/filepath"

    "golang.org/x/crypto/ssh"
    "github.com/pkg/sftp"
)

// TransferFile transfers a file to the remote host using SFTP and ensures the target directory exists
func TransferFile(client *ssh.Client, localFilePath, remoteFilePath string) error {
    // Create an SFTP client
    sftpClient, err := sftp.NewClient(client)
    if err != nil {
        return err
    }
    defer sftpClient.Close()

    // Ensure the remote directory exists
    remoteDir := filepath.Dir(remoteFilePath)

    if err := sftpClient.MkdirAll(remoteDir); err != nil {
        return err
    }

    // Open the local file
    localFile, err := os.Open(localFilePath)
    if err != nil {
        return err
    }
    defer localFile.Close()

    // Create the remote file
    remoteFile, err := sftpClient.Create(remoteFilePath)
    if err != nil {
        return err
    }
    defer remoteFile.Close()

    // Copy the local file content to the remote file
    if _, err := remoteFile.ReadFrom(localFile); err != nil {
        return fmt.Errorf("failed to copy file content to remote file %s: %v", remoteFilePath, err)
    }

    logger.Infof("File %s transferred to %s:%s", localFilePath, client.RemoteAddr(), remoteFilePath)
    return nil
}

func TransferFileWithRoot(client *ssh.Client, localFilePath, remoteFilePath string) error {
	// Detect filename
	fileName := filepath.Base(localFilePath)

    // Create an SFTP client
    sftpClient, err := sftp.NewClient(client)
    if err != nil {
        return err
    }
    defer sftpClient.Close()

    // Ensure the remote directory exists
    remoteDir := filepath.Dir(remoteFilePath)

	err = RunRemoteCommand(client, fmt.Sprintf("sudo mkdir -p %s", remoteDir))
	if err != nil {
		return err
	}

    // Open the local file
    localFile, err := os.Open(localFilePath)
    if err != nil {
        return err
    }
    defer localFile.Close()

    // Create the remote file
    remoteFile, err := sftpClient.Create("/tmp/"+fileName)
    if err != nil {
        return err
    }
    defer remoteFile.Close()

    // Copy the local file content to the remote file
    if _, err := remoteFile.ReadFrom(localFile); err != nil {
        return err
    }

	// Move the file to the desired location with root privileges
	err = RunRemoteCommand(client, fmt.Sprintf("sudo mv /tmp/%s %s", fileName, remoteFilePath))
	if err != nil {
		return err
	}

    logger.Infof("File %s transferred to %s:%s", localFilePath, client.RemoteAddr(), remoteFilePath)
    return nil
}