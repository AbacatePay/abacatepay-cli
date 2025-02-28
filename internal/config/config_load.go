package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

const (
	configName = ".abacatepay"
	configType = "toml"
)

// Structure of the TOML file
// -> LiveMode is a boolean that determines if the app is in live mode
type Config struct {
	Token     string `mapstructure:"token"`
	LiveMode  bool   `mapstructure:"live_mode"`
	Endpoints struct {
		API        string `mapstructure:"api"`
		WebSocket  string `mapstructure:"websocket"`
		ForwardURL string `mapstructure:"forward_url"`
	} `mapstructure:"endpoints"`
}

// Initialize the Viper
func Init() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	viper.AddConfigPath(home)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	// Default values
	viper.SetDefault("endpoints.api", "https://api.abacatepay.com")
	viper.SetDefault("endpoints.websocket", "wss://ws.abacatepay.com/ws")
	viper.SetDefault("endpoints.forward_url", "http://localhost:8080/webhook")

	configPath := filepath.Join(home, configName+"."+configType)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := viper.WriteConfigAs(configPath); err != nil {
			return fmt.Errorf("failed to create config file: %v", err)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}

	return nil
}

func Get() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	// Convert the struct to a map[string]interface{}
	var rawConfig map[string]interface{}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "mapstructure",
		Result:  &rawConfig,
	})
	if err != nil {
		return fmt.Errorf("failed to create decoder: %v", err)
	}
	if err := decoder.Decode(cfg); err != nil {
		return fmt.Errorf("failed to decode config: %v", err)
	}

	// Set all values in Viper recursively
	setViperKeys("", rawConfig)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %v", err)
	}
	return nil
}

// Set all values in Viper recursively
func setViperKeys(prefix string, data map[string]interface{}) {
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		// If the value is a map, process recursively
		if nested, ok := value.(map[string]interface{}); ok {
			setViperKeys(fullKey, nested)
		} else {
			viper.Set(fullKey, value)
		}
	}
}
