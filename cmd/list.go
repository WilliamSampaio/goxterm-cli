/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"goxterm-cli/internal/config"
	"goxterm-cli/internal/sshclient"
	"goxterm-cli/internal/store"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all credentials",
	Run: func(cmd *cobra.Command, args []string) {
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

		if len(db.SshSessions) == 0 {
			fmt.Println("No credentials found in the store.")
			os.Exit(0)
		}

		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}:",
			Active:   "🟢 {{ .Name | cyan }}",
			Inactive: "🟤 {{ .Name | cyan }}",
			Selected: "- CREDENTIAL: {{ .Name | cyan }}",
			Details: `
--------- SSH ----------
{{ "Name:" | faint }} {{ .Name }}
{{ "Host:" | faint }} {{ .Host }}
{{ "Port:" | faint }} {{ .Port }}
{{ "User:" | faint }} {{ .User }}
			`,
		}

		searcher := func(input string, index int) bool {
			credential := db.SshSessions[index]
			name := strings.ReplaceAll(strings.ToLower(credential.Name), " ", "")
			input = strings.ReplaceAll(strings.ToLower(input), " ", "")

			return strings.Contains(name, input)
		}

		prompt := promptui.Select{
			Label:     "Credentials",
			Items:     db.SshSessions,
			Templates: templates,
			Size:      4,
			Searcher:  searcher,
		}

		i, _, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		sshclient.ConnectAndRun(db.SshSessions[i])
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
