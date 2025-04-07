/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"steward/pkg/common"
	"steward/pkg/run"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the configuration to the system",
	Long: `Apply the configuration to the system. This command will read the configuration
file and apply the settings to the system. It will also provide real-time updates
to the user about the progress of the application process.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := common.LoadConfig("steward-config/config.yaml")
		if err != nil {
			logger.Errorf("Failed to load steward config: %v", err)
		}
		logger.Infof("Steward config loaded successfully")

		// Apply the configuration
		// err = run.ApplyConfig(config)
		run.ApplyConfigWithProgress(config)
		// if err != nil {
		// 	logger.Errorf("Failed to apply configuration: %v", err)
		// }
		logger.Infof("Configuration applied successfully")
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
