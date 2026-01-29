package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// InsightsActions groups LLM insights endpoints.
type InsightsActions struct {
	c *Client
}

// Insights helper accessor.
func (c *Client) Insights() *InsightsActions { return &InsightsActions{c: c} }

// LLMInsights retrieves AI-generated sleep insights for a date range.
func (i *InsightsActions) LLMInsights(ctx context.Context, from, to string) (any, error) {
	if err := i.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/llm-insights", i.c.UserID)
	q := url.Values{"from": []string{from}, "to": []string{to}}
	var res any
	err := i.c.do(ctx, http.MethodGet, path, q, nil, &res)
	return res, err
}

// CreateLLMInsightsBatch creates a batch of AI insights.
func (i *InsightsActions) CreateLLMInsightsBatch(ctx context.Context, body map[string]any) error {
	if err := i.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/llm-insights/batch", i.c.UserID)
	return i.c.do(ctx, http.MethodPost, path, nil, body, nil)
}

// LLMInsightsSettings retrieves AI insights settings.
func (i *InsightsActions) LLMInsightsSettings(ctx context.Context) (any, error) {
	if err := i.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/llm-insights/settings", i.c.UserID)
	var res any
	err := i.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// UpdateLLMInsightsSettings updates AI insights settings.
func (i *InsightsActions) UpdateLLMInsightsSettings(ctx context.Context, body map[string]any) error {
	if err := i.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/llm-insights/settings", i.c.UserID)
	return i.c.do(ctx, http.MethodPut, path, nil, body, nil)
}

// SubmitLLMInsightFeedback submits feedback for a specific AI insight.
func (i *InsightsActions) SubmitLLMInsightFeedback(ctx context.Context, insightID string, body map[string]any) error {
	if err := i.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/llm-insights/%s/feedback", i.c.UserID, insightID)
	return i.c.do(ctx, http.MethodPost, path, nil, body, nil)
}
