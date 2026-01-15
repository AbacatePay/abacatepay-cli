package config

import (
	"os"
	"testing"
	"time"
)

func TestNewSimpleConfig(t *testing.T) {
	builder := NewSimpleConfig()
	if builder == nil {
		t.Fatal("NewSimpleConfig() returned nil")
	}

	// Test default values without building
	if builder.config.apiURL != "https://api.abacatepay.com" {
		t.Errorf("Expected default APIURL to be 'https://api.abacatepay.com', got '%s'", builder.config.apiURL)
	}

	if builder.config.environment != "production" {
		t.Errorf("Expected default Environment to be 'production', got '%s'", builder.config.environment)
	}

	if builder.config.debug != false {
		t.Errorf("Expected default Debug to be false, got %t", builder.config.debug)
	}

	if builder.config.timeout != 30 {
		t.Errorf("Expected default Timeout to be 30, got %d", builder.config.timeout)
	}

	// Test successful build
	config, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	if config.apiURL != "https://api.abacatepay.com" {
		t.Errorf("Expected APIURL to be 'https://api.abacatepay.com', got '%s'", config.apiURL)
	}

	if config.environment != "production" {
		t.Errorf("Expected Environment to be 'production', got '%s'", config.environment)
	}

	if config.debug != false {
		t.Errorf("Expected Debug to be false, got %t", config.debug)
	}

	if config.timeout != 30 {
		t.Errorf("Expected Timeout to be 30, got %d", config.timeout)
	}
}

func TestSimpleConfigBuilder_Environment(t *testing.T) {
	tests := []struct {
		env         string
		expectedURL string
		expectedDbg bool
	}{
		{"production", "https://api.abacatepay.com", false},
		{"prod", "https://api.abacatepay.com", false},
		{"staging", "https://staging-api.abacatepay.com", true},
		{"stage", "https://staging-api.abacatepay.com", true},
		{"local", "http://191.252.202.128:8080", true},
		{"development", "http://191.252.202.128:8080", true},
		{"dev", "http://191.252.202.128:8080", true},
		{"unknown", "https://api.abacatepay.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			config := NewSimpleConfig().
				Environment(tt.env).
				MustBuild()

			if config.environment != tt.env {
				t.Errorf("Expected Environment to be '%s', got '%s'", tt.env, config.environment)
			}

			if config.apiURL != tt.expectedURL {
				t.Errorf("Expected APIURL to be '%s', got '%s'", tt.expectedURL, config.apiURL)
			}

			if config.debug != tt.expectedDbg {
				t.Errorf("Expected Debug to be %t, got %t", tt.expectedDbg, config.debug)
			}
		})
	}
}

func TestSimpleConfigBuilder_APIURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://api.test.com", "https://api.test.com"},
		{"https://api.test.com/", "https://api.test.com"},
		{"http://localhost:3000", "http://localhost:3000"},
		{"http://localhost:3000/", "http://localhost:3000"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			config := NewSimpleConfig().
				APIURL(tt.input).
				MustBuild()

			if config.apiURL != tt.expected {
				t.Errorf("Expected APIURL to be '%s', got '%s'", tt.expected, config.apiURL)
			}
		})
	}
}

func TestSimpleConfigBuilder_Debug(t *testing.T) {
	tests := []struct {
		debug bool
	}{
		{true},
		{false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			config := NewSimpleConfig().
				Debug(tt.debug).
				MustBuild()

			if config.debug != tt.debug {
				t.Errorf("Expected Debug to be %t, got %t", tt.debug, config.debug)
			}
		})
	}
}

func TestSimpleConfigBuilder_Timeout(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{10, 10},
		{60, 60},
		{0, 30},  // Should use default
		{-5, 30}, // Should use default
		{120, 120},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			config := NewSimpleConfig().
				Timeout(tt.input).
				MustBuild()

			if config.timeout != tt.expected {
				t.Errorf("Expected Timeout to be %d, got %d", tt.expected, config.timeout)
			}
		})
	}
}

func TestSimpleConfig_ValidateTimeout(t *testing.T) {
	t.Run("valid positive timeout", func(t *testing.T) {
		config := NewSimpleConfig().
			Timeout(60).
			MustBuild()

		if config.timeout != 60 {
			t.Errorf("Expected timeout validation to pass for positive value, got %d", config.timeout)
		}
	})

	t.Run("zero timeout validation", func(t *testing.T) {
		builder := NewSimpleConfig()
		builder.config.timeout = 0 // Manually set to trigger validation

		_, err := builder.Build()

		if err == nil {
			t.Error("Expected validation error for timeout <= 0")
		}
		if err != nil && err.Error() != "timeout must be greater than 0" {
			t.Errorf("Expected 'timeout must be greater than 0', got '%s'", err.Error())
		}
	})
}

