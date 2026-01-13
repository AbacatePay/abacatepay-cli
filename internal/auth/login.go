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

	DeviceCode string `json:"deviceCode"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type LoginParams struct {
	Config      *config.Config
	Client      *resty.Client
	Store       TokenStore
	Context     context.Context
	APIKey      string
	ProfileName string
	OpenBrowser func(string) error
}

func Login(params *LoginParams) error {
	profile := params.ProfileName

	if profile == "" {
		profile = fmt.Sprintf("profile-%d", time.Now().Unix()%10000)
	}

	if params.APIKey != "" {
		user, err := ValidateToken(params.Client, params.Config.APIBaseURL, params.APIKey)

		if err != nil {
			return err
		}

		if err := params.Store.SaveNamed(profile, params.APIKey); err != nil {
			return fmt.Errorf("falha ao salvar API Key: %w", err)
		}

		if err := params.Store.SetActiveProfile(profile); err != nil {
			return fmt.Errorf("falha ao definir perfil ativo: %w", err)
		}

		fmt.Printf("Bem-vindo, %s! API Key salva no perfil: %s\n", user.Name, profile)
		return nil

	}

	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}

	var result DeviceLoginResponse

	resp, err := params.Client.R().
		SetContext(params.Context).
		SetBody(map[string]string{"host": host}).
		SetResult(&result).
		Post(params.Config.APIBaseURL + "/device-login")

	if err != nil {
		return fmt.Errorf("falha na requisição de login: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("login falhou com status %d", resp.StatusCode())
	}

	tryOpen := func() bool {
		if params.OpenBrowser == nil {
			return false
		}
		if err := params.OpenBrowser(result.VerificationURI); err != nil {
			slog.Debug("Não foi possível abrir o navegador automaticamente", "error", err)
			return false
		}
		fmt.Printf("Tentando abrir o navegador em: %s\n", result.VerificationURI)
		return true
	}

	if !tryOpen() {
		fmt.Println("Por favor, abra o seguinte link no navegador para autenticar:")
		fmt.Printf("%s\n", result.VerificationURI)
	}

	token, err := pollForToken(params.Context, params.Config, params.Client, result.DeviceCode)
	if err != nil {
		return err
	}

	user, err := ValidateToken(params.Client, params.Config.APIBaseURL, token)
	if err != nil {
		return err
	}

	if err := params.Store.SaveNamed(profile, token); err != nil {
		return fmt.Errorf("falha ao salvar token: %w", err)
	}

	if err := params.Store.SetActiveProfile(profile); err != nil {
		return fmt.Errorf("falha ao definir perfil ativo: %w", err)
	}

	fmt.Printf("Bem-vindo, %s! Login realizado no perfil: %s\n", user.Name, profile)
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
