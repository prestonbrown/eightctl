package model

import (
	"encoding/json"
	"errors"
	"strings"
)

// Side represents left or right side of the bed.
type Side int

const (
	Left Side = iota + 1
	Right
)

func (s Side) String() string {
	switch s {
	case Left:
		return "left"
	case Right:
		return "right"
	default:
		return ""
	}
}

func ParseSide(s string) (Side, error) {
	switch strings.ToLower(s) {
	case "left":
		return Left, nil
	case "right":
		return Right, nil
	default:
		return 0, errors.New("invalid side: must be 'left' or 'right'")
	}
}

func (s Side) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Side) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	parsed, err := ParseSide(str)
	if err != nil {
		return err
	}
	*s = parsed
	return nil
}

// PowerState represents the pod's power mode.
type PowerState int

const (
	PowerOff PowerState = iota
	PowerSmart
	PowerManual
)

func (p PowerState) String() string {
	switch p {
	case PowerOff:
		return "off"
	case PowerSmart:
		return "smart"
	case PowerManual:
		return "manual"
	default:
		return "off"
	}
}

func ParsePowerState(s string) PowerState {
	switch strings.ToLower(s) {
	case "smart":
		return PowerSmart
	case "manual":
		return PowerManual
	default:
		return PowerOff
	}
}
