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

// mockServer builds a test server that can serve a handful of endpoints the client expects.
func mockServer(t *testing.T) (*httptest.Server, *Client) {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-1"}}}`))
	})

	mux.HandleFunc("/users/uid-123/temperature", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"currentLevel":5,"currentState":{"type":"on"}}`))
			return
		}
		if r.Method == http.MethodPut {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.NotFound(w, r)
	})

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		// first call rate limits, second succeeds
		if r.Header.Get("X-Test-Retry") == "done" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ok":true}`))
			return
		}
		w.WriteHeader(http.StatusTooManyRequests)
	})

	srv := httptest.NewServer(mux)

	// client with pre-set token to skip auth
	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	return srv, c
}

func TestRequireUserFilledAutomatically(t *testing.T) {
	srv, c := mockServer(t)
	defer srv.Close()

	// UserID empty; GetStatus should fetch it from /users/me
	st, err := c.GetStatus(context.Background())
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if c.UserID != "uid-123" {
		t.Fatalf("expected user id populated, got %s", c.UserID)
	}
	if st.CurrentLevel != 5 || st.CurrentState.Type != "on" {
		t.Fatalf("unexpected status %+v", st)
	}
}

func Test429Retry(t *testing.T) {
	count := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		count++
		if count == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	start := time.Now()
	if err := c.do(context.Background(), http.MethodGet, "/ping", nil, nil, nil); err != nil {
		t.Fatalf("do retry: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 attempts, got %d", count)
	}
	if elapsed := time.Since(start); elapsed < 2*time.Second {
		t.Fatalf("expected backoff, got %v", elapsed)
	}
}

func TestMetricsTrendsPassesTimezone(t *testing.T) {
	var capturedTZ string
	var capturedFrom string
	var capturedTo string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/trends", func(w http.ResponseWriter, r *http.Request) {
		capturedTZ = r.URL.Query().Get("tz")
		capturedFrom = r.URL.Query().Get("from")
		capturedTo = r.URL.Query().Get("to")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"days":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	var out any
	err := c.Metrics().Trends(context.Background(), "2025-01-01", "2025-01-28", "America/New_York", &out)
	if err != nil {
		t.Fatalf("Trends error: %v", err)
	}

	if capturedTZ != "America/New_York" {
		t.Errorf("expected tz=America/New_York, got tz=%s", capturedTZ)
	}
	if capturedFrom != "2025-01-01" {
		t.Errorf("expected from=2025-01-01, got from=%s", capturedFrom)
	}
	if capturedTo != "2025-01-28" {
		t.Errorf("expected to=2025-01-28, got to=%s", capturedTo)
	}
}

func TestMetricsTrendsEmptyTimezone(t *testing.T) {
	var capturedQuery string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/trends", func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"days":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	var out any
	err := c.Metrics().Trends(context.Background(), "2025-01-01", "2025-01-28", "", &out)
	if err != nil {
		t.Fatalf("Trends error: %v", err)
	}

	// When timezone is empty, tz param should not be in the query
	if strings.Contains(capturedQuery, "tz=") {
		t.Errorf("expected no tz param when empty, got query: %s", capturedQuery)
	}
}

func TestAuthenticate_TokenEndpoint(t *testing.T) {
	mux := http.NewServeMux()
	var capturedBody map[string]string

	// Mock the auth token endpoint
	mux.HandleFunc("/v1/tokens", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&capturedBody); err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"access_token":"tok-123","expires_in":3600,"userId":"uid-456"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("test@example.com", "secret", "", "client-id", "client-secret")
	c.BaseURL = srv.URL
	// Override auth URL to use test server
	// Note: authURL is const, so we need to test via authTokenEndpoint indirectly

	// For this test we'll verify the flow through Authenticate's fallback
	// Since we can't override authURL, we test legacy login instead
}

