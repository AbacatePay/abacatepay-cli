package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type SimpleConfig struct {
	apiURL      string
	environment string
	debug       bool
	apiKey      string
	appID       string
	timeout     int
}

type SimpleConfigBuilder struct {
	config *SimpleConfig
}

func NewSimpleConfig() *SimpleConfigBuilder {
	return &SimpleConfigBuilder{
		config: &SimpleConfig{
			apiURL:      "https://api.abacatepay.com",
			environment: "production",
			debug:       false,
			timeout:     30,
		},
	}
}

func (b *SimpleConfigBuilder) Environment(env string) *SimpleConfigBuilder {
	env = strings.ToLower(env)
	b.config.environment = env

	switch env {
	case "production", "prod":
		b.config.apiURL = "https://api.abacatepay.com"
		b.config.debug = false
	case "staging", "stage":
		b.config.apiURL = "https://staging-api.abacatepay.com"
		b.config.debug = true
	case "local", "development", "dev":
		b.config.apiURL = "http://191.252.202.128:8080"
		b.config.debug = true
	}

	return b
}

func (b *SimpleConfigBuilder) APIURL(url string) *SimpleConfigBuilder {
	b.config.apiURL = strings.TrimSuffix(url, "/")
	return b
}

func (b *SimpleConfigBuilder) Debug(debug bool) *SimpleConfigBuilder {
	b.config.debug = debug
	return b
}

func (b *SimpleConfigBuilder) APIKey(key string) *SimpleConfigBuilder {
	b.config.apiKey = key
	return b
}

func (b *SimpleConfigBuilder) AppID(id string) *SimpleConfigBuilder {
	b.config.appID = id
	return b
}

func (b *SimpleConfigBuilder) Timeout(seconds int) *SimpleConfigBuilder {
	if seconds > 0 {
		b.config.timeout = seconds
	}
	return b
}

func (b *SimpleConfigBuilder) Build() (*SimpleConfig, error) {
	if err := b.config.validate(); err != nil {
		return nil, err
	}
	return b.config, nil
}

func (b *SimpleConfigBuilder) MustBuild() *SimpleConfig {
	config, err := b.Build()
	if err != nil {
		panic(fmt.Sprintf("SimpleConfig validation failed: %v", err))
	}
	return config
}

func LoadFromEnvironment() *SimpleConfigBuilder {
	builder := NewSimpleConfig()

	envHandlers := map[string]func(string){
		SimpleEnvEnvironment: func(v string) { builder.Environment(v) },
		SimpleEnvAPIKey:      func(v string) { builder.APIKey(v) },
		SimpleEnvAppID:       func(v string) { builder.AppID(v) },
		SimpleEnvAPIURL:      func(v string) { builder.APIURL(v) },
	}

	for envVar, handler := range envHandlers {
		if value := os.Getenv(envVar); value != "" {
			handler(value)
		}
	}

	if debugStr := os.Getenv(SimpleEnvDebug); debugStr != "" {
		if debug, err := strconv.ParseBool(debugStr); err == nil {
			builder.Debug(debug)
		}
	}

	if timeoutStr := os.Getenv(SimpleEnvTimeout); timeoutStr != "" {
		if timeout, err := strconv.Atoi(timeoutStr); err == nil && timeout > 0 {
			builder.Timeout(timeout)
		}
	}

	return builder
}

func SimpleDefault() *SimpleConfig {
	return NewSimpleConfig().
		Environment("production").
		MustBuild()
}

func SimpleLocal() *SimpleConfig {
	return NewSimpleConfig().
		Environment("local").
		MustBuild()
}

func SimpleDebug() *SimpleConfig {
	return NewSimpleConfig().
		Environment("local").
		Debug(true).
		MustBuild()
}

func (c *SimpleConfig) validate() error {
	if c.apiURL == "" {
		return fmt.Errorf("API URL cannot be empty")
	}

	if !strings.HasPrefix(c.apiURL, "http://") && !strings.HasPrefix(c.apiURL, "https://") {
		return fmt.Errorf("API URL must start with http:// or https://")
	}

	if c.timeout <= 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}

	return nil
}

func (c *SimpleConfig) HasAPIKey() bool {
	return c.apiKey != ""
}

func (c *SimpleConfig) HasAppID() bool {
	return c.appID != ""
}

func (c *SimpleConfig) IsProduction() bool {
	return c.environment == "production" || c.environment == "prod"
}

func (c *SimpleConfig) IsStaging() bool {
	return c.environment == "staging" || c.environment == "stage"
}

func (c *SimpleConfig) IsDevelopment() bool {
	return c.environment == "local" || c.environment == "development" || c.environment == "dev"
}

func (c *SimpleConfig) GetEnvironment() string {
	return c.environment
}

func (c *SimpleConfig) GetAPIURL() string {
	return c.apiURL
}

func (c *SimpleConfig) GetDebug() bool {
	return c.debug
}

func (c *SimpleConfig) GetTimeout() int {
	return c.timeout
}

func (c *SimpleConfig) String() string {
	return fmt.Sprintf("SimpleConfig{Environment: %s, APIURL: %s, Debug: %v, Timeout: %ds}",
		c.environment, c.apiURL, c.debug, c.timeout)
}

func (c *SimpleConfig) Clone() *SimpleConfig {
	return &SimpleConfig{
		apiURL:      c.apiURL,
		environment: c.environment,
		debug:       c.debug,
		apiKey:      c.apiKey,
		appID:       c.appID,
		timeout:     c.timeout,
	}
}

func (c *SimpleConfig) ToLegacy() *Config {
	timeoutDuration := time.Duration(c.timeout) * time.Second

	return &Config{
		APIBaseURL:        c.apiURL,
		WebSocketBaseURL:  strings.Replace(c.apiURL, "http", "ws", 1) + "/ws",
		HTTPTimeout:       timeoutDuration,
		ServiceName:       "abacatepay-cli",
		TokenKey:          "auth-token",
		DefaultForwardURL: "http://localhost:3000/webhooks",
		MaxRetries:        30,
		RetryCount:        3,
		RetryWaitTime:     time.Second,
		PollInterval:      2 * time.Second,
		Verbose:           c.debug,
	}
}

const (
	SimpleEnvEnvironment = "ABACATE_ENV"
	SimpleEnvAPIKey      = "ABACATE_API_KEY"
	SimpleEnvAppID       = "ABACATE_APP_ID"
	SimpleEnvAPIURL      = "ABACATE_API_URL"
	SimpleEnvDebug       = "ABACATE_DEBUG"
	SimpleEnvTimeout     = "ABACATE_TIMEOUT"
)
