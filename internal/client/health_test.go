package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHealthActions_HealthSurvey(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health-survey/test-drive", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"survey":{"status":"active"}}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	var res map[string]any
	err := c.Health().HealthSurvey(context.Background(), &res)
	if err != nil {
		t.Fatalf("HealthSurvey error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
	survey, ok := res["survey"].(map[string]any)
	if !ok {
		t.Error("expected survey in response")
	} else if survey["status"] != "active" {
		t.Errorf("expected status=active, got %v", survey["status"])
	}
}

func TestHealthActions_UpdateHealthSurvey(t *testing.T) {
	var capturedBody map[string]any
	var capturedQuery string

	mux := http.NewServeMux()
	mux.HandleFunc("/health-survey/test-drive", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		capturedQuery = r.URL.RawQuery
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

	body := map[string]any{"answer": "yes"}
	err := c.Health().UpdateHealthSurvey(context.Background(), body)
	if err != nil {
		t.Fatalf("UpdateHealthSurvey error: %v", err)
	}
	if capturedBody["answer"] != "yes" {
		t.Errorf("expected answer=yes, got %v", capturedBody["answer"])
	}
	if !strings.Contains(capturedQuery, "enableValidation=true") {
		t.Errorf("expected enableValidation=true in query, got %s", capturedQuery)
	}
}

func TestHealthActions_HealthCheckpoints(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/health-integrations/sources/apple-health/checkpoints", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"checkpoints":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	var res map[string]any
	err := c.Health().HealthCheckpoints(context.Background(), "apple-health", &res)
	if err != nil {
		t.Fatalf("HealthCheckpoints error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
	if capturedPath != "/users/uid-123/health-integrations/sources/apple-health/checkpoints" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestHealthActions_UploadHealthData(t *testing.T) {
	var capturedPath string
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/health-integrations/sources/apple-health", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
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

	body := map[string]any{"metrics": []string{"steps", "heartRate"}}
	err := c.Health().UploadHealthData(context.Background(), "apple-health", body)
	if err != nil {
		t.Fatalf("UploadHealthData error: %v", err)
	}
	if capturedPath != "/users/uid-123/health-integrations/sources/apple-health" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if capturedBody["metrics"] == nil {
		t.Error("expected metrics in body")
	}
}
