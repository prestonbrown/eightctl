package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAutopilotActions_Details(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/autopilotDetails", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"enabled":true,"mode":"auto"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Autopilot().Details(context.Background())
	if err != nil {
		t.Fatalf("Details error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestAutopilotActions_History(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/autopilot-history", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"history":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Autopilot().History(context.Background())
	if err != nil {
		t.Fatalf("History error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestAutopilotActions_Recap(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/autopilotDetails/autopilotRecap", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"recap":{}}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Autopilot().Recap(context.Background())
	if err != nil {
		t.Fatalf("Recap error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestAutopilotActions_SetLevelSuggestions(t *testing.T) {
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/level-suggestions-mode", func(w http.ResponseWriter, r *http.Request) {
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
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Autopilot().SetLevelSuggestions(context.Background(), true)
	if err != nil {
		t.Fatalf("SetLevelSuggestions error: %v", err)
	}
	if capturedBody["enabled"] != true {
		t.Errorf("expected enabled=true, got %v", capturedBody["enabled"])
	}
}

func TestAutopilotActions_SetSnoreMitigation(t *testing.T) {
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/autopilotDetails/snoringMitigation", func(w http.ResponseWriter, r *http.Request) {
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

	err := c.Autopilot().SetSnoreMitigation(context.Background(), false)
	if err != nil {
		t.Fatalf("SetSnoreMitigation error: %v", err)
	}
	if capturedBody["enabled"] != false {
		t.Errorf("expected enabled=false, got %v", capturedBody["enabled"])
	}
}
