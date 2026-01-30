package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type TagsActions struct{ c *Client }

func (c *Client) Tags() *TagsActions { return &TagsActions{c: c} }

// GetTags retrieves sleep session tags for a date range.
// from/to should be ISO date strings (e.g., "2024-01-01").
func (t *TagsActions) GetTags(ctx context.Context, from, to string) (any, error) {
	if err := t.c.requireUser(ctx); err != nil {
		return nil, err
	}
	q := url.Values{}
	if from != "" {
		q.Set("from", from)
	}
	if to != "" {
		q.Set("to", to)
	}
	path := fmt.Sprintf("/users/%s/tags", t.c.UserID)
	var res any
	err := t.c.do(ctx, http.MethodGet, path, q, nil, &res)
	return res, err
}

// SaveTags saves tags for a specific day.
// day should be an ISO date string (e.g., "2024-01-01").
func (t *TagsActions) SaveTags(ctx context.Context, day string, tags any) error {
	if err := t.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/days/%s/tags", t.c.UserID, day)
	return t.c.do(ctx, http.MethodPut, path, nil, tags, nil)
}

// Presence tags (ground truth for sleep detection calibration)

// GetPresenceTags retrieves presence/truth tags for the user.
func (t *TagsActions) GetPresenceTags(ctx context.Context) (any, error) {
	if err := t.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/truth-tags", t.c.UserID)
	var res any
	err := t.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// CreatePresenceTag creates a new presence tag.
func (t *TagsActions) CreatePresenceTag(ctx context.Context, tag any) (any, error) {
	if err := t.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/truth-tags", t.c.UserID)
	var res any
	err := t.c.do(ctx, http.MethodPost, path, nil, tag, &res)
	return res, err
}

// UpdatePresenceTag updates an existing presence tag.
func (t *TagsActions) UpdatePresenceTag(ctx context.Context, tagID string, tag any) (any, error) {
	if err := t.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/truth-tags/%s", t.c.UserID, tagID)
	var res any
	err := t.c.do(ctx, http.MethodPut, path, nil, tag, &res)
	return res, err
}

// DeletePresenceTag deletes a presence tag.
func (t *TagsActions) DeletePresenceTag(ctx context.Context, tagID string) error {
	if err := t.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/truth-tags/%s", t.c.UserID, tagID)
	return t.c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}
