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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := common.LoadConfig("steward-config/config.yaml")
		if err != nil {
			logger.Errorf("Failed to load steward config: %v", err)
			return err
		}
		logger.Infof("Steward config loaded successfully")

		// Apply the configuration
		err = run.ApplyConfig(config)
		if err != nil {
			logger.Errorf("Failed to apply configuration: %v", err)
			return err
		}
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