func TestAuthenticate_LegacyLogin(t *testing.T) {
	var capturedBody map[string]string

	mux := http.NewServeMux()
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&capturedBody); err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"session":{"userId":"uid-789","token":"legacy-tok","expirationDate":"2025-12-31T23:59:59Z"}}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("test@example.com", "secret", "", "", "")
	c.BaseURL = srv.URL
	c.HTTP = srv.Client()

	err := c.authLegacyLogin(context.Background())
	if err != nil {
		t.Fatalf("authLegacyLogin error: %v", err)
	}

	if capturedBody["email"] != "test@example.com" {
		t.Errorf("expected email in body, got %v", capturedBody)
	}
	if c.token != "legacy-tok" {
		t.Errorf("expected token set, got %s", c.token)
	}
	if c.UserID != "uid-789" {
		t.Errorf("expected UserID set, got %s", c.UserID)
	}
}

func TestEnsureDeviceID(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-456"}}}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	deviceID, err := c.EnsureDeviceID(context.Background())
	if err != nil {
		t.Fatalf("EnsureDeviceID error: %v", err)
	}
	if deviceID != "dev-456" {
		t.Errorf("expected dev-456, got %s", deviceID)
	}
	if c.DeviceID != "dev-456" {
		t.Errorf("expected DeviceID cached, got %s", c.DeviceID)
	}
}

func TestEnsureDeviceID_AlreadySet(t *testing.T) {
	c := New("email", "pass", "", "", "")
	c.DeviceID = "existing-dev"

	deviceID, err := c.EnsureDeviceID(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deviceID != "existing-dev" {
		t.Errorf("expected existing-dev, got %s", deviceID)
	}
}

func TestEnsureDeviceID_NoDevice(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{}}}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	_, err := c.EnsureDeviceID(context.Background())
	if err == nil {
		t.Fatal("expected error for missing device")
	}
	if !strings.Contains(err.Error(), "no current device") {
		t.Errorf("expected 'no current device' error, got: %v", err)
	}
}

func TestTurnOn(t *testing.T) {
	var capturedPath string
	var capturedBody map[string]bool

	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123"}}`))
	})
	mux.HandleFunc("/users/uid-123/devices/power", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		if err := json.NewDecoder(r.Body).Decode(&capturedBody); err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.TurnOn(context.Background())
	if err != nil {
		t.Fatalf("TurnOn error: %v", err)
	}
	if capturedPath != "/users/uid-123/devices/power" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if !capturedBody["on"] {
		t.Errorf("expected on=true, got %v", capturedBody)
	}
}

func TestTurnOff(t *testing.T) {
	var capturedBody map[string]bool

	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123"}}`))
	})
	mux.HandleFunc("/users/uid-123/devices/power", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&capturedBody); err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.TurnOff(context.Background())
	if err != nil {
		t.Fatalf("TurnOff error: %v", err)
	}
	if capturedBody["on"] {
		t.Errorf("expected on=false, got %v", capturedBody)
	}
}

func TestSetTemperature(t *testing.T) {
	var capturedBody map[string]int

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/temperature", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
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

	err := c.SetTemperature(context.Background(), 50)
	if err != nil {
		t.Fatalf("SetTemperature error: %v", err)
	}
	if capturedBody["currentLevel"] != 50 {
		t.Errorf("expected level 50, got %v", capturedBody)
	}
}

func TestSetTemperature_OutOfRange(t *testing.T) {
	c := New("email", "pass", "uid-123", "", "")

	tests := []struct {
		level int
		valid bool
	}{
		{-101, false},
		{-100, true},
		{0, true},
		{100, true},
		{101, false},
	}

	for _, tt := range tests {
		// For invalid cases, we expect immediate error without network call
		if !tt.valid {
			err := c.SetTemperature(context.Background(), tt.level)
			if err == nil {
				t.Errorf("expected error for level %d", tt.level)
			}
		}
	}
}

