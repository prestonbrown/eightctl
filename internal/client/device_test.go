package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDeviceActions_Info(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-456"}}}`))
	})
	mux.HandleFunc("/devices/dev-456", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"dev-456","model":"Pod 3"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Device().Info(context.Background())
	if err != nil {
		t.Fatalf("Info error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestDeviceActions_Peripherals(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-456"}}}`))
	})
	mux.HandleFunc("/devices/dev-456/peripherals", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"peripherals":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Device().Peripherals(context.Background())
	if err != nil {
		t.Fatalf("Peripherals error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestDeviceActions_Owner(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-456"}}}`))
	})
	mux.HandleFunc("/devices/dev-456/owner", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"owner":{"email":"owner@example.com"}}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Device().Owner(context.Background())
	if err != nil {
		t.Fatalf("Owner error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestDeviceActions_Warranty(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-456"}}}`))
	})
	mux.HandleFunc("/devices/dev-456/warranty", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"warranty":{"expires":"2026-01-01"}}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Device().Warranty(context.Background())
	if err != nil {
		t.Fatalf("Warranty error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestDeviceActions_Online(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-456"}}}`))
	})
	mux.HandleFunc("/devices/dev-456/online", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"online":true}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Device().Online(context.Background())
	if err != nil {
		t.Fatalf("Online error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestDeviceActions_PrimingTasks(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-456"}}}`))
	})
	mux.HandleFunc("/devices/dev-456/priming/tasks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"tasks":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Device().PrimingTasks(context.Background())
	if err != nil {
		t.Fatalf("PrimingTasks error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestDeviceActions_PrimingSchedule(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-456"}}}`))
	})
	mux.HandleFunc("/devices/dev-456/priming/schedule", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"schedule":{}}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Device().PrimingSchedule(context.Background())
	if err != nil {
		t.Fatalf("PrimingSchedule error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestDeviceActions_Update(t *testing.T) {
	var capturedPath string
	var capturedMethod string
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-456"}}}`))
	})
	mux.HandleFunc("/devices/dev-456", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedMethod = r.Method
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &capturedBody)
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Device().Update(context.Background(), map[string]any{"name": "My Pod"})
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if capturedMethod != http.MethodPut {
		t.Errorf("expected PUT, got %s", capturedMethod)
	}
	if capturedPath != "/devices/dev-456" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if capturedBody["name"] != "My Pod" {
		t.Errorf("unexpected body: %v", capturedBody)
	}
}

func TestDeviceActions_SetOwner(t *testing.T) {
	var capturedPath string
	var capturedMethod string
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-456"}}}`))
	})
	mux.HandleFunc("/devices/dev-456/owner", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedMethod = r.Method
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &capturedBody)
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Device().SetOwner(context.Background(), map[string]any{"userId": "new-owner"})
	if err != nil {
		t.Fatalf("SetOwner error: %v", err)
	}
	if capturedMethod != http.MethodPut {
		t.Errorf("expected PUT, got %s", capturedMethod)
	}
	if capturedPath != "/devices/dev-456/owner" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if capturedBody["userId"] != "new-owner" {
		t.Errorf("unexpected body: %v", capturedBody)
	}
}

func TestDeviceActions_SetPeripherals(t *testing.T) {
	var capturedPath string
	var capturedMethod string
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-456"}}}`))
	})
	mux.HandleFunc("/devices/dev-456/peripherals", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedMethod = r.Method
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &capturedBody)
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Device().SetPeripherals(context.Background(), map[string]any{"peripherals": []string{"p1"}})
	if err != nil {
		t.Fatalf("SetPeripherals error: %v", err)
	}
	if capturedMethod != http.MethodPut {
		t.Errorf("expected PUT, got %s", capturedMethod)
	}
	if capturedPath != "/devices/dev-456/peripherals" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if capturedBody["peripherals"] == nil {
		t.Errorf("unexpected body: %v", capturedBody)
	}
}

func TestDeviceActions_AddPeripheral(t *testing.T) {
	var capturedPath string
	var capturedMethod string
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-456"}}}`))
	})
	mux.HandleFunc("/devices/dev-456/peripherals", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedMethod = r.Method
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &capturedBody)
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Device().AddPeripheral(context.Background(), map[string]any{"peripheral": "new-p"})
	if err != nil {
		t.Fatalf("AddPeripheral error: %v", err)
	}
	if capturedMethod != http.MethodPatch {
		t.Errorf("expected PATCH, got %s", capturedMethod)
	}
	if capturedPath != "/devices/dev-456/peripherals" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if capturedBody["peripheral"] != "new-p" {
		t.Errorf("unexpected body: %v", capturedBody)
	}
}

func TestDeviceActions_GetBLEKey(t *testing.T) {
	var capturedPath string
	var capturedMethod string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-456"}}}`))
	})
	mux.HandleFunc("/devices/dev-456/security/key", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedMethod = r.Method
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"key":"ble-key-123"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	var result any
	err := c.Device().GetBLEKey(context.Background(), &result)
	if err != nil {
		t.Fatalf("GetBLEKey error: %v", err)
	}
	if capturedMethod != http.MethodPost {
		t.Errorf("expected POST, got %s", capturedMethod)
	}
	if capturedPath != "/devices/dev-456/security/key" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if result == nil {
		t.Error("expected non-nil response")
	}
}

func TestDeviceActions_UpdatePrimingSchedule(t *testing.T) {
	var capturedPath string
	var capturedMethod string
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-456"}}}`))
	})
	mux.HandleFunc("/devices/dev-456/priming/schedule", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedMethod = r.Method
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &capturedBody)
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Device().UpdatePrimingSchedule(context.Background(), map[string]any{"enabled": true})
	if err != nil {
		t.Fatalf("UpdatePrimingSchedule error: %v", err)
	}
	if capturedMethod != http.MethodPut {
		t.Errorf("expected PUT, got %s", capturedMethod)
	}
	if capturedPath != "/devices/dev-456/priming/schedule" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if capturedBody["enabled"] != true {
		t.Errorf("unexpected body: %v", capturedBody)
	}
}

func TestDeviceActions_CreatePrimingTask(t *testing.T) {
	var capturedPath string
	var capturedMethod string
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-456"}}}`))
	})
	mux.HandleFunc("/devices/dev-456/priming/tasks", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedMethod = r.Method
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &capturedBody)
		w.WriteHeader(http.StatusCreated)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Device().CreatePrimingTask(context.Background(), map[string]any{"type": "prime"})
	if err != nil {
		t.Fatalf("CreatePrimingTask error: %v", err)
	}
	if capturedMethod != http.MethodPost {
		t.Errorf("expected POST, got %s", capturedMethod)
	}
	if capturedPath != "/devices/dev-456/priming/tasks" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if capturedBody["type"] != "prime" {
		t.Errorf("unexpected body: %v", capturedBody)
	}
}

func TestDeviceActions_CancelPrimingTask(t *testing.T) {
	var capturedPath string
	var capturedMethod string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-456"}}}`))
	})
	mux.HandleFunc("/devices/dev-456/priming/tasks", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Device().CancelPrimingTask(context.Background())
	if err != nil {
		t.Fatalf("CancelPrimingTask error: %v", err)
	}
	if capturedMethod != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", capturedMethod)
	}
	if capturedPath != "/devices/dev-456/priming/tasks" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}
