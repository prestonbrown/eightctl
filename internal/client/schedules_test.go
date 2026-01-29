package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListSchedules(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/temperature/schedules", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"schedules":[{"id":"s1","startTime":"22:00","level":20,"daysOfWeek":[1,2,3,4,5],"enabled":true}]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	schedules, err := c.ListSchedules(context.Background())
	if err != nil {
		t.Fatalf("ListSchedules error: %v", err)
	}
	if len(schedules) != 1 {
		t.Errorf("expected 1 schedule, got %d", len(schedules))
	}
	if schedules[0].StartTime != "22:00" {
		t.Errorf("expected startTime 22:00, got %s", schedules[0].StartTime)
	}
}

func TestCreateSchedule(t *testing.T) {
	var capturedBody TemperatureSchedule

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/temperature/schedules", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&capturedBody); err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"schedule":{"id":"new-sched","startTime":"23:00","level":30,"daysOfWeek":[0,6],"enabled":true}}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	sched := TemperatureSchedule{
		StartTime:  "23:00",
		Level:      30,
		DaysOfWeek: []int{0, 6},
		Enabled:    true,
	}
	created, err := c.CreateSchedule(context.Background(), sched)
	if err != nil {
		t.Fatalf("CreateSchedule error: %v", err)
	}
	if created.ID != "new-sched" {
		t.Errorf("expected id new-sched, got %s", created.ID)
	}
	if capturedBody.StartTime != "23:00" {
		t.Errorf("expected captured startTime 23:00, got %s", capturedBody.StartTime)
	}
}

func TestUpdateSchedule(t *testing.T) {
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/temperature/schedules/sched-1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&capturedBody); err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"schedule":{"id":"sched-1","startTime":"22:00","level":50,"daysOfWeek":[1,2,3,4,5],"enabled":true}}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	patch := map[string]any{"level": 50}
	updated, err := c.UpdateSchedule(context.Background(), "sched-1", patch)
	if err != nil {
		t.Fatalf("UpdateSchedule error: %v", err)
	}
	if updated.Level != 50 {
		t.Errorf("expected level 50, got %d", updated.Level)
	}
}

func TestDeleteSchedule(t *testing.T) {
	var capturedMethod string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/temperature/schedules/sched-to-delete", func(w http.ResponseWriter, r *http.Request) {
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

	err := c.DeleteSchedule(context.Background(), "sched-to-delete")
	if err != nil {
		t.Fatalf("DeleteSchedule error: %v", err)
	}
	if capturedMethod != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", capturedMethod)
	}
}
