package client

import (
	"context"
	"fmt"
	"net/http"
)

// AwayModeActions groups away mode endpoints.
type AwayModeActions struct{ c *Client }

// AwayMode helper accessor.
func (c *Client) AwayMode() *AwayModeActions { return &AwayModeActions{c: c} }

// AwayModeStatus represents the away mode state.
type AwayModeStatus struct {
	Enabled bool `json:"enabled"`
}

// Get retrieves the current away mode status.
func (a *AwayModeActions) Get(ctx context.Context) (*AwayModeStatus, error) {
	if err := a.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/away-mode", a.c.UserID)
	var res AwayModeStatus
	if err := a.c.doAppAPI(ctx, http.MethodGet, path, nil, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// Set enables or disables away mode.
func (a *AwayModeActions) Set(ctx context.Context, enabled bool) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/away-mode", a.c.UserID)
	body := map[string]any{"enabled": enabled}
	return a.c.doAppAPI(ctx, http.MethodPut, path, nil, body, nil)
}

// Enable enables away mode.
func (a *AwayModeActions) Enable(ctx context.Context) error {
	return a.Set(ctx, true)
}

// Disable disables away mode.
func (a *AwayModeActions) Disable(ctx context.Context) error {
	return a.Set(ctx, false)
}
