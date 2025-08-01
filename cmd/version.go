package cmd

import (
	"fmt"
	"goxterm-cli/internal/constants"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(constants.AppName, "version", constants.AppVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