func TestSimpleConfigBuilder_Build(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *SimpleConfigBuilder
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid config",
			setup: func() *SimpleConfigBuilder {
				return NewSimpleConfig()
			},
			expectErr: false,
		},
		{
			name: "empty API URL",
			setup: func() *SimpleConfigBuilder {
				return NewSimpleConfig().APIURL("")
			},
			expectErr: true,
			errMsg:    "API URL cannot be empty",
		},
		{
			name: "invalid API URL - no protocol",
			setup: func() *SimpleConfigBuilder {
				return NewSimpleConfig().APIURL("api.example.com")
			},
			expectErr: true,
			errMsg:    "API URL must start with http:// or https://",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.setup().Build()

			if tt.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				} else if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestSimpleConfigBuilder_MustBuild(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config := NewSimpleConfig().MustBuild()
		if config == nil {
			t.Error("Expected config to be non-nil")
		}
	})

	t.Run("invalid config panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected MustBuild to panic on invalid config")
			}
		}()
		NewSimpleConfig().APIURL("invalid-url").MustBuild()
	})
}

func TestLoadFromEnvironment(t *testing.T) {
	clearSimpleEnvVars()

	tests := []struct {
		name     string
		setupEnv func()
		expected func(*SimpleConfig) bool
	}{
		{
			name:     "default values",
			setupEnv: func() {},
			expected: func(c *SimpleConfig) bool {
				return c.environment == "production" && !c.debug
			},
		},
		{
			name: "environment override",
			setupEnv: func() {
				os.Setenv(SimpleEnvEnvironment, "staging")
			},
			expected: func(c *SimpleConfig) bool {
				return c.environment == "staging" && c.debug && c.apiURL == "https://staging-api.abacatepay.com"
			},
		},
		{
			name: "full environment",
			setupEnv: func() {
				os.Setenv(SimpleEnvAPIKey, "env-key")
				os.Setenv(SimpleEnvAppID, "env-app")
				os.Setenv(SimpleEnvAPIURL, "https://env-api.com")
				os.Setenv(SimpleEnvDebug, "true")
				os.Setenv(SimpleEnvTimeout, "45")
			},
			expected: func(c *SimpleConfig) bool {
				return c.apiKey == "env-key" &&
					c.appID == "env-app" &&
					c.apiURL == "https://env-api.com" &&
					c.debug == true &&
					c.timeout == 45
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearSimpleEnvVars()
			tt.setupEnv()

			config := LoadFromEnvironment().MustBuild()

			if !tt.expected(config) {
				t.Errorf("Environment configuration was not applied correctly. Got: APIKey=%s, AppID=%s, APIURL=%s, Debug=%t, Timeout=%d",
					config.apiKey, config.appID, config.apiURL, config.debug, config.timeout)
			}
		})
	}
}

func TestLoadFromEnvironment_InvalidValues(t *testing.T) {
	clearSimpleEnvVars()

	tests := []struct {
		name string
		env  map[string]string
	}{
		{
			name: "invalid debug",
			env: map[string]string{
				SimpleEnvDebug: "not-a-boolean",
			},
		},
		{
			name: "invalid timeout",
			env: map[string]string{
				SimpleEnvTimeout: "not-a-number",
			},
		},
		{
			name: "negative timeout",
			env: map[string]string{
				SimpleEnvTimeout: "-5",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearSimpleEnvVars()
			for k, v := range tt.env {
				os.Setenv(k, v)
			}

			config := LoadFromEnvironment().MustBuild()

			if tt.env[SimpleEnvDebug] == "not-a-boolean" && config.debug != false {
				t.Error("Invalid debug value should not change default")
			}

			if tt.env[SimpleEnvTimeout] == "not-a-number" && config.timeout != 30 {
				t.Error("Invalid timeout value should not change default")
			}

			if tt.env[SimpleEnvTimeout] == "-5" && config.timeout != 30 {
				t.Error("Negative timeout value should not change default")
			}
		})
	}
}

func TestSimpleDefault(t *testing.T) {
	config := SimpleDefault()

	if config == nil {
		t.Fatal("SimpleDefault() returned nil")
	}

	if !config.IsProduction() {
		t.Errorf("Expected production environment, got %s", config.environment)
	}

	if config.apiURL != "https://api.abacatepay.com" {
		t.Errorf("Expected production API URL, got %s", config.apiURL)
	}

	if config.debug != false {
		t.Errorf("Expected debug to be false in production, got %v", config.debug)
	}

	if config.HasAPIKey() {
		t.Error("SimpleDefault() should not have API key by default")
	}

	if config.HasAppID() {
		t.Error("SimpleDefault() should not have App ID by default")
	}
}

