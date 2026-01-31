// Package hubitat provides an HTTP server adapter for Hubitat smart home integration.
// Hubitat uses HTTP to communicate with local devices, and this server exposes
// endpoints for Hubitat Maker API / custom drivers to call.
package hubitat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/steipete/eightctl/internal/adapter"
	"github.com/steipete/eightctl/internal/model"
	"github.com/steipete/eightctl/internal/state"
)

// Compile-time check that Adapter implements adapter.Adapter.
var _ adapter.Adapter = (*Adapter)(nil)

// Adapter implements the adapter.Adapter interface for Hubitat integration.
type Adapter struct {
	server       *http.Server
	stateManager *state.Manager
	port         int
	pollInterval time.Duration
}

// New creates a new Hubitat adapter.
func New(stateManager *state.Manager, port int, pollInterval time.Duration) *Adapter {
	return &Adapter{
		stateManager: stateManager,
		port:         port,
		pollInterval: pollInterval,
	}
}

// Start begins the HTTP server for Hubitat integration.
func (a *Adapter) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/status", a.handleStatus)
	mux.HandleFunc("/left/status", a.handleSideStatus(model.Left))
	mux.HandleFunc("/right/status", a.handleSideStatus(model.Right))
	mux.HandleFunc("/left/on", a.handleSideOn(model.Left))
	mux.HandleFunc("/right/on", a.handleSideOn(model.Right))
	mux.HandleFunc("/left/off", a.handleSideOff(model.Left))
	mux.HandleFunc("/right/off", a.handleSideOff(model.Right))
	mux.HandleFunc("/left/temperature", a.handleSideTemperature(model.Left))
	mux.HandleFunc("/right/temperature", a.handleSideTemperature(model.Right))

	a.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", a.port),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait briefly to check for immediate startup errors
	select {
	case err := <-errChan:
		return fmt.Errorf("failed to start HTTP server: %w", err)
	case <-time.After(100 * time.Millisecond):
		// Server started successfully
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// HandleCommand processes a command from the smart home platform.
func (a *Adapter) HandleCommand(ctx context.Context, cmd adapter.Command) error {
	switch cmd.Action {
	case adapter.ActionOn:
		return a.stateManager.TurnOn(ctx, cmd.Side)
	case adapter.ActionOff:
		return a.stateManager.TurnOff(ctx, cmd.Side)
	case adapter.ActionSetTemp:
		if cmd.Temperature == nil {
			return fmt.Errorf("temperature level required for set_temperature action")
		}
		return a.stateManager.SetTemperature(ctx, cmd.Side, *cmd.Temperature)
	default:
		return fmt.Errorf("unknown action: %s", cmd.Action)
	}
}

// Stop gracefully shuts down the HTTP server.
func (a *Adapter) Stop() error {
	if a.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return a.server.Shutdown(ctx)
}

// StatusResponse represents the full device status response.
type StatusResponse struct {
	ID    string      `json:"id"`
	Left  *SideStatus `json:"left,omitempty"`
	Right *SideStatus `json:"right,omitempty"`
}

// SideStatus represents the status of one side of the bed.
type SideStatus struct {
	On             bool    `json:"on"`
	Level          int     `json:"level"`
	BedTemperature float64 `json:"bed_temperature"`
}

// handleStatus returns the full device state as JSON.
func (a *Adapter) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	deviceState, err := a.stateManager.GetState(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get state: %v", err), http.StatusInternalServerError)
		return
	}

	resp := StatusResponse{
		ID: deviceState.ID,
	}

	if deviceState.LeftUser != nil {
		resp.Left = &SideStatus{
			On:             deviceState.LeftUser.IsOn(),
			Level:          deviceState.LeftUser.TargetLevel,
			BedTemperature: deviceState.LeftUser.BedTemperature,
		}
	}

	if deviceState.RightUser != nil {
		resp.Right = &SideStatus{
			On:             deviceState.RightUser.IsOn(),
			Level:          deviceState.RightUser.TargetLevel,
			BedTemperature: deviceState.RightUser.BedTemperature,
		}
	}

	writeJSON(w, resp)
}

// handleSideStatus returns a handler for a specific side's status.
func (a *Adapter) handleSideStatus(side model.Side) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		deviceState, err := a.stateManager.GetState(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get state: %v", err), http.StatusInternalServerError)
			return
		}

		userState := deviceState.GetSide(side)
		if userState == nil {
			http.Error(w, fmt.Sprintf("no user assigned to %s side", side), http.StatusNotFound)
			return
		}

		resp := SideStatus{
			On:             userState.IsOn(),
			Level:          userState.TargetLevel,
			BedTemperature: userState.BedTemperature,
		}

		writeJSON(w, resp)
	}
}

// handleSideOn returns a handler to turn on a specific side.
func (a *Adapter) handleSideOn(side model.Side) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cmd := adapter.Command{
			Action: adapter.ActionOn,
			Side:   side,
		}

		if err := a.HandleCommand(r.Context(), cmd); err != nil {
			http.Error(w, fmt.Sprintf("failed to turn on: %v", err), http.StatusInternalServerError)
			return
		}

		writeJSON(w, map[string]string{"status": "ok"})
	}
}

// handleSideOff returns a handler to turn off a specific side.
func (a *Adapter) handleSideOff(side model.Side) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cmd := adapter.Command{
			Action: adapter.ActionOff,
			Side:   side,
		}

		if err := a.HandleCommand(r.Context(), cmd); err != nil {
			http.Error(w, fmt.Sprintf("failed to turn off: %v", err), http.StatusInternalServerError)
			return
		}

		writeJSON(w, map[string]string{"status": "ok"})
	}
}

// handleSideTemperature returns a handler to set temperature for a specific side.
func (a *Adapter) handleSideTemperature(side model.Side) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		levelStr := r.URL.Query().Get("level")
		if levelStr == "" {
			http.Error(w, "level parameter required", http.StatusBadRequest)
			return
		}

		level, err := strconv.Atoi(strings.TrimSpace(levelStr))
		if err != nil {
			http.Error(w, "invalid level: must be an integer", http.StatusBadRequest)
			return
		}

		if level < -100 || level > 100 {
			http.Error(w, "invalid level: must be between -100 and 100", http.StatusBadRequest)
			return
		}

		cmd := adapter.Command{
			Action:      adapter.ActionSetTemp,
			Side:        side,
			Temperature: &level,
		}

		if err := a.HandleCommand(r.Context(), cmd); err != nil {
			http.Error(w, fmt.Sprintf("failed to set temperature: %v", err), http.StatusInternalServerError)
			return
		}

		writeJSON(w, map[string]any{"status": "ok", "level": level})
	}
}

// writeJSON writes a JSON response with proper content type.
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
