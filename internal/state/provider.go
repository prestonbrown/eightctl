package state

import (
	"context"

	"github.com/steipete/eightctl/internal/model"
)

// StateProvider defines the interface for getting and setting device state.
type StateProvider interface {
	// GetState returns the current device state.
	GetState(ctx context.Context) (*model.DeviceState, error)

	// SetTemperature sets the target temperature level for a side.
	SetTemperature(ctx context.Context, side model.Side, level int) error

	// TurnOn powers on a side.
	TurnOn(ctx context.Context, side model.Side) error

	// TurnOff powers off a side.
	TurnOff(ctx context.Context, side model.Side) error
}

// StateChange represents a change in device state.
type StateChange struct {
	Old *model.DeviceState
	New *model.DeviceState
}

// PresenceChange represents a change in user presence.
type PresenceChange struct {
	Side    model.Side
	Present bool
	User    *model.UserState
}

// Observer receives state change notifications.
type Observer interface {
	// OnStateChange is called when device state changes.
	OnStateChange(change StateChange)

	// OnPresenceChange is called when presence changes for a side.
	OnPresenceChange(change PresenceChange)
}
