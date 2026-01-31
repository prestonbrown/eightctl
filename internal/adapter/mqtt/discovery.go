// Package mqtt provides Home Assistant MQTT integration for Eight Sleep Pods.
package mqtt

import "fmt"

// DeviceInfo represents Home Assistant device registry information.
type DeviceInfo struct {
	Identifiers  []string `json:"identifiers"`
	Name         string   `json:"name"`
	Manufacturer string   `json:"manufacturer"`
	Model        string   `json:"model"`
}

// ClimateDiscovery represents a Home Assistant MQTT climate entity configuration.
type ClimateDiscovery struct {
	Name     string     `json:"name"`
	UniqueID string     `json:"unique_id"`
	Device   DeviceInfo `json:"device"`

	// State topics
	CurrentTemperatureTopic string `json:"current_temperature_topic"`
	TemperatureStateTopic   string `json:"temperature_state_topic"`
	ModeStateTopic          string `json:"mode_state_topic"`
	AvailabilityTopic       string `json:"availability_topic"`

	// Command topics
	TemperatureCommandTopic string `json:"temperature_command_topic"`
	ModeCommandTopic        string `json:"mode_command_topic"`

	// Temperature range (Eight Sleep uses -100 to +100 levels)
	MinTemp  float64 `json:"min_temp"`
	MaxTemp  float64 `json:"max_temp"`
	TempStep float64 `json:"temp_step"`
	TempUnit string  `json:"temperature_unit"`

	// Modes
	Modes []string `json:"modes"`
}

// DiscoveryTopic returns the MQTT discovery topic for a climate entity.
// Format: {topicPrefix}/climate/{deviceID}_{side}/config
func DiscoveryTopic(topicPrefix, deviceID, side string) string {
	return fmt.Sprintf("%s/climate/%s_%s/config", topicPrefix, deviceID, side)
}

// GenerateDiscoveryConfigs creates Home Assistant discovery payloads for a device.
// topicPrefix is typically "homeassistant" for default HA setup.
// deviceID is the Eight Sleep device ID.
// deviceName is a human-readable name like "Bedroom Pod".
// Returns a map with keys "left" and "right" containing the discovery configs.
func GenerateDiscoveryConfigs(topicPrefix, deviceID, deviceName string) map[string]ClimateDiscovery {
	sides := []string{"left", "right"}
	configs := make(map[string]ClimateDiscovery, 2)

	device := DeviceInfo{
		Identifiers:  []string{deviceID},
		Name:         deviceName,
		Manufacturer: "Eight Sleep",
		Model:        "Pod",
	}

	for _, side := range sides {
		configs[side] = ClimateDiscovery{
			Name:     fmt.Sprintf("%s %s", deviceName, side),
			UniqueID: fmt.Sprintf("eightsleep_%s_%s", deviceID, side),
			Device:   device,

			// State topics
			CurrentTemperatureTopic: fmt.Sprintf("eightsleep/%s/%s/temperature", deviceID, side),
			TemperatureStateTopic:   fmt.Sprintf("eightsleep/%s/%s/temperature", deviceID, side),
			ModeStateTopic:          fmt.Sprintf("eightsleep/%s/%s/mode", deviceID, side),
			AvailabilityTopic:       fmt.Sprintf("eightsleep/%s/availability", deviceID),

			// Command topics
			TemperatureCommandTopic: fmt.Sprintf("eightsleep/%s/%s/set_temperature", deviceID, side),
			ModeCommandTopic:        fmt.Sprintf("eightsleep/%s/%s/set_mode", deviceID, side),

			// Temperature range (Eight Sleep levels: -100 to +100)
			MinTemp:  -100,
			MaxTemp:  100,
			TempStep: 1,
			TempUnit: "C", // HA requires a unit, but we use levels

			// Modes: off, heat (positive levels), cool (negative levels)
			Modes: []string{"off", "heat", "cool"},
		}
	}

	return configs
}
