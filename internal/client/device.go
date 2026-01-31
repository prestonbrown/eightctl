package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type DeviceActions struct{ c *Client }

func (c *Client) Device() *DeviceActions { return &DeviceActions{c: c} }

func (d *DeviceActions) Info(ctx context.Context) (any, error) {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s", id)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (d *DeviceActions) Peripherals(ctx context.Context) (any, error) {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s/peripherals", id)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (d *DeviceActions) Owner(ctx context.Context) (any, error) {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s/owner", id)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (d *DeviceActions) Warranty(ctx context.Context) (any, error) {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s/warranty", id)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (d *DeviceActions) Online(ctx context.Context) (any, error) {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s/online", id)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (d *DeviceActions) PrimingTasks(ctx context.Context) (any, error) {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s/priming/tasks", id)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (d *DeviceActions) PrimingSchedule(ctx context.Context) (any, error) {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s/priming/schedule", id)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (d *DeviceActions) Update(ctx context.Context, body map[string]any) error {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/devices/%s", id)
	return d.c.do(ctx, http.MethodPut, path, nil, body, nil)
}

func (d *DeviceActions) SetOwner(ctx context.Context, body map[string]any) error {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/devices/%s/owner", id)
	return d.c.do(ctx, http.MethodPut, path, nil, body, nil)
}

func (d *DeviceActions) SetPeripherals(ctx context.Context, body map[string]any) error {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/devices/%s/peripherals", id)
	return d.c.do(ctx, http.MethodPut, path, nil, body, nil)
}

func (d *DeviceActions) AddPeripheral(ctx context.Context, body map[string]any) error {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/devices/%s/peripherals", id)
	return d.c.do(ctx, http.MethodPatch, path, nil, body, nil)
}

func (d *DeviceActions) GetBLEKey(ctx context.Context, out any) error {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/devices/%s/security/key", id)
	return d.c.do(ctx, http.MethodPost, path, nil, nil, out)
}

func (d *DeviceActions) UpdatePrimingSchedule(ctx context.Context, body map[string]any) error {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/devices/%s/priming/schedule", id)
	return d.c.do(ctx, http.MethodPut, path, nil, body, nil)
}

func (d *DeviceActions) CreatePrimingTask(ctx context.Context, body map[string]any) error {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/devices/%s/priming/tasks", id)
	return d.c.do(ctx, http.MethodPost, path, nil, body, nil)
}

func (d *DeviceActions) CancelPrimingTask(ctx context.Context) error {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/devices/%s/priming/tasks", id)
	return d.c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// DeviceWithUsers contains device info including user assignments.
type DeviceWithUsers struct {
	ID              string  `json:"id"`
	LeftUserID      string  `json:"leftUserId"`
	RightUserID     string  `json:"rightUserId"`
	RoomTemperature float64 `json:"roomTemperature"`
	WaterLevel      int     `json:"waterLevel"`
	IsPriming       bool    `json:"isPriming"`
	NeedsPriming    bool    `json:"needsPriming"`
}

// GetWithUsers fetches device info with left/right user assignments.
func (d *DeviceActions) GetWithUsers(ctx context.Context) (*DeviceWithUsers, error) {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s", id)
	query := url.Values{}
	query.Set("filter", "leftUserId,rightUserId,awaySides")

	var res struct {
		Result struct {
			ID              string  `json:"id"`
			LeftUserID      string  `json:"leftUserId"`
			RightUserID     string  `json:"rightUserId"`
			RoomTemperature float64 `json:"roomTemperature"`
			WaterLevel      int     `json:"waterLevel"`
			Priming         struct {
				Status string `json:"status"`
			} `json:"priming"`
		} `json:"result"`
	}
	err = d.c.do(ctx, http.MethodGet, path, query, nil, &res)
	if err != nil {
		return nil, err
	}

	return &DeviceWithUsers{
		ID:              res.Result.ID,
		LeftUserID:      res.Result.LeftUserID,
		RightUserID:     res.Result.RightUserID,
		RoomTemperature: res.Result.RoomTemperature,
		WaterLevel:      res.Result.WaterLevel,
		IsPriming:       res.Result.Priming.Status == "priming",
		NeedsPriming:    res.Result.Priming.Status == "needed",
	}, nil
}
