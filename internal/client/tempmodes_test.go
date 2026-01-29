package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestTempModes_NapActivate(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/temperature/nap-mode/activate", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
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

	err := c.TempModes().NapActivate(context.Background())
	if err != nil {
		t.Fatalf("NapActivate error: %v", err)
	}
	if capturedPath != "/users/uid-123/temperature/nap-mode/activate" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestTempModes_NapDeactivate(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/temperature/nap-mode/deactivate", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.TempModes().NapDeactivate(context.Background())
	if err != nil {
		t.Fatalf("NapDeactivate error: %v", err)
	}
	if capturedPath != "/users/uid-123/temperature/nap-mode/deactivate" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestTempModes_NapExtend(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/temperature/nap-mode/extend", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.TempModes().NapExtend(context.Background())
	if err != nil {
		t.Fatalf("NapExtend error: %v", err)
	}
	if capturedPath != "/users/uid-123/temperature/nap-mode/extend" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestTempModes_NapStatus(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/temperature/nap-mode/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"active":false}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	var out any
	err := c.TempModes().NapStatus(context.Background(), &out)
	if err != nil {
		t.Fatalf("NapStatus error: %v", err)
	}
}

func TestTempModes_HotFlashActivate(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/temperature/hot-flash-mode/activate", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.TempModes().HotFlashActivate(context.Background())
	if err != nil {
		t.Fatalf("HotFlashActivate error: %v", err)
	}
	if capturedPath != "/users/uid-123/temperature/hot-flash-mode/activate" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestTempModes_HotFlashDeactivate(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/temperature/hot-flash-mode/deactivate", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.TempModes().HotFlashDeactivate(context.Background())
	if err != nil {
		t.Fatalf("HotFlashDeactivate error: %v", err)
	}
	if capturedPath != "/users/uid-123/temperature/hot-flash-mode/deactivate" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestTempModes_HotFlashStatus(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/temperature/hot-flash-mode", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"active":false}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	var out any
	err := c.TempModes().HotFlashStatus(context.Background(), &out)
	if err != nil {
		t.Fatalf("HotFlashStatus error: %v", err)
	}
}

func TestTempModes_TempEvents(t *testing.T) {
	var capturedQuery string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/temp-events", func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"events":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	var out any
	err := c.TempModes().TempEvents(context.Background(), "2025-01-01", "2025-01-28", &out)
	if err != nil {
		t.Fatalf("TempEvents error: %v", err)
	}
	if !strings.Contains(capturedQuery, "from=2025-01-01") {
		t.Errorf("expected from in query, got %s", capturedQuery)
	}
	if !strings.Contains(capturedQuery, "to=2025-01-28") {
		t.Errorf("expected to in query, got %s", capturedQuery)
	}
}

func TestTempModes_TempEvents_EmptyDates(t *testing.T) {
	var capturedQuery string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/temp-events", func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"events":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	var out any
	err := c.TempModes().TempEvents(context.Background(), "", "", &out)
	if err != nil {
		t.Fatalf("TempEvents error: %v", err)
	}
	if strings.Contains(capturedQuery, "from=") {
		t.Errorf("expected no from when empty, got %s", capturedQuery)
	}
}
