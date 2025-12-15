package auth

import (
	"fmt"
	"os"

	"github.com/99designs/keyring"
)

type TokenStore interface {
	Save(token string) error
	Get() (string, error)
	Delete() error
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
		return nil, fmt.Errorf("falha ao obter diretório home: %w", err)
	}

	return keyring.Open(keyring.Config{
		ServiceName:              k.serviceName,
		FilePasswordFunc:         keyring.TerminalPrompt,
		FileDir:                  homeDir,
		KeychainTrustApplication: true,
	})
}

func (k *KeyringStore) Save(token string) error {
	ring, err := k.getKeyring()
	if err != nil {
		return fmt.Errorf("falha ao abrir keyring: %w", err)
	}

	if err := ring.Set(keyring.Item{
		Key:  k.tokenKey,
		Data: []byte(token),
	}); err != nil {
		return fmt.Errorf("falha ao salvar token: %w", err)
	}

	return nil
}

func (k *KeyringStore) Get() (string, error) {
	ring, err := k.getKeyring()
	if err != nil {
		return "", fmt.Errorf("falha ao abrir keyring: %w", err)
	}

	item, err := ring.Get(k.tokenKey)
	if err != nil {
		if err == keyring.ErrKeyNotFound {
			return "", nil
		}
		return "", fmt.Errorf("falha ao recuperar token: %w", err)
	}

	return string(item.Data), nil
}

func (k *KeyringStore) Delete() error {
	ring, err := k.getKeyring()
	if err != nil {
		return fmt.Errorf("falha ao abrir keyring: %w", err)
	}

	if err := ring.Remove(k.tokenKey); err != nil && err != keyring.ErrKeyNotFound {
		return fmt.Errorf("falha ao remover token: %w", err)
	}

	return nil
}

type MemoryStore struct {
	token string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) Save(token string) error {
	m.token = token
	return nil
}

func (m *MemoryStore) Get() (string, error) {
	return m.token, nil
}

func (m *MemoryStore) Delete() error {
	m.token = ""
	return nil
}
