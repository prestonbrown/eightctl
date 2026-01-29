package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListAlarms(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/alarms", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"alarms":[{"id":"a1","enabled":true,"time":"07:00","daysOfWeek":[1,2,3,4,5],"vibration":true}]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	alarms, err := c.ListAlarms(context.Background())
	if err != nil {
		t.Fatalf("ListAlarms error: %v", err)
	}
	if len(alarms) != 1 {
		t.Errorf("expected 1 alarm, got %d", len(alarms))
	}
	if alarms[0].Time != "07:00" {
		t.Errorf("expected time 07:00, got %s", alarms[0].Time)
	}
}

func TestCreateAlarm(t *testing.T) {
	var capturedBody Alarm

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/alarms", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&capturedBody); err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"alarm":{"id":"new-alarm","enabled":true,"time":"08:00","daysOfWeek":[1,2,3,4,5],"vibration":true}}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	alarm := Alarm{
		Enabled:    true,
		Time:       "08:00",
		DaysOfWeek: []int{1, 2, 3, 4, 5},
		Vibration:  true,
	}
	created, err := c.CreateAlarm(context.Background(), alarm)
	if err != nil {
		t.Fatalf("CreateAlarm error: %v", err)
	}
	if created.ID != "new-alarm" {
		t.Errorf("expected id new-alarm, got %s", created.ID)
	}
	if capturedBody.Time != "08:00" {
		t.Errorf("expected captured time 08:00, got %s", capturedBody.Time)
	}
}

func TestUpdateAlarm(t *testing.T) {
	var capturedBody map[string]any
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/alarms/alarm-1", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&capturedBody); err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"alarm":{"id":"alarm-1","enabled":false,"time":"09:00","daysOfWeek":[1,2,3,4,5],"vibration":true}}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	patch := map[string]any{"enabled": false}
	updated, err := c.UpdateAlarm(context.Background(), "alarm-1", patch)
	if err != nil {
		t.Fatalf("UpdateAlarm error: %v", err)
	}
	if capturedPath != "/users/uid-123/alarms/alarm-1" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
	if updated.Enabled {
		t.Error("expected enabled=false")
	}
}

func TestDeleteAlarm(t *testing.T) {
	var capturedPath string
	var capturedMethod string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/alarms/alarm-to-delete", func(w http.ResponseWriter, r *http.Request) {
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

	err := c.DeleteAlarm(context.Background(), "alarm-to-delete")
	if err != nil {
		t.Fatalf("DeleteAlarm error: %v", err)
	}
	if capturedMethod != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", capturedMethod)
	}
	if capturedPath != "/users/uid-123/alarms/alarm-to-delete" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}
