package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-resty/resty/v2"

	"abacatepay-cli/internal/config"
	"abacatepay-cli/internal/style"

	//"abacatepay-cli/internal/style"
	"abacatepay-cli/internal/types"
)

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

		existingToken, _ := params.Store.GetNamed(profile)
		if existingToken != "" {
			slog.Info("Updating existing profile", "name", profile)
		}

		if err := params.Store.SaveNamed(profile, params.APIKey); err != nil {
			return fmt.Errorf("failed to store API key: %w", err)
		}

		if err := params.Store.SetActiveProfile(profile); err != nil {
			return fmt.Errorf("failed to activate profile: %w", err)
		}

		fmt.Printf("Welcome back, %s\nProfile: %s\n", user.Name, profile)

		return nil

	}

	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}

	var result types.DeviceLoginResponse

	resp, err := params.Client.R().
		SetContext(params.Context).
		SetBody(map[string]string{"host": host}).
		SetResult(&result).
		Post(params.Config.APIBaseURL + "/device-login")
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("login failed (status %d)", resp.StatusCode())
	}

	tryOpen := func() bool {
		if params.OpenBrowser == nil {
			return false
		}

		if err := params.OpenBrowser(result.VerificationURI); err != nil {
			slog.Debug("Unable to open browser automatically", "error", err)

			return false
		}

		fmt.Printf("Opening browser: %s\n", result.VerificationURI)

		return true
	}

	if !tryOpen() {
		fmt.Println("Open the link below to continue authentication:")
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

	existingToken, _ := params.Store.GetNamed(profile)
	if existingToken != "" {
		slog.Info("Updating existing profile", "name", profile)
	}

	if err := params.Store.SaveNamed(profile, token); err != nil {
		return fmt.Errorf("failed to store API key: %w", err)
	}

	if err := params.Store.SetActiveProfile(profile); err != nil {
		return fmt.Errorf("failed to activate profile: %w", err)
	}

	fmt.Printf("Authenticated as %s\nProfile: %s\n", user.Name, profile)

	return nil
}

func Logout(store TokenStore) error {
	if err := store.Delete(); err != nil {
		return err
	}

	slog.Info("Signed out")

	return nil
}

func pollForToken(ctx context.Context, cfg *config.Config, client *resty.Client, deviceCode string) (string, error) {
	s := style.Spinner()
	defer s.Stop()

	ticker := time.NewTicker(2 * time.Second)

	defer ticker.Stop()

	for range 150 {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-ticker.C:
		}

		var result types.TokenResponse

		resp, err := client.R().
			SetContext(ctx).
			SetBody(map[string]string{"deviceCode": deviceCode}).
			SetResult(&result).
			Post(cfg.APIBaseURL + "/token")
		if err != nil {
			slog.Debug("Token request failed", "error", err)

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
			return "", fmt.Errorf("invalid device code")
		case http.StatusUnauthorized:
			return "", fmt.Errorf("authorization denied")
		case http.StatusInternalServerError:
			slog.Warn("Server error, retrying...")
			continue
		default:
			slog.Debug("Unexpected response", "status", resp.StatusCode())
		}
	}

	return "", fmt.Errorf("authorization timed out")
}

func retryConnection(cfg *config.Config) {
}
