package cmd

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
    "steward/pkg/common"
)

var configFilePath string // Variable to store the custom configuration file path

// initCmd represents the init command
var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize a new configuration file",
    Long: `Initialize a new configuration file. This command will create a new configuration
file in the specified directory. The new configuration file will be created with default values.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        // Use the default path if no custom path is provided
        if configFilePath == "" {
            configFilePath = "./config.yaml"
        }

        // Check if file path includes .yaml or .yml
        if filepath.Ext(configFilePath) != ".yaml" && filepath.Ext(configFilePath) != ".yml" {
            err := fmt.Errorf("Invalid file name extension: %s. Please use .yaml or .yml", configFilePath)
            logger.Errorf("Failed to generate steward config: %v", err)
            return err
        }

        // Check if the configuration file already exists
        if _, err := os.Stat(configFilePath); err == nil {
            logger.Infof("Configuration file already exists at %s. Skipping initialization.", configFilePath)
            return nil
        } else if !os.IsNotExist(err) {
            // If there's an error other than "file does not exist", return it
            logger.Errorf("Failed to check if configuration file exists: %v", err)
            return err
        }

        // Generate a new configuration file
        err := common.GenerateStewardConfig(configFilePath)
        if err != nil {
            logger.Errorf("Failed to generate steward config: %v", err)
            return err
        }
        logger.Infof("Steward config generated successfully at %s", configFilePath)

        // Create additional directories: "config/output" and "config/template"
        outputDir := "output"
        templateDir := "template"

        for _, dir := range []string{outputDir, templateDir} {
            if err := os.MkdirAll(dir, os.ModePerm); err != nil {
                logger.Errorf("Failed to create directory %s: %v", dir, err)
                return err
            }
            logger.Infof("Directory created: %s", dir)
        }

        return nil
    },
}

func init() {
    rootCmd.AddCommand(initCmd)

    // Add a flag for specifying the configuration file path
    initCmd.Flags().StringVarP(&configFilePath, "config", "c", "", "Path to the configuration file (default: steward-config/config.yaml)")
}