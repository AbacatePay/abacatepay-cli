package payments

import "github.com/AbacatePay/abacatepay-cli/internal/utils"

func ExecutePaymentAction(local, verbose bool, action func(*Service) error) error {
	deps, err := utils.SetupClient(local, verbose)
	if err != nil {
		return err
	}

	pixService := New(deps.Client, deps.Config.APIBaseURL, verbose)
	return action(pixService)
}
