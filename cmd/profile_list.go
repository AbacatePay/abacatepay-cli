package cmd

import (
	"fmt"

	"abacatepay-cli/internal/output"
	"abacatepay-cli/internal/style"
	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var listProfilesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listProfiles()
	},
}

func init() {
	profileCmd.AddCommand(listProfilesCmd)
}

func listProfiles() error {
	deps := utils.SetupDependencies(Local, Verbose)

	profiles, err := deps.Store.List()
	if err != nil {
		return fmt.Errorf("error listing profiles: %w", err)
	}

	active, err := deps.Store.GetActiveProfile()
	if err != nil {
		active = ""
	}

	if len(profiles) == 0 {
		output.Print(output.Result{
			Title: "No profiles found",
			Fields: map[string]string{
				"Hint": "Use 'abacatepay login' to create one.",
			},
			Data: map[string]any{
				"profiles": []string{},
				"active":   "",
			},
		})
		return nil
	}

	profileData := make([]map[string]any, 0, len(profiles))
	rows := make([][]string, 0, len(profiles))

	for _, p := range profiles {
		token, _ := deps.Store.GetNamed(p)
		isActive := p == active

		shortKey := ""
		if token != "" && len(token) > 10 {
			shortKey = token[:10] + "..."
		}

		activeMarker := ""
		if isActive {
			activeMarker = "Yes"
		}

		rows = append(rows, []string{p, shortKey, activeMarker})
		profileData = append(profileData, map[string]any{
			"name":   p,
			"active": isActive,
		})
	}

	if output.GetFormat() == output.FormatText {
		profileMap := make(map[string]string)
		for _, p := range profiles {
			token, _ := deps.Store.GetNamed(p)
			profileMap[p] = token
		}
		style.ProfileSimpleList(profileMap, active)
		return nil
	}

	output.Print(output.Result{
		Title:   "Profiles",
		Headers: []string{"Name", "API Key", "Active"},
		Rows:    rows,
		Data: map[string]any{
			"profiles": profileData,
			"active":   active,
		},
	})

	return nil
}
