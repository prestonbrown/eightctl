package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserActions_GetMe(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123"}}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	var result any
	err := c.Users().GetMe(context.Background(), &result)
	if err != nil {
		t.Fatalf("GetMe error: %v", err)
	}
	if capturedPath != "/users/me" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestUserActions_GetUser(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"userId":"uid-123"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	var result any
	err := c.Users().GetUser(context.Background(), "uid-123", &result)
	if err != nil {
		t.Fatalf("GetUser error: %v", err)
	}
	if capturedPath != "/users/uid-123" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestUserActions_UpdateUser(t *testing.T) {
	var capturedPath string
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &capturedBody)
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Users().UpdateUser(context.Background(), map[string]any{"name": "Test User"})
	if err != nil {
		t.Fatalf("UpdateUser error: %v", err)
	}
	if capturedPath != "/users/uid-123" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if capturedBody["name"] != "Test User" {
		t.Errorf("unexpected body: %v", capturedBody)
	}
}

func TestUserActions_UpdateEmail(t *testing.T) {
	var capturedPath string
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/email", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &capturedBody)
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Users().UpdateEmail(context.Background(), map[string]any{"email": "new@example.com"})
	if err != nil {
		t.Fatalf("UpdateEmail error: %v", err)
	}
	if capturedPath != "/users/uid-123/email" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if capturedBody["email"] != "new@example.com" {
		t.Errorf("unexpected body: %v", capturedBody)
	}
}

func TestUserActions_PasswordReset(t *testing.T) {
	var capturedPath string
	var capturedBody map[string]string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/password-reset", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &capturedBody)
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Users().PasswordReset(context.Background(), "user@example.com")
	if err != nil {
		t.Fatalf("PasswordReset error: %v", err)
	}
	if capturedPath != "/users/password-reset" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if capturedBody["email"] != "user@example.com" {
		t.Errorf("unexpected body: %v", capturedBody)
	}
}
