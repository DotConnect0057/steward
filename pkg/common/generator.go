package common

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"steward/utils"
	"text/template"

	"gopkg.in/yaml.v3"
)

var logger = utils.SetupLogging(false)

func GenerateStewardConfig(configPath string) error {
	// Define the configuration structure
	config := Config{
		Common: struct {
			Application   Application             `yaml:"application" json:"application"`
			Configuration []ConfigurationTemplate `yaml:"configuration" json:"configuration"`
			Commands      []Command               `yaml:"command" json:"command"`
		}{
			Application: Application{
				Core: []CoreApp{
					{
						Manager: "apt",
						Name:    "apt-transport-https",
						Version: "",
					},
					{
						Name:    "ca-certificates",
						Manager: "apt",
						Version: "",
					},
					{
						Name:    "curl",
						Manager: "apt",
						Version: "",
					},
					{
						Name:    "gpg",
						Manager: "apt",
						Version: "",
					},
				},
				External: []ExternalApp{
					{
						Name:      "kubernetes",
						GPGKeyURL: "https://pkgs.k8s.io/core:/stable:/v1.32/deb/Release.key",
						Repo:      "https://pkgs.k8s.io/core:/stable:/v1.32/deb/",
						Manager:   "apt",
						Version:   "",
					},
					{
						Name:      "kubernetes",
						GPGKeyURL: "https://pkgs.k8s.io/core:/stable:/v1.32/deb/Release.key",
						Repo:      "https://pkgs.k8s.io/core:/stable:/v1.32/deb/",
						Manager:   "apt",
						Version:   "",
					},
					{
						Name:      "kubeadm",
						GPGKeyURL: "https://pkgs.k8s.io/core:/stable:/v1.32/deb/Release.key",
						Repo:      "https://pkgs.k8s.io/core:/stable:/v1.32/deb/",
						Manager:   "apt",
						Version:   "",
					},
				},
			},
			Configuration: []ConfigurationTemplate{
				{
					Name:         "haproxy",
					TemplateFile: "./template/haproxy_template.cfg",
					OutputFile:   "./output/haproxy.cfg",
					RemoteFile:   "/etc/haproxy/haproxy.cfg",
					Sudo:         true,
					Data: map[string]interface{}{
						"backends": []map[string]interface{}{
							{"name": "web1", "address": "192.168.1.101", "port": 80},
							{"name": "web2", "address": "192.168.1.102", "port": 80},
							{"name": "web3", "address": "192.168.1.103", "port": 80},
						},
					},
				},
			},
			Commands: []Command{
				{
					Name:           "update-packages",
					Command:        "apt-get update",
					ExpectedOutput: "",
					Sudo:           true,
				},
			},
		},
		Hosts: []Host{
			{
				Host:     "ubuntu-ssh-service",
				Port:     "22",
				User:     "myuser",
				Password: "password",
				Application: Application{
					Core: []CoreApp{
						{
							Name:    "nginx",
							Manager: "",
							Version: "",
						},
					},
					External: []ExternalApp{
						{
							Name:      "cri-o",
							GPGKeyURL: "https://pkgs.k8s.io/addons:/cri-o:/stable:/v1.32/deb/Release.key",
							Repo:      "https://pkgs.k8s.io/addons:/cri-o:/stable:/v1.32/deb/",
							Manager:   "apt",
							Version:   "",
						},
					},
				},
				Configuration: []ConfigurationTemplate{
					{
						Name:         "keepalived",
						TemplateFile: "./template/keepalived_template.cfg",
						OutputFile:   "./output/keepalived.cfg",
						RemoteFile:   "/etc/keepalived/keepalived.conf",
						Sudo:         true,
						Data: map[string]interface{}{
							"virtual_ip": "192.168.1.100",
							"interface":  "eth0",
							"priority":   100,
							"state":      "MASTER",
						},
					},
				},
				Commands: []Command{
					{
						Name:           "install-nginx",
						Command:        "apt-get install -y nginx",
						ExpectedOutput: "",
						Sudo:           true,
					},
				},
			},
		},
	}

	// Generate YAML file
	err := writeYAMLFile(configPath, config)
	if err != nil {
		logger.Errorf("Failed to generate YAML file: %v", err)
		return err
	}
	logger.Infof("YAML file generated at: %s", configPath)

	// Generate JSON file
	configDir := filepath.Dir(configPath)
	jsonFile := filepath.Join(configDir, "config.json")
	err = writeJSONFile(jsonFile, config)
	if err != nil {
		logger.Errorf("Failed to generate JSON file: %v", err)
		return err
	}
	logger.Infof("JSON file generated at: %s", jsonFile)

	return nil
}

// GenerateConfig generates a configuration file based on the provided template and data
func GenerateConfig(templatePath string, outputPath string, data interface{}) error {

	// Read the template file
	templateData, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return err
	}

	// Parse the template
	tmpl, err := template.New("config").Parse(string(templateData))
	if err != nil {
		return err
	}

	// Create the output file
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Execute the template with the provided data
	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}

	logger.Infof("Config file generated at: %s", outputPath) // Corrected log format
	return nil
}

func writeYAMLFile(filePath string, data interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	return encoder.Encode(data)
}

func writeJSONFile(filePath string, data interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(data)
}

// DebugData prints the data structure for debugging purposes
func DebugData(data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	logger.Infof("Template Data: %s", string(jsonData)) // Corrected log format
}
