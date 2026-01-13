package auth

import (
	"fmt"
	"os"

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

func (k *KeyringStore) getKeyring() (keyring.Keyring, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return nil, fmt.Errorf("failed to resolve home directory: %w", err)
	}

	return keyring.Open(keyring.Config{
		ServiceName:              k.serviceName,
		FilePasswordFunc:         keyring.TerminalPrompt,
		FileDir:                  homeDir,
		KeychainTrustApplication: true,
	})
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

	if err != nil {
		if err == keyring.ErrKeyNotFound {
			return "", nil
		}

		return "", fmt.Errorf("failed to read from keyring: %w", err)
	}

	return string(item.Data), nil
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
