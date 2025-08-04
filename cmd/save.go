package cmd

import (
	"bufio"
	"fmt"
	"goxterm-cli/internal/config"
	"goxterm-cli/internal/store"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// saveCmd represents the open command
var saveCmd = &cobra.Command{
	Use:   "save [connection]",
	Short: "Update or add a new connection if it doesn't exist",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if len(args) < 1 {
			fmt.Println("Please provide a connection string.")
			fmt.Println("Example: goxterm open user@host:port")
			fmt.Println("")
			cmd.Help()
			os.Exit(1)
		}

		connString := args[0]

		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		if !store.Exists(cfg.StorePath) {
			store.Initialize(cfg.StorePath)
		}

		db, err := store.Load(cfg.StorePath)
		if err != nil {
			fmt.Printf("Error loading store: %v\n", err)
			os.Exit(1)
		}

		if name == "" {
			fmt.Print("Name: ")
			nameFromUser, _ := reader.ReadString('\n')
			name = strings.TrimSpace(nameFromUser)
		}

		if password == "" {
			fmt.Print("SSH Password: ")
			passwordFromUser, _ := reader.ReadString('\n')
			password = strings.TrimSpace(passwordFromUser)
		}

		split1 := strings.Split(connString, "@")
		if len(split1) != 2 || split1[0] == "" || split1[1] == "" {
			fmt.Println("Invalid connection string format. Use 'user@host:port'.")
			os.Exit(1)
		}

		user := split1[0]
		host := ""
		port := 22 // Default SSH port

		split2 := strings.Split(split1[1], ":")
		if len(split2) == 2 {
			port, err = strconv.Atoi(split2[1])
			if err != nil {
				fmt.Printf("Invalid port number: %s\n", split2[1])
				os.Exit(1)
			}
		}

		host = split2[0]

		credential := store.SshSession{
			Session: store.Session{
				Name: name,
			},
			Host:     host,
			Port:     port,
			User:     user,
			Password: password,
		}

		db.AddSshSession(credential)

		if err := store.Save(cfg.StorePath, &db); err != nil {
			fmt.Printf("Error saving store: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Connection '%s' added successfully.\n", name)

	},
}

func init() {
	saveCmd.Flags().StringVarP(&name, "name", "n", "", "Enter the name of the connection")
	saveCmd.Flags().StringVarP(&password, "password", "p", "", "Enter the password")
	rootCmd.AddCommand(saveCmd)
}
