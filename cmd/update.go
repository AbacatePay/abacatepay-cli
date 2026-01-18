package cmd

import (
	"context"
	"fmt"
	"os"

	"abacatepay-cli/internal/version"

	"github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upgrade"},
	Short:   "Update CLI to new version if it is available",
	RunE: func(cmd *cobra.Command, args []string) error {
		return update()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func update() error {
	ctx := context.Background()

	latest, found, err := version.CheckUpdate(ctx, rootCmd.Version)
	if err != nil {
		return fmt.Errorf("couldn’t check for updates: %w", err)
	}

	if !found {
		fmt.Printf("You’re already on the latest version (%s)\n", rootCmd.Version)
		return nil
	}

	fmt.Printf("Update available: %s\n", latest.Version())
	fmt.Println("Downloading and installing...")

	exe, _ := os.Executable()

	if err := selfupdate.UpdateTo(ctx, latest.AssetURL, latest.AssetName, exe); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	fmt.Println("Update complete ✨")

	return nil
}
