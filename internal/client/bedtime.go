package client

import (
	"context"
	"fmt"
	"net/http"
)

// BedtimeActions groups bedtime scheduling endpoints.
type BedtimeActions struct{ c *Client }

// Bedtime helper accessor.
func (c *Client) Bedtime() *BedtimeActions { return &BedtimeActions{c: c} }

// BedtimeSchedule represents the bedtime temperature schedule.
type BedtimeSchedule struct {
	Enabled bool `json:"enabled"`
	Time    struct {
		Start string `json:"start"` // HH:MM format
		End   string `json:"end"`   // HH:MM format
	} `json:"time"`
	Stages []BedtimeStage `json:"stages,omitempty"`
}

// BedtimeStage represents a temperature stage in the bedtime schedule.
type BedtimeStage struct {
	Level    int    `json:"level"`    // -100 to 100
	Duration int    `json:"duration"` // minutes
	Type     string `json:"type,omitempty"`
}

// Get retrieves the current bedtime schedule.
func (b *BedtimeActions) Get(ctx context.Context) (*BedtimeSchedule, error) {
	if err := b.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/temperature", b.c.UserID)
	var res struct {
		Bedtime BedtimeSchedule `json:"bedtime"`
	}
	if err := b.c.doAppAPI(ctx, http.MethodGet, path, nil, nil, &res); err != nil {
		return nil, err
	}
	return &res.Bedtime, nil
}

// Set updates the bedtime schedule.
func (b *BedtimeActions) Set(ctx context.Context, schedule *BedtimeSchedule) error {
	if err := b.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/bedtime", b.c.UserID)
	return b.c.doAppAPI(ctx, http.MethodPut, path, nil, schedule, nil)
}

// Enable enables the bedtime schedule.
func (b *BedtimeActions) Enable(ctx context.Context) error {
	if err := b.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/bedtime", b.c.UserID)
	body := map[string]any{"enabled": true}
	return b.c.doAppAPI(ctx, http.MethodPut, path, nil, body, nil)
}

// Disable disables the bedtime schedule.
func (b *BedtimeActions) Disable(ctx context.Context) error {
	if err := b.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/bedtime", b.c.UserID)
	body := map[string]any{"enabled": false}
	return b.c.doAppAPI(ctx, http.MethodPut, path, nil, body, nil)
}
