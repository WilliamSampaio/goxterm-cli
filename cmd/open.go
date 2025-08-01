/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"goxterm-cli/internal/config"
	"goxterm-cli/internal/sshclient"
	"goxterm-cli/internal/store"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open a SSH connection",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if name == "" {
			fmt.Print("Name: ")
			nameFromUser, _ := reader.ReadString('\n')
			name = strings.TrimSpace(nameFromUser)
		}

		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		if !store.Exists(cfg.StorePath) {
			fmt.Printf("Store does not exist or is not located: %v\n", err)
			os.Exit(1)
		}

		db, err := store.Load(cfg.StorePath)
		if err != nil {
			fmt.Printf("Error loading store: %v\n", err)
			os.Exit(1)
		}

		credential, exists := db.Credentials[name]
		if !exists {
			fmt.Printf("Connection '%s' not found in the store.\n", name)
			fmt.Println("Please run 'goxterm save' to add a new connection.")
			os.Exit(1)
		}

		sshclient.ConnectAndRun(credential)
	},
}

func init() {
	openCmd.Flags().StringVarP(&name, "name", "n", "", "Enter the name of the connection")
	rootCmd.AddCommand(openCmd)
}
