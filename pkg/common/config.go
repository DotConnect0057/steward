package common

import (
    "encoding/json"
    "fmt"
    "os"

    "gopkg.in/yaml.v3"
)

// Config represents the structure of the configuration file
type Config struct {
    Common struct {
        Packages struct {
            Standard   []string `yaml:"standard" json:"standard"`
            ThirdParty []struct {
                Name       string   `yaml:"name" json:"name"`
                GPGKeyURL  string   `yaml:"gpg_key_url" json:"gpg_key_url"`
                Repo       string   `yaml:"repo" json:"repo"`
                Packages   []string `yaml:"packages" json:"packages"`
            } `yaml:"third_party" json:"third_party"`
        } `yaml:"packages" json:"packages"`
        Templates []struct {
            Name         string `yaml:"name" json:"name"`
            TemplateFile string `yaml:"template_file" json:"template_file"`
            OutputFile   string `yaml:"output_file" json:"output_file"`
            RemoteFile   string `yaml:"remote_file" json:"remote_file"`
            Sudo         bool   `yaml:"sudo" json:"sudo"`
            Data         any    `yaml:"data" json:"data"`
        } `yaml:"templates" json:"templates"`
    } `yaml:"common" json:"common"`
    Hosts []struct {
        Host              string `yaml:"host" json:"host"`
        User              string `yaml:"user" json:"user"`
        Password          string `yaml:"password" json:"password"`
        SSHKey            string `yaml:"ssh_key,omitempty" json:"ssh_key,omitempty"`
        Packages          struct {
            Standard          []string `yaml:"standard" json:"standard"`
            ThirdPartyPackages []struct {
                Name       string   `yaml:"name" json:"name"`
                GPGKeyURL  string   `yaml:"gpg_key_url" json:"gpg_key_url"`
                Repo       string   `yaml:"repo" json:"repo"`
                Packages   []string `yaml:"packages" json:"packages"`
            } `yaml:"third_party_packages" json:"third_party_packages"`
        } `yaml:"packages" json:"packages"`
        Templates []struct {
            Name         string `yaml:"name" json:"name"`
            TemplateFile string `yaml:"template_file" json:"template_file"`
            OutputFile   string `yaml:"output_file" json:"output_file"`
            RemoteFile   string `yaml:"remote_file" json:"remote_file"`
            Sudo         bool   `yaml:"sudo" json:"sudo"`
            Data         any    `yaml:"data" json:"data"`
        } `yaml:"templates" json:"templates"`
    } `yaml:"hosts" json:"hosts"`
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

// isYAML checks if the file is a YAML file based on its extension
func isYAML(filePath string) bool {
    return len(filePath) > 5 && (filePath[len(filePath)-5:] == ".yaml" || filePath[len(filePath)-4:] == ".yml")
}

// isJSON checks if the file is a JSON file based on its extension
func isJSON(filePath string) bool {
    return len(filePath) > 5 && filePath[len(filePath)-5:] == ".json"
}