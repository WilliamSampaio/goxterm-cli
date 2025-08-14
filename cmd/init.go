package cmd

import (
	"bufio"
	"fmt"
	"goxterm-cli/internal/config"
	"goxterm-cli/internal/constants"
	"goxterm-cli/internal/database"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var overwrite bool
var noMigrate bool

// initCmd represents the setup command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the application",
	Run: func(cmd *cobra.Command, args []string) {
		Initialize()
	},
}

func init() {
	initCmd.Flags().BoolVarP(&overwrite, "overwrite", "O", false, "Overwrite on")
	initCmd.Flags().BoolVarP(&noMigrate, "no-migrate", "X", false, "Ignore migration")
	rootCmd.AddCommand(initCmd)
}

func Initialize() {
	reader := bufio.NewReader(os.Stdin)

	if config.Exists() && !overwrite {
		fmt.Println("Configuration already exists. Use 'goxterm init --overwrite' to overwrite.")
		return
	}

	prompt := promptui.Select{
		Label: "Select Store Type",
		Items: []string{"sqlite3", "bbolt (not implemented)"},
	}

	_, storeType, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	defaultStoreDir := config.ConfigDir()

	fmt.Print("Store Dir (default ", defaultStoreDir, "): ")
	storeDir, _ := reader.ReadString('\n')
	storeDir = strings.TrimSpace(storeDir)

	if storeDir == "" {
		storeDir = defaultStoreDir
	}

	cfg := config.Config{
		Version:   constants.AppVersion,
		StoreType: storeType,
		StorePath: filepath.Join(storeDir, "database."+storeType),
	}

	if err := config.Save(cfg); err != nil {
		fmt.Printf("Failed to save configuration: %v\n", err)
		os.Exit(1)
	}

	if !noMigrate {
		if err := database.RunMigration(database.PerformDown, 0); err != nil {
			fmt.Printf("Failed to run migrations: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("Configuration saved successfully.")
	os.Exit(0)
}
