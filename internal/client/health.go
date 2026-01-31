package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// HealthActions groups health integration endpoints.
type HealthActions struct {
	c *Client
}

// Health helper accessor.
func (c *Client) Health() *HealthActions { return &HealthActions{c: c} }

// HealthSurvey fetches the health survey test-drive data.
func (h *HealthActions) HealthSurvey(ctx context.Context, out any) error {
	return h.c.do(ctx, http.MethodGet, "/health-survey/test-drive", nil, nil, out)
}

// UpdateHealthSurvey updates the health survey with validation enabled.
func (h *HealthActions) UpdateHealthSurvey(ctx context.Context, body map[string]any) error {
	q := url.Values{"enableValidation": []string{"true"}}
	return h.c.do(ctx, http.MethodPatch, "/health-survey/test-drive", q, body, nil)
}

// HealthCheckpoints fetches health integration checkpoints for a source.
func (h *HealthActions) HealthCheckpoints(ctx context.Context, sourceID string, out any) error {
	if err := h.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/health-integrations/sources/%s/checkpoints", h.c.UserID, sourceID)
	return h.c.do(ctx, http.MethodGet, path, nil, nil, out)
}

// UploadHealthData uploads health data for a source.
func (h *HealthActions) UploadHealthData(ctx context.Context, sourceID string, body map[string]any) error {
	if err := h.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/health-integrations/sources/%s", h.c.UserID, sourceID)
	return h.c.do(ctx, http.MethodPost, path, nil, body, nil)
}
