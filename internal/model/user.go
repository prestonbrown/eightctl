package model

import "time"

const PresenceTimeout = 10 * time.Minute

// UserState represents the state of one side of the bed.
type UserState struct {
	ID                string     `json:"id"`
	Email             string     `json:"email"`
	Side              Side       `json:"side"`
	BedTemperature    float64    `json:"bed_temperature"`
	TargetLevel       int        `json:"target_level"`
	State             PowerState `json:"state"`
	SleepStage        SleepStage `json:"sleep_stage"`
	HeartRate         float64    `json:"heart_rate"`
	HRV               float64    `json:"hrv"`
	BreathRate        float64    `json:"breath_rate"`
	LastHeartRateTime time.Time  `json:"last_heart_rate_time"`
}

// IsPresent returns true if the user appears to be in bed.
// Based on pyEight: presence is determined by heart rate data within the last 10 minutes.
func (u *UserState) IsPresent() bool {
	if u.LastHeartRateTime.IsZero() {
		return false
	}
	return time.Since(u.LastHeartRateTime) < PresenceTimeout
}

// IsOn returns true if the side is actively heating/cooling.
func (u *UserState) IsOn() bool {
	return u.State == PowerSmart || u.State == PowerManual
}
