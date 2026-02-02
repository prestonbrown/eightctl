package mqtt

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/steipete/eightctl/internal/adapter"
	"github.com/steipete/eightctl/internal/client"
	"github.com/steipete/eightctl/internal/model"
	"github.com/steipete/eightctl/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestAdapter creates a test adapter with mocked Eight Sleep API.
func setupTestAdapter(t *testing.T, setTempCalls *[]setTempCall, turnOnCalls, turnOffCalls *[]string) (*Adapter, *httptest.Server) {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/login" && r.Method == http.MethodPost:
			// Mock auth endpoint
			json.NewEncoder(w).Encode(map[string]any{
				"session": map[string]any{
					"token":     "test-token",
					"expiresAt": "2099-01-01T00:00:00Z",
					"userId":    "user-123",
				},
			})
		case r.URL.Path == "/users/me" && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]any{
				"user": map[string]any{
					"userId": "user-123",
					"currentDevice": map[string]any{
						"id": "device-123",
					},
				},
			})
		case r.URL.Path == "/devices/device-123" && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]any{
				"result": map[string]any{
					"id":              "device-123",
					"leftUserId":      "left-user",
					"rightUserId":     "right-user",
					"roomTemperature": 70.0,
					"waterLevel":      100,
				},
			})
		case r.URL.Path == "/users/left-user/temperature" && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]any{
				"currentLevel": -20,
				"currentState": map[string]any{"type": "smart"},
			})
		case r.URL.Path == "/users/right-user/temperature" && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]any{
				"currentLevel": 10,
				"currentState": map[string]any{"type": "off"},
			})
		case r.URL.Path == "/users/left-user/temperature" && r.Method == http.MethodPut:
			if setTempCalls != nil {
				var body map[string]any
				json.NewDecoder(r.Body).Decode(&body)
				if level, ok := body["currentLevel"].(float64); ok {
					*setTempCalls = append(*setTempCalls, setTempCall{"left-user", int(level)})
				}
			}
			w.WriteHeader(http.StatusOK)
		case r.URL.Path == "/users/right-user/temperature" && r.Method == http.MethodPut:
			if setTempCalls != nil {
				var body map[string]any
				json.NewDecoder(r.Body).Decode(&body)
				if level, ok := body["currentLevel"].(float64); ok {
					*setTempCalls = append(*setTempCalls, setTempCall{"right-user", int(level)})
				}
			}
			w.WriteHeader(http.StatusOK)
		case r.URL.Path == "/users/left-user/devices/power" && r.Method == http.MethodPost:
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			if on, ok := body["on"].(bool); ok {
				if on && turnOnCalls != nil {
					*turnOnCalls = append(*turnOnCalls, "left-user")
				}
				if !on && turnOffCalls != nil {
					*turnOffCalls = append(*turnOffCalls, "left-user")
				}
			}
			w.WriteHeader(http.StatusOK)
		case r.URL.Path == "/users/right-user/devices/power" && r.Method == http.MethodPost:
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			if on, ok := body["on"].(bool); ok {
				if on && turnOnCalls != nil {
					*turnOnCalls = append(*turnOnCalls, "right-user")
				}
				if !on && turnOffCalls != nil {
					*turnOffCalls = append(*turnOffCalls, "right-user")
				}
			}
			w.WriteHeader(http.StatusOK)
		default:
			t.Logf("unhandled request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	c := client.New("test@test.com", "pass", "", "", "")
	c.BaseURL = srv.URL

	mgr := state.NewManager(c, "device-123", state.WithCacheTTL(5*time.Second))

	// Pre-populate state
	_, err := mgr.GetState(context.Background())
	require.NoError(t, err)

	cfg := Config{
		BrokerURL:    "tcp://localhost:1883",
		TopicPrefix:  "homeassistant",
		DeviceID:     "device-123",
		DeviceName:   "Test Pod",
		PollInterval: 60 * time.Second,
		ClientID:     "test-client",
	}

	a := New(cfg, mgr)
	return a, srv
}

type setTempCall struct {
	userID string
	level  int
}

