package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var port int

var rootCmd = &cobra.Command{
	Use:   "goxterm",
	Short: "Efficiently manage your SSH connections in one place using your favorite terminal.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			os.Exit(1)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
