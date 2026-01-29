package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// SettingsActions groups user settings and preferences endpoints.
type SettingsActions struct{ c *Client }

// Settings helper accessor.
func (c *Client) Settings() *SettingsActions { return &SettingsActions{c: c} }

// TapSettings returns tap gesture settings for the specified device.
func (s *SettingsActions) TapSettings(ctx context.Context, deviceID string) (any, error) {
	if err := s.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/devices/%s/tap-settings", s.c.UserID, deviceID)
	var res any
	err := s.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// UpdateTapSettings updates tap gesture settings for the specified device.
func (s *SettingsActions) UpdateTapSettings(ctx context.Context, deviceID string, body map[string]any) error {
	if err := s.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/devices/%s/tap-settings", s.c.UserID, deviceID)
	return s.c.do(ctx, http.MethodPut, path, nil, body, nil)
}

// TapHistory returns the user's tap gesture history.
func (s *SettingsActions) TapHistory(ctx context.Context, from string) (any, error) {
	if err := s.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/tap-history", s.c.UserID)
	q := url.Values{"from": []string{from}}
	var res any
	err := s.c.do(ctx, http.MethodGet, path, q, nil, &res)
	return res, err
}

// LevelSuggestions returns temperature level suggestions for the user.
func (s *SettingsActions) LevelSuggestions(ctx context.Context) (any, error) {
	if err := s.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/level-suggestions", s.c.UserID)
	var res any
	err := s.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// BlanketRecommendations returns blanket temperature recommendations.
func (s *SettingsActions) BlanketRecommendations(ctx context.Context) (any, error) {
	if err := s.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/recommendations/blanket", s.c.UserID)
	var res any
	err := s.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// Perks returns member perks for the user.
func (s *SettingsActions) Perks(ctx context.Context) (any, error) {
	if err := s.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/perks", s.c.UserID)
	var res any
	err := s.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// GetReferralLink generates/retrieves the user's personal referral link (v2 API).
func (s *SettingsActions) GetReferralLink(ctx context.Context) (any, error) {
	if err := s.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/v2/users/%s/referral/personal-referral-link", s.c.UserID)
	var res any
	err := s.c.do(ctx, http.MethodPut, path, nil, nil, &res)
	return res, err
}

// ReferralCampaigns returns available referral campaigns (v2 API).
func (s *SettingsActions) ReferralCampaigns(ctx context.Context) (any, error) {
	if err := s.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/v2/users/%s/referral/campaigns", s.c.UserID)
	var res any
	err := s.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// Purchases returns purchase tracker information.
func (s *SettingsActions) Purchases(ctx context.Context) (any, error) {
	var res any
	err := s.c.do(ctx, http.MethodGet, "/purchase-tracker", nil, nil, &res)
	return res, err
}

// MaintenanceInsertStatus returns device maintenance insert status.
func (s *SettingsActions) MaintenanceInsertStatus(ctx context.Context) (any, error) {
	if err := s.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/user/%s/device_maintenance/maintenance_insert", s.c.UserID)
	q := url.Values{"v": []string{"2"}}
	var res any
	err := s.c.do(ctx, http.MethodGet, path, q, nil, &res)
	return res, err
}
