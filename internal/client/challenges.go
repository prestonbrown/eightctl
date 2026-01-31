package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// ChallengesActions groups challenges endpoints.
type ChallengesActions struct {
	c *Client
}

// Challenges helper accessor.
func (c *Client) Challenges() *ChallengesActions { return &ChallengesActions{c: c} }

// Challenges retrieves challenges for the user, optionally filtered by state.
func (a *ChallengesActions) Challenges(ctx context.Context, state string, out any) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	var q url.Values
	if state != "" {
		q = url.Values{}
		q.Set("state", state)
	}
	path := fmt.Sprintf("/users/%s/challenges", a.c.UserID)
	return a.c.do(ctx, http.MethodGet, path, q, nil, out)
}
