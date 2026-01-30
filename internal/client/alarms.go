package client

import (
	"context"
	"fmt"
	"net/http"
)

// AlarmActions groups alarm endpoints.
type AlarmActions struct {
	c *Client
}

// Alarms helper accessor.
func (c *Client) Alarms() *AlarmActions { return &AlarmActions{c: c} }

// Snooze snoozes an alarm for the default duration (9 minutes).
// Uses the dedicated alarm endpoint per APK.
func (a *AlarmActions) Snooze(ctx context.Context, alarmID string) error {
	return a.SnoozeWithDuration(ctx, alarmID, 9)
}

// SnoozeWithDuration snoozes an alarm for a specified number of minutes.
func (a *AlarmActions) SnoozeWithDuration(ctx context.Context, alarmID string, minutes int) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/alarms/%s/snooze", a.c.UserID, alarmID)
	body := map[string]any{"snoozeMinutes": minutes}
	return a.c.doAppAPI(ctx, http.MethodPut, path, nil, body, nil)
}

// Dismiss dismisses a specific alarm.
// Uses the dedicated alarm endpoint per APK.
func (a *AlarmActions) Dismiss(ctx context.Context, alarmID string) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/alarms/%s/dismiss", a.c.UserID, alarmID)
	return a.c.doAppAPI(ctx, http.MethodPut, path, nil, map[string]any{}, nil)
}

// DismissAll dismisses all active alarms.
// Uses the dedicated alarm endpoint per APK.
func (a *AlarmActions) DismissAll(ctx context.Context) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/alarms/active/dismiss-all", a.c.UserID)
	return a.c.doAppAPI(ctx, http.MethodPut, path, nil, map[string]any{}, nil)
}

func (a *AlarmActions) VibrationTest(ctx context.Context) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/vibration-test", a.c.UserID)
	return a.c.doAppAPI(ctx, http.MethodPost, path, nil, map[string]any{}, nil)
}
