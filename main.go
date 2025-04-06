/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	// "os"
	// "fmt"

	"steward/cmd"
	"steward/utils"
)

var logger = utils.SetupLogging(true)

func main() {
	// Execute the main command
	cmd.Execute()
	// if err := cmd.Execute(); err != nil {
	// 	// fmt.Fprintln(os.Stderr, err)
	// 	logger.Fatalf("Error executing command: %v", err)
	// 	os.Exit(1)
	// }
}
