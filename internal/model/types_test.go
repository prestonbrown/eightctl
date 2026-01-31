package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSide_String(t *testing.T) {
	assert.Equal(t, "left", Left.String())
	assert.Equal(t, "right", Right.String())
}

func TestParseSide(t *testing.T) {
	tests := []struct {
		input    string
		expected Side
		wantErr  bool
	}{
		{"left", Left, false},
		{"LEFT", Left, false},
		{"right", Right, false},
		{"RIGHT", Right, false},
		{"invalid", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseSide(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestSide_JSON(t *testing.T) {
	type wrapper struct {
		Side Side `json:"side"`
	}
	w := wrapper{Side: Left}
	data, err := json.Marshal(w)
	require.NoError(t, err)
	assert.Equal(t, `{"side":"left"}`, string(data))

	var w2 wrapper
	err = json.Unmarshal([]byte(`{"side":"right"}`), &w2)
	require.NoError(t, err)
	assert.Equal(t, Right, w2.Side)
}

func TestPowerState_String(t *testing.T) {
	assert.Equal(t, "off", PowerOff.String())
	assert.Equal(t, "smart", PowerSmart.String())
	assert.Equal(t, "manual", PowerManual.String())
}

func TestParsePowerState(t *testing.T) {
	tests := []struct {
		input    string
		expected PowerState
	}{
		{"off", PowerOff},
		{"smart", PowerSmart},
		{"manual", PowerManual},
		{"unknown", PowerOff}, // Default to off
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, ParsePowerState(tt.input))
		})
	}
}

func TestSleepStage_String(t *testing.T) {
	assert.Equal(t, "awake", StageAwake.String())
	assert.Equal(t, "light", StageLight.String())
	assert.Equal(t, "deep", StageDeep.String())
	assert.Equal(t, "rem", StageREM.String())
	assert.Equal(t, "unknown", StageUnknown.String())
}

func TestParseSleepStage(t *testing.T) {
	tests := []struct {
		input    string
		expected SleepStage
	}{
		{"awake", StageAwake},
		{"light", StageLight},
		{"deep", StageDeep},
		{"rem", StageREM},
		{"REM", StageREM},
		{"invalid", StageUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, ParseSleepStage(tt.input))
		})
	}
}
