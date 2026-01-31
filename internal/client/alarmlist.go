package client

import (
	"context"
	"fmt"
	"net/http"
)

// AlarmRepeat represents the repeat configuration for an alarm.
type AlarmRepeat struct {
	Enabled  bool            `json:"enabled"`
	WeekDays map[string]bool `json:"weekDays,omitempty"`
}

// AlarmVibration represents vibration settings for an alarm.
type AlarmVibration struct {
	Enabled    bool   `json:"enabled"`
	PowerLevel int    `json:"powerLevel,omitempty"`
	Pattern    string `json:"pattern,omitempty"`
}

// AlarmThermal represents thermal wake settings for an alarm.
type AlarmThermal struct {
	Enabled bool `json:"enabled"`
	Level   int  `json:"level,omitempty"`
}

// Alarm represents alarm payload (app-api format).
type Alarm struct {
	ID             string         `json:"id,omitempty"`
	Enabled        bool           `json:"enabled"`
	Time           string         `json:"time"` // HH:MM:SS format
	Repeat         AlarmRepeat    `json:"repeat,omitempty"`
	Vibration      AlarmVibration `json:"vibration,omitempty"`
	Thermal        AlarmThermal   `json:"thermal,omitempty"`
	Snoozing       bool           `json:"snoozing,omitempty"`
	DismissedUntil string         `json:"dismissedUntil,omitempty"`
	SkippedUntil   string         `json:"skippedUntil,omitempty"`
	SnoozedUntil   string         `json:"snoozedUntil,omitempty"`
	SkipNext       bool           `json:"skipNext,omitempty"`
	NextTimestamp  string         `json:"nextTimestamp,omitempty"`
	StartTimestamp string         `json:"startTimestamp,omitempty"`
	EndTimestamp   string         `json:"endTimestamp,omitempty"`
}

// ListAlarms retrieves alarms using the v2 API endpoint per APK.
func (c *Client) ListAlarms(ctx context.Context) ([]Alarm, error) {
	if err := c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/v2/users/%s/alarms", c.UserID)
	var res struct {
		Alarms []Alarm `json:"alarms"`
	}
	if err := c.doAppAPI(ctx, http.MethodGet, path, nil, nil, &res); err != nil {
		return nil, err
	}
	return res.Alarms, nil
}

func (c *Client) CreateAlarm(ctx context.Context, alarm Alarm) (*Alarm, error) {
	if err := c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/alarms", c.UserID)
	var res struct {
		Alarm Alarm `json:"alarm"`
	}
	if err := c.doAppAPI(ctx, http.MethodPost, path, nil, alarm, &res); err != nil {
		return nil, err
	}
	return &res.Alarm, nil
}

func (c *Client) UpdateAlarm(ctx context.Context, alarmID string, patch map[string]any) (*Alarm, error) {
	if err := c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/alarms/%s", c.UserID, alarmID)
	var res struct {
		Alarm Alarm `json:"alarm"`
	}
	if err := c.doAppAPI(ctx, http.MethodPut, path, nil, patch, &res); err != nil {
		return nil, err
	}
	return &res.Alarm, nil
}

func (c *Client) DeleteAlarm(ctx context.Context, alarmID string) error {
	if err := c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/alarms/%s", c.UserID, alarmID)
	return c.doAppAPI(ctx, http.MethodDelete, path, nil, nil, nil)
}

// ListAlarmsV1 retrieves alarms using the legacy v1 API endpoint.
func (c *Client) ListAlarmsV1(ctx context.Context, out any) error {
	if err := c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/alarms", c.UserID)
	return c.doAppAPI(ctx, http.MethodGet, path, nil, nil, out)
}
