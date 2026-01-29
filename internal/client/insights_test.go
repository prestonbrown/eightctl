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

func TestInsightsActions_LLMInsights(t *testing.T) {
	var capturedQuery string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/llm-insights", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		capturedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"insights":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Insights().LLMInsights(context.Background(), "2024-01-01", "2024-01-31")
	if err != nil {
		t.Fatalf("LLMInsights error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
	if !strings.Contains(capturedQuery, "from=2024-01-01") {
		t.Errorf("expected from=2024-01-01 in query, got %s", capturedQuery)
	}
	if !strings.Contains(capturedQuery, "to=2024-01-31") {
		t.Errorf("expected to=2024-01-31 in query, got %s", capturedQuery)
	}
}

func TestInsightsActions_CreateLLMInsightsBatch(t *testing.T) {
	var capturedBody map[string]any
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/llm-insights/batch", func(w http.ResponseWriter, r *http.Request) {
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

	body := map[string]any{"dates": []string{"2024-01-01", "2024-01-02"}}
	err := c.Insights().CreateLLMInsightsBatch(context.Background(), body)
	if err != nil {
		t.Fatalf("CreateLLMInsightsBatch error: %v", err)
	}
	if capturedPath != "/users/uid-123/llm-insights/batch" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if capturedBody["dates"] == nil {
		t.Error("expected dates in body")
	}
}

func TestInsightsActions_LLMInsightsSettings(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/llm-insights/settings", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"enabled":true}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Insights().LLMInsightsSettings(context.Background())
	if err != nil {
		t.Fatalf("LLMInsightsSettings error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
	if capturedPath != "/users/uid-123/llm-insights/settings" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestInsightsActions_UpdateLLMInsightsSettings(t *testing.T) {
	var capturedBody map[string]any
	var capturedMethod string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/llm-insights/settings", func(w http.ResponseWriter, r *http.Request) {
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

	body := map[string]any{"enabled": false}
	err := c.Insights().UpdateLLMInsightsSettings(context.Background(), body)
	if err != nil {
		t.Fatalf("UpdateLLMInsightsSettings error: %v", err)
	}
	if capturedMethod != http.MethodPut {
		t.Errorf("expected PUT, got %s", capturedMethod)
	}
	if capturedBody["enabled"] != false {
		t.Errorf("expected enabled=false, got %v", capturedBody["enabled"])
	}
}

func TestInsightsActions_SubmitLLMInsightFeedback(t *testing.T) {
	var capturedPath string
	var capturedBody map[string]any
	var capturedMethod string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/llm-insights/insight-456/feedback", func(w http.ResponseWriter, r *http.Request) {
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

	body := map[string]any{"rating": "helpful"}
	err := c.Insights().SubmitLLMInsightFeedback(context.Background(), "insight-456", body)
	if err != nil {
		t.Fatalf("SubmitLLMInsightFeedback error: %v", err)
	}
	if capturedMethod != http.MethodPost {
		t.Errorf("expected POST, got %s", capturedMethod)
	}
	if capturedPath != "/users/uid-123/llm-insights/insight-456/feedback" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if capturedBody["rating"] != "helpful" {
		t.Errorf("expected rating=helpful, got %v", capturedBody["rating"])
	}
}
