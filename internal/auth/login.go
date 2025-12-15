package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/go-resty/resty/v2"

	"abacatepay-cli/internal/config"
)

type DeviceLoginResponse struct {
	VerificationURI string `json:"verificationUri"`
	DeviceCode      string `json:"deviceCode"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func Login(ctx context.Context, cfg *config.Config, client *resty.Client, store TokenStore) error {
	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}

	var result DeviceLoginResponse
	resp, err := client.R().
		SetContext(ctx).
		SetBody(map[string]string{"host": host}).
		SetResult(&result).
		Post(cfg.APIBaseURL + "/device-login")

	if err != nil {
		return fmt.Errorf("falha na requisição de login: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("login falhou com status %d", resp.StatusCode())
	}

	if result.VerificationURI == "" {
		return fmt.Errorf("resposta inválida: URI de verificação ausente")
	}

	if result.DeviceCode == "" {
		return fmt.Errorf("resposta inválida: código do dispositivo ausente")
	}

	fmt.Println("Abra o seguinte link no navegador para autenticar:")
	fmt.Printf("%s\n", result.VerificationURI)

	token, err := pollForToken(ctx, cfg, client, result.DeviceCode)
	if err != nil {
		return err
	}

	if err := store.Save(token); err != nil {
		return fmt.Errorf("falha ao salvar token: %w", err)
	}

	slog.Info("Login realizado com sucesso")
	return nil
}

func Logout(store TokenStore) error {
	if err := store.Delete(); err != nil {
		return err
	}

	slog.Info("Logout realizado com sucesso")
	return nil
}

func pollForToken(ctx context.Context, cfg *config.Config, client *resty.Client, deviceCode string) (string, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Aguardando autorização..."
	s.Start()
	defer s.Stop()

	ticker := time.NewTicker(cfg.PollInterval)
	defer ticker.Stop()

	for retries := 0; retries < cfg.MaxRetries; retries++ {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-ticker.C:
		}

		var result TokenResponse
		resp, err := client.R().
			SetContext(ctx).
			SetBody(map[string]string{"deviceCode": deviceCode}).
			SetResult(&result).
			Post(cfg.APIBaseURL + "/token")

		if err != nil {
			slog.Debug("Falha na requisição de token", "error", err)
			continue
		}

		switch resp.StatusCode() {
		case http.StatusOK:
			if result.Token != "" {
				return result.Token, nil
			}
		case http.StatusAccepted:
			continue
		case http.StatusBadRequest:
			return "", fmt.Errorf("código do dispositivo inválido")
		case http.StatusUnauthorized:
			return "", fmt.Errorf("não autorizado")
		case http.StatusInternalServerError:
			slog.Warn("Erro no servidor, tentando novamente...")
			continue
		default:
			slog.Debug("Status inesperado", "status", resp.StatusCode())
		}
	}

	return "", fmt.Errorf("tempo de autorização expirado. Tente novamente")
}
