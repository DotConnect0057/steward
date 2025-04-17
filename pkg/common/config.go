package common

import (
    "encoding/json"
    "fmt"
    "os"

    "gopkg.in/yaml.v3"
)

// Application represents the core and external applications
type Application struct {
    Core     map[string]string `yaml:"core" json:"core"`
    External []ExternalApp     `yaml:"external" json:"external"`
}

// ExternalApp represents an external application with GPG key, repo, and packages
type ExternalApp struct {
    Name      string            `yaml:"name" json:"name"`
    GPGKeyURL string            `yaml:"gpg_key_url" json:"gpg_key_url"`
    Repo      string            `yaml:"repo" json:"repo"`
    Packages  map[string]string `yaml:"packages" json:"packages"`
}

// ConfigurationTemplate represents a configuration template
type ConfigurationTemplate struct {
    Name         string      `yaml:"name" json:"name"`
    TemplateFile string      `yaml:"template_file" json:"template_file"`
    OutputFile   string      `yaml:"output_file" json:"output_file"`
    RemoteFile   string      `yaml:"remote_file" json:"remote_file"`
    Sudo         bool        `yaml:"sudo" json:"sudo"`
    Data         interface{} `yaml:"data" json:"data"`
}

// Command represents a custom command to execute
type Command struct {
    Name           string `yaml:"name" json:"name"`
    Command        string `yaml:"command" json:"command"`
    ExpectedOutput string `yaml:"expected_output" json:"expected_output"`
    Sudo           bool   `yaml:"sudo" json:"sudo"`
}

// Host represents a host configuration
type Host struct {
    Host              string                   `yaml:"host" json:"host"`
    Port              string                   `yaml:"port" json:"port"`
    User              string                   `yaml:"user" json:"user"`
    Password          string                   `yaml:"password" json:"password"`
    SSHKey            string                   `yaml:"ssh_key,omitempty" json:"ssh_key,omitempty"`
    Application       Application              `yaml:"application" json:"application"`
    Configuration     []ConfigurationTemplate  `yaml:"configuration" json:"configuration"`
    Commands          []Command                `yaml:"command" json:"command"`
}

// Config represents the structure of the configuration file
type Config struct {
    Common struct {
        Application      Application              `yaml:"application" json:"application"`
        Configuration    []ConfigurationTemplate  `yaml:"configuration" json:"configuration"`
        Commands         []Command                `yaml:"command" json:"command"`
    } `yaml:"common" json:"common"`
    Hosts []Host `yaml:"hosts" json:"hosts"`
}

// LoadConfig loads a configuration file (YAML or JSON) into the Config struct
func LoadConfig(filePath string) (*Config, error) {
    // Open the file
    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    // Read the file content
    var config Config
    if isYAML(filePath) {
        decoder := yaml.NewDecoder(file)
        if err := decoder.Decode(&config); err != nil {
            return nil, fmt.Errorf("failed to parse YAML config: %w", err)
        }
    } else if isJSON(filePath) {
        decoder := json.NewDecoder(file)
        if err := decoder.Decode(&config); err != nil {
            return nil, fmt.Errorf("failed to parse JSON config: %w", err)
        }
    } else {
        return nil, fmt.Errorf("unsupported config file format: %s", filePath)
    }

    return &config, nil
}

// This function overwrites the existing file with the new data
func UpdateConfigFile(filePath string, config *Config) error {
    // Open the file for writing
    file, err := os.Create(filePath+".tmp")
    if err != nil {
        return fmt.Errorf("failed to open config file for writing: %w", err)
    }
    defer file.Close()

    // Write the config data to the file
    if isYAML(filePath) {
        encoder := yaml.NewEncoder(file)
        defer encoder.Close()
        if err := encoder.Encode(config); err != nil {
            return fmt.Errorf("failed to write YAML config: %w", err)
        }
    } else if isJSON(filePath) {
        encoder := json.NewEncoder(file)
        if err := encoder.Encode(config); err != nil {
            return fmt.Errorf("failed to write JSON config: %w", err)
        }
    } else {
        return fmt.Errorf("unsupported config file format: %s", filePath)
    }

    return nil
}

// ValidateConfig validates the configuration file
func ValidateConfig(config *Config) error {
    // Implement validation logic here
    // For example, check if required fields are present and valid
    for _, host := range config.Hosts {
        if host.Host == "" {
            return fmt.Errorf("host is required")
        }
        if host.User == "" {
            return fmt.Errorf("user is required for host %s", host.Host)
        }
    }

    return nil
}

// isYAML checks if the file is a YAML file based on its extension
func isYAML(filePath string) bool {
    return len(filePath) > 5 && (filePath[len(filePath)-5:] == ".yaml" || filePath[len(filePath)-4:] == ".yml")
}

// isJSON checks if the file is a JSON file based on its extension
func isJSON(filePath string) bool {
    return len(filePath) > 5 && filePath[len(filePath)-5:] == ".json"
}