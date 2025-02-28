package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.Run("Cria arquivo com valores padrão", func(t *testing.T) {
		// Configure temporary directory
		dir := t.TempDir()
		os.Setenv("HOME", dir)
		defer os.Unsetenv("HOME")

		// Test initialization
		err := Init()
		assert.NoError(t, err)

		// Verify if file was created
		cfgPath := filepath.Join(dir, ".abacatepay.toml")
		_, err = os.Stat(cfgPath)
		assert.NoError(t, err)

		// Verify default values
		assert.Equal(t, "https://api.abacatepay.com", viper.GetString("endpoints.api"))
		assert.Equal(t, "http://localhost:8080/webhook", viper.GetString("endpoints.forward_url"))
	})
}

// Test if the configuration exists and loads correctly
func TestGet(t *testing.T) {
	t.Run("Load existing configuration", func(t *testing.T) {
		// Setup
		dir := t.TempDir()
		os.Setenv("HOME", dir)
		viper.Reset()

		// Create manual file
		cfgContent := `
			token = "test_token"
			live_mode = true

			[endpoints]
			api = "https://custom.api"
			websocket = "wss://custom.ws"
			forward_url = "http://localhost:8080/webhook"
			`
		cfgPath := filepath.Join(dir, ".abacatepay.toml")
		os.WriteFile(cfgPath, []byte(cfgContent), 0644)

		// Test
		err := Init()
		assert.NoError(t, err)

		cfg, err := Get()
		assert.NoError(t, err)

		// Verifications
		assert.Equal(t, "test_token", cfg.Token)
		assert.True(t, cfg.LiveMode)
		assert.Equal(t, "https://custom.api", cfg.Endpoints.API)
		assert.Equal(t, "wss://custom.ws", cfg.Endpoints.WebSocket)
		assert.Equal(t, "http://localhost:8080/webhook", cfg.Endpoints.ForwardURL)
	})
}

func TestSave(t *testing.T) {
	t.Run("Persist changes in file", func(t *testing.T) {
		// Setup
		dir := t.TempDir()
		os.Setenv("HOME", dir)
		viper.Reset()

		err := Init()
		assert.NoError(t, err)

		// Modify configuration
		newCfg := &Config{
			Token:    "new_token",
			LiveMode: true,
		}
		newCfg.Endpoints.API = "https://new.api"
		newCfg.Endpoints.WebSocket = "wss://new.ws"
		newCfg.Endpoints.ForwardURL = "http://localhost:8080/new_webhook"

		// Save
		err = Save(newCfg)
		assert.NoError(t, err)

		// Reload to verify
		viper.Reset()
		err = Init()
		assert.NoError(t, err)

		cfg, err := Get()
		assert.NoError(t, err)

		assert.Equal(t, "new_token", cfg.Token)
		assert.Equal(t, "https://new.api", cfg.Endpoints.API)
		assert.Equal(t, "wss://new.ws", cfg.Endpoints.WebSocket)
		assert.Equal(t, "http://localhost:8080/new_webhook", cfg.Endpoints.ForwardURL)
	})
}
