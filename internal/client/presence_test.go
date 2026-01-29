package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetPresence(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/presence", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"presence":true}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	present, err := c.GetPresence(context.Background())
	if err != nil {
		t.Fatalf("GetPresence error: %v", err)
	}
	if !present {
		t.Error("expected presence=true")
	}
}

func TestGetPresence_NotPresent(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/presence", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"presence":false}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	present, err := c.GetPresence(context.Background())
	if err != nil {
		t.Fatalf("GetPresence error: %v", err)
	}
	if present {
		t.Error("expected presence=false")
	}
}
