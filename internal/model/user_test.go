package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserState_IsPresent_WithRecentHeartRate(t *testing.T) {
	u := &UserState{
		LastHeartRateTime: time.Now().Add(-5 * time.Minute),
	}
	assert.True(t, u.IsPresent())
}

func TestUserState_IsPresent_WithStaleHeartRate(t *testing.T) {
	u := &UserState{
		LastHeartRateTime: time.Now().Add(-15 * time.Minute),
	}
	assert.False(t, u.IsPresent())
}

func TestUserState_IsPresent_WithZeroTime(t *testing.T) {
	u := &UserState{}
	assert.False(t, u.IsPresent())
}

func TestUserState_IsOn(t *testing.T) {
	tests := []struct {
		state    PowerState
		expected bool
	}{
		{PowerOff, false},
		{PowerSmart, true},
		{PowerManual, true},
	}
	for _, tt := range tests {
		t.Run(tt.state.String(), func(t *testing.T) {
			u := &UserState{State: tt.state}
			assert.Equal(t, tt.expected, u.IsOn())
		})
	}
}
