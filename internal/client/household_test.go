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

func TestHouseholdActions_SetCurrentSet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/users/uid-123/current-set", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Household().SetCurrentSet(context.Background(), map[string]any{"setId": "set-1"})
	if err != nil {
		t.Fatalf("SetCurrentSet error: %v", err)
	}
}

func TestHouseholdActions_ClearCurrentSet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/users/uid-123/current-set", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Household().ClearCurrentSet(context.Background())
	if err != nil {
		t.Fatalf("ClearCurrentSet error: %v", err)
	}
}

func TestHouseholdActions_SetReturnDate(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/users/uid-123/schedule", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Household().SetReturnDate(context.Background(), map[string]any{"returnDate": "2024-01-15"})
	if err != nil {
		t.Fatalf("SetReturnDate error: %v", err)
	}
}

func TestHouseholdActions_RemoveReturnDate(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/users/uid-123/schedule/set-456", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Household().RemoveReturnDate(context.Background(), "set-456")
	if err != nil {
		t.Fatalf("RemoveReturnDate error: %v", err)
	}
}

func TestHouseholdActions_AddDevice(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/households/hh-123/devices", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Household().AddDevice(context.Background(), "hh-123", map[string]any{"deviceId": "dev-1"})
	if err != nil {
		t.Fatalf("AddDevice error: %v", err)
	}
}

func TestHouseholdActions_UpdateDevice(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/devices/dev-123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Household().UpdateDevice(context.Background(), "dev-123", map[string]any{"name": "My Pod"})
	if err != nil {
		t.Fatalf("UpdateDevice error: %v", err)
	}
}

func TestHouseholdActions_RemoveDevice(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/devices/dev-123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Household().RemoveDevice(context.Background(), "dev-123")
	if err != nil {
		t.Fatalf("RemoveDevice error: %v", err)
	}
}

func TestHouseholdActions_InviteUser(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/households/hh-123/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Household().InviteUser(context.Background(), "hh-123", map[string]any{"email": "guest@example.com"})
	if err != nil {
		t.Fatalf("InviteUser error: %v", err)
	}
}

func TestHouseholdActions_RespondToInvitation(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/households/hh-123/users/uid-456", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Household().RespondToInvitation(context.Background(), "hh-123", "uid-456", map[string]any{"accept": true})
	if err != nil {
		t.Fatalf("RespondToInvitation error: %v", err)
	}
}

func TestHouseholdActions_RemoveGuest(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/households/hh-123/users/uid-456", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Household().RemoveGuest(context.Background(), "hh-123", "uid-456")
	if err != nil {
		t.Fatalf("RemoveGuest error: %v", err)
	}
}

func TestHouseholdActions_AddGuests(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/households/hh-123/devices/dev-456/guests", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Household().AddGuests(context.Background(), "hh-123", "dev-456", map[string]any{"userIds": []string{"uid-789"}})
	if err != nil {
		t.Fatalf("AddGuests error: %v", err)
	}
}

func TestHouseholdActions_UpdateDeviceSet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/households/hh-123/sets/set-456", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Household().UpdateDeviceSet(context.Background(), "hh-123", "set-456", map[string]any{"name": "Living Room"})
	if err != nil {
		t.Fatalf("UpdateDeviceSet error: %v", err)
	}
}

func TestHouseholdActions_RemoveDeviceSet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/households/hh-123/sets/set-456", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Household().RemoveDeviceSet(context.Background(), "hh-123", "set-456")
	if err != nil {
		t.Fatalf("RemoveDeviceSet error: %v", err)
	}
}

func TestHouseholdActions_RemoveDeviceAssignment(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/household/devices/dev-123/assignment/users/uid-456", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Household().RemoveDeviceAssignment(context.Background(), "dev-123", "uid-456")
	if err != nil {
		t.Fatalf("RemoveDeviceAssignment error: %v", err)
	}
}
