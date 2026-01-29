package client

import (
	"context"
	"fmt"
	"net/http"
)

// SubscriptionsActions groups subscription endpoints.
type SubscriptionsActions struct {
	c *Client
}

// Subscriptions helper accessor.
func (c *Client) Subscriptions() *SubscriptionsActions { return &SubscriptionsActions{c: c} }

// Subscriptions returns the user's subscriptions.
// Uses v3 API: GET v3/users/{userId}/subscriptions
func (s *SubscriptionsActions) Subscriptions(ctx context.Context, out any) error {
	if err := s.c.requireUser(ctx); err != nil {
		return err
	}
	return s.c.doV3(ctx, http.MethodGet, fmt.Sprintf("/users/%s/subscriptions", s.c.UserID), nil, nil, out)
}

// CreateTemporarySubscription creates a temporary subscription.
// Uses v3 API: POST v3/users/{userId}/subscriptions/temporary
func (s *SubscriptionsActions) CreateTemporarySubscription(ctx context.Context, body map[string]any) error {
	if err := s.c.requireUser(ctx); err != nil {
		return err
	}
	return s.c.doV3(ctx, http.MethodPost, fmt.Sprintf("/users/%s/subscriptions/temporary", s.c.UserID), nil, body, nil)
}

// RedeemSubscription redeems a subscription code.
// Uses v3 API: POST v3/users/{userId}/subscriptions/redeem
func (s *SubscriptionsActions) RedeemSubscription(ctx context.Context, body map[string]any) error {
	if err := s.c.requireUser(ctx); err != nil {
		return err
	}
	return s.c.doV3(ctx, http.MethodPost, fmt.Sprintf("/users/%s/subscriptions/redeem", s.c.UserID), nil, body, nil)
}
