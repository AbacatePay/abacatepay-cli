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
	pixSendAmount      int
	pixSendKeyType     string
	pixSendKey         string
	pixSendExternalId  string
	pixSendOtp         string
	pixSendDescription string
)

var pixSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a PIX to any key",
	RunE: func(cmd *cobra.Command, args []string) error {
		if pixSendOtp == "" {
			return fmt.Errorf("a flag --otp é obrigatória")
		}
		if pixSendExternalId == "" {
			return fmt.Errorf("a flag --external-id é obrigatória")
		}
		if pixSendAmount < 100 {
			return fmt.Errorf("o valor mínimo para envio de PIX é de R$ 1,00 (100 centavos)")
		}

		validPixTypes := map[string]bool{
			"EMAIL":  true,
			"CPF":    true,
			"CNPJ":   true,
			"PHONE":  true,
			"RANDOM": true,
		}
		if !validPixTypes[strings.ToUpper(pixSendKeyType)] {
			return fmt.Errorf("tipo de chave PIX inválido. Valores aceitos: EMAIL, CPF, CNPJ, PHONE, RANDOM")
		}

		if pixSendKey == "" {
			return fmt.Errorf("a chave PIX (--pix-key) não pode ser vazia")
		}

		deps, err := utils.SetupClient(Local, Verbose)
		if err != nil {
			return err
		}

		payload := map[string]interface{}{
			"externalId":  pixSendExternalId,
			"amount":      pixSendAmount,
			"pixKeyType":  strings.ToUpper(pixSendKeyType),
			"pixKey":      pixSendKey,
			"otp":         pixSendOtp,
		}
		if pixSendDescription != "" {
			payload["description"] = pixSendDescription
		}

		// Utilizando o endpoint HTTP autenticado no relay
		endpoint := deps.Config.APIBaseURL + "/cli/pix/send"
		
		fmt.Println("Enviando PIX...")
		
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
				style.PrintError(fmt.Sprintf("Erro ao enviar PIX: %s", errResp.Message))
				return fmt.Errorf("falha no envio de PIX: %s", errResp.Message)
			}
			style.PrintError(fmt.Sprintf("Erro ao enviar PIX (Status %d): %s", resp.StatusCode(), string(resp.Body())))
			return fmt.Errorf("falha no envio de PIX (Status %d)", resp.StatusCode())
		}

		style.PrintSuccess("PIX enviado com sucesso!", map[string]string{
			"Valor":      fmt.Sprintf("R$ %.2f", float64(pixSendAmount)/100),
			"Chave":      pixSendKey,
			"Destino":    pixSendKeyType,
			"External ID": pixSendExternalId,
		})
		return nil
	},
}

func init() {
	pixSendCmd.Flags().IntVar(&pixSendAmount, "amount", 0, "Valor do PIX em centavos (mínimo 100)")
	pixSendCmd.Flags().StringVar(&pixSendKeyType, "pix-key-type", "", "Tipo de chave PIX (EMAIL, CPF, CNPJ, PHONE, RANDOM)")
	pixSendCmd.Flags().StringVar(&pixSendKey, "pix-key", "", "Chave PIX de destino")
	pixSendCmd.Flags().StringVar(&pixSendExternalId, "external-id", "", "ID externo da transação (obrigatório)")
	pixSendCmd.Flags().StringVar(&pixSendOtp, "otp", "", "Código OTP de confirmação (obrigatório)")
	pixSendCmd.Flags().StringVar(&pixSendDescription, "description", "", "Descrição opcional do envio")

	_ = pixSendCmd.MarkFlagRequired("amount")
	_ = pixSendCmd.MarkFlagRequired("pix-key-type")
	_ = pixSendCmd.MarkFlagRequired("pix-key")
	_ = pixSendCmd.MarkFlagRequired("external-id")
	_ = pixSendCmd.MarkFlagRequired("otp")

	pixCmd.AddCommand(pixSendCmd)
}
