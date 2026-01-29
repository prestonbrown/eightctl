package client

import (
	"context"
	"fmt"
	"net/http"
)

// PushActions groups push notification endpoints.
type PushActions struct {
	c *Client
}

// Push helper accessor.
func (c *Client) Push() *PushActions { return &PushActions{c: c} }

func (p *PushActions) UpdatePushToken(ctx context.Context, deviceID string, body map[string]any) error {
	path := fmt.Sprintf("/users/me/push-targets/%s", deviceID)
	return p.c.do(ctx, http.MethodPut, path, nil, body, nil)
}

func (p *PushActions) DeletePushToken(ctx context.Context, token string) error {
	path := fmt.Sprintf("/users/me/push-targets/token/%s", token)
	return p.c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}
