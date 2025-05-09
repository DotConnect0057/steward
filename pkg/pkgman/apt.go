package pkgman

import (
    "fmt"
    "steward/pkg/exec"
    "golang.org/x/crypto/ssh"
    "steward/utils"
    "strings"
)

var logger = utils.SetupLogging(false)

// AptManager provides methods to manage apt packages on a remote server
type AptManager struct {
    Client *ssh.Client
}

// NewAptManager creates a new instance of AptManager
func NewAptManager(client *ssh.Client) *AptManager {
    return &AptManager{Client: client}
}

// UpdateRepo updates the apt package repository
func (a *AptManager) UpdateRepo(sudoPass string) error {
    command := fmt.Sprintf("sudo apt update")
    // return exec.RunRemoteCommand(a.Client, command)
    return exec.RunRemoteCommandWithSudo(a.Client, command, sudoPass)
}

// InstallPackage installs a package using apt
func (a *AptManager) InstallPackage(sudoPass string, packageName string) error {
    command := fmt.Sprintf("sudo apt install -y %s", packageName)
    return exec.RunRemoteCommandWithSudo(a.Client, command, sudoPass)
}

// AddRepository adds a third-party repository to the system
func (a *AptManager) AddRepository(sudoPass string, repoName string, repoUrl string) error {
    // Check if the repository is already added
    checkCommand := fmt.Sprintf("grep -h '^deb .*%s' /etc/apt/sources.list /etc/apt/sources.list.d/*.list || true", repoUrl)
    output, err := exec.RunRemoteCommandWithOutput(a.Client, checkCommand)
    if err != nil {
        return fmt.Errorf("failed to check repository: %w", err)
    }

    if strings.Contains(output, repoUrl) {
        logger.Infof("Repository '%s' is already added. Skipping.", repoName)
        return nil
    }

    // Add the repository if not already added
    command := fmt.Sprintf("echo 'deb [signed-by=/etc/apt/keyrings/%s-apt-keyring.gpg] %s /' | sudo tee /etc/apt/sources.list.d/%s.list && sudo apt update", repoName, repoUrl, repoName)
    return exec.RunRemoteCommandWithSudo(a.Client, command, sudoPass)
}

// InstallGPGKey installs a GPG key from a URL
func (a *AptManager) InstallGPGKey(sudoPass string, keyName string, keyURL string) error {
    // Check if the GPG key is already installed
    checkCommand := fmt.Sprintf("test -f /etc/apt/keyrings/%s-apt-keyring.gpg && echo 'exists' || true", keyName)
    output, err := exec.RunRemoteCommandWithOutput(a.Client, checkCommand)
    if err != nil {
        return fmt.Errorf("failed to check GPG key: %w", err)
    }

    if strings.Contains(output, "exists") {
        logger.Infof("GPG key '%s' is already installed. Skipping.", keyName)
        return nil
    }

    // Create Directory for keyrings if it doesn't exist
    createDirCommand := "sudo mkdir -p /etc/apt/keyrings"
    if err := exec.RunRemoteCommandWithSudo(a.Client, createDirCommand, sudoPass); err != nil {
        return fmt.Errorf("failed to create keyrings directory: %w", err)
    }

    // Install the GPG key if not already installed
    command := fmt.Sprintf("curl -fsSL %s | sudo gpg --dearmor -o /etc/apt/keyrings/%s-apt-keyring.gpg", keyURL, keyName)
    return exec.RunRemoteCommandWithSudo(a.Client, command, sudoPass)
}

// Check if a package is installed
func (a *AptManager) IsPackageInstalled(packageName string) (bool, error) {
    command := fmt.Sprintf("dpkg -l | grep '^ii' | grep '%s'", packageName)
    output, err := exec.RunRemoteCommandWithOutput(a.Client, command)
    if err != nil {
        return false, fmt.Errorf("failed to check package: %w", err)
    }

    if strings.Contains(output, packageName) {
        return true, nil
    }
    return false, nil
}

// FetchInstalledVersion fetches the installed version of a package and updates a configuration file
func (a *AptManager) FetchInstalledVersion(packageName string) (string, error) {
    // Check if the package is installed
    isInstalled, err := a.IsPackageInstalled(packageName)
    if err != nil {
        return "", fmt.Errorf("failed to check if package is installed: %w", err)
    }

    if !isInstalled {
        return "", fmt.Errorf("package '%s' is not installed", packageName)
    }

    // Fetch the installed version of the package
    command := fmt.Sprintf("dpkg -l | grep '^ii' | grep '  %s  ' | awk '{print $3}'", packageName)
    version, err := exec.RunRemoteCommandWithOutput(a.Client, command)
    if err != nil {
        return "", fmt.Errorf("failed to fetch installed version of package '%s': %w", packageName, err)
    }

    version = strings.TrimSpace(version) // Remove any trailing whitespace

    // Return new package name with version which installed
    return fmt.Sprintf("%s", version), nil 
}