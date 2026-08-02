package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/AbacatePay/abacatepay-cli/internal/style"
	"github.com/AbacatePay/abacatepay-cli/internal/utils"
	"github.com/spf13/cobra"
)

var (
	withdrawAmount      int
	withdrawPixKeyType  string
	withdrawPixKey      string
	withdrawOtp         string
	withdrawDescription string
)

var withdrawCmd = &cobra.Command{
	Use:   "withdraw",
	Short: "Request a payout/withdrawal via PIX",
	RunE: func(cmd *cobra.Command, args []string) error {
		if withdrawOtp == "" {
			return fmt.Errorf("a flag --otp é obrigatória")
		}

		if withdrawAmount < 350 {
			return fmt.Errorf("o valor mínimo para saque é de R$ 3,50 (350 centavos)")
		}

		validPixTypes := map[string]bool{
			"EMAIL":  true,
			"CPF":    true,
			"CNPJ":   true,
			"PHONE":  true,
			"RANDOM": true,
		}
		if !validPixTypes[strings.ToUpper(withdrawPixKeyType)] {
			return fmt.Errorf("tipo de chave PIX inválido. Valores aceitos: EMAIL, CPF, CNPJ, PHONE, RANDOM")
		}

		if withdrawPixKey == "" {
			return fmt.Errorf("a chave PIX (--pix-key) não pode ser vazia")
		}

		deps, err := utils.SetupClient(Local, Verbose)
		if err != nil {
			return err
		}

		payload := map[string]interface{}{
			"amount":      withdrawAmount,
			"pixKeyType":  strings.ToUpper(withdrawPixKeyType),
			"pixKey":      withdrawPixKey,
			"otp":         withdrawOtp,
		}
		if withdrawDescription != "" {
			payload["description"] = withdrawDescription
		}

		// Utilizando o endpoint HTTP autenticado no relay (api.abacatepay.com/cli/withdraw)
		endpoint := deps.Config.APIBaseURL + "/cli/withdraw"
		
		fmt.Println("Solicitando saque...")
		
		resp, err := deps.Client.R().
			SetBody(payload).
			Post(endpoint)

		if err != nil {
			style.PrintError(fmt.Sprintf("Falha na comunicação com a API: %v", err))
			return err
		}

		if resp.IsError() {
			var errResp struct {
				Error   string `json:"error"`
				Message string `json:"message"`
			}
			if err := json.Unmarshal(resp.Body(), &errResp); err == nil && errResp.Message != "" {
				style.PrintError(fmt.Sprintf("Erro ao solicitar saque: %s", errResp.Message))
				return fmt.Errorf("falha no saque: %s", errResp.Message)
			}
			style.PrintError(fmt.Sprintf("Erro ao solicitar saque (Status %d): %s", resp.StatusCode(), string(resp.Body())))
			return fmt.Errorf("falha no saque (Status %d)", resp.StatusCode())
		}

		style.PrintSuccess("Saque solicitado com sucesso!", map[string]string{
			"Valor":    fmt.Sprintf("R$ %.2f", float64(withdrawAmount)/100),
			"Chave":    withdrawPixKey,
			"Destino":  withdrawPixKeyType,
		})
		return nil
	},
}

func init() {
	withdrawCmd.Flags().IntVar(&withdrawAmount, "amount", 0, "Valor do saque em centavos (mínimo 350)")
	withdrawCmd.Flags().StringVar(&withdrawPixKeyType, "pix-key-type", "", "Tipo de chave PIX (EMAIL, CPF, CNPJ, PHONE, RANDOM)")
	withdrawCmd.Flags().StringVar(&withdrawPixKey, "pix-key", "", "Chave PIX de destino")
	withdrawCmd.Flags().StringVar(&withdrawOtp, "otp", "", "Código OTP de confirmação (obrigatório)")
	withdrawCmd.Flags().StringVar(&withdrawDescription, "description", "", "Descrição opcional do saque")

	_ = withdrawCmd.MarkFlagRequired("amount")
	_ = withdrawCmd.MarkFlagRequired("pix-key-type")
	_ = withdrawCmd.MarkFlagRequired("pix-key")
	_ = withdrawCmd.MarkFlagRequired("otp")

	rootCmd.AddCommand(withdrawCmd)
}
