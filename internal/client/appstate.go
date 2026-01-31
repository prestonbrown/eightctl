package client

import (
	"context"
	"fmt"
	"net/http"
)

// AppStateActions groups app-state endpoints.
type AppStateActions struct {
	c *Client
}

// AppState helper accessor.
func (c *Client) AppState() *AppStateActions { return &AppStateActions{c: c} }

// MessagesState retrieves the user's app-state messages.
func (a *AppStateActions) MessagesState(ctx context.Context, out any) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/app-state/messages", a.c.UserID)
	return a.c.do(ctx, http.MethodGet, path, nil, nil, out)
}

// UpdateMessagesState replaces the user's app-state messages (PUT).
func (a *AppStateActions) UpdateMessagesState(ctx context.Context, body map[string]any) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/app-state/messages", a.c.UserID)
	return a.c.do(ctx, http.MethodPut, path, nil, body, nil)
}

// PatchMessagesState partially updates the user's app-state messages (PATCH).
func (a *AppStateActions) PatchMessagesState(ctx context.Context, body map[string]any) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/app-state/messages", a.c.UserID)
	return a.c.do(ctx, http.MethodPatch, path, nil, body, nil)
}
