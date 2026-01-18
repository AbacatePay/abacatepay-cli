// Package utils...
package utils

import (
	"context"
	"fmt"
	"log/slog"

	"abacatepay-cli/internal/auth"
	"abacatepay-cli/internal/client"
	"abacatepay-cli/internal/config"
	"abacatepay-cli/internal/logger"
	"abacatepay-cli/internal/store"

	"github.com/go-resty/resty/v2"
)

type StartListenerParams struct {
	Context    context.Context
	Config     *config.Config
	Client     *resty.Client
	Store      store.TokenStore
	Token      string
	ForwardURL string
	Version    string
	Mock       bool
}

type Dependencies struct {
	Config *config.Config
	Client *resty.Client
	Store  store.TokenStore
}

func SetupTransactionLogger() (*slog.Logger, error) {
	logCfg, err := logger.DefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to configure logger: %w", err)
	}

	return logger.NewTransactionLogger(logCfg)
}

func GetConfig(local bool) *config.Config {
	if !local {
		return config.Default()
	}

	return config.Local()
}

func GetStore(cfg *config.Config) store.TokenStore {
	return store.NewKeyringStore(cfg.ServiceName, cfg.TokenKey)
}

func SetupDependencies(local bool, verbose bool) *Dependencies {
	cfg := GetConfig(local)
	cfg.Verbose = verbose

	cli := client.New(cfg)
	store := GetStore(cfg)

	return &Dependencies{
		Config: cfg,
		Client: cli,
		Store:  store,
	}
}

func SetupClient(local, verbose bool) (*Dependencies, error) {
	if !IsOnline() {
		return nil, fmt.Errorf("you’re offline — check your connection and try again")
	}

	deps := SetupDependencies(local, verbose)
	activeProfile, err := deps.Store.GetActiveProfile()
	if err != nil {
		return nil, fmt.Errorf("failed to get active profile: %w", err)
	}

	token, err := deps.Store.GetNamed(activeProfile)
	if err != nil || token == "" {
		return nil, fmt.Errorf("token not found for active profile: %s", activeProfile)
	}

	_, err = auth.ValidateToken(deps.Client, deps.Config.APIBaseURL, token)
	if err != nil {
		return nil, fmt.Errorf("session expired for profile %s: %w", activeProfile, err)
	}

	deps.Client.SetAuthToken(token)
	return deps, nil
}
