package cmd

import (
	"context"
	"fmt"

	"github.com/steipete/eightctl/internal/client"
	"github.com/steipete/eightctl/internal/model"
)

// getUserIDForSide fetches the user ID for a given side.
func getUserIDForSide(ctx context.Context, cl *client.Client, side model.Side) (string, error) {
	device, err := cl.Device().GetWithUsers(ctx)
	if err != nil {
		return "", err
	}
	switch side {
	case model.Left:
		if device.LeftUserID == "" {
			return "", fmt.Errorf("no user assigned to left side")
		}
		return device.LeftUserID, nil
	case model.Right:
		if device.RightUserID == "" {
			return "", fmt.Errorf("no user assigned to right side")
		}
		return device.RightUserID, nil
	default:
		return "", fmt.Errorf("invalid side")
	}
}
