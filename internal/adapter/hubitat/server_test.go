package hubitat

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

type setTempCall struct {
	userID string
	level  int
}

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
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			if setTempCalls != nil {
				if level, ok := body["currentLevel"].(float64); ok {
					*setTempCalls = append(*setTempCalls, setTempCall{"left-user", int(level)})
				}
			}
			w.WriteHeader(http.StatusOK)
		case r.URL.Path == "/users/right-user/temperature" && r.Method == http.MethodPut:
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			if setTempCalls != nil {
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

	a := New(mgr, 0, 60*time.Second) // Port 0 for testing
	return a, srv
}

func TestNew(t *testing.T) {
	c := client.New("test@test.com", "pass", "", "", "")
	mgr := state.NewManager(c, "device-123")

	a := New(mgr, 8380, 30*time.Second)

	assert.NotNil(t, a)
	assert.Equal(t, 8380, a.port)
	assert.Equal(t, 30*time.Second, a.pollInterval)
}

func TestAdapter_HandleStatus(t *testing.T) {
	a, srv := setupTestAdapter(t, nil, nil, nil)
	defer srv.Close()

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()

	a.handleStatus(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp StatusResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, "device-123", resp.ID)
	require.NotNil(t, resp.Left)
	assert.True(t, resp.Left.On) // smart mode = on
	assert.Equal(t, -20, resp.Left.Level)
	require.NotNil(t, resp.Right)
	assert.False(t, resp.Right.On) // off mode
	assert.Equal(t, 10, resp.Right.Level)
}

func TestAdapter_HandleStatus_MethodNotAllowed(t *testing.T) {
	a, srv := setupTestAdapter(t, nil, nil, nil)
	defer srv.Close()

	req := httptest.NewRequest(http.MethodPost, "/status", nil)
	w := httptest.NewRecorder()

	a.handleStatus(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestAdapter_HandleSideStatus_Left(t *testing.T) {
	a, srv := setupTestAdapter(t, nil, nil, nil)
	defer srv.Close()

	handler := a.handleSideStatus(model.Left)

	req := httptest.NewRequest(http.MethodGet, "/left/status", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp SideStatus
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.True(t, resp.On)
	assert.Equal(t, -20, resp.Level)
}

func TestAdapter_HandleSideStatus_Right(t *testing.T) {
	a, srv := setupTestAdapter(t, nil, nil, nil)
	defer srv.Close()

	handler := a.handleSideStatus(model.Right)

	req := httptest.NewRequest(http.MethodGet, "/right/status", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp SideStatus
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.False(t, resp.On) // off mode
	assert.Equal(t, 10, resp.Level)
}

func TestAdapter_HandleSideOn(t *testing.T) {
	var turnOnCalls []string
	a, srv := setupTestAdapter(t, nil, &turnOnCalls, nil)
	defer srv.Close()

	handler := a.handleSideOn(model.Left)

	req := httptest.NewRequest(http.MethodPut, "/left/on", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	require.Len(t, turnOnCalls, 1)
	assert.Equal(t, "left-user", turnOnCalls[0])
}

func TestAdapter_HandleSideOn_MethodNotAllowed(t *testing.T) {
	a, srv := setupTestAdapter(t, nil, nil, nil)
	defer srv.Close()

	handler := a.handleSideOn(model.Left)

	req := httptest.NewRequest(http.MethodGet, "/left/on", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestAdapter_HandleSideOff(t *testing.T) {
	var turnOffCalls []string
	a, srv := setupTestAdapter(t, nil, nil, &turnOffCalls)
	defer srv.Close()

	handler := a.handleSideOff(model.Right)

	req := httptest.NewRequest(http.MethodPut, "/right/off", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	require.Len(t, turnOffCalls, 1)
	assert.Equal(t, "right-user", turnOffCalls[0])
}

func TestAdapter_HandleSideTemperature(t *testing.T) {
	var setTempCalls []setTempCall
	a, srv := setupTestAdapter(t, &setTempCalls, nil, nil)
	defer srv.Close()

	handler := a.handleSideTemperature(model.Left)

	req := httptest.NewRequest(http.MethodPut, "/left/temperature?level=-30", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	require.Len(t, setTempCalls, 1)
	assert.Equal(t, "left-user", setTempCalls[0].userID)
	assert.Equal(t, -30, setTempCalls[0].level)
}

func TestAdapter_HandleSideTemperature_MissingLevel(t *testing.T) {
	a, srv := setupTestAdapter(t, nil, nil, nil)
	defer srv.Close()

	handler := a.handleSideTemperature(model.Left)

	req := httptest.NewRequest(http.MethodPut, "/left/temperature", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "level parameter required")
}

func TestAdapter_HandleSideTemperature_InvalidLevel(t *testing.T) {
	a, srv := setupTestAdapter(t, nil, nil, nil)
	defer srv.Close()

	handler := a.handleSideTemperature(model.Left)

	req := httptest.NewRequest(http.MethodPut, "/left/temperature?level=abc", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid level")
}

func TestAdapter_HandleSideTemperature_LevelOutOfRange(t *testing.T) {
	a, srv := setupTestAdapter(t, nil, nil, nil)
	defer srv.Close()

	tests := []struct {
		level    string
		expected string
	}{
		{"-150", "must be between -100 and 100"},
		{"150", "must be between -100 and 100"},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			handler := a.handleSideTemperature(model.Left)

			req := httptest.NewRequest(http.MethodPut, "/left/temperature?level="+tt.level, nil)
			w := httptest.NewRecorder()

			handler(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), tt.expected)
		})
	}
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

	level := 50
	cmd := adapter.Command{
		Action:      adapter.ActionSetTemp,
		Side:        model.Left,
		Temperature: &level,
	}

	err := a.HandleCommand(context.Background(), cmd)
	require.NoError(t, err)

	require.Len(t, setTempCalls, 1)
	assert.Equal(t, "left-user", setTempCalls[0].userID)
	assert.Equal(t, 50, setTempCalls[0].level)
}

func TestAdapter_HandleCommand_SetTemp_MissingLevel(t *testing.T) {
	a, srv := setupTestAdapter(t, nil, nil, nil)
	defer srv.Close()

	cmd := adapter.Command{
		Action: adapter.ActionSetTemp,
		Side:   model.Left,
		// Temperature is nil
	}

	err := a.HandleCommand(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "temperature level required")
}

func TestAdapter_HandleCommand_UnknownAction(t *testing.T) {
	a, srv := setupTestAdapter(t, nil, nil, nil)
	defer srv.Close()

	cmd := adapter.Command{
		Action: "invalid",
		Side:   model.Left,
	}

	err := a.HandleCommand(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown action")
}

func TestAdapter_Stop_NilServer(t *testing.T) {
	c := client.New("test@test.com", "pass", "", "", "")
	mgr := state.NewManager(c, "device-123")
	a := New(mgr, 8380, 60*time.Second)

	// Server is nil before Start is called
	err := a.Stop()
	assert.NoError(t, err)
}

// Verify compile-time interface compliance
var _ adapter.Adapter = (*Adapter)(nil)
