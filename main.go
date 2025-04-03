/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"steward/cmd"
	"steward/utils"
)

var logger = utils.SetupLogging(true)

func main() {
	// Execute the main command
	if err := cmd.Execute(); err != nil {
		logger.Fatalf("Error executing command: %v", err)
	}
}
