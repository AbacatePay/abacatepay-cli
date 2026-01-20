package utils

import (
	"fmt"
	"net"
	"net/url"
	"os/exec"
	"runtime"
	"time"

	"abacatepay-cli/internal/style"
)

const DefaultForwardURL = "http://localhost:3000/webhooks/abacatepay"

func validateForwardURL(s string) error {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL must start with http:// or https://")
	}
	if u.Host == "" {
		return fmt.Errorf("URL must include a valid host")
	}
	return nil
}

func promptForwardURL(defaultURL string, result *string) error {
	if err := style.Input("Forward events to", defaultURL, result, validateForwardURL); err != nil {
		return err
	}
	if *result == "" {
		*result = defaultURL
	}
	return nil
}

func GetForwardURL(flagValue, defaultURL string) (string, error) {
	if flagValue != "" {
		if err := validateForwardURL(flagValue); err != nil {
			return "", err
		}
		return flagValue, nil
	}

	var result string
	if err := promptForwardURL(defaultURL, &result); err != nil {
		return "", err
	}
	return result, nil
}

func IsOnline() bool {
	const googleDNS = "8.8.8.8:53"
	const timeout = 2 * time.Second

	conn, err := net.DialTimeout("tcp", googleDNS, timeout)
	if err != nil {
		return false
	}

	defer conn.Close()

	return true
}

func OpenBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return exec.Command(cmd, args...).Start()
}
