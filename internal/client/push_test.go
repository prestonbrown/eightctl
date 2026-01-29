package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPushActions_UpdatePushToken(t *testing.T) {
	var capturedPath string
	var capturedMethod string
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/me/push-targets/device-123", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedMethod = r.Method
		if err := json.NewDecoder(r.Body).Decode(&capturedBody); err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	body := map[string]any{"token": "fcm-token-abc", "platform": "android"}
	err := c.Push().UpdatePushToken(context.Background(), "device-123", body)
	if err != nil {
		t.Fatalf("UpdatePushToken error: %v", err)
	}
	if capturedMethod != http.MethodPut {
		t.Errorf("expected PUT, got %s", capturedMethod)
	}
	if capturedPath != "/users/me/push-targets/device-123" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if capturedBody["token"] != "fcm-token-abc" {
		t.Errorf("expected token=fcm-token-abc, got %v", capturedBody["token"])
	}
	if capturedBody["platform"] != "android" {
		t.Errorf("expected platform=android, got %v", capturedBody["platform"])
	}
}

func TestPushActions_DeletePushToken(t *testing.T) {
	var capturedPath string
	var capturedMethod string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/me/push-targets/token/fcm-token-xyz", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Push().DeletePushToken(context.Background(), "fcm-token-xyz")
	if err != nil {
		t.Fatalf("DeletePushToken error: %v", err)
	}
	if capturedMethod != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", capturedMethod)
	}
	if capturedPath != "/users/me/push-targets/token/fcm-token-xyz" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}
