package mqtt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoveryTopic(t *testing.T) {
	tests := []struct {
		name        string
		topicPrefix string
		deviceID    string
		side        string
		expected    string
	}{
		{
			name:        "left side with homeassistant prefix",
			topicPrefix: "homeassistant",
			deviceID:    "device-123",
			side:        "left",
			expected:    "homeassistant/climate/device-123_left/config",
		},
		{
			name:        "right side with homeassistant prefix",
			topicPrefix: "homeassistant",
			deviceID:    "device-456",
			side:        "right",
			expected:    "homeassistant/climate/device-456_right/config",
		},
		{
			name:        "custom prefix",
			topicPrefix: "custom/prefix",
			deviceID:    "pod-1",
			side:        "left",
			expected:    "custom/prefix/climate/pod-1_left/config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DiscoveryTopic(tt.topicPrefix, tt.deviceID, tt.side)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateDiscoveryConfigs(t *testing.T) {
	configs := GenerateDiscoveryConfigs("homeassistant", "device-123", "Bedroom Pod")

	require.Len(t, configs, 2)
	require.Contains(t, configs, "left")
	require.Contains(t, configs, "right")

	// Check left config
	left := configs["left"]
	assert.Equal(t, "Bedroom Pod left", left.Name)
	assert.Equal(t, "eightsleep_device-123_left", left.UniqueID)
	assert.Equal(t, float64(-100), left.MinTemp)
	assert.Equal(t, float64(100), left.MaxTemp)
	assert.Equal(t, float64(1), left.TempStep)
	assert.Contains(t, left.Modes, "off")
	assert.Contains(t, left.Modes, "heat")
	assert.Contains(t, left.Modes, "cool")

	// Check right config
	right := configs["right"]
	assert.Equal(t, "Bedroom Pod right", right.Name)
	assert.Equal(t, "eightsleep_device-123_right", right.UniqueID)
}

func TestGenerateDiscoveryConfigs_Topics(t *testing.T) {
	configs := GenerateDiscoveryConfigs("homeassistant", "pod-1", "Test Pod")

	left := configs["left"]

	// State topics
	assert.Equal(t, "eightsleep/pod-1/left/temperature", left.CurrentTemperatureTopic)
	assert.Equal(t, "eightsleep/pod-1/left/temperature", left.TemperatureStateTopic)
	assert.Equal(t, "eightsleep/pod-1/left/mode", left.ModeStateTopic)
	assert.Equal(t, "eightsleep/pod-1/availability", left.AvailabilityTopic)

	// Command topics
	assert.Equal(t, "eightsleep/pod-1/left/set_temperature", left.TemperatureCommandTopic)
	assert.Equal(t, "eightsleep/pod-1/left/set_mode", left.ModeCommandTopic)
}

func TestGenerateDiscoveryConfigs_DeviceInfo(t *testing.T) {
	configs := GenerateDiscoveryConfigs("homeassistant", "device-abc", "Living Room Pod")

	left := configs["left"]
	device := left.Device

	assert.Equal(t, "Living Room Pod", device.Name)
	assert.Equal(t, "Eight Sleep", device.Manufacturer)
	assert.Equal(t, "Pod", device.Model)
	require.Len(t, device.Identifiers, 1)
	assert.Equal(t, "device-abc", device.Identifiers[0])
}

func TestClimateDiscovery_Modes(t *testing.T) {
	configs := GenerateDiscoveryConfigs("homeassistant", "device-1", "Pod")

	for side, config := range configs {
		t.Run(side, func(t *testing.T) {
			// All sides should have the same modes
			require.Len(t, config.Modes, 3)
			assert.Equal(t, "off", config.Modes[0])
			assert.Equal(t, "heat", config.Modes[1])
			assert.Equal(t, "cool", config.Modes[2])
		})
	}
}
