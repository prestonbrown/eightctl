package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestChallengesActions_Challenges(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/challenges", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"challenges":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	var res map[string]any
	err := c.Challenges().Challenges(context.Background(), "", &res)
	if err != nil {
		t.Fatalf("Challenges error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestChallengesActions_ChallengesWithState(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/challenges", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != "active" {
			t.Errorf("expected state=active, got %q", r.URL.Query().Get("state"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"challenges":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	var res map[string]any
	err := c.Challenges().Challenges(context.Background(), "active", &res)
	if err != nil {
		t.Fatalf("Challenges error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}
