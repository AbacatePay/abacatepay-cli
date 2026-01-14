package cmd

import (
	"fmt"

	"abacatepay-cli/internal/auth"
	"abacatepay-cli/internal/mock"
	"abacatepay-cli/internal/payments/pix"
	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var createInteractive bool

var createPaymentCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new payment charge",
	Long: `Create a new payment charge.
By default, creates a mock payment with random data.
Use -i to enter interactive mode and specify details.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return createPayment(cmd)
	},
}

func init() {
	paymentsCmd.AddCommand(createPaymentCmd)

	createPaymentCmd.Flags().BoolVarP(&createInteractive, "interactive", "i", false, "Enable interactive mode")
}

func createPayment(cmd *cobra.Command) error {
	if !utils.IsOnline() {
		return fmt.Errorf("you’re offline — check your connection and try again")
	}
	deps := utils.SetupDependencies(Local, Verbose)
	activeProfile, err := deps.Store.GetActiveProfile()
	if err != nil {
		return fmt.Errorf("failed to get active profile: %w", err)
	}

	token, err := deps.Store.GetNamed(activeProfile)
	if err != nil || token == "" {
		return fmt.Errorf("token not found for active profile: %s", activeProfile)
	}

	_, err = auth.ValidateToken(deps.Client, deps.Config.APIBaseURL, token)
	if err != nil {
		return fmt.Errorf("session expired for profile %s: %w", activeProfile, err)
	}

	deps.Client.SetAuthToken(token)

	if createInteractive {
		// TODO: Implement interactive mode if needed or delegate
		return fmt.Errorf("interactive mode not yet implemented in this simplified setup")
	}

	b := mock.CreatePixQRCodeMock()

	params := &pix.CreatePixQRCodeParams{
		Client:  deps.Client,
		Token:   token,
		Body:    b,
		BaseURL: deps.Config.APIBaseURL,
	}
	salve := pix.CreateQRCode(params)

	return salve
}

