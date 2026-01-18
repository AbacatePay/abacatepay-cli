package store

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/99designs/keyring"
)

type TokenStore interface {
	Save(token string) error
	SaveNamed(name, token string) error
	Get() (string, error)
	GetNamed(name string) (string, error)
	Delete() error
	DeleteNamed(name string) error
	SetActiveProfile(name string) error
	GetActiveProfile() (string, error)
	List() ([]string, error)
}

type KeyringStore struct {
	serviceName string
	tokenKey    string
}

func NewKeyringStore(serviceName, tokenKey string) *KeyringStore {
	return &KeyringStore{
		serviceName: serviceName,
		tokenKey:    tokenKey,
	}
}

func deriveKeyringPassword() string {
	var parts []string

	if u, err := user.Current(); err == nil {
		parts = append(parts, u.Username, u.Uid)
	}

	if machineID := getMachineID(); machineID != "" {
		parts = append(parts, machineID)
	}

	if hostname, err := os.Hostname(); err == nil {
		parts = append(parts, hostname)
	}

	if home, err := os.UserHomeDir(); err == nil {
		parts = append(parts, home)
	}

	combined := strings.Join(parts, "|")
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

func (k *KeyringStore) List() ([]string, error) {
	ring, err := k.getKeyring()
	if err != nil {
		return nil, fmt.Errorf("failed to open keyring: %w", err)
	}

	keys, err := ring.Keys()
	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %w", err)
	}

	var profiles []string
	for _, key := range keys {
		if key == activeProfileKey {
			continue
		}

		profiles = append(profiles, key)
	}
	return profiles, nil
}

func (k *KeyringStore) Save(token string) error {
	return k.SaveNamed(k.tokenKey, token)
}

func (k *KeyringStore) SaveNamed(name, token string) error {
	ring, err := k.getKeyring()
	if err != nil {
		return fmt.Errorf("failed to open keyring: %w", err)
	}

	if err := ring.Set(keyring.Item{
		Key:  name,
		Data: []byte(token),
	}); err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}

	return nil
}

func (k *KeyringStore) Get() (string, error) {
	return k.GetNamed(k.tokenKey)
}

func (k *KeyringStore) GetNamed(name string) (string, error) {
	ring, err := k.getKeyring()
	if err != nil {
		return "", fmt.Errorf("failed to open keyring: %w", err)
	}

	item, err := ring.Get(name)
	if err == nil {
		return string(item.Data), nil
	}

	if err == keyring.ErrKeyNotFound {
		return "", nil
	}

	return "", fmt.Errorf("failed to read from keyring: %w", err)
}

func (k *KeyringStore) Delete() error {
	return k.DeleteNamed(k.tokenKey)
}

func (k *KeyringStore) DeleteNamed(name string) error {
	ring, err := k.getKeyring()
	if err != nil {
		return fmt.Errorf("failed to open keyring: %w", err)
	}

	if err := ring.Remove(name); err != nil && err != keyring.ErrKeyNotFound {
		return fmt.Errorf("failed to remove token: %w", err)
	}

	return nil
}

const activeProfileKey = "active-profile-name"

func (k *KeyringStore) SetActiveProfile(name string) error {
	return k.SaveNamed(activeProfileKey, name)
}

func (k *KeyringStore) GetActiveProfile() (string, error) {
	return k.GetNamed(activeProfileKey)
}
