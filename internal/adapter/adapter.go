// Package adapter provides the framework for smart home integrations.
// Adapters connect Eight Sleep Pods to smart home platforms like
// Home Assistant, HomeKit, and Google Home.
package adapter

import (
	"context"
	"time"

	"github.com/steipete/eightctl/internal/model"
)

// Action represents a command action from a smart home platform.
type Action string

const (
	ActionOn      Action = "on"
	ActionOff     Action = "off"
	ActionSetTemp Action = "set_temperature"
)

// Command represents a command received from a smart home platform.
type Command struct {
	Action      Action
	Side        model.Side
	Temperature *int // only for ActionSetTemp
}

// Adapter defines the interface for smart home platform integrations.
type Adapter interface {
	// Start begins the adapter's main loop (connect, publish discovery, etc.)
	Start(ctx context.Context) error

	// HandleCommand processes a command from the smart home platform
	HandleCommand(ctx context.Context, cmd Command) error

	// Stop gracefully shuts down the adapter
	Stop() error
}

// Config holds common configuration for adapters.
type Config struct {
	PollInterval time.Duration
	DeviceID     string
	DeviceName   string
}
