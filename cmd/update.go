package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update CLI to new version if is available",
	RunE: func(cmd *cobra.Command, args []string) error {
		return update()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func update() error {
	const slug string = "AbacatePay/abacatepay-cli"

	latest, found, err := selfupdate.DetectLatest(context.Background(), selfupdate.ParseSlug(slug))
	if err != nil {
		return fmt.Errorf("error to update cli's version: %w", err)
	}

	if !found || latest.LessOrEqual(rootCmd.Version) {
		return nil
	}

	exe, _ := os.Executable()

	return selfupdate.UpdateTo(context.Background(), latest.AssetURL, latest.AssetName, exe)
}
