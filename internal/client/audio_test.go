package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAudioActions_Tracks(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/audio/tracks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"tracks":[{"id":"t1","title":"Rain","type":"nature"}]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	tracks, err := c.Audio().Tracks(context.Background())
	if err != nil {
		t.Fatalf("Tracks error: %v", err)
	}
	if len(tracks) != 1 {
		t.Errorf("expected 1 track, got %d", len(tracks))
	}
}

func TestAudioActions_Categories(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/audio/categories", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"categories":["nature","ambient"]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Audio().Categories(context.Background())
	if err != nil {
		t.Fatalf("Categories error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestAudioActions_PlayerState(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/audio/player/state", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"state":"playing","trackId":"t1"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Audio().PlayerState(context.Background())
	if err != nil {
		t.Fatalf("PlayerState error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestAudioActions_Play(t *testing.T) {
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/audio/player", func(w http.ResponseWriter, r *http.Request) {
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

	err := c.Audio().Play(context.Background(), "track-123")
	if err != nil {
		t.Fatalf("Play error: %v", err)
	}
	if capturedBody["action"] != "play" {
		t.Errorf("expected action=play, got %v", capturedBody["action"])
	}
	if capturedBody["trackId"] != "track-123" {
		t.Errorf("expected trackId=track-123, got %v", capturedBody["trackId"])
	}
}

func TestAudioActions_Play_NoTrackID(t *testing.T) {
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/audio/player", func(w http.ResponseWriter, r *http.Request) {
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

	err := c.Audio().Play(context.Background(), "")
	if err != nil {
		t.Fatalf("Play error: %v", err)
	}
	if _, exists := capturedBody["trackId"]; exists {
		t.Error("expected no trackId when empty")
	}
}

func TestAudioActions_Pause(t *testing.T) {
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/audio/player", func(w http.ResponseWriter, r *http.Request) {
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

	err := c.Audio().Pause(context.Background())
	if err != nil {
		t.Fatalf("Pause error: %v", err)
	}
	if capturedBody["action"] != "pause" {
		t.Errorf("expected action=pause, got %v", capturedBody["action"])
	}
}

func TestAudioActions_Seek(t *testing.T) {
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/audio/player/seek", func(w http.ResponseWriter, r *http.Request) {
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

	err := c.Audio().Seek(context.Background(), 30000)
	if err != nil {
		t.Fatalf("Seek error: %v", err)
	}
	if capturedBody["position"] != float64(30000) {
		t.Errorf("expected position=30000, got %v", capturedBody["position"])
	}
}

func TestAudioActions_Volume(t *testing.T) {
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/audio/player/volume", func(w http.ResponseWriter, r *http.Request) {
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

	err := c.Audio().Volume(context.Background(), 75)
	if err != nil {
		t.Fatalf("Volume error: %v", err)
	}
	if capturedBody["level"] != float64(75) {
		t.Errorf("expected level=75, got %v", capturedBody["level"])
	}
}

func TestAudioActions_Pair(t *testing.T) {
	var capturedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-456"}}}`))
	})
	mux.HandleFunc("/devices/dev-456/audio/player/pair", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	err := c.Audio().Pair(context.Background())
	if err != nil {
		t.Fatalf("Pair error: %v", err)
	}
	if capturedPath != "/devices/dev-456/audio/player/pair" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}

func TestAudioActions_RecommendedNext(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/audio/tracks/recommended-next-track", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"trackId":"next-track"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Audio().RecommendedNext(context.Background())
	if err != nil {
		t.Fatalf("RecommendedNext error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestAudioActions_Favorites(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/audio/tracks/favorites", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"favorites":["t1","t2"]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	res, err := c.Audio().Favorites(context.Background())
	if err != nil {
		t.Fatalf("Favorites error: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil response")
	}
}

func TestAudioActions_AddFavorite(t *testing.T) {
	var capturedBody map[string]any

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/audio/tracks/favorites", func(w http.ResponseWriter, r *http.Request) {
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

	err := c.Audio().AddFavorite(context.Background(), "track-123")
	if err != nil {
		t.Fatalf("AddFavorite error: %v", err)
	}
	if capturedBody["trackId"] != "track-123" {
		t.Errorf("expected trackId=track-123, got %v", capturedBody["trackId"])
	}
}

func TestAudioActions_RemoveFavorite(t *testing.T) {
	var capturedPath string
	var capturedMethod string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/uid-123/audio/tracks/favorites/track-123", func(w http.ResponseWriter, r *http.Request) {
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

	err := c.Audio().RemoveFavorite(context.Background(), "track-123")
	if err != nil {
		t.Fatalf("RemoveFavorite error: %v", err)
	}
	if capturedMethod != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", capturedMethod)
	}
	if capturedPath != "/users/uid-123/audio/tracks/favorites/track-123" {
		t.Errorf("unexpected path: %s", capturedPath)
	}
}
