package common

import (
    "os"
    "testing"
)

func TestLoadConfigYAML(t *testing.T) {
    // Create a temporary YAML configuration file
    yamlContent := `
common:
  packages:
    standard:
      - apt-transport-https
      - ca-certificates
    third_party:
      - name: "kubernetes"
        gpg_key_url: "https://pkgs.k8s.io/core:/stable:/v1.32/deb/Release.key"
        repo: "https://pkgs.k8s.io/core:/stable:/v1.32/deb/"
        packages:
          - kubeadm
          - kubelet
  templates:
    - name: "haproxy"
      template_file: "haproxy_template.cfg"
      output_file: "haproxy.cfg"
      remote_file: "/etc/haproxy/haproxy.cfg"
      sudo: true
      data:
        backends:
          - name: "web1"
            address: "192.168.1.101"
            port: 80
hosts:
  - host: "192.168.100.14"
    user: "admin"
    password: "admin"
    packages:
      standard:
        - nginx
      third_party_packages:
        - name: "cri-o"
          gpg_key_url: "https://pkgs.k8s.io/addons:/cri-o:/stable:/v1.32/deb/Release.key"
          repo: "https://pkgs.k8s.io/addons:/cri-o:/stable:/v1.32/deb/"
          packages:
            - cri-o
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
    if len(config.Common.Packages.Standard) != 2 {
        t.Errorf("Expected 2 standard packages, got %d", len(config.Common.Packages.Standard))
    }
    if config.Common.Packages.ThirdParty[0].Name != "kubernetes" {
        t.Errorf("Expected third-party package name 'kubernetes', got '%s'", config.Common.Packages.ThirdParty[0].Name)
    }
    if config.Hosts[0].Host != "192.168.100.14" {
        t.Errorf("Expected host '192.168.100.14', got '%s'", config.Hosts[0].Host)
    }
}

func TestLoadConfigJSON(t *testing.T) {
    // Create a temporary JSON configuration file
    jsonContent := `{
        "common": {
            "packages": {
                "standard": ["apt-transport-https", "ca-certificates"],
                "third_party": [
                    {
                        "name": "kubernetes",
                        "gpg_key_url": "https://pkgs.k8s.io/core:/stable:/v1.32/deb/Release.key",
                        "repo": "https://pkgs.k8s.io/core:/stable:/v1.32/deb/",
                        "packages": ["kubeadm", "kubelet"]
                    }
                ]
            },
            "templates": [
                {
                    "name": "haproxy",
                    "template_file": "haproxy_template.cfg",
                    "output_file": "haproxy.cfg",
                    "remote_file": "/etc/haproxy/haproxy.cfg",
                    "sudo": true,
                    "data": {
                        "backends": [
                            {"name": "web1", "address": "192.168.1.101", "port": 80}
                        ]
                    }
                }
            ]
        },
        "hosts": [
            {
                "host": "192.168.100.14",
                "user": "admin",
                "password": "admin",
                "packages": {
                    "standard": ["nginx"],
                    "third_party_packages": [
                        {
                            "name": "cri-o",
                            "gpg_key_url": "https://pkgs.k8s.io/addons:/cri-o:/stable:/v1.32/deb/Release.key",
                            "repo": "https://pkgs.k8s.io/addons:/cri-o:/stable:/v1.32/deb/",
                            "packages": ["cri-o"]
                        }
                    ]
                }
            }
        ]
    }`

    tmpFile, err := os.CreateTemp("", "config-*.json")
    if err != nil {
        t.Fatalf("Failed to create temporary file: %v", err)
    }
    defer os.Remove(tmpFile.Name()) // Clean up the file after the test

    _, err = tmpFile.Write([]byte(jsonContent))
    if err != nil {
        t.Fatalf("Failed to write to temporary file: %v", err)
    }
    tmpFile.Close()

    // Load the configuration
    config, err := LoadConfig(tmpFile.Name())
    if err != nil {
        t.Fatalf("Failed to load JSON config: %v", err)
    }

    // Validate the loaded configuration
    if len(config.Common.Packages.Standard) != 2 {
        t.Errorf("Expected 2 standard packages, got %d", len(config.Common.Packages.Standard))
    }
    if config.Common.Packages.ThirdParty[0].Name != "kubernetes" {
        t.Errorf("Expected third-party package name 'kubernetes', got '%s'", config.Common.Packages.ThirdParty[0].Name)
    }
    if config.Hosts[0].Host != "192.168.100.14" {
        t.Errorf("Expected host '192.168.100.14', got '%s'", config.Hosts[0].Host)
    }
}

func TestLoadConfigFromYAMLFile(t *testing.T) {
    // Path to the YAML configuration file
    configFilePath := "../..//config/steward/config.yaml"

    // Load the configuration
    config, err := LoadConfig(configFilePath)
    if err != nil {
        t.Fatalf("Failed to load YAML config from file: %v", err)
    }

    // Validate the loaded configuration
    if len(config.Common.Packages.Standard) == 0 {
        t.Errorf("Expected standard packages, but got none")
    }
    if len(config.Common.Packages.ThirdParty) == 0 {
        t.Errorf("Expected third-party packages, but got none")
    }
    if len(config.Hosts) == 0 {
        t.Errorf("Expected hosts, but got none")
    }

    // Example validation for a specific host
    if config.Hosts[0].Host != "192.168.100.14" {
        t.Errorf("Expected host '192.168.100.14', got '%s'", config.Hosts[0].Host)
    }
}

func TestLoadConfigFromJSONFile(t *testing.T) {
    // Path to the JSON configuration file
    configFilePath := "../../config/steward/config.json"

    // Load the configuration
    config, err := LoadConfig(configFilePath)
    if err != nil {
        t.Fatalf("Failed to load JSON config from file: %v", err)
    }

    // Validate the loaded configuration
    if len(config.Common.Packages.Standard) == 0 {
        t.Errorf("Expected standard packages, but got none")
    }
    if len(config.Common.Packages.ThirdParty) == 0 {
        t.Errorf("Expected third-party packages, but got none")
    }
    if len(config.Hosts) == 0 {
        t.Errorf("Expected hosts, but got none")
    }

    // Example validation for a specific host
    if config.Hosts[0].Host != "192.168.100.14" {
        t.Errorf("Expected host '192.168.100.14', got '%s'", config.Hosts[0].Host)
    }
}