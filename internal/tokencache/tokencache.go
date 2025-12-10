package tokencache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/99designs/keyring"
)

const (
	serviceName = "eightctl"
	tokenKey    = "oauth-token"
)

type CachedToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    string    `json:"user_id,omitempty"`
}

var openKeyring = defaultOpenKeyring

func defaultOpenKeyring() (keyring.Keyring, error) {
	home, _ := os.UserHomeDir()
	return keyring.Open(keyring.Config{
		ServiceName: serviceName,
		AllowedBackends: []keyring.BackendType{
			keyring.KeychainBackend,
			keyring.SecretServiceBackend,
			keyring.WinCredBackend,
			keyring.FileBackend,
		},
		FileDir:          filepath.Join(home, ".config", "eightctl", "keyring"),
		FilePasswordFunc: filePassword,
	})
}

func filePassword(_ string) (string, error) {
	return serviceName + "-fallback", nil
}

func Save(token string, expiresAt time.Time, userID string) error {
	ring, err := openKeyring()
	if err != nil {
		return err
	}
	data, err := json.Marshal(CachedToken{
		Token:     token,
		ExpiresAt: expiresAt,
		UserID:    userID,
	})
	if err != nil {
		return err
	}
	return ring.Set(keyring.Item{
		Key:  tokenKey,
		Data: data,
	})
}

func Load() (*CachedToken, error) {
	ring, err := openKeyring()
	if err != nil {
		return nil, err
	}
	item, err := ring.Get(tokenKey)
	if err != nil {
		return nil, err
	}
	var cached CachedToken
	if err := json.Unmarshal(item.Data, &cached); err != nil {
		return nil, err
	}
	if time.Now().After(cached.ExpiresAt) {
		return nil, keyring.ErrKeyNotFound
	}
	return &cached, nil
}

func Clear() error {
	ring, err := openKeyring()
	if err != nil {
		return err
	}
	return ring.Remove(tokenKey)
}
