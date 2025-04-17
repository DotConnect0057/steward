package cmd

import (
    "os"

    "github.com/spf13/cobra"
    "steward/pkg/common"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize a new configuration file",
    Long: `Initialize a new configuration file. This command will create a new configuration
file in the specified directory. The new configuration file will be created with default values.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        // Define the default configuration file path
        configFilePath := "steward-config/config.yaml"

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
        err := common.GenerateStewardConfig()
        if err != nil {
            logger.Errorf("Failed to generate steward config: %v", err)
            return err
        }

        logger.Infof("Steward config generated successfully at %s", configFilePath)
        return nil
    },
}

func init() {
    rootCmd.AddCommand(initCmd)

    // Here you will define your flags and configuration settings.

    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // initCmd.PersistentFlags().String("foo", "", "A help for foo")

    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    // initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}