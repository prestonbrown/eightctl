package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPairingActions_StartAutoPairing(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-456"}}}`))
	})
	mux.HandleFunc("/devices/dev-456/auto-pairing/start", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"pairingId":"pair-789","status":"started"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Pairing().StartAutoPairing(context.Background())
	if err != nil {
		t.Fatalf("StartAutoPairing error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestPairingActions_GetPairingStatus(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-456"}}}`))
	})
	mux.HandleFunc("/devices/dev-456/auto-pairing/status/pair-789", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"pairingId":"pair-789","status":"completed"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	var out any
	err := c.Pairing().GetPairingStatus(context.Background(), "pair-789", &out)
	if err != nil {
		t.Fatalf("GetPairingStatus error: %v", err)
	}
	if out == nil {
		t.Error("expected non-nil response")
	}
}
