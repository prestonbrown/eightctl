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

func TestSettingsActions_TapSettings(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/devices/dev-456/tap-settings", func(w http.ResponseWriter, r *http.Request) {
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

	res, err := c.Settings().TapSettings(context.Background(), "dev-456")
	if err != nil {
		t.Fatalf("TapSettings error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
	if capturedPath != "/users/uid-123/devices/dev-456/tap-settings" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestSettingsActions_UpdateTapSettings(t *testing.T) {
	var capturedBody map[string]any
	var capturedMethod string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/devices/dev-456/tap-settings", func(w http.ResponseWriter, r *http.Request) {
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
	err := c.Settings().UpdateTapSettings(context.Background(), "dev-456", body)
	if err != nil {
		t.Fatalf("UpdateTapSettings error: %v", err)
	}
	if capturedMethod != http.MethodPut {
		t.Errorf("expected PUT, got %s", capturedMethod)
	}
	if capturedBody["enabled"] != false {
		t.Errorf("expected enabled=false, got %v", capturedBody["enabled"])
	}
}

func TestSettingsActions_TapHistory(t *testing.T) {
	var capturedQuery string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/tap-history", func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
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

	res, err := c.Settings().TapHistory(context.Background(), "2024-01-01")
	if err != nil {
		t.Fatalf("TapHistory error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
	if !strings.Contains(capturedQuery, "from=2024-01-01") {
		t.Errorf("expected from=2024-01-01, got %s", capturedQuery)
	}
}

func TestSettingsActions_LevelSuggestions(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/level-suggestions", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"suggestions":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Settings().LevelSuggestions(context.Background())
	if err != nil {
		t.Fatalf("LevelSuggestions error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
	if capturedPath != "/users/uid-123/level-suggestions" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestSettingsActions_BlanketRecommendations(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/recommendations/blanket", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"recommendations":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Settings().BlanketRecommendations(context.Background())
	if err != nil {
		t.Fatalf("BlanketRecommendations error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
	if capturedPath != "/users/uid-123/recommendations/blanket" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestSettingsActions_Perks(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/perks", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"perks":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Settings().Perks(context.Background())
	if err != nil {
		t.Fatalf("Perks error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
	if capturedPath != "/users/uid-123/perks" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestSettingsActions_GetReferralLink(t *testing.T) {
	var capturedPath string
	var capturedMethod string

	mux := http.NewServeMux()
	mux.HandleFunc("/v2/users/uid-123/referral/personal-referral-link", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedMethod = r.Method
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"link":"https://eight.sl/abc123"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Settings().GetReferralLink(context.Background())
	if err != nil {
		t.Fatalf("GetReferralLink error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
	if capturedMethod != http.MethodPut {
		t.Errorf("expected PUT, got %s", capturedMethod)
	}
	if capturedPath != "/v2/users/uid-123/referral/personal-referral-link" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestSettingsActions_ReferralCampaigns(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/v2/users/uid-123/referral/campaigns", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"campaigns":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Settings().ReferralCampaigns(context.Background())
	if err != nil {
		t.Fatalf("ReferralCampaigns error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
	if capturedPath != "/v2/users/uid-123/referral/campaigns" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestSettingsActions_Purchases(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/purchase-tracker", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"purchases":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Settings().Purchases(context.Background())
	if err != nil {
		t.Fatalf("Purchases error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
	if capturedPath != "/purchase-tracker" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestSettingsActions_MaintenanceInsertStatus(t *testing.T) {
	var capturedPath string
	var capturedQuery string

	mux := http.NewServeMux()
	mux.HandleFunc("/user/uid-123/device_maintenance/maintenance_insert", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedQuery = r.URL.RawQuery
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Settings().MaintenanceInsertStatus(context.Background())
	if err != nil {
		t.Fatalf("MaintenanceInsertStatus error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
	if capturedPath != "/user/uid-123/device_maintenance/maintenance_insert" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if !strings.Contains(capturedQuery, "v=2") {
		t.Errorf("expected v=2, got %s", capturedQuery)
	}
}
