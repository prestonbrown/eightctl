package client

import (
	"context"
	"fmt"
	"net/http"
)

type BaseActions struct{ c *Client }

func (c *Client) Base() *BaseActions { return &BaseActions{c: c} }

func (b *BaseActions) Info(ctx context.Context) (any, error) {
	if err := b.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/base", b.c.UserID)
	var res any
	err := b.c.doAppAPI(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// SetAngle sets the adjustable base angles.
// torsoAngle: head/torso elevation angle
// legAngle: leg/foot elevation angle
func (b *BaseActions) SetAngle(ctx context.Context, torsoAngle, legAngle int) error {
	if err := b.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/base/angle", b.c.UserID)
	body := map[string]any{"torsoAngle": torsoAngle, "legAngle": legAngle}
	return b.c.doAppAPI(ctx, http.MethodPost, path, nil, body, nil)
}

// SetAngleWithSnoreMitigation sets base angles with snore mitigation enabled.
func (b *BaseActions) SetAngleWithSnoreMitigation(ctx context.Context, torsoAngle, legAngle int) error {
	if err := b.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/base/angle", b.c.UserID)
	body := map[string]any{"torsoAngle": torsoAngle, "legAngle": legAngle, "snoreMitigation": true}
	return b.c.doAppAPI(ctx, http.MethodPost, path, nil, body, nil)
}

// StopMovement stops base movement.
func (b *BaseActions) StopMovement(ctx context.Context) error {
	if err := b.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/base/angle", b.c.UserID)
	return b.c.doAppAPI(ctx, http.MethodDelete, path, nil, nil, nil)
}

// Presets retrieves available base presets.
// Uses v2 endpoint per APK.
func (b *BaseActions) Presets(ctx context.Context) (any, error) {
	if err := b.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/v2/users/%s/base/presets", b.c.UserID)
	var res any
	err := b.c.doAppAPI(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (b *BaseActions) RunPreset(ctx context.Context, name string) error {
	if err := b.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/base/presets", b.c.UserID)
	body := map[string]any{"name": name}
	return b.c.doAppAPI(ctx, http.MethodPost, path, nil, body, nil)
}

func (b *BaseActions) VibrationTest(ctx context.Context) error {
	deviceID, err := b.c.EnsureDeviceID(ctx)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/devices/%s/vibration-test", deviceID)
	return b.c.doAppAPI(ctx, http.MethodPost, path, nil, map[string]any{}, nil)
}

// StopVibrationTest stops the vibration test.
func (b *BaseActions) StopVibrationTest(ctx context.Context) error {
	deviceID, err := b.c.EnsureDeviceID(ctx)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/v2/devices/%s/vibration-test/stop", deviceID)
	return b.c.doAppAPI(ctx, http.MethodPut, path, nil, nil, nil)
}

// PairFirstFound pairs the first found base via Bluetooth.
func (b *BaseActions) PairFirstFound(ctx context.Context) error {
	deviceID, err := b.c.EnsureDeviceID(ctx)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/devices/%s/base/pairfirstfoundbase", deviceID)
	return b.c.doAppAPI(ctx, http.MethodPost, path, nil, nil, nil)
}

// Unpair unpairs the adjustable base from the device.
func (b *BaseActions) Unpair(ctx context.Context) error {
	deviceID, err := b.c.EnsureDeviceID(ctx)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/devices/%s/base", deviceID)
	return b.c.doAppAPI(ctx, http.MethodDelete, path, nil, nil, nil)
}