func TestSimpleLocal(t *testing.T) {
	config := SimpleLocal()

	if config == nil {
		t.Fatal("SimpleLocal() returned nil")
	}

	if !config.IsDevelopment() {
		t.Errorf("Expected development environment, got %s", config.environment)
	}

	if config.apiURL != "http://191.252.202.128:8080" {
		t.Errorf("Expected local API URL, got %s", config.apiURL)
	}

	if config.debug != true {
		t.Errorf("Expected debug to be true in local, got %v", config.debug)
	}

	if config.HasAPIKey() {
		t.Error("SimpleLocal() should not have API key by default")
	}

	if config.HasAppID() {
		t.Error("SimpleLocal() should not have App ID by default")
	}
}

func TestSimpleDebug(t *testing.T) {
	config := SimpleDebug()

	if config == nil {
		t.Fatal("SimpleDebug() returned nil")
	}

	if !config.IsDevelopment() {
		t.Errorf("Expected development environment, got %s", config.environment)
	}

	if config.apiURL != "http://191.252.202.128:8080" {
		t.Errorf("Expected local API URL, got %s", config.apiURL)
	}

	if config.debug != true {
		t.Errorf("Expected debug to be true, got %v", config.debug)
	}

	if config.HasAPIKey() {
		t.Error("SimpleDebug() should not have API key by default")
	}

	if config.HasAppID() {
		t.Error("SimpleDebug() should not have App ID by default")
	}
}

func TestSimpleConfig_IsProduction(t *testing.T) {
	tests := []struct {
		env      string
		expected bool
	}{
		{"production", true},
		{"prod", true},
		{"staging", false},
		{"local", false},
		{"development", false},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			config := NewSimpleConfig().
				Environment(tt.env).
				MustBuild()

			if config.IsProduction() != tt.expected {
				t.Errorf("Expected IsProduction() to be %t for environment '%s', got %t", tt.expected, tt.env, config.IsProduction())
			}
		})
	}
}

func TestSimpleConfig_IsStaging(t *testing.T) {
	tests := []struct {
		env      string
		expected bool
	}{
		{"production", false},
		{"staging", true},
		{"stage", true},
		{"local", false},
		{"development", false},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			config := NewSimpleConfig().
				Environment(tt.env).
				MustBuild()

			if config.IsStaging() != tt.expected {
				t.Errorf("Expected IsStaging() to be %t for environment '%s', got %t", tt.expected, tt.env, config.IsStaging())
			}
		})
	}
}

func TestSimpleConfig_IsDevelopment(t *testing.T) {
	tests := []struct {
		env      string
		expected bool
	}{
		{"production", false},
		{"staging", false},
		{"local", true},
		{"development", true},
		{"dev", true},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			config := NewSimpleConfig().
				Environment(tt.env).
				MustBuild()

			if config.IsDevelopment() != tt.expected {
				t.Errorf("Expected IsDevelopment() to be %t for environment '%s', got %t", tt.expected, tt.env, config.IsDevelopment())
			}
		})
	}
}

func TestSimpleConfig_String(t *testing.T) {
	config := NewSimpleConfig().
		APIKey("secret-key").
		AppID("test-app").
		Environment("production").
		Debug(true).
		MustBuild()

	str := config.String()

	if str == "" {
		t.Error("String() returned empty string")
	}

	if str[0] != 'S' {
		t.Error("String() should start with 'SimpleConfig{'")
	}

	if str[len(str)-1] != '}' {
		t.Error("String() should end with '}'")
	}

	// Ensure sensitive data is not in the string
	if len(str) > len("secret-key") {
		for i := 0; i <= len(str)-len("secret-key"); i++ {
			match := true
			for j := 0; j < len("secret-key"); j++ {
				if str[i+j] != "secret-key"[j] {
					match = false
					break
				}
			}
			if match {
				t.Error("String() should not contain API key")
				break
			}
		}
	}
}

