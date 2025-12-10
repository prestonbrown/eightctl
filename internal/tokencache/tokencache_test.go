package tokencache

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"github.com/99designs/keyring"
)

func testKeyring(t *testing.T) keyring.Keyring {
	t.Helper()
	tmpDir := t.TempDir()
	ring, err := keyring.Open(keyring.Config{
		ServiceName:      serviceName + "-test",
		AllowedBackends:  []keyring.BackendType{keyring.FileBackend},
		FileDir:          filepath.Join(tmpDir, "keyring"),
		FilePasswordFunc: func(_ string) (string, error) { return "test-pass", nil },
	})
	if err != nil {
		t.Fatalf("failed to open test keyring: %v", err)
	}
	return ring
}

func TestSaveAndLoad(t *testing.T) {
	ring := testKeyring(t)

	token := "test-token-123"
	expiresAt := time.Now().Add(time.Hour)
	userID := "user-456"

	// Save token to keyring
	data, err := json.Marshal(CachedToken{
		Token:     token,
		ExpiresAt: expiresAt,
		UserID:    userID,
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := ring.Set(keyring.Item{Key: tokenKey, Data: data}); err != nil {
		t.Fatalf("set: %v", err)
	}

	// Load token from keyring
	item, err := ring.Get(tokenKey)
	if err != nil {
		t.Fatalf("get: %v", err)
	}

	var cached CachedToken
	if err := json.Unmarshal(item.Data, &cached); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if cached.Token != token {
		t.Errorf("token: got %s, want %s", cached.Token, token)
	}
	if cached.UserID != userID {
		t.Errorf("userID: got %s, want %s", cached.UserID, userID)
	}
	if !cached.ExpiresAt.Equal(expiresAt) {
		t.Errorf("expiresAt: got %v, want %v", cached.ExpiresAt, expiresAt)
	}
}

func TestExpiredToken(t *testing.T) {
	ring := testKeyring(t)

	// Save expired token
	expiredTime := time.Now().Add(-time.Hour)
	data, err := json.Marshal(CachedToken{
		Token:     "expired-token",
		ExpiresAt: expiredTime,
		UserID:    "user-123",
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := ring.Set(keyring.Item{Key: tokenKey, Data: data}); err != nil {
		t.Fatalf("set: %v", err)
	}

	// Try to load - should get error for expired token
	item, err := ring.Get(tokenKey)
	if err != nil {
		t.Fatalf("get: %v", err)
	}

	var cached CachedToken
	if err := json.Unmarshal(item.Data, &cached); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !time.Now().After(cached.ExpiresAt) {
		t.Errorf("token should be expired")
	}
}

func TestClear(t *testing.T) {
	ring := testKeyring(t)

	// Save a token
	data, err := json.Marshal(CachedToken{
		Token:     "test-token",
		ExpiresAt: time.Now().Add(time.Hour),
		UserID:    "user-123",
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := ring.Set(keyring.Item{Key: tokenKey, Data: data}); err != nil {
		t.Fatalf("set: %v", err)
	}

	// Verify token exists
	if _, err := ring.Get(tokenKey); err != nil {
		t.Fatalf("get before clear: %v", err)
	}

	// Clear the token
	if err := ring.Remove(tokenKey); err != nil {
		t.Fatalf("remove: %v", err)
	}

	// Verify token is gone
	if _, err := ring.Get(tokenKey); err != keyring.ErrKeyNotFound {
		t.Errorf("expected ErrKeyNotFound after clear, got: %v", err)
	}
}

func TestIntegrationSaveLoadClear(t *testing.T) {
	// Override openKeyring for this test
	tmpDir := t.TempDir()
	origOpenKeyring := openKeyring
	openKeyring = func() (keyring.Keyring, error) {
		return keyring.Open(keyring.Config{
			ServiceName:      serviceName + "-integration",
			AllowedBackends:  []keyring.BackendType{keyring.FileBackend},
			FileDir:          filepath.Join(tmpDir, "keyring"),
			FilePasswordFunc: func(_ string) (string, error) { return "test-pass", nil },
		})
	}
	t.Cleanup(func() { openKeyring = origOpenKeyring })

	token := "integration-token"
	expiresAt := time.Now().Add(2 * time.Hour)
	userID := "integration-user"

	// Test Save
	if err := Save(token, expiresAt, userID); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Test Load
	cached, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cached.Token != token {
		t.Errorf("token: got %s, want %s", cached.Token, token)
	}
	if cached.UserID != userID {
		t.Errorf("userID: got %s, want %s", cached.UserID, userID)
	}

	// Test Clear
	if err := Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}

	// Verify cleared
	if _, err := Load(); err != keyring.ErrKeyNotFound {
		t.Errorf("expected ErrKeyNotFound after Clear, got: %v", err)
	}
}

func TestLoadExpiredTokenReturnsError(t *testing.T) {
	tmpDir := t.TempDir()
	origOpenKeyring := openKeyring
	openKeyring = func() (keyring.Keyring, error) {
		return keyring.Open(keyring.Config{
			ServiceName:      serviceName + "-expired",
			AllowedBackends:  []keyring.BackendType{keyring.FileBackend},
			FileDir:          filepath.Join(tmpDir, "keyring"),
			FilePasswordFunc: func(_ string) (string, error) { return "test-pass", nil },
		})
	}
	t.Cleanup(func() { openKeyring = origOpenKeyring })

	// Save expired token
	expiredTime := time.Now().Add(-time.Minute)
	if err := Save("expired-token", expiredTime, "user-id"); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Load should return ErrKeyNotFound for expired token
	if _, err := Load(); err != keyring.ErrKeyNotFound {
		t.Errorf("expected ErrKeyNotFound for expired token, got: %v", err)
	}
}

func TestFilePasswordFunc(t *testing.T) {
	pw, err := filePassword("any-string")
	if err != nil {
		t.Fatalf("filePassword: %v", err)
	}
	if pw != serviceName+"-fallback" {
		t.Errorf("password: got %s, want %s", pw, serviceName+"-fallback")
	}
}
