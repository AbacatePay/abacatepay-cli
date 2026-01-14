package cmd

import (
	"abacatepay-cli/internal/utils"
	"fmt"
	"os"
	"text/tabwriter"

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
		fmt.Println("No profiles found. Use 'abacatepay login' to create one.")

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
