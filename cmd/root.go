package cmd

import (
	"fmt"
	"goxterm-cli/internal/config"
	"os"

	"github.com/spf13/cobra"
)

var name string
var password string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goxterm",
	Short: "Efficiently manage your SSH connections in one place using your favorite terminal.",
	Run: func(cmd *cobra.Command, args []string) {
		if !config.Exists() {
			fmt.Println("No configuration found. Please run 'goxterm setup' to initialize or goxterm --help.")
			os.Exit(1)
		}

		if len(args) < 1 {
			cmd.Help()
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goxterm-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
