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
	"abacatepay-cli/internal/store"
	"abacatepay-cli/internal/style"

	"abacatepay-cli/internal/types"
)

type LoginParams struct {
	Config      *config.Config
	Client      *resty.Client
	Store       store.TokenStore
	Context     context.Context
	APIKey      string
	ProfileName string
	OpenBrowser func(string) error
}

func Login(params *LoginParams) error {
	if params.APIKey != "" {
		return saveAPIKey(params)
	}

	// For local mode, mock the token
	if params.Config.APIBaseURL == "http://191.252.202.128:8080" {
		params.APIKey = "mock_token_local_dev"
		return saveAPIKey(params)
	}

	return deviceCodeFlow(params)
}

	if params.APIKey != "" {
		return loginWithAPIKey(params)
	}

	return loginWithDeviceFlow(params)
}

func loginWithAPIKey(params *LoginParams) error {
	user, err := ValidateToken(params.Client, params.Config.APIBaseURL, params.APIKey)
	if err != nil {
		return err
	}

	if err := saveAndActivateProfile(params.Store, params.ProfileName, params.APIKey); err != nil {
		return err
	}

	fmt.Printf("Welcome back, %s\nProfile: %s\n", user.Name, params.ProfileName)

	return nil
}

func loginWithDeviceFlow(params *LoginParams) error {
	host := getHostname()

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

	if !tryOpenBrowser(params, result.VerificationURI) {
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

	if err := saveAndActivateProfile(params.Store, params.ProfileName, token); err != nil {
		return err
	}

	fmt.Printf("Authenticated as %s\nProfile: %s\n", user.Name, params.ProfileName)

	return nil
}

func saveAndActivateProfile(st store.TokenStore, profile, token string) error {
	existingToken, _ := st.GetNamed(profile)
	if existingToken != "" {
		slog.Info("Updating existing profile", "name", profile)
	}

	if err := st.SaveNamed(profile, token); err != nil {
		return fmt.Errorf("failed to store API key: %w", err)
	}

	if err := st.SetActiveProfile(profile); err != nil {
		return fmt.Errorf("failed to activate profile: %w", err)
	}

	return nil
}

func tryOpenBrowser(params *LoginParams, verificationURI string) bool {
	if params.OpenBrowser == nil {
		return false
	}

	if err := params.OpenBrowser(verificationURI); err != nil {
		slog.Debug("Unable to open browser automatically", "error", err)

		return false
	}

	fmt.Printf("Opening browser: %s\n", verificationURI)

	return true
}

func getHostname() string {
	host, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return host
}

func Logout(st store.TokenStore) (string, error) {
	activeProfile, err := st.GetActiveProfile()
	if err != nil {
		return "", fmt.Errorf("failed to get active profile: %w", err)
	}

	if activeProfile == "" {
		return "", fmt.Errorf("no active profile found")
	}

	if err := st.DeleteNamed(activeProfile); err != nil {
		return "", fmt.Errorf("failed to delete token: %w", err)
	}

	profiles, _ := st.List()
	if len(profiles) > 0 {
		_ = st.SetActiveProfile(profiles[0])

		slog.Info("Signed out", "profile", activeProfile, "switched_to", profiles[0])
		return activeProfile, nil
	}

	_ = st.SetActiveProfile("")
	slog.Info("Signed out", "profile", activeProfile)

	return activeProfile, nil
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
