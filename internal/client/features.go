package client

import (
	"context"
	"fmt"
	"net/http"
)

type FeaturesActions struct{ c *Client }

func (c *Client) Features() *FeaturesActions { return &FeaturesActions{c: c} }

// GetFeatureFlags retrieves feature flags/release features for the user.
// This can help understand which features are enabled for the account.
func (f *FeaturesActions) GetFeatureFlags(ctx context.Context) (any, error) {
	if err := f.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/release-features", f.c.UserID)
	var res any
	err := f.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}
