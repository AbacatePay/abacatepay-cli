package auth

import (
	"fmt"
	"os"
	"path/filepath"

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

func (k *KeyringStore) getKeyring() (keyring.Keyring, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("falha ao obter diretório home: %w", err)
	}

	storeDir := filepath.Join(homeDir, ".abacatepay", "keyring")
	if err := os.MkdirAll(storeDir, 0o700); err != nil {
		return nil, fmt.Errorf("falha ao criar diretório do keyring: %w", err)
	}

	return keyring.Open(keyring.Config{
		ServiceName: k.serviceName,
		FilePasswordFunc: func(_ string) (string, error) {
			return "abacatepay-cli-auto-unlock", nil
		},
		FileDir:                  storeDir,
		KeychainTrustApplication: true,
	})
}

func (k *KeyringStore) List() ([]string, error) {
	ring, err := k.getKeyring()
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir keyring: %w", err)
	}

	keys, err := ring.Keys()
	if err != nil {
		return nil, fmt.Errorf("falha ao listar chaves: %w", err)
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
		return fmt.Errorf("falha ao abrir keyring: %w", err)
	}

	if err := ring.Set(keyring.Item{
		Key:  name,
		Data: []byte(token),
	}); err != nil {
		return fmt.Errorf("falha ao salvar no keyring: %w", err)
	}

	return nil
}

func (k *KeyringStore) Get() (string, error) {
	return k.GetNamed(k.tokenKey)
}

func (k *KeyringStore) GetNamed(name string) (string, error) {
	ring, err := k.getKeyring()
	if err != nil {
		return "", fmt.Errorf("falha ao abrir keyring: %w", err)
	}

	item, err := ring.Get(name)
	if err != nil {
		if err == keyring.ErrKeyNotFound {
			return "", nil
		}
		return "", fmt.Errorf("falha ao recuperar do keyring: %w", err)
	}

	return string(item.Data), nil
}

func (k *KeyringStore) Delete() error {
	return k.DeleteNamed(k.tokenKey)
}

func (k *KeyringStore) DeleteNamed(name string) error {
	ring, err := k.getKeyring()
	if err != nil {
		return fmt.Errorf("falha ao abrir keyring: %w", err)
	}

	if err := ring.Remove(name); err != nil && err != keyring.ErrKeyNotFound {
		return fmt.Errorf("falha ao remover do keyring: %w", err)
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
