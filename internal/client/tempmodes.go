package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type TempModes struct{ c *Client }

func (c *Client) TempModes() *TempModes { return &TempModes{c: c} }

// Nap mode controls
func (t *TempModes) NapActivate(ctx context.Context) error {
	return t.simplePost(ctx, "/temperature/nap-mode/activate")
}

func (t *TempModes) NapDeactivate(ctx context.Context) error {
	return t.simplePost(ctx, "/temperature/nap-mode/deactivate")
}

func (t *TempModes) NapExtend(ctx context.Context) error {
	return t.simplePost(ctx, "/temperature/nap-mode/extend")
}

func (t *TempModes) NapStatus(ctx context.Context, out any) error {
	return t.simpleGet(ctx, "/temperature/nap-mode/status", out)
}

func (t *TempModes) NapSettings(ctx context.Context, out any) error {
	return t.simpleGet(ctx, "/temperature/nap-mode", out)
}

func (t *TempModes) NapAlarmSettings(ctx context.Context, out any) error {
	return t.simpleGet(ctx, "/temporary-mode/nap-mode", out)
}

func (t *TempModes) UpdateNapAlarmSettings(ctx context.Context, body map[string]any) error {
	return t.simplePut(ctx, "/temporary-mode/nap-mode", body)
}

// Hot flash controls
func (t *TempModes) HotFlashActivate(ctx context.Context) error {
	return t.simplePost(ctx, "/temperature/hot-flash-mode/activate")
}

func (t *TempModes) HotFlashDeactivate(ctx context.Context) error {
	return t.simplePost(ctx, "/temperature/hot-flash-mode/deactivate")
}

func (t *TempModes) HotFlashStatus(ctx context.Context, out any) error {
	return t.simpleGet(ctx, "/temperature/hot-flash-mode", out)
}

func (t *TempModes) UpdateHotFlashSettings(ctx context.Context, body map[string]any) error {
	return t.simplePut(ctx, "/temperature/hot-flash-mode", body)
}

func (t *TempModes) DeleteHotFlashSettings(ctx context.Context) error {
	return t.simpleDelete(ctx, "/temperature/hot-flash-mode")
}

// Temp events history
func (t *TempModes) TempEvents(ctx context.Context, from, to string, out any) error {
	if err := t.c.requireUser(ctx); err != nil {
		return err
	}
	q := url.Values{}
	if from != "" {
		q.Set("from", from)
	}
	if to != "" {
		q.Set("to", to)
	}
	path := fmt.Sprintf("/users/%s/temp-events", t.c.UserID)
	return t.c.do(ctx, http.MethodGet, path, q, nil, out)
}

func (t *TempModes) simplePost(ctx context.Context, suffix string) error {
	if err := t.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s%s", t.c.UserID, suffix)
	return t.c.do(ctx, http.MethodPost, path, nil, map[string]string{}, nil)
}

func (t *TempModes) simpleGet(ctx context.Context, suffix string, out any) error {
	if err := t.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s%s", t.c.UserID, suffix)
	return t.c.do(ctx, http.MethodGet, path, nil, nil, out)
}

func (t *TempModes) simplePut(ctx context.Context, suffix string, body any) error {
	if err := t.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s%s", t.c.UserID, suffix)
	return t.c.do(ctx, http.MethodPut, path, nil, body, nil)
}

func (t *TempModes) simpleDelete(ctx context.Context, suffix string) error {
	if err := t.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s%s", t.c.UserID, suffix)
	return t.c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}
