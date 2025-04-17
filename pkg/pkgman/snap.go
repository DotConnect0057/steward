package pkgman

import (
    "fmt"
    "steward/pkg/exec"
    "golang.org/x/crypto/ssh"
    "steward/utils"
    "strings"
)

var snapLogger = utils.SetupLogging(false)

// SnapManager provides methods to manage Snap packages on a remote server
type SnapManager struct {
    Client *ssh.Client
}

// NewSnapManager creates a new instance of SnapManager
func NewSnapManager(client *ssh.Client) *SnapManager {
    return &SnapManager{Client: client}
}

// InstallPackage installs a Snap package
func (s *SnapManager) InstallPackage(sudoPass string, packageName string) error {
    command := fmt.Sprintf("sudo snap install %s", packageName)
    return exec.RunRemoteCommandWithSudo(s.Client, command, sudoPass)
}

// RemovePackage removes a Snap package
func (s *SnapManager) RemovePackage(sudoPass string, packageName string) error {
    command := fmt.Sprintf("sudo snap remove %s", packageName)
    return exec.RunRemoteCommandWithSudo(s.Client, command, sudoPass)
}

// RefreshPackages refreshes Snap packages
func (s *SnapManager) RefreshPackages(sudoPass string) error {
    command := "sudo snap refresh"
    return exec.RunRemoteCommandWithSudo(s.Client, command, sudoPass)
}

// AddRepository adds a third-party Snap repository to the system
func (s *SnapManager) AddRepository(sudoPass string, assertionFilePath string) error {
    // Import the assertion file
    command := fmt.Sprintf("sudo snap ack %s", assertionFilePath)
    err := exec.RunRemoteCommandWithSudo(s.Client, command, sudoPass)
    if err != nil {
        return fmt.Errorf("failed to add Snap repository: %w", err)
    }

    snapLogger.Infof("Successfully added Snap repository using assertion file: %s", assertionFilePath)
    return nil
}

// IsPackageInstalled checks if a Snap package is installed
func (s *SnapManager) IsPackageInstalled(packageName string) (bool, error) {
    command := fmt.Sprintf("snap list | grep '^%s '", packageName)
    output, err := exec.RunRemoteCommandWithOutput(s.Client, command)
    if err != nil {
        return false, fmt.Errorf("failed to check Snap package: %w", err)
    }

    if strings.Contains(output, packageName) {
        return true, nil
    }
    return false, nil
}

// ListInstalledPackages lists all installed Snap packages
func (s *SnapManager) ListInstalledPackages() ([]string, error) {
    command := "snap list --all"
    output, err := exec.RunRemoteCommandWithOutput(s.Client, command)
    if err != nil {
        return nil, fmt.Errorf("failed to list Snap packages: %w", err)
    }

    var packages []string
    lines := strings.Split(output, "\n")
    for _, line := range lines {
        fields := strings.Fields(line)
        if len(fields) > 0 && fields[0] != "Name" { // Skip the header line
            packages = append(packages, fields[0])
        }
    }
    return packages, nil
}