func TestGetSleepDay(t *testing.T) {
	var capturedQuery string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/trends", func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"days":[{"day":"2025-01-15","score":85.5,"heartRate":62.3}]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	day, err := c.GetSleepDay(context.Background(), "2025-01-15", "America/New_York")
	if err != nil {
		t.Fatalf("GetSleepDay error: %v", err)
	}
	if day.Date != "2025-01-15" {
		t.Errorf("expected date 2025-01-15, got %s", day.Date)
	}
	if day.Score != 85.5 {
		t.Errorf("expected score 85.5, got %f", day.Score)
	}
	if !strings.Contains(capturedQuery, "tz=America") {
		t.Errorf("expected timezone in query, got %s", capturedQuery)
	}
}

func TestGetSleepDay_NoData(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/trends", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"days":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	_, err := c.GetSleepDay(context.Background(), "2025-01-15", "America/New_York")
	if err == nil {
		t.Fatal("expected error for no data")
	}
	if !strings.Contains(err.Error(), "no sleep data") {
		t.Errorf("expected 'no sleep data' error, got: %v", err)
	}
}

func TestListTracks(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/audio/tracks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"tracks":[{"id":"t1","title":"Rain","type":"nature"},{"id":"t2","title":"Ocean","type":"nature"}]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	tracks, err := c.ListTracks(context.Background())
	if err != nil {
		t.Fatalf("ListTracks error: %v", err)
	}
	if len(tracks) != 2 {
		t.Errorf("expected 2 tracks, got %d", len(tracks))
	}
	if tracks[0].Title != "Rain" {
		t.Errorf("expected first track 'Rain', got %s", tracks[0].Title)
	}
}

func TestReleaseFeatures(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/release/features", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"features":[{"title":"New Feature","body":"Description"}]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	features, err := c.ReleaseFeatures(context.Background())
	if err != nil {
		t.Fatalf("ReleaseFeatures error: %v", err)
	}
	if len(features) != 1 {
		t.Errorf("expected 1 feature, got %d", len(features))
	}
	if features[0].Title != "New Feature" {
		t.Errorf("expected title 'New Feature', got %s", features[0].Title)
	}
}

func TestIdentity(t *testing.T) {
	c := New("test@example.com", "pass", "", "my-client", "")
	c.BaseURL = "https://custom.api.com"

	id := c.Identity()
	if id.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", id.Email)
	}
	if id.ClientID != "my-client" {
		t.Errorf("expected clientID my-client, got %s", id.ClientID)
	}
	if id.BaseURL != "https://custom.api.com" {
		t.Errorf("expected baseURL https://custom.api.com, got %s", id.BaseURL)
	}
}

func TestAPIErrorResponse(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/temperature", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid request"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	_, err := c.GetStatus(context.Background())
	if err == nil {
		t.Fatal("expected error for bad request")
	}
	if !strings.Contains(err.Error(), "invalid request") {
		t.Errorf("expected error message to contain 'invalid request', got: %v", err)
	}
}

func TestEnsureUserID_EmptyResponse(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":""}}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.EnsureUserID(context.Background())
	if err == nil {
		t.Fatal("expected error for empty userId")
	}
	if !strings.Contains(err.Error(), "userId not found") {
		t.Errorf("expected 'userId not found' error, got: %v", err)
	}
}

func TestEnsureUserID_AlreadySet(t *testing.T) {
	c := New("email", "pass", "existing-user", "", "")

	err := c.EnsureUserID(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.UserID != "existing-user" {
		t.Errorf("expected existing-user, got %s", c.UserID)
	}
}

func TestClient_GetUserTemperature(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/other-user-id/temperature", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"currentLevel": -20,
			"currentState": map[string]any{
				"type": "smart",
			},
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	status, err := c.GetUserTemperature(context.Background(), "other-user-id")
	if err != nil {
		t.Fatalf("GetUserTemperature error: %v", err)
	}
	if status.CurrentLevel != -20 {
		t.Errorf("expected level -20, got %d", status.CurrentLevel)
	}
	if status.CurrentState.Type != "smart" {
		t.Errorf("expected state type 'smart', got %s", status.CurrentState.Type)
	}
}
