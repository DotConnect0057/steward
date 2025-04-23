package cmd

import (
    "github.com/spf13/cobra"
    "steward/pkg/common"
    "steward/pkg/run"
)

var configPath string // Variable to store the configuration file path

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
    Use:   "apply",
    Short: "Apply the configuration to the system",
    Long: `Apply the configuration to the system. This command will read the configuration
file and apply the settings to the system. It will also provide real-time updates
to the user about the progress of the application process.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        // Load the configuration file
        config, err := common.LoadConfig(configPath)
        if err != nil {
            logger.Errorf("Failed to load steward config from %s: %v", configPath, err)
            return err
        }
        logger.Infof("Debug %s", config)

		// Merge common parameters into host-specific configurations
		mergedConfig, err := common.MergeCommonToHosts(config)
		if err != nil {
			logger.Errorf("Error merging common parameters: %v\n", err)
			return err
		}
        logger.Infof("Steward config loaded successfully from %s", configPath)
        logger.Infof("Merged Config %s", mergedConfig)

        // Apply the configuration
        updatedConfig := run.ApplyConfigWithProgress(mergedConfig)
		logger.Infof("Applying configuration... %s", updatedConfig)
        if updatedConfig == nil {
            logger.Errorf("Failed to apply configuration")
            return err
        }
        logger.Infof("Configuration applied successfully")
        logger.Debugf("Updated configuration: %v", updatedConfig)

        // Update the configuration file with updatedConfig
        err = common.UpdateConfigFile(configPath, updatedConfig)
        if err != nil {
            logger.Errorf("Failed to update config file at %s: %v", configPath, err)
            return err
        }
        logger.Infof("Configuration file updated successfully at %s", configPath)
        return nil
    },
}

func init() {
    rootCmd.AddCommand(applyCmd)

    // Add a flag for specifying the configuration file path
    applyCmd.Flags().StringVarP(&configPath, "config", "c", "./config.yaml", "Path to the configuration file")
}