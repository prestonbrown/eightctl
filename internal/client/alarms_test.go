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
	var capturedMethod string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/routines", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.AppAPIBaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Alarms().Snooze(context.Background(), "a1")
	if err != nil {
		t.Fatalf("Snooze error: %v", err)
	}
	if capturedPath != "/users/uid-123/routines" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if capturedMethod != http.MethodPut {
		t.Errorf("expected PUT, got %s", capturedMethod)
	}
}

func TestAlarmActions_Dismiss(t *testing.T) {
	var capturedPath string
	var capturedMethod string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/routines", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.AppAPIBaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Alarms().Dismiss(context.Background(), "a1")
	if err != nil {
		t.Fatalf("Dismiss error: %v", err)
	}
	if capturedPath != "/users/uid-123/routines" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if capturedMethod != http.MethodPut {
		t.Errorf("expected PUT, got %s", capturedMethod)
	}
}

func TestAlarmActions_DismissAll(t *testing.T) {
	var capturedPath string
	var capturedMethod string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/routines", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.AppAPIBaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Alarms().DismissAll(context.Background())
	if err != nil {
		t.Fatalf("DismissAll error: %v", err)
	}
	if capturedPath != "/users/uid-123/routines" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if capturedMethod != http.MethodPut {
		t.Errorf("expected PUT, got %s", capturedMethod)
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
	c.AppAPIBaseURL = srv.URL
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
