package common

import (
    "encoding/json"
    "io/ioutil"
    "os"
    "text/template"
    "steward/utils"
)

var logger = utils.SetupLogging(true)

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

// DebugData prints the data structure for debugging purposes
func DebugData(data interface{}) {
    jsonData, _ := json.MarshalIndent(data, "", "  ")
    logger.Infof("Template Data: %s", string(jsonData)) // Corrected log format
}