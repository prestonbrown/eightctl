package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBaseActions_Info(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/base", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"model":"Pod Pro","firmware":"1.2.3"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.AppAPIBaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Base().Info(context.Background())
	if err != nil {
		t.Fatalf("Info error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestBaseActions_SetAngle(t *testing.T) {
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/base/angle", func(w http.ResponseWriter, r *http.Request) {
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
	c.AppAPIBaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Base().SetAngle(context.Background(), 30, 15)
	if err != nil {
		t.Fatalf("SetAngle error: %v", err)
	}
	if capturedBody["torsoAngle"] != float64(30) {
		t.Errorf("expected torsoAngle=30, got %v", capturedBody["torsoAngle"])
	}
	if capturedBody["legAngle"] != float64(15) {
		t.Errorf("expected legAngle=15, got %v", capturedBody["legAngle"])
	}
}

func TestBaseActions_Presets(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v2/users/uid-123/base/presets", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"presets":[{"name":"reading"},{"name":"tv"}]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.AppAPIBaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Base().Presets(context.Background())
	if err != nil {
		t.Fatalf("Presets error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestBaseActions_RunPreset(t *testing.T) {
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/base/presets", func(w http.ResponseWriter, r *http.Request) {
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
	c.AppAPIBaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Base().RunPreset(context.Background(), "reading")
	if err != nil {
		t.Fatalf("RunPreset error: %v", err)
	}
	if capturedBody["name"] != "reading" {
		t.Errorf("expected name=reading, got %v", capturedBody["name"])
	}
}

func TestBaseActions_VibrationTest(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-456"}}}`))
	})
	mux.HandleFunc("/devices/dev-456/vibration-test", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.AppAPIBaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Base().VibrationTest(context.Background())
	if err != nil {
		t.Fatalf("VibrationTest error: %v", err)
	}
	if capturedPath != "/devices/dev-456/vibration-test" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}
