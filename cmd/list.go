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

type SshCredentialWithName struct {
	Name string
	store.SshCredential
}

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

		if len(db.Credentials) == 0 {
			fmt.Println("No credentials found in the store.")
			os.Exit(0)
		}

		var list []SshCredentialWithName

		for key, value := range db.Credentials {
			list = append(list, SshCredentialWithName{
				Name: key,
				SshCredential: store.SshCredential{
					Host:     value.Host,
					Port:     value.Port,
					User:     value.User,
					Password: value.Password,
				},
			})
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
			credential := list[index]
			name := strings.ReplaceAll(strings.ToLower(credential.Name), " ", "")
			input = strings.ReplaceAll(strings.ToLower(input), " ", "")

			return strings.Contains(name, input)
		}

		prompt := promptui.Select{
			Label:     "Credentials",
			Items:     list,
			Templates: templates,
			Size:      4,
			Searcher:  searcher,
		}

		i, _, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		sshclient.ConnectAndRun(store.SshCredential{
			Host:     list[i].Host,
			Port:     list[i].Port,
			User:     list[i].User,
			Password: list[i].Password,
		})
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
