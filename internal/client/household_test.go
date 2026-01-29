package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHouseholdActions_Summary(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/users/uid-123/summary", func(w http.ResponseWriter, r *http.Request) {
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

	res, err := c.Household().Summary(context.Background())
	if err != nil {
		t.Fatalf("Summary error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestHouseholdActions_Schedule(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/users/uid-123/schedule", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"schedule":{}}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Household().Schedule(context.Background())
	if err != nil {
		t.Fatalf("Schedule error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestHouseholdActions_CurrentSet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/users/uid-123/current-set", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"currentSet":{}}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Household().CurrentSet(context.Background())
	if err != nil {
		t.Fatalf("CurrentSet error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestHouseholdActions_Invitations(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/users/uid-123/invitations", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"invitations":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Household().Invitations(context.Background())
	if err != nil {
		t.Fatalf("Invitations error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestHouseholdActions_Devices(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/users/uid-123/devices", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"devices":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Household().Devices(context.Background())
	if err != nil {
		t.Fatalf("Devices error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestHouseholdActions_Users(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/users/uid-123/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"users":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Household().Users(context.Background())
	if err != nil {
		t.Fatalf("Users error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestHouseholdActions_Guests(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/users/uid-123/guests", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"guests":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Household().Guests(context.Background())
	if err != nil {
		t.Fatalf("Guests error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}
