package common

import (
    "encoding/json"
    "io/ioutil"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
    "text/template"
    "steward/utils"
)

var logger = utils.SetupLogging(false)

func GenerateStewardConfig() error {
    // Define the data structure
    data := map[string]interface{}{
        "common": map[string]interface{}{
            "packages": map[string]interface{}{
                "standard": []string{
                    "apt-transport-https",
                    "ca-certificates",
                    "curl",
                    "gpg",
                },
                "third_party": []map[string]interface{}{
                    {
                        "name":       "kubernetes",
                        "gpg_key_url": "https://pkgs.k8s.io/core:/stable:/v1.32/deb/Release.key",
                        "repo":       "https://pkgs.k8s.io/core:/stable:/v1.32/deb/",
                        "packages": []string{
                            "kubeadm",
                            "kubelet",
                            "kubectl",
                        },
                    },
                },
            },
            "templates": []map[string]interface{}{
                {
                    "name":          "haproxy",
                    "template_file": "./config/template/haproxy_template.cfg",
                    "output_file":   "./config/output/haproxy.cfg",
                    "remote_file":   "/etc/haproxy/haproxy.cfg",
                    "sudo":          true,
                    "data": map[string]interface{}{
                        "backends": []map[string]interface{}{
                            {"name": "web1", "address": "192.168.1.101", "port": 80},
                            {"name": "web2", "address": "192.168.1.102", "port": 80},
                            {"name": "web3", "address": "192.168.1.103", "port": 80},
                        },
                    },
                },
            },
        },
        "hosts": []map[string]interface{}{
            {
                "host":     "192.168.100.14",
                "user":     "admin",
                "password": "admin",
                "packages": map[string]interface{}{
                    "standard": []string{"nginx"},
                    "third_party_packages": []map[string]interface{}{
                        {
                            "name":       "cri-o",
                            "gpg_key_url": "https://pkgs.k8s.io/addons:/cri-o:/stable:/v1.32/deb/Release.key",
                            "repo":       "https://pkgs.k8s.io/addons:/cri-o:/stable:/v1.32/deb/",
                            "packages":   []string{"cri-o"},
                        },
                    },
                },
                "templates": []map[string]interface{}{
                    {
                        "name":          "keepalived",
                        "template_file": "./config/template/keepalived_template.cfg",
                        "output_file":   "./config/output/keepalived.cfg",
                        "remote_file":   "/etc/keepalived/keepalived.conf",
                        "sudo":          true,
                        "data": map[string]interface{}{
                            "virtual_ip": "192.168.1.100",
                            "interface":  "eth0",
                            "priority":   100,
                            "state":      "MASTER",
                        },
                    },
                },
            },
        },
    }

    // Create a temporary directory for output
    outputDir := "./steward-config"
    err := os.MkdirAll(outputDir, 0755)
    if err != nil {
        logger.Errorf("Failed to create dir: %v", err)
        return err
    }

    // Generate YAML file
    yamlFile := filepath.Join(outputDir, "config.yaml")
    err = writeYAMLFile(yamlFile, data)
    if err != nil {
        logger.Errorf("Failed to generate YAML file: %v", err)
        return err
    }
    logger.Infof("YAML file generated at: %s", yamlFile)

    // Generate JSON file
    jsonFile := filepath.Join(outputDir, "config.json")
    err = writeJSONFile(jsonFile, data)
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