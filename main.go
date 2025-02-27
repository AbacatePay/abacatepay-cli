package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

const (
	apiBaseURL       = "https://api.abacatepay.com"
	configFileName   = ".abacatepay.json"
	pollInterval     = 2 * time.Second
	maxRetries       = 30
	websocketBaseURL = "wss://ws.abacatepay.com/ws"
)

// Config structure
type Config struct {
	Token string `json:"token"`
}

// ReadConfig reads the stored configuration
func ReadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return &Config{}, nil // Return an empty config if file doesn't exist
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// WriteConfig writes configuration to file
func WriteConfig(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, configFileName), nil
}

// PollForToken polls the API until the token is received
func PollForToken(deviceCode string) (string, error) {
	for retries := 0; retries < maxRetries; retries++ {
		time.Sleep(pollInterval)

		body, _ := json.Marshal(map[string]string{"deviceCode": deviceCode})
		resp, err := http.Post(apiBaseURL+"/token", "application/json", bytes.NewBuffer(body))
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			var result map[string]string
			json.NewDecoder(resp.Body).Decode(&result)
			if token, ok := result["token"]; ok {
				return token, nil
			}
		}
	}

	return "", fmt.Errorf("authorization timed out. Please try again")
}

// Login handles device authentication
func Login() {
	host, _ := os.Hostname()

	body, _ := json.Marshal(map[string]string{"host": host})
	resp, err := http.Post(apiBaseURL+"/device-login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("❌ Login failed: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	fmt.Println("🔗 Open the following link in your browser to authenticate:")
	fmt.Printf("👉 %s\n", result["verificationUri"])

	fmt.Println("⌛ Waiting for authorization...")
	token, err := PollForToken(result["deviceCode"])
	if err != nil {
		log.Fatalf("❌ %v", err)
	}

	config := &Config{Token: token}
	if err := WriteConfig(config); err != nil {
		log.Fatalf("❌ Failed to save login token: %v", err)
	}

	fmt.Println("✅ Logged in successfully.")
}

// Logout removes the stored authentication token
func Logout() error {
	config := &Config{Token: ""}
	if err := WriteConfig(config); err != nil {
		return fmt.Errorf("failed to clear login token: %v", err)
	}
	return nil
}

// Listen connects to WebSocket and forwards webhooks
func Listen(forwardURL string, token string) {
	fmt.Println("🌐 Starting webhook listener...")

	conn, _, err := websocket.DefaultDialer.Dial(websocketBaseURL+"?token="+token, nil)
	if err != nil {
		log.Fatalf("❌ Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	fmt.Println("✅ Webhook listener started")
	fmt.Printf("🔄 Forwarding events to: %s\n", forwardURL)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("❌ Webhook listener closed unexpectedly: %v", err)
		}

		fmt.Println("\n🌟 Received Webhook:")
		fmt.Println(string(message))

		// Forward webhook
		resp, err := http.Post(forwardURL, "application/json", bytes.NewBuffer(message))
		if err != nil {
			log.Printf("❌ Failed to forward webhook: %v", err)
			continue
		}
		resp.Body.Close()
	}
}

func main() {
	var forwardURL string

	rootCmd := &cobra.Command{Use: "abacatepay-cli"}

	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to AbacatePay",
		Run: func(cmd *cobra.Command, args []string) {
			Login()
		},
	}

	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out from AbacatePay",
		Run: func(cmd *cobra.Command, args []string) {
			if err := Logout(); err != nil {
				log.Fatalf("❌ %v", err)
			}
			fmt.Println("✅ Logged out successfully")
		},
	}

	listenCmd := &cobra.Command{
		Use:   "listen",
		Short: "Listen for webhooks and forward them to your local server",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := ReadConfig()
			if err != nil {
				log.Fatalf("❌ Failed to read config: %v", err)
			}

			if config.Token == "" {
				log.Fatal("❌ You are not logged in. Please run `abacatepay-cli login` first.")
			}
			Listen(forwardURL, config.Token)
		},
	}

	listenCmd.Flags().StringVarP(&forwardURL, "forward", "f", "", "Local server URL to forward webhooks")
	listenCmd.MarkFlagRequired("forward")

	rootCmd.AddCommand(loginCmd, logoutCmd, listenCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