func TestSimpleConfig_Clone(t *testing.T) {
	original := NewSimpleConfig().
		APIKey("test-key").
		AppID("test-app").
		Environment("staging").
		Debug(true).
		Timeout(60).
		MustBuild()

	clone := original.Clone()

	if clone == original {
		t.Error("Clone() should return a different instance")
	}

	if clone.apiKey != original.apiKey {
		t.Error("Clone() should preserve API key")
	}

	if clone.appID != original.appID {
		t.Error("Clone() should preserve App ID")
	}

	if clone.environment != original.environment {
		t.Error("Clone() should preserve Environment")
	}

	if clone.debug != original.debug {
		t.Error("Clone() should preserve Debug setting")
	}

	if clone.timeout != original.timeout {
		t.Error("Clone() should preserve Timeout")
	}
}

func TestSimpleConfig_ToLegacy(t *testing.T) {
	simple := NewSimpleConfig().
		APIKey("test-key").
		AppID("test-app").
		Environment("production").
		Debug(true).
		Timeout(45).
		APIURL("https://api.custom.com").
		MustBuild()

	legacy := simple.ToLegacy()

	if legacy == nil {
		t.Fatal("ToLegacy() returned nil")
	}

	if legacy.APIBaseURL != "https://api.custom.com" {
		t.Errorf("Expected APIBaseURL to be 'https://api.custom.com', got '%s'", legacy.APIBaseURL)
	}

	if legacy.WebSocketBaseURL != "wss://api.custom.com/ws" {
		t.Errorf("Expected WebSocketBaseURL to be 'wss://api.custom.com/ws', got '%s'", legacy.WebSocketBaseURL)
	}

	if legacy.HTTPTimeout != 45*time.Second {
		t.Errorf("Expected HTTPTimeout to be 45s, got %v", legacy.HTTPTimeout)
	}

	if legacy.ServiceName != "abacatepay-cli" {
		t.Errorf("Expected ServiceName to be 'abacatepay-cli', got '%s'", legacy.ServiceName)
	}

	if legacy.Verbose != true {
		t.Error("Expected Verbose to be true")
	}

	if legacy.TokenKey != "auth-token" {
		t.Errorf("Expected TokenKey to be 'auth-token', got '%s'", legacy.TokenKey)
	}
}

func TestSimpleConfig_HasCredentials(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *SimpleConfig
		hasKey   bool
		hasAppID bool
	}{
		{
			name: "both credentials present",
			setup: func() *SimpleConfig {
				return NewSimpleConfig().
					APIKey("test-key").
					AppID("test-app").
					MustBuild()
			},
			hasKey:   true,
			hasAppID: true,
		},
		{
			name: "only API key present",
			setup: func() *SimpleConfig {
				return NewSimpleConfig().
					APIKey("test-key").
					MustBuild()
			},
			hasKey:   true,
			hasAppID: false,
		},
		{
			name: "only App ID present",
			setup: func() *SimpleConfig {
				return NewSimpleConfig().
					AppID("test-app").
					MustBuild()
			},
			hasKey:   false,
			hasAppID: true,
		},
		{
			name: "no credentials present",
			setup: func() *SimpleConfig {
				return NewSimpleConfig().MustBuild()
			},
			hasKey:   false,
			hasAppID: false,
		},
		{
			name: "empty API key",
			setup: func() *SimpleConfig {
				return NewSimpleConfig().
					APIKey("").
					MustBuild()
			},
			hasKey:   false,
			hasAppID: false,
		},
		{
			name: "empty App ID",
			setup: func() *SimpleConfig {
				return NewSimpleConfig().
					AppID("").
					MustBuild()
			},
			hasKey:   false,
			hasAppID: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.setup()

			if config.HasAPIKey() != tt.hasKey {
				t.Errorf("HasAPIKey() = %v, want %v", config.HasAPIKey(), tt.hasKey)
			}

			if config.HasAppID() != tt.hasAppID {
				t.Errorf("HasAppID() = %v, want %v", config.HasAppID(), tt.hasAppID)
			}
		})
	}
}

func TestSimpleConfig_ToLegacyWithoutCredentials(t *testing.T) {
	config := NewSimpleConfig().MustBuild()
	legacy := config.ToLegacy()

	if legacy == nil {
		t.Fatal("ToLegacy() returned nil")
	}

	if legacy.APIBaseURL != config.apiURL {
		t.Errorf("APIBaseURL = %v, want %v", legacy.APIBaseURL, config.apiURL)
	}

	if legacy.Verbose != config.debug {
		t.Errorf("Verbose = %v, want %v", legacy.Verbose, config.debug)
	}
}

func clearSimpleEnvVars() {
	envVars := []string{
		SimpleEnvEnvironment,
		SimpleEnvAPIKey,
		SimpleEnvAppID,
		SimpleEnvAPIURL,
		SimpleEnvDebug,
		SimpleEnvTimeout,
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}
