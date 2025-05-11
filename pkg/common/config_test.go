package common

import (
	"os"
	"testing"
)

func TestLoadConfigYAML(t *testing.T) {
	// Create a temporary YAML configuration file
	yamlContent := `
common:
  application:
    core:
      - name: "keepalived"
        manager: "apt"
        version: "2.0.20"
    external:
      - name: "cri-o"
        gpg_key_url: "https://pkgs.k8s.io/addons:/cri-o:/stable:/v1.32/deb/Release.key"
        repo: "https://pkgs.k8s.io/addons:/cri-o:/stable:/v1.32/deb/"
        manager: "apt"
        version: "1.32"
hosts:
  - host: "192.168.100.14"
    port: "22"
    user: "admin"
    password: "admin"
    application:
      core:
        - name: "keepalived"
          manager: "apt"
          version: "2.0.20"
      external:
        - name: "cri-o"
          gpg_key_url: "https://pkgs.k8s.io/addons:/cri-o:/stable:/v1.32/deb/Release.key"
          repo: "https://pkgs.k8s.io/addons:/cri-o:/stable:/v1.32/deb/"
          manager: "apt"
          version: "1.32"
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the file after the test

	_, err = tmpFile.Write([]byte(yamlContent))
	if err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	tmpFile.Close()

	// Load the configuration
	config, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load YAML config: %v", err)
	}

	// Validate the loaded configuration
	if len(config.Common.Application.Core) != 1 {
		t.Errorf("Expected 1 core application, got %d", len(config.Common.Application.Core))
	}
	if config.Common.Application.Core[0].Name != "keepalived" {
		t.Errorf("Expected core application name 'keepalived', got '%s'", config.Common.Application.Core[0].Name)
	}
	if len(config.Common.Application.External) != 1 {
		t.Errorf("Expected 1 external application, got %d", len(config.Common.Application.External))
	}
	if config.Common.Application.External[0].Name != "cri-o" {
		t.Errorf("Expected external application name 'cri-o', got '%s'", config.Common.Application.External[0].Name)
	}
	if config.Hosts[0].Host != "192.168.100.14" {
		t.Errorf("Expected host '192.168.100.14', got '%s'", config.Hosts[0].Host)
	}
	if config.Hosts[0].Application.Core[0].Name != "keepalived" {
		t.Errorf("Expected host core application name 'keepalived', got '%s'", config.Hosts[0].Application.Core[0].Name)
	}
	if config.Hosts[0].Application.External[0].Name != "cri-o" {
		t.Errorf("Expected host external application name 'cri-o', got '%s'", config.Hosts[0].Application.External[0].Name)
	}
}
