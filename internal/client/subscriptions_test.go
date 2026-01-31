package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSubscriptionsActions_Subscriptions(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v3/users/uid-123/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"subscriptions":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL + "/v1"
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	var res any
	err := c.Subscriptions().Subscriptions(context.Background(), &res)
	if err != nil {
		t.Fatalf("Subscriptions error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestSubscriptionsActions_CreateTemporarySubscription(t *testing.T) {
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/v3/users/uid-123/subscriptions/temporary", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&capturedBody); err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL + "/v1"
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	body := map[string]any{"duration": 30}
	err := c.Subscriptions().CreateTemporarySubscription(context.Background(), body)
	if err != nil {
		t.Fatalf("CreateTemporarySubscription error: %v", err)
	}
	if capturedBody["duration"] != float64(30) {
		t.Errorf("expected duration=30, got %v", capturedBody["duration"])
	}
}

func TestSubscriptionsActions_RedeemSubscription(t *testing.T) {
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/v3/users/uid-123/subscriptions/redeem", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&capturedBody); err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL + "/v1"
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	body := map[string]any{"code": "PROMO123"}
	err := c.Subscriptions().RedeemSubscription(context.Background(), body)
	if err != nil {
		t.Fatalf("RedeemSubscription error: %v", err)
	}
	if capturedBody["code"] != "PROMO123" {
		t.Errorf("expected code=PROMO123, got %v", capturedBody["code"])
	}
}
