package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAppStateActions_MessagesState(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/app-state/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"messages":["welcome"]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	var res map[string]any
	err := c.AppState().MessagesState(context.Background(), &res)
	if err != nil {
		t.Fatalf("MessagesState error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
	if msgs, ok := res["messages"].([]any); !ok || len(msgs) != 1 {
		t.Errorf("expected messages array with 1 element, got %v", res["messages"])
	}
}

func TestAppStateActions_UpdateMessagesState(t *testing.T) {
	var capturedBody map[string]any
	var capturedMethod string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/app-state/messages", func(w http.ResponseWriter, r *http.Request) {
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

	body := map[string]any{"dismissed": []string{"msg-1"}}
	err := c.AppState().UpdateMessagesState(context.Background(), body)
	if err != nil {
		t.Fatalf("UpdateMessagesState error: %v", err)
	}
	if capturedMethod != http.MethodPut {
		t.Errorf("expected PUT, got %s", capturedMethod)
	}
	if capturedBody["dismissed"] == nil {
		t.Error("expected dismissed field in body")
	}
}

func TestAppStateActions_PatchMessagesState(t *testing.T) {
	var capturedBody map[string]any
	var capturedMethod string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/app-state/messages", func(w http.ResponseWriter, r *http.Request) {
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

	body := map[string]any{"read": []string{"msg-2"}}
	err := c.AppState().PatchMessagesState(context.Background(), body)
	if err != nil {
		t.Fatalf("PatchMessagesState error: %v", err)
	}
	if capturedMethod != http.MethodPatch {
		t.Errorf("expected PATCH, got %s", capturedMethod)
	}
	if capturedBody["read"] == nil {
		t.Error("expected read field in body")
	}
}
