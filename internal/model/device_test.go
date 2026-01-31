package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeviceState_GetSide(t *testing.T) {
	left := &UserState{ID: "left-user", Side: Left}
	right := &UserState{ID: "right-user", Side: Right}
	d := &DeviceState{
		LeftUser:  left,
		RightUser: right,
	}

	assert.Equal(t, left, d.GetSide(Left))
	assert.Equal(t, right, d.GetSide(Right))
}

func TestDeviceState_GetSide_Nil(t *testing.T) {
	d := &DeviceState{}
	assert.Nil(t, d.GetSide(Left))
	assert.Nil(t, d.GetSide(Right))
}

func TestDeviceState_HasBothSides(t *testing.T) {
	tests := []struct {
		name     string
		left     *UserState
		right    *UserState
		expected bool
	}{
		{"both", &UserState{}, &UserState{}, true},
		{"left only", &UserState{}, nil, false},
		{"right only", nil, &UserState{}, false},
		{"neither", nil, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DeviceState{LeftUser: tt.left, RightUser: tt.right}
			assert.Equal(t, tt.expected, d.HasBothSides())
		})
	}
}

func TestDeviceState_JSON(t *testing.T) {
	d := &DeviceState{
		ID:              "device-123",
		RoomTemperature: 68.5,
		HasWater:        true,
		LeftUser:        &UserState{ID: "left", Side: Left, TargetLevel: -20},
	}

	// Just verify it marshals without error
	data, err := json.Marshal(d)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"room_temperature":68.5`)
	assert.Contains(t, string(data), `"target_level":-20`)
}
