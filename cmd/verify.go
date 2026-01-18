package cmd

import (
	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var verifySecret, verifySignature, verifyPayload string

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Args:  cobra.ExactArgs(0),
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return verify()
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}

func verify() error {
	deps, err := utils.SetupClient(Local, Verbose)
	if err != nil {
		return err
	}

	return nil
}
