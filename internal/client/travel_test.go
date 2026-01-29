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

func TestTravelActions_Trips(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/travel/trips", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"trips":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Travel().Trips(context.Background())
	if err != nil {
		t.Fatalf("Trips error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestTravelActions_CreateTrip(t *testing.T) {
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/travel/trips", func(w http.ResponseWriter, r *http.Request) {
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
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	body := map[string]any{"destination": "NYC"}
	err := c.Travel().CreateTrip(context.Background(), body)
	if err != nil {
		t.Fatalf("CreateTrip error: %v", err)
	}
	if capturedBody["destination"] != "NYC" {
		t.Errorf("expected destination=NYC, got %v", capturedBody["destination"])
	}
}

func TestTravelActions_CreatePlan(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/travel/trips/trip-1/plans", func(w http.ResponseWriter, r *http.Request) {
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

	err := c.Travel().CreatePlan(context.Background(), "trip-1", map[string]any{"type": "flight"})
	if err != nil {
		t.Fatalf("CreatePlan error: %v", err)
	}
	if capturedPath != "/users/uid-123/travel/trips/trip-1/plans" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestTravelActions_UpdatePlan(t *testing.T) {
	var capturedPath string
	var capturedMethod string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/travel/plans/plan-1", func(w http.ResponseWriter, r *http.Request) {
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

	err := c.Travel().UpdatePlan(context.Background(), "plan-1", map[string]any{"status": "confirmed"})
	if err != nil {
		t.Fatalf("UpdatePlan error: %v", err)
	}
	if capturedMethod != http.MethodPatch {
		t.Errorf("expected PATCH, got %s", capturedMethod)
	}
	if capturedPath != "/users/uid-123/travel/plans/plan-1" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestTravelActions_DeleteTrip(t *testing.T) {
	var capturedMethod string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/travel/trips/trip-to-delete", func(w http.ResponseWriter, r *http.Request) {
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

	err := c.Travel().DeleteTrip(context.Background(), "trip-to-delete")
	if err != nil {
		t.Fatalf("DeleteTrip error: %v", err)
	}
	if capturedMethod != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", capturedMethod)
	}
}

func TestTravelActions_Plans(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/travel/trips/trip-1/plans", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"plans":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Travel().Plans(context.Background(), "trip-1")
	if err != nil {
		t.Fatalf("Plans error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
	if capturedPath != "/users/uid-123/travel/trips/trip-1/plans" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestTravelActions_PlanTasks(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/travel/plans/plan-1/tasks", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"tasks":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Travel().PlanTasks(context.Background(), "plan-1")
	if err != nil {
		t.Fatalf("PlanTasks error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
	if capturedPath != "/users/uid-123/travel/plans/plan-1/tasks" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestTravelActions_AirportSearch(t *testing.T) {
	var capturedQuery string

	mux := http.NewServeMux()
	mux.HandleFunc("/travel/airport-search", func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"airports":[]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Travel().AirportSearch(context.Background(), "JFK")
	if err != nil {
		t.Fatalf("AirportSearch error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
	if !strings.Contains(capturedQuery, "query=JFK") {
		t.Errorf("expected query=JFK, got %s", capturedQuery)
	}
}

func TestTravelActions_FlightStatus(t *testing.T) {
	var capturedQuery string

	mux := http.NewServeMux()
	mux.HandleFunc("/travel/flight-status", func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"on-time"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Travel().FlightStatus(context.Background(), "AA123")
	if err != nil {
		t.Fatalf("FlightStatus error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
	if !strings.Contains(capturedQuery, "flightNumber=AA123") {
		t.Errorf("expected flightNumber=AA123, got %s", capturedQuery)
	}
}
