/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"steward/pkg/common"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new configuration file",
	Long: `Initialize a new configuration file. This command will create a new configuration
file in the specified directory. The new configuration file will be created with default values.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := common.GenerateStewardConfig()
		if err != nil {
			logger.Errorf("Failed to generate steward config: %v", err)
		}
		logger.Infof("Steward config generated successfully")
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