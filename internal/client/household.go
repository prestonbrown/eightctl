package client

import (
	"context"
	"fmt"
	"net/http"
)

type HouseholdActions struct{ c *Client }

func (c *Client) Household() *HouseholdActions { return &HouseholdActions{c: c} }

func (h *HouseholdActions) Summary(ctx context.Context) (any, error) {
	if err := h.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/household/users/%s/summary", h.c.UserID)
	var res any
	err := h.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (h *HouseholdActions) Schedule(ctx context.Context) (any, error) {
	if err := h.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/household/users/%s/schedule", h.c.UserID)
	var res any
	err := h.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (h *HouseholdActions) CurrentSet(ctx context.Context) (any, error) {
	if err := h.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/household/users/%s/current-set", h.c.UserID)
	var res any
	err := h.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (h *HouseholdActions) Invitations(ctx context.Context) (any, error) {
	if err := h.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/household/users/%s/invitations", h.c.UserID)
	var res any
	err := h.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (h *HouseholdActions) Devices(ctx context.Context) (any, error) {
	if err := h.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/household/users/%s/devices", h.c.UserID)
	var res any
	err := h.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (h *HouseholdActions) Users(ctx context.Context) (any, error) {
	if err := h.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/household/users/%s/users", h.c.UserID)
	var res any
	err := h.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (h *HouseholdActions) Guests(ctx context.Context) (any, error) {
	if err := h.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/household/users/%s/guests", h.c.UserID)
	var res any
	err := h.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (h *HouseholdActions) SetCurrentSet(ctx context.Context, body map[string]any) error {
	if err := h.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/household/users/%s/current-set", h.c.UserID)
	return h.c.do(ctx, http.MethodPut, path, nil, body, nil)
}

func (h *HouseholdActions) ClearCurrentSet(ctx context.Context) error {
	if err := h.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/household/users/%s/current-set", h.c.UserID)
	return h.c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

func (h *HouseholdActions) SetReturnDate(ctx context.Context, body map[string]any) error {
	if err := h.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/household/users/%s/schedule", h.c.UserID)
	return h.c.do(ctx, http.MethodPost, path, nil, body, nil)
}

func (h *HouseholdActions) RemoveReturnDate(ctx context.Context, setID string) error {
	if err := h.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/household/users/%s/schedule/%s", h.c.UserID, setID)
	return h.c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

func (h *HouseholdActions) AddDevice(ctx context.Context, householdID string, body map[string]any) error {
	if err := h.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/household/households/%s/devices", householdID)
	return h.c.do(ctx, http.MethodPost, path, nil, body, nil)
}

func (h *HouseholdActions) UpdateDevice(ctx context.Context, deviceID string, body map[string]any) error {
	if err := h.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/household/devices/%s", deviceID)
	return h.c.do(ctx, http.MethodPut, path, nil, body, nil)
}

func (h *HouseholdActions) RemoveDevice(ctx context.Context, deviceID string) error {
	if err := h.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/household/devices/%s", deviceID)
	return h.c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

func (h *HouseholdActions) InviteUser(ctx context.Context, householdID string, body map[string]any) error {
	if err := h.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/household/households/%s/users", householdID)
	return h.c.do(ctx, http.MethodPost, path, nil, body, nil)
}

func (h *HouseholdActions) RespondToInvitation(ctx context.Context, householdID, userID string, body map[string]any) error {
	if err := h.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/household/households/%s/users/%s", householdID, userID)
	return h.c.do(ctx, http.MethodPost, path, nil, body, nil)
}

func (h *HouseholdActions) RemoveGuest(ctx context.Context, householdID, userID string) error {
	if err := h.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/household/households/%s/users/%s", householdID, userID)
	return h.c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

func (h *HouseholdActions) AddGuests(ctx context.Context, householdID, deviceID string, body map[string]any) error {
	if err := h.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/household/households/%s/devices/%s/guests", householdID, deviceID)
	return h.c.do(ctx, http.MethodPost, path, nil, body, nil)
}

func (h *HouseholdActions) UpdateDeviceSet(ctx context.Context, householdID, setID string, body map[string]any) error {
	if err := h.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/household/households/%s/sets/%s", householdID, setID)
	return h.c.do(ctx, http.MethodPut, path, nil, body, nil)
}

func (h *HouseholdActions) RemoveDeviceSet(ctx context.Context, householdID, setID string) error {
	if err := h.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/household/households/%s/sets/%s", householdID, setID)
	return h.c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

func (h *HouseholdActions) RemoveDeviceAssignment(ctx context.Context, deviceID, userID string) error {
	if err := h.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/household/devices/%s/assignment/users/%s", deviceID, userID)
	return h.c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}
