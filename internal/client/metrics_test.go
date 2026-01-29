package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestMetricsActions_Intervals(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/intervals/session-abc", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"intervals":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	var out any
	err := c.Metrics().Intervals(context.Background(), "session-abc", &out)
	if err != nil {
		t.Fatalf("Intervals error: %v", err)
	}
	if capturedPath != "/users/uid-123/intervals/session-abc" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestMetricsActions_Summary(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/metrics/summary", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"summary":{}}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	var out any
	err := c.Metrics().Summary(context.Background(), &out)
	if err != nil {
		t.Fatalf("Summary error: %v", err)
	}
}

func TestMetricsActions_Aggregate(t *testing.T) {
	var capturedQuery string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/metrics/aggregate", func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"aggregate":{}}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	var out any
	err := c.Metrics().Aggregate(context.Background(), &out)
	if err != nil {
		t.Fatalf("Aggregate error: %v", err)
	}
	if !strings.Contains(capturedQuery, "v2=true") {
		t.Errorf("expected v2=true in query, got %s", capturedQuery)
	}
}

func TestMetricsActions_Insights(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/insights", func(w http.ResponseWriter, r *http.Request) {
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

	var out any
	err := c.Metrics().Insights(context.Background(), &out)
	if err != nil {
		t.Fatalf("Insights error: %v", err)
	}
}

func TestMetricsActions_UpdateSession(t *testing.T) {
	var capturedPath, capturedMethod string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/intervals/session-abc", func(w http.ResponseWriter, r *http.Request) {
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

	body := map[string]any{"note": "test note"}
	err := c.Metrics().UpdateSession(context.Background(), "session-abc", body)
	if err != nil {
		t.Fatalf("UpdateSession error: %v", err)
	}
	if capturedPath != "/users/uid-123/intervals/session-abc" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if capturedMethod != http.MethodPut {
		t.Errorf("expected PUT, got %s", capturedMethod)
	}
}

func TestMetricsActions_DeleteSession(t *testing.T) {
	var capturedPath, capturedMethod string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/intervals/session-abc", func(w http.ResponseWriter, r *http.Request) {
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

	err := c.Metrics().DeleteSession(context.Background(), "session-abc")
	if err != nil {
		t.Fatalf("DeleteSession error: %v", err)
	}
	if capturedPath != "/users/uid-123/intervals/session-abc" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if capturedMethod != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", capturedMethod)
	}
}

func TestMetricsActions_SendFeedback(t *testing.T) {
	var capturedPath, capturedMethod string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/feedback", func(w http.ResponseWriter, r *http.Request) {
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

	body := map[string]any{"rating": 5, "comment": "Great!"}
	err := c.Metrics().SendFeedback(context.Background(), body)
	if err != nil {
		t.Fatalf("SendFeedback error: %v", err)
	}
	if capturedPath != "/users/uid-123/feedback" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if capturedMethod != http.MethodPost {
		t.Errorf("expected POST, got %s", capturedMethod)
	}
}
