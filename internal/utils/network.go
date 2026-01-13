package utils

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"time"
)

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
