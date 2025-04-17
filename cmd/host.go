package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "steward/pkg/common"
)

// hostCmd represents the host command
var hostCmd = &cobra.Command{
    Use:   "host",
    Short: "Manage hosts in the configuration",
    Long:  `Add, update, or delete hosts in the configuration file.`,
}

// addHostCmd represents the add subcommand
var addHostCmd = &cobra.Command{
    Use:   "add",
    Short: "Add a new host to the configuration",
    RunE: func(cmd *cobra.Command, args []string) error {
        username, _ := cmd.Flags().GetString("username")
        password, _ := cmd.Flags().GetString("password")
        key, _ := cmd.Flags().GetString("key")
        host, _ := cmd.Flags().GetString("host")

        if host == "" || username == "" {
            // fmt.Println("Error: Host and username are required")
            return fmt.Errorf("Error: Host and username are required")
        }

		if password == "" && key == "" {
			// fmt.Println("Error: Either password or SSH key must be provided")
			return fmt.Errorf("Error: Either password or SSH key must be provided")
		}

        config, err := common.LoadConfig("steward-config/config.yaml")
        if err != nil {
            // logger.Errorf("Failed to load configuration: %v", err)
            return fmt.Errorf("Failed to load configuration: %v", err)
        }

        // Add the new host
        newHost := common.Host{
            Host:     host,
            User:     username,
            Password: password,
            SSHKey:      key,
        }
        config.Hosts = append(config.Hosts, newHost)

        // Save the updated configuration
        err = common.UpdateConfigFile("steward-config/config.yaml", config)
        if err != nil {
            // logger.Errorf("Failed to update configuration file: %v", err)
            return fmt.Errorf("Failed to update configuration file: %v", err)
        }
        logger.Infof("Host %s added successfully", host)

		return nil
    },
}

// updateHostCmd represents the update subcommand
var updateHostCmd = &cobra.Command{
    Use:   "update",
    Short: "Update an existing host in the configuration",
    RunE: func(cmd *cobra.Command, args []string) error {
        username, _ := cmd.Flags().GetString("username")
        password, _ := cmd.Flags().GetString("password")
        key, _ := cmd.Flags().GetString("key")
        host, _ := cmd.Flags().GetString("host")

        if host == "" {
            // fmt.Println("Error: Host is required")
            return fmt.Errorf("Error: Host is required")
        }

        config, err := common.LoadConfig("steward-config/config.yaml")
        if err != nil {
            // logger.Errorf("Failed to load configuration: %v", err)
            return fmt.Errorf("Failed to load configuration: %v", err)
        }

        // Update the host
        updated := false
        for i, h := range config.Hosts {
            if h.Host == host {
                if username != "" {
                    config.Hosts[i].User = username
                }
                if password != "" {
                    config.Hosts[i].Password = password
                }
                if key != "" {
                    config.Hosts[i].SSHKey = key
                }
                updated = true
                break
            }
        }

        if !updated {
            // logger.Errorf("Host %s not found", host)
            return fmt.Errorf("Host %s not found", host)
        }

        // Save the updated configuration
        err = common.UpdateConfigFile("steward-config/config.yaml", config)
        if err != nil {
            // logger.Errorf("Failed to update configuration file: %v", err)
            return fmt.Errorf("Failed to update configuration file: %v", err)
        }
        logger.Infof("Host %s updated successfully", host)

		return nil
    },
}

// deleteHostCmd represents the delete subcommand
var deleteHostCmd = &cobra.Command{
    Use:   "delete",
    Short: "Delete a host from the configuration",
    RunE: func(cmd *cobra.Command, args []string) error {
        host, _ := cmd.Flags().GetString("host")

        if host == "" {
            // fmt.Println("Error: Host is required")
            return fmt.Errorf("Error: Host is required")
        }

        config, err := common.LoadConfig("steward-config/config.yaml")
        if err != nil {
            // logger.Errorf("Failed to load configuration: %v", err)
            return fmt.Errorf("Failed to load configuration: %v", err)
        }

        // Delete the host
        var updatedHosts []common.Host
        deleted := false
        for _, h := range config.Hosts {
            if h.Host != host {
                updatedHosts = append(updatedHosts, h)
            } else {
                deleted = true
            }
        }

        if !deleted {
            // logger.Errorf("Host %s not found", host)
            return fmt.Errorf("Host %s not found", host)
        }

        config.Hosts = updatedHosts

        // Save the updated configuration
        err = common.UpdateConfigFile("steward-config/config.yaml", config)
        if err != nil {
            // logger.Errorf("Failed to update configuration file: %v", err)
            return fmt.Errorf("Failed to update configuration file: %v", err)
        }
        logger.Infof("Host %s deleted successfully", host)

		return nil
    },
}

func init() {
    rootCmd.AddCommand(hostCmd)

    // Add subcommands to the host command
    hostCmd.AddCommand(addHostCmd)
    hostCmd.AddCommand(updateHostCmd)
    hostCmd.AddCommand(deleteHostCmd)

    // Add flags for the subcommands
    addHostCmd.Flags().StringP("host", "H", "", "Host address (required)")
    addHostCmd.Flags().StringP("username", "u", "", "Username for the host (required)")
    addHostCmd.Flags().StringP("password", "p", "", "Password for the host")
    addHostCmd.Flags().StringP("key", "k", "", "SSH key for the host")

    updateHostCmd.Flags().StringP("host", "H", "", "Host address (required)")
    updateHostCmd.Flags().StringP("username", "u", "", "New username for the host")
    updateHostCmd.Flags().StringP("password", "p", "", "New password for the host")
    updateHostCmd.Flags().StringP("key", "k", "", "New SSH key for the host")

    deleteHostCmd.Flags().StringP("host", "H", "", "Host address (required)")
}