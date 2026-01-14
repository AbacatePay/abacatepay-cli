package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:     "profile",
	Aliases: []string{"profiles"},
	Short:   "Manage saved authentication profiles",
	Long:    "Allows listing, switching, and removing AbacatePay access profiles configured on this machine.",
}

var listProfilesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listProfiles()
	},
}

var deleteProfileCmd = &cobra.Command{
	Use:     "delete [name]",
	Aliases: []string{"remove", "rm"},
	Short:   "Remove a profile",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return deleteProfile(args[0])
	},
}

func init() {
	profileCmd.AddCommand(listProfilesCmd)
	profileCmd.AddCommand(deleteProfileCmd)

	rootCmd.AddCommand(profileCmd)
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
		fmt.Println("No profiles found. Use 'abacatepay-cli login' to create one.")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "  NAME\tSTATUS")
	for _, p := range profiles {
		status := ""
		prefix := "  "

		if p == active {
			status = "(active)"
			prefix = "* "
		}
		fmt.Fprintf(w, "%s%s\t%s\n", prefix, p, status)
	}
	w.Flush()

	return nil
}

func deleteProfile(name string) error {
	deps := utils.SetupDependencies(Local, Verbose)

	token, err := deps.Store.GetNamed(name)
	if err != nil {
		return fmt.Errorf("error verifying profile: %w", err)
	}
	if token == "" {
		return fmt.Errorf("profile '%s' not found", name)
	}

	active, _ := deps.Store.GetActiveProfile()
	if active == name {
		return fmt.Errorf("cannot delete the active profile. Switch to another profile first using 'profile use'")
	}

	if err := deps.Store.DeleteNamed(name); err != nil {
		return fmt.Errorf("error deleting profile: %w", err)
	}

	fmt.Printf("Profile '%s' successfully removed.\n", name)
	return nil
}