func TestNew(t *testing.T) {
	c := client.New("test@test.com", "pass", "", "", "")
	mgr := state.NewManager(c, "device-123")

	cfg := Config{
		BrokerURL:    "tcp://localhost:1883",
		TopicPrefix:  "homeassistant",
		DeviceID:     "device-123",
		DeviceName:   "Test Pod",
		PollInterval: 60 * time.Second,
		ClientID:     "test-client",
	}

	a := New(cfg, mgr)

	assert.NotNil(t, a)
	assert.Equal(t, cfg.BrokerURL, a.cfg.BrokerURL)
	assert.Equal(t, cfg.DeviceID, a.cfg.DeviceID)
}

func TestAdapter_HandleCommand_TurnOn(t *testing.T) {
	var turnOnCalls []string
	a, srv := setupTestAdapter(t, nil, &turnOnCalls, nil)
	defer srv.Close()

	cmd := adapter.Command{
		Action: adapter.ActionOn,
		Side:   model.Left,
	}

	err := a.HandleCommand(context.Background(), cmd)
	require.NoError(t, err)

	require.Len(t, turnOnCalls, 1)
	assert.Equal(t, "left-user", turnOnCalls[0])
}

func TestAdapter_HandleCommand_TurnOff(t *testing.T) {
	var turnOffCalls []string
	a, srv := setupTestAdapter(t, nil, nil, &turnOffCalls)
	defer srv.Close()

	cmd := adapter.Command{
		Action: adapter.ActionOff,
		Side:   model.Right,
	}

	err := a.HandleCommand(context.Background(), cmd)
	require.NoError(t, err)

	require.Len(t, turnOffCalls, 1)
	assert.Equal(t, "right-user", turnOffCalls[0])
}

func TestAdapter_HandleCommand_SetTemp(t *testing.T) {
	var setTempCalls []setTempCall
	a, srv := setupTestAdapter(t, &setTempCalls, nil, nil)
	defer srv.Close()

	level := -30
	cmd := adapter.Command{
		Action:      adapter.ActionSetTemp,
		Side:        model.Left,
		Temperature: &level,
	}

	err := a.HandleCommand(context.Background(), cmd)
	require.NoError(t, err)

	require.Len(t, setTempCalls, 1)
	assert.Equal(t, "left-user", setTempCalls[0].userID)
	assert.Equal(t, -30, setTempCalls[0].level)
}

func TestAdapter_HandleCommand_SetTemp_MissingTemperature(t *testing.T) {
	a, srv := setupTestAdapter(t, nil, nil, nil)
	defer srv.Close()

	cmd := adapter.Command{
		Action: adapter.ActionSetTemp,
		Side:   model.Left,
		// Temperature is nil
	}

	err := a.HandleCommand(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "temperature required")
}

func TestAdapter_HandleCommand_UnknownAction(t *testing.T) {
	a, srv := setupTestAdapter(t, nil, nil, nil)
	defer srv.Close()

	cmd := adapter.Command{
		Action: "invalid_action",
		Side:   model.Left,
	}

	err := a.HandleCommand(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown action")
}

func TestAdapter_powerStateToMode(t *testing.T) {
	a := &Adapter{}

	tests := []struct {
		state    model.PowerState
		level    int
		expected string
	}{
		{model.PowerOff, 0, "off"},
		{model.PowerSmart, 20, "heat"},
		{model.PowerSmart, -20, "cool"},
		{model.PowerManual, 50, "heat"},
		{model.PowerManual, -50, "cool"},
		{model.PowerSmart, 0, "heat"}, // Zero treated as heat
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := a.powerStateToMode(tt.state, tt.level)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_OptionalFields(t *testing.T) {
	cfg := Config{
		BrokerURL:    "tcp://localhost:1883",
		TopicPrefix:  "homeassistant",
		DeviceID:     "pod-1",
		DeviceName:   "Test",
		PollInterval: 30 * time.Second,
		ClientID:     "client-1",
		Username:     "user",
		Password:     "pass",
	}

	assert.Equal(t, "user", cfg.Username)
	assert.Equal(t, "pass", cfg.Password)
}

// Verify compile-time interface compliance
var _ adapter.Adapter = (*Adapter)(nil)
