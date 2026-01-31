package client

import (
	"context"
	"fmt"
	"net/http"
)

// UserActions groups user endpoints.
type UserActions struct {
	c *Client
}

// Users helper accessor.
func (c *Client) Users() *UserActions { return &UserActions{c: c} }

// GetMe fetches the current user's profile.
func (u *UserActions) GetMe(ctx context.Context, out any) error {
	return u.c.do(ctx, http.MethodGet, "/users/me", nil, nil, out)
}

// GetUser fetches a user by ID.
func (u *UserActions) GetUser(ctx context.Context, userID string, out any) error {
	path := fmt.Sprintf("/users/%s", userID)
	return u.c.do(ctx, http.MethodGet, path, nil, nil, out)
}

// UpdateUser updates the current user's profile.
func (u *UserActions) UpdateUser(ctx context.Context, body map[string]any) error {
	if err := u.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s", u.c.UserID)
	return u.c.do(ctx, http.MethodPut, path, nil, body, nil)
}

// UpdateEmail updates the current user's email address.
func (u *UserActions) UpdateEmail(ctx context.Context, body map[string]any) error {
	if err := u.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/email", u.c.UserID)
	return u.c.do(ctx, http.MethodPost, path, nil, body, nil)
}

// PasswordReset initiates a password reset for the given email.
func (u *UserActions) PasswordReset(ctx context.Context, email string) error {
	body := map[string]string{"email": email}
	return u.c.do(ctx, http.MethodPost, "/users/password-reset", nil, body, nil)
}
