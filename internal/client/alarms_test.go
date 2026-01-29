package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAlarmActions_Snooze(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/alarms/a1/snooze", func(w http.ResponseWriter, r *http.Request) {
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

	err := c.Alarms().Snooze(context.Background(), "a1")
	if err != nil {
		t.Fatalf("Snooze error: %v", err)
	}
	if capturedPath != "/users/uid-123/alarms/a1/snooze" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestAlarmActions_Dismiss(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/alarms/a1/dismiss", func(w http.ResponseWriter, r *http.Request) {
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

	err := c.Alarms().Dismiss(context.Background(), "a1")
	if err != nil {
		t.Fatalf("Dismiss error: %v", err)
	}
	if capturedPath != "/users/uid-123/alarms/a1/dismiss" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestAlarmActions_DismissAll(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/alarms/active/dismiss-all", func(w http.ResponseWriter, r *http.Request) {
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

	err := c.Alarms().DismissAll(context.Background())
	if err != nil {
		t.Fatalf("DismissAll error: %v", err)
	}
	if capturedPath != "/users/uid-123/alarms/active/dismiss-all" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestAlarmActions_VibrationTest(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/vibration-test", func(w http.ResponseWriter, r *http.Request) {
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

	err := c.Alarms().VibrationTest(context.Background())
	if err != nil {
		t.Fatalf("VibrationTest error: %v", err)
	}
	if capturedPath != "/users/uid-123/vibration-test" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}
