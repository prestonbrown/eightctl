package model

// DeviceState represents the complete state of an Eight Sleep pod.
type DeviceState struct {
	ID              string     `json:"id"`
	RoomTemperature float64    `json:"room_temperature"`
	HasWater        bool       `json:"has_water"`
	IsPriming       bool       `json:"is_priming"`
	NeedsPriming    bool       `json:"needs_priming"`
	LeftUser        *UserState `json:"left,omitempty"`
	RightUser       *UserState `json:"right,omitempty"`
}

// GetSide returns the UserState for the specified side.
func (d *DeviceState) GetSide(side Side) *UserState {
	switch side {
	case Left:
		return d.LeftUser
	case Right:
		return d.RightUser
	default:
		return nil
	}
}

// HasBothSides returns true if both left and right users are configured.
func (d *DeviceState) HasBothSides() bool {
	return d.LeftUser != nil && d.RightUser != nil
}
