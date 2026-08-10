package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/AbacatePay/abacatepay-cli/internal/output"
	"github.com/AbacatePay/abacatepay-cli/internal/version"
	"github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the CLI to the latest version, if one is available",
	RunE: func(cmd *cobra.Command, args []string) error {
		return upgrade()
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

func upgrade() error {
	ctx := context.Background()

	return output.RunTask("Checking for updates...", func() (output.Result, error) {
		latest, found, err := version.CheckUpdate(ctx, version.Version)
		if err != nil {
			return output.Result{}, fmt.Errorf("could not check for updates: %w", err)
		}

		if !found {
			return output.Result{
				Title:  "Already up to date",
				Fields: map[string]string{"Version": currentVersionLabel()},
			}, nil
		}

		exe, err := os.Executable()
		if err != nil {
			return output.Result{}, fmt.Errorf("could not resolve executable path: %w", err)
		}

		if err := selfupdate.UpdateTo(ctx, latest.AssetURL, latest.AssetName, exe); err != nil {
			return output.Result{}, fmt.Errorf("update failed: %w", err)
		}

		return output.Result{
			Title:  "Update complete ✨",
			Fields: map[string]string{"Version": latest.Version()},
		}, nil
	})
}

func currentVersionLabel() string {
	if version.Version == "" {
		return "dev"
	}
	return version.Version
}
