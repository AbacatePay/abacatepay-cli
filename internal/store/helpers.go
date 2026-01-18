package store

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/99designs/keyring"
)

func getMachineID() string {
	var paths []string

	switch runtime.GOOS {
	case "linux":
		paths = []string{"/etc/machine-id", "/var/lib/dbus/machine-id"}
	case "darwin":
		return ""
	case "windows":
		return ""
	}

	for _, path := range paths {
		if data, err := os.ReadFile(path); err == nil {
			return strings.TrimSpace(string(data))
		}
	}

	return ""
}

func (k *KeyringStore) getKeyring() (keyring.Keyring, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve home directory: %w", err)
	}

	storeDir := filepath.Join(homeDir, ".abacatepay", "keyring")
	if err := os.MkdirAll(storeDir, 0o700); err != nil {
		return nil, fmt.Errorf("failed to create keyring directory: %w", err)
	}

	return keyring.Open(keyring.Config{
		ServiceName: k.serviceName,
		FilePasswordFunc: func(_ string) (string, error) {
			return deriveKeyringPassword(), nil
		},
		FileDir:                  storeDir,
		KeychainTrustApplication: true,
	})
}
