package adapter

import (
	"testing"

	"github.com/steipete/eightctl/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestAction_String(t *testing.T) {
	tests := []struct {
		action   Action
		expected string
	}{
		{ActionOn, "on"},
		{ActionOff, "off"},
		{ActionSetTemp, "set_temperature"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.action))
		})
	}
}

func TestCommand_Fields(t *testing.T) {
	level := -20
	cmd := Command{
		Action:      ActionSetTemp,
		Side:        model.Left,
		Temperature: &level,
	}

	assert.Equal(t, ActionSetTemp, cmd.Action)
	assert.Equal(t, model.Left, cmd.Side)
	assert.NotNil(t, cmd.Temperature)
	assert.Equal(t, -20, *cmd.Temperature)
}

func TestCommand_NilTemperature(t *testing.T) {
	cmd := Command{
		Action: ActionOn,
		Side:   model.Right,
	}

	assert.Nil(t, cmd.Temperature)
}

func TestConfig_Fields(t *testing.T) {
	cfg := Config{
		PollInterval: 60,
		DeviceID:     "device-123",
		DeviceName:   "Bedroom Pod",
	}

	assert.Equal(t, 60, int(cfg.PollInterval))
	assert.Equal(t, "device-123", cfg.DeviceID)
	assert.Equal(t, "Bedroom Pod", cfg.DeviceName)
}
