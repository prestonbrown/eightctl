# Smart Home Integration Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add Home Assistant (MQTT) and Hubitat integrations to eightctl with dual-side bed support.

**Architecture:** Domain models define strongly-typed state; StateManager provides caching and observer notifications; Platform adapters (MQTT, Hubitat) implement common interface and use StateManager for all operations.

**Tech Stack:** Go 1.21+, Eclipse Paho MQTT, net/http for Hubitat server, Groovy for Hubitat drivers.

---

## Phase 1: Domain Models

### Task 1.1: Create Side and PowerState Types

**Files:**
- Create: `internal/model/types.go`
- Create: `internal/model/types_test.go`

**Step 1: Write the failing test**

```go
// internal/model/types_test.go
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
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/model/... -v`
Expected: FAIL - package does not exist

**Step 3: Write minimal implementation**

```go
// internal/model/types.go
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
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/model/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/model/types.go internal/model/types_test.go
git commit -m "feat(model): add Side and PowerState types"
```

---

### Task 1.2: Create SleepStage Type

**Files:**
- Modify: `internal/model/types.go`
- Modify: `internal/model/types_test.go`

**Step 1: Write the failing test**

Add to `internal/model/types_test.go`:

```go
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
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/model/... -v -run SleepStage`
Expected: FAIL - undefined: StageAwake

**Step 3: Write minimal implementation**

Add to `internal/model/types.go`:

```go
// SleepStage represents the current sleep stage.
type SleepStage int

const (
	StageUnknown SleepStage = iota
	StageAwake
	StageLight
	StageDeep
	StageREM
)

func (s SleepStage) String() string {
	switch s {
	case StageAwake:
		return "awake"
	case StageLight:
		return "light"
	case StageDeep:
		return "deep"
	case StageREM:
		return "rem"
	default:
		return "unknown"
	}
}

func ParseSleepStage(s string) SleepStage {
	switch strings.ToLower(s) {
	case "awake":
		return StageAwake
	case "light":
		return StageLight
	case "deep":
		return StageDeep
	case "rem":
		return StageREM
	default:
		return StageUnknown
	}
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/model/... -v -run SleepStage`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/model/types.go internal/model/types_test.go
git commit -m "feat(model): add SleepStage type"
```

---

### Task 1.3: Create UserState Model

**Files:**
- Create: `internal/model/user.go`
- Create: `internal/model/user_test.go`

**Step 1: Write the failing test**

```go
// internal/model/user_test.go
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
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/model/... -v -run UserState`
Expected: FAIL - undefined: UserState

**Step 3: Write minimal implementation**

```go
// internal/model/user.go
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
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/model/... -v -run UserState`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/model/user.go internal/model/user_test.go
git commit -m "feat(model): add UserState with presence detection"
```

---

### Task 1.4: Create DeviceState Model

**Files:**
- Create: `internal/model/device.go`
- Create: `internal/model/device_test.go`

**Step 1: Write the failing test**

```go
// internal/model/device_test.go
package model

import (
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
```

Add import at top of test file:
```go
import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/model/... -v -run DeviceState`
Expected: FAIL - undefined: DeviceState

**Step 3: Write minimal implementation**

```go
// internal/model/device.go
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
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/model/... -v -run DeviceState`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/model/device.go internal/model/device_test.go
git commit -m "feat(model): add DeviceState model"
```

---

## Phase 2: Client Enhancements

### Task 2.1: Add GetDeviceWithUsers Method

**Files:**
- Modify: `internal/client/device.go`
- Modify: `internal/client/device_test.go`

**Step 1: Write the failing test**

Add to `internal/client/device_test.go`:

```go
func TestDeviceActions_GetWithUsers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/users/me" && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]any{
				"user": map[string]any{
					"userId": "user-123",
					"currentDevice": map[string]any{
						"id": "device-456",
					},
				},
			})
		case strings.HasPrefix(r.URL.Path, "/devices/device-456"):
			// Verify filter query param
			assert.Contains(t, r.URL.Query().Get("filter"), "leftUserId")
			assert.Contains(t, r.URL.Query().Get("filter"), "rightUserId")
			json.NewEncoder(w).Encode(map[string]any{
				"result": map[string]any{
					"id":              "device-456",
					"roomTemperature": 68.5,
					"leftUserId":      "left-user-id",
					"rightUserId":     "right-user-id",
					"priming": map[string]any{
						"status": "ready",
					},
					"waterLevel": 100,
				},
			})
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(404)
		}
	}))
	defer srv.Close()

	cl := New("test@test.com", "pass", "", "", "")
	cl.BaseURL = srv.URL
	cl.token = "fake-token"
	cl.tokenExp = time.Now().Add(time.Hour)

	info, err := cl.Device().GetWithUsers(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "device-456", info.ID)
	assert.Equal(t, "left-user-id", info.LeftUserID)
	assert.Equal(t, "right-user-id", info.RightUserID)
	assert.InDelta(t, 68.5, info.RoomTemperature, 0.01)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/client/... -v -run GetWithUsers`
Expected: FAIL - cl.Device().GetWithUsers undefined

**Step 3: Write minimal implementation**

Add to `internal/client/device.go`:

```go
// DeviceWithUsers contains device info including user assignments.
type DeviceWithUsers struct {
	ID              string  `json:"id"`
	LeftUserID      string  `json:"leftUserId"`
	RightUserID     string  `json:"rightUserId"`
	RoomTemperature float64 `json:"roomTemperature"`
	WaterLevel      int     `json:"waterLevel"`
	IsPriming       bool    `json:"isPriming"`
	NeedsPriming    bool    `json:"needsPriming"`
}

// GetWithUsers fetches device info with left/right user assignments.
func (d *DeviceActions) GetWithUsers(ctx context.Context) (*DeviceWithUsers, error) {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s", id)
	query := url.Values{}
	query.Set("filter", "leftUserId,rightUserId,awaySides")

	var res struct {
		Result struct {
			ID              string  `json:"id"`
			LeftUserID      string  `json:"leftUserId"`
			RightUserID     string  `json:"rightUserId"`
			RoomTemperature float64 `json:"roomTemperature"`
			WaterLevel      int     `json:"waterLevel"`
			Priming         struct {
				Status string `json:"status"`
			} `json:"priming"`
		} `json:"result"`
	}
	err = d.c.do(ctx, http.MethodGet, path, query, nil, &res)
	if err != nil {
		return nil, err
	}

	return &DeviceWithUsers{
		ID:              res.Result.ID,
		LeftUserID:      res.Result.LeftUserID,
		RightUserID:     res.Result.RightUserID,
		RoomTemperature: res.Result.RoomTemperature,
		WaterLevel:      res.Result.WaterLevel,
		IsPriming:       res.Result.Priming.Status == "priming",
		NeedsPriming:    res.Result.Priming.Status == "needed",
	}, nil
}
```

Add import: `"net/url"`

**Step 4: Run test to verify it passes**

Run: `go test ./internal/client/... -v -run GetWithUsers`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/client/device.go internal/client/device_test.go
git commit -m "feat(client): add GetWithUsers to fetch device with user assignments"
```

---

### Task 2.2: Add GetUserTemperature Method for Any User

**Files:**
- Modify: `internal/client/eightsleep.go`
- Modify: `internal/client/eightsleep_test.go`

**Step 1: Write the failing test**

Add to `internal/client/eightsleep_test.go`:

```go
func TestClient_GetUserTemperature(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/users/other-user-id/temperature" && r.Method == http.MethodGet {
			json.NewEncoder(w).Encode(map[string]any{
				"currentLevel": -20,
				"currentState": map[string]any{
					"type": "smart",
				},
			})
			return
		}
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(404)
	}))
	defer srv.Close()

	cl := New("test@test.com", "pass", "", "", "")
	cl.BaseURL = srv.URL
	cl.token = "fake-token"
	cl.tokenExp = time.Now().Add(time.Hour)

	status, err := cl.GetUserTemperature(context.Background(), "other-user-id")
	require.NoError(t, err)
	assert.Equal(t, -20, status.CurrentLevel)
	assert.Equal(t, "smart", status.CurrentState.Type)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/client/... -v -run GetUserTemperature`
Expected: FAIL - cl.GetUserTemperature undefined

**Step 3: Write minimal implementation**

Add to `internal/client/eightsleep.go`:

```go
// GetUserTemperature fetches temperature status for a specific user ID.
// Unlike GetStatus which uses the authenticated user, this allows querying any user.
func (c *Client) GetUserTemperature(ctx context.Context, userID string) (*TempStatus, error) {
	if err := c.ensureToken(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%s/temperature", userID)
	var res TempStatus
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/client/... -v -run GetUserTemperature`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/client/eightsleep.go internal/client/eightsleep_test.go
git commit -m "feat(client): add GetUserTemperature for any user ID"
```

---

### Task 2.3: Add SetUserTemperature Method for Any User

**Files:**
- Modify: `internal/client/eightsleep.go`
- Modify: `internal/client/eightsleep_test.go`

**Step 1: Write the failing test**

Add to `internal/client/eightsleep_test.go`:

```go
func TestClient_SetUserTemperature(t *testing.T) {
	var receivedBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/users/other-user-id/temperature" && r.Method == http.MethodPut {
			json.NewDecoder(r.Body).Decode(&receivedBody)
			w.WriteHeader(200)
			return
		}
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(404)
	}))
	defer srv.Close()

	cl := New("test@test.com", "pass", "", "", "")
	cl.BaseURL = srv.URL
	cl.token = "fake-token"
	cl.tokenExp = time.Now().Add(time.Hour)

	err := cl.SetUserTemperature(context.Background(), "other-user-id", -30)
	require.NoError(t, err)
	assert.Equal(t, float64(-30), receivedBody["currentLevel"])
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/client/... -v -run SetUserTemperature`
Expected: FAIL - cl.SetUserTemperature undefined

**Step 3: Write minimal implementation**

Add to `internal/client/eightsleep.go`:

```go
// SetUserTemperature sets temperature for a specific user ID.
func (c *Client) SetUserTemperature(ctx context.Context, userID string, level int) error {
	if err := c.ensureToken(ctx); err != nil {
		return err
	}
	if level < -100 || level > 100 {
		return fmt.Errorf("level must be between -100 and 100")
	}
	path := fmt.Sprintf("/users/%s/temperature", userID)
	body := map[string]int{"currentLevel": level}
	return c.do(ctx, http.MethodPut, path, nil, body, nil)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/client/... -v -run SetUserTemperature`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/client/eightsleep.go internal/client/eightsleep_test.go
git commit -m "feat(client): add SetUserTemperature for any user ID"
```

---

### Task 2.4: Add TurnOnUser and TurnOffUser Methods

**Files:**
- Modify: `internal/client/eightsleep.go`
- Modify: `internal/client/eightsleep_test.go`

**Step 1: Write the failing test**

Add to `internal/client/eightsleep_test.go`:

```go
func TestClient_TurnOnUser(t *testing.T) {
	var receivedBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/users/other-user-id/temperature" && r.Method == http.MethodPut {
			json.NewDecoder(r.Body).Decode(&receivedBody)
			w.WriteHeader(200)
			return
		}
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(404)
	}))
	defer srv.Close()

	cl := New("test@test.com", "pass", "", "", "")
	cl.BaseURL = srv.URL
	cl.token = "fake-token"
	cl.tokenExp = time.Now().Add(time.Hour)

	err := cl.TurnOnUser(context.Background(), "other-user-id")
	require.NoError(t, err)

	state, ok := receivedBody["currentState"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "smart", state["type"])
}

func TestClient_TurnOffUser(t *testing.T) {
	var receivedBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/users/other-user-id/temperature" && r.Method == http.MethodPut {
			json.NewDecoder(r.Body).Decode(&receivedBody)
			w.WriteHeader(200)
			return
		}
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(404)
	}))
	defer srv.Close()

	cl := New("test@test.com", "pass", "", "", "")
	cl.BaseURL = srv.URL
	cl.token = "fake-token"
	cl.tokenExp = time.Now().Add(time.Hour)

	err := cl.TurnOffUser(context.Background(), "other-user-id")
	require.NoError(t, err)

	state, ok := receivedBody["currentState"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "off", state["type"])
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/client/... -v -run "TurnOnUser|TurnOffUser"`
Expected: FAIL - undefined methods

**Step 3: Write minimal implementation**

Add to `internal/client/eightsleep.go`:

```go
// TurnOnUser turns on the pod for a specific user ID.
func (c *Client) TurnOnUser(ctx context.Context, userID string) error {
	if err := c.ensureToken(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/temperature", userID)
	body := map[string]any{
		"currentState": map[string]string{
			"type": "smart",
		},
	}
	return c.do(ctx, http.MethodPut, path, nil, body, nil)
}

// TurnOffUser turns off the pod for a specific user ID.
func (c *Client) TurnOffUser(ctx context.Context, userID string) error {
	if err := c.ensureToken(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/temperature", userID)
	body := map[string]any{
		"currentState": map[string]string{
			"type": "off",
		},
	}
	return c.do(ctx, http.MethodPut, path, nil, body, nil)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/client/... -v -run "TurnOnUser|TurnOffUser"`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/client/eightsleep.go internal/client/eightsleep_test.go
git commit -m "feat(client): add TurnOnUser and TurnOffUser methods"
```

---

## Phase 3: State Manager

### Task 3.1: Create StateProvider Interface

**Files:**
- Create: `internal/state/provider.go`
- Create: `internal/state/provider_test.go`

**Step 1: Write the failing test**

```go
// internal/state/provider_test.go
package state

import (
	"context"
	"testing"

	"github.com/steipete/eightctl/internal/model"
	"github.com/stretchr/testify/assert"
)

// Verify interface compliance at compile time
var _ StateProvider = (*mockProvider)(nil)

type mockProvider struct {
	state *model.DeviceState
	err   error
}

func (m *mockProvider) GetState(ctx context.Context) (*model.DeviceState, error) {
	return m.state, m.err
}

func (m *mockProvider) SetTemperature(ctx context.Context, side model.Side, level int) error {
	return m.err
}

func (m *mockProvider) TurnOn(ctx context.Context, side model.Side) error {
	return m.err
}

func (m *mockProvider) TurnOff(ctx context.Context, side model.Side) error {
	return m.err
}

func TestMockProvider_ImplementsInterface(t *testing.T) {
	// If this compiles, the test passes
	var p StateProvider = &mockProvider{}
	assert.NotNil(t, p)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/state/... -v`
Expected: FAIL - package does not exist

**Step 3: Write minimal implementation**

```go
// internal/state/provider.go
package state

import (
	"context"

	"github.com/steipete/eightctl/internal/model"
)

// StateProvider defines the interface for accessing and controlling Eight Sleep state.
type StateProvider interface {
	// GetState returns the current device state including both sides.
	GetState(ctx context.Context) (*model.DeviceState, error)

	// SetTemperature sets the temperature level for a specific side.
	SetTemperature(ctx context.Context, side model.Side, level int) error

	// TurnOn activates the pod for a specific side.
	TurnOn(ctx context.Context, side model.Side) error

	// TurnOff deactivates the pod for a specific side.
	TurnOff(ctx context.Context, side model.Side) error
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/state/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/state/provider.go internal/state/provider_test.go
git commit -m "feat(state): add StateProvider interface"
```

---

### Task 3.2: Create Observer Interface

**Files:**
- Modify: `internal/state/provider.go`
- Modify: `internal/state/provider_test.go`

**Step 1: Write the failing test**

Add to `internal/state/provider_test.go`:

```go
var _ Observer = (*mockObserver)(nil)

type mockObserver struct {
	stateChanges    []stateChange
	presenceChanges []presenceChange
}

type stateChange struct {
	old, new *model.DeviceState
}

type presenceChange struct {
	side    model.Side
	present bool
}

func (m *mockObserver) OnStateChange(old, new *model.DeviceState) {
	m.stateChanges = append(m.stateChanges, stateChange{old, new})
}

func (m *mockObserver) OnPresenceChange(side model.Side, present bool) {
	m.presenceChanges = append(m.presenceChanges, presenceChange{side, present})
}

func TestMockObserver_ImplementsInterface(t *testing.T) {
	var o Observer = &mockObserver{}
	assert.NotNil(t, o)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/state/... -v -run Observer`
Expected: FAIL - undefined: Observer

**Step 3: Write minimal implementation**

Add to `internal/state/provider.go`:

```go
// Observer receives notifications about state changes.
type Observer interface {
	// OnStateChange is called when the device state changes.
	OnStateChange(old, new *model.DeviceState)

	// OnPresenceChange is called when presence status changes for a side.
	OnPresenceChange(side model.Side, present bool)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/state/... -v -run Observer`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/state/provider.go internal/state/provider_test.go
git commit -m "feat(state): add Observer interface"
```

---

### Task 3.3: Create StateManager with Caching

**Files:**
- Create: `internal/state/manager.go`
- Create: `internal/state/manager_test.go`

**Step 1: Write the failing test**

```go
// internal/state/manager_test.go
package state

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/steipete/eightctl/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockClient simulates the Eight Sleep API client
type mockClient struct {
	deviceInfo      *DeviceInfo
	userTemps       map[string]*UserTemperature
	fetchCount      atomic.Int32
	setTempCalls    []setTempCall
	turnOnCalls     []string
	turnOffCalls    []string
}

type setTempCall struct {
	userID string
	level  int
}

type DeviceInfo struct {
	ID              string
	LeftUserID      string
	RightUserID     string
	RoomTemperature float64
}

type UserTemperature struct {
	CurrentLevel int
	StateType    string
}

func (m *mockClient) GetDeviceWithUsers(ctx context.Context) (*DeviceInfo, error) {
	m.fetchCount.Add(1)
	return m.deviceInfo, nil
}

func (m *mockClient) GetUserTemperature(ctx context.Context, userID string) (*UserTemperature, error) {
	return m.userTemps[userID], nil
}

func (m *mockClient) SetUserTemperature(ctx context.Context, userID string, level int) error {
	m.setTempCalls = append(m.setTempCalls, setTempCall{userID, level})
	return nil
}

func (m *mockClient) TurnOnUser(ctx context.Context, userID string) error {
	m.turnOnCalls = append(m.turnOnCalls, userID)
	return nil
}

func (m *mockClient) TurnOffUser(ctx context.Context, userID string) error {
	m.turnOffCalls = append(m.turnOffCalls, userID)
	return nil
}

func TestStateManager_GetState_CachesResult(t *testing.T) {
	mock := &mockClient{
		deviceInfo: &DeviceInfo{
			ID:              "device-1",
			LeftUserID:      "left-user",
			RightUserID:     "right-user",
			RoomTemperature: 70.0,
		},
		userTemps: map[string]*UserTemperature{
			"left-user":  {CurrentLevel: -20, StateType: "smart"},
			"right-user": {CurrentLevel: 10, StateType: "off"},
		},
	}

	mgr := NewManager(mock, WithCacheTTL(5*time.Second))
	ctx := context.Background()

	// First call fetches from API
	state1, err := mgr.GetState(ctx)
	require.NoError(t, err)
	assert.Equal(t, "device-1", state1.ID)
	assert.Equal(t, int32(1), mock.fetchCount.Load())

	// Second call uses cache
	state2, err := mgr.GetState(ctx)
	require.NoError(t, err)
	assert.Equal(t, state1, state2)
	assert.Equal(t, int32(1), mock.fetchCount.Load()) // Still 1
}

func TestStateManager_SetTemperature(t *testing.T) {
	mock := &mockClient{
		deviceInfo: &DeviceInfo{
			ID:          "device-1",
			LeftUserID:  "left-user",
			RightUserID: "right-user",
		},
		userTemps: map[string]*UserTemperature{
			"left-user":  {CurrentLevel: 0, StateType: "off"},
			"right-user": {CurrentLevel: 0, StateType: "off"},
		},
	}

	mgr := NewManager(mock)
	ctx := context.Background()

	// Must fetch state first to know user IDs
	_, err := mgr.GetState(ctx)
	require.NoError(t, err)

	err = mgr.SetTemperature(ctx, model.Left, -30)
	require.NoError(t, err)

	require.Len(t, mock.setTempCalls, 1)
	assert.Equal(t, "left-user", mock.setTempCalls[0].userID)
	assert.Equal(t, -30, mock.setTempCalls[0].level)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/state/... -v -run StateManager`
Expected: FAIL - undefined: NewManager

**Step 3: Write minimal implementation**

```go
// internal/state/manager.go
package state

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/steipete/eightctl/internal/model"
)

// ClientInterface defines the methods needed from the Eight Sleep client.
type ClientInterface interface {
	GetDeviceWithUsers(ctx context.Context) (*DeviceInfo, error)
	GetUserTemperature(ctx context.Context, userID string) (*UserTemperature, error)
	SetUserTemperature(ctx context.Context, userID string, level int) error
	TurnOnUser(ctx context.Context, userID string) error
	TurnOffUser(ctx context.Context, userID string) error
}

// DeviceInfo contains device data with user assignments.
type DeviceInfo struct {
	ID              string
	LeftUserID      string
	RightUserID     string
	RoomTemperature float64
	HasWater        bool
	IsPriming       bool
	NeedsPriming    bool
}

// UserTemperature contains temperature state for a user.
type UserTemperature struct {
	CurrentLevel int
	StateType    string
}

// Manager provides cached access to Eight Sleep state.
type Manager struct {
	client   ClientInterface
	cacheTTL time.Duration

	mu          sync.RWMutex
	cache       *model.DeviceState
	cacheTime   time.Time
	leftUserID  string
	rightUserID string

	observers []Observer
}

// Option configures the Manager.
type Option func(*Manager)

// WithCacheTTL sets the cache time-to-live.
func WithCacheTTL(ttl time.Duration) Option {
	return func(m *Manager) {
		m.cacheTTL = ttl
	}
}

// NewManager creates a new state manager.
func NewManager(client ClientInterface, opts ...Option) *Manager {
	m := &Manager{
		client:   client,
		cacheTTL: 60 * time.Second, // Default 60s
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// GetState returns the current device state, using cache if valid.
func (m *Manager) GetState(ctx context.Context) (*model.DeviceState, error) {
	m.mu.RLock()
	if m.cache != nil && time.Since(m.cacheTime) < m.cacheTTL {
		defer m.mu.RUnlock()
		return m.cache, nil
	}
	m.mu.RUnlock()

	return m.fetchState(ctx)
}

func (m *Manager) fetchState(ctx context.Context) (*model.DeviceState, error) {
	device, err := m.client.GetDeviceWithUsers(ctx)
	if err != nil {
		return nil, err
	}

	state := &model.DeviceState{
		ID:              device.ID,
		RoomTemperature: device.RoomTemperature,
		HasWater:        device.HasWater,
		IsPriming:       device.IsPriming,
		NeedsPriming:    device.NeedsPriming,
	}

	// Fetch left user state
	if device.LeftUserID != "" {
		leftTemp, err := m.client.GetUserTemperature(ctx, device.LeftUserID)
		if err == nil {
			state.LeftUser = &model.UserState{
				ID:          device.LeftUserID,
				Side:        model.Left,
				TargetLevel: leftTemp.CurrentLevel,
				State:       model.ParsePowerState(leftTemp.StateType),
			}
		}
	}

	// Fetch right user state
	if device.RightUserID != "" {
		rightTemp, err := m.client.GetUserTemperature(ctx, device.RightUserID)
		if err == nil {
			state.RightUser = &model.UserState{
				ID:          device.RightUserID,
				Side:        model.Right,
				TargetLevel: rightTemp.CurrentLevel,
				State:       model.ParsePowerState(rightTemp.StateType),
			}
		}
	}

	m.mu.Lock()
	m.cache = state
	m.cacheTime = time.Now()
	m.leftUserID = device.LeftUserID
	m.rightUserID = device.RightUserID
	m.mu.Unlock()

	return state, nil
}

func (m *Manager) getUserID(side model.Side) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	switch side {
	case model.Left:
		if m.leftUserID == "" {
			return "", errors.New("left user ID not known; call GetState first")
		}
		return m.leftUserID, nil
	case model.Right:
		if m.rightUserID == "" {
			return "", errors.New("right user ID not known; call GetState first")
		}
		return m.rightUserID, nil
	default:
		return "", errors.New("invalid side")
	}
}

// SetTemperature sets the temperature for a side.
func (m *Manager) SetTemperature(ctx context.Context, side model.Side, level int) error {
	userID, err := m.getUserID(side)
	if err != nil {
		return err
	}
	err = m.client.SetUserTemperature(ctx, userID, level)
	if err == nil {
		m.invalidateCache()
	}
	return err
}

// TurnOn activates the pod for a side.
func (m *Manager) TurnOn(ctx context.Context, side model.Side) error {
	userID, err := m.getUserID(side)
	if err != nil {
		return err
	}
	err = m.client.TurnOnUser(ctx, userID)
	if err == nil {
		m.invalidateCache()
	}
	return err
}

// TurnOff deactivates the pod for a side.
func (m *Manager) TurnOff(ctx context.Context, side model.Side) error {
	userID, err := m.getUserID(side)
	if err != nil {
		return err
	}
	err = m.client.TurnOffUser(ctx, userID)
	if err == nil {
		m.invalidateCache()
	}
	return err
}

func (m *Manager) invalidateCache() {
	m.mu.Lock()
	m.cacheTime = time.Time{}
	m.mu.Unlock()
}

// Subscribe adds an observer for state changes.
func (m *Manager) Subscribe(o Observer) {
	m.mu.Lock()
	m.observers = append(m.observers, o)
	m.mu.Unlock()
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/state/... -v -run StateManager`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/state/manager.go internal/state/manager_test.go
git commit -m "feat(state): add StateManager with caching"
```

---

## Phase 4: CLI Commands (Dual-Side Support)

### Task 4.1: Add `eightctl sides` Command

**Files:**
- Create: `internal/cmd/sides.go`
- Modify: `internal/cmd/root.go` (add to init)

**Step 1: Write the failing test**

Create a simple integration test:

```go
// internal/cmd/sides_test.go
package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSidesCmd_Exists(t *testing.T) {
	// Verify the command is registered
	cmd := rootCmd
	found := false
	for _, c := range cmd.Commands() {
		if c.Name() == "sides" {
			found = true
			break
		}
	}
	assert.True(t, found, "sides command should be registered")
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cmd/... -v -run SidesCmd`
Expected: FAIL - sides command not found

**Step 3: Write minimal implementation**

```go
// internal/cmd/sides.go
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/steipete/eightctl/internal/client"
	"github.com/steipete/eightctl/internal/output"
)

var sidesCmd = &cobra.Command{
	Use:   "sides",
	Short: "Show which users are assigned to which side of the bed",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		cl := client.New(
			viper.GetString("email"),
			viper.GetString("password"),
			viper.GetString("user_id"),
			viper.GetString("client_id"),
			viper.GetString("client_secret"),
		)

		device, err := cl.Device().GetWithUsers(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get device info: %w", err)
		}

		rows := []map[string]any{
			{"side": "left", "user_id": device.LeftUserID},
			{"side": "right", "user_id": device.RightUserID},
		}

		fields := viper.GetStringSlice("fields")
		rows = output.FilterFields(rows, fields)
		headers := fields
		if len(headers) == 0 {
			headers = []string{"side", "user_id"}
		}
		return output.Print(output.Format(viper.GetString("output")), headers, rows)
	},
}

func init() {
	rootCmd.AddCommand(sidesCmd)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cmd/... -v -run SidesCmd`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/cmd/sides.go internal/cmd/sides_test.go
git commit -m "feat(cmd): add sides command to show user-side mappings"
```

---

### Task 4.2: Add `--side` Flag to Status Command

**Files:**
- Modify: `internal/cmd/status.go`

**Step 1: Design**

Add a `--side` flag that accepts "left" or "right". When specified, show status for that side. When omitted, show status for the authenticated user's side (current behavior).

**Step 2: Implementation**

Modify `internal/cmd/status.go`:

```go
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/steipete/eightctl/internal/client"
	"github.com/steipete/eightctl/internal/model"
	"github.com/steipete/eightctl/internal/output"
)

var statusSide string

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show device status",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		cl := client.New(
			viper.GetString("email"),
			viper.GetString("password"),
			viper.GetString("user_id"),
			viper.GetString("client_id"),
			viper.GetString("client_secret"),
		)

		var st *client.TempStatus
		var err error

		if statusSide != "" {
			// Fetch specific side
			side, parseErr := model.ParseSide(statusSide)
			if parseErr != nil {
				return parseErr
			}

			device, devErr := cl.Device().GetWithUsers(context.Background())
			if devErr != nil {
				return fmt.Errorf("failed to get device info: %w", devErr)
			}

			var userID string
			switch side {
			case model.Left:
				userID = device.LeftUserID
			case model.Right:
				userID = device.RightUserID
			}

			if userID == "" {
				return fmt.Errorf("no user assigned to %s side", side)
			}

			st, err = cl.GetUserTemperature(context.Background(), userID)
		} else {
			// Default: authenticated user's side
			st, err = cl.GetStatus(context.Background())
		}

		if err != nil {
			return err
		}

		row := map[string]any{"mode": st.CurrentState.Type, "level": st.CurrentLevel}
		if statusSide != "" {
			row["side"] = statusSide
		}

		fields := viper.GetStringSlice("fields")
		rows := output.FilterFields([]map[string]any{row}, fields)
		headers := fields
		if len(headers) == 0 {
			if statusSide != "" {
				headers = []string{"side", "mode", "level"}
			} else {
				headers = []string{"mode", "level"}
			}
		}
		return output.Print(output.Format(viper.GetString("output")), headers, rows)
	},
}

func init() {
	statusCmd.Flags().StringVar(&statusSide, "side", "", "Show status for specific side (left or right)")
	rootCmd.AddCommand(statusCmd)
}
```

**Step 3: Test manually**

```bash
go build -o bin/eightctl ./cmd/eightctl
./bin/eightctl status --help  # Should show --side flag
```

**Step 4: Commit**

```bash
git add internal/cmd/status.go
git commit -m "feat(cmd): add --side flag to status command"
```

---

### Task 4.3: Add `--side` Flag to Temp Command

**Files:**
- Modify: `internal/cmd/temp.go`

**Step 1: Read current implementation**

Read `internal/cmd/temp.go` to understand current structure.

**Step 2: Modify to add side flag**

Similar pattern to status command - add `--side` flag and use `SetUserTemperature` when specified.

**Step 3: Test manually**

```bash
./bin/eightctl temp --help  # Should show --side flag
```

**Step 4: Commit**

```bash
git add internal/cmd/temp.go
git commit -m "feat(cmd): add --side flag to temp command"
```

---

### Task 4.4: Add `--side` Flag to On/Off Commands

**Files:**
- Modify: `internal/cmd/on.go`
- Modify: `internal/cmd/off.go`

**Step 1: Modify on.go**

Add `--side` flag, use `TurnOnUser` when specified.

**Step 2: Modify off.go**

Add `--side` flag, use `TurnOffUser` when specified.

**Step 3: Test manually**

```bash
./bin/eightctl on --help   # Should show --side flag
./bin/eightctl off --help  # Should show --side flag
```

**Step 4: Commit**

```bash
git add internal/cmd/on.go internal/cmd/off.go
git commit -m "feat(cmd): add --side flag to on/off commands"
```

---

## Phase 5: Adapter Interface

### Task 5.1: Create Adapter Interface

**Files:**
- Create: `internal/adapter/adapter.go`
- Create: `internal/adapter/adapter_test.go`

**Step 1: Write the failing test**

```go
// internal/adapter/adapter_test.go
package adapter

import (
	"context"
	"testing"

	"github.com/steipete/eightctl/internal/model"
	"github.com/steipete/eightctl/internal/state"
	"github.com/stretchr/testify/assert"
)

// Verify interface compliance
var _ Adapter = (*mockAdapter)(nil)

type mockAdapter struct {
	started bool
	stopped bool
}

func (m *mockAdapter) Start(ctx context.Context, provider state.StateProvider) error {
	m.started = true
	return nil
}

func (m *mockAdapter) HandleCommand(cmd Command) error {
	return nil
}

func (m *mockAdapter) Stop() error {
	m.stopped = true
	return nil
}

func TestMockAdapter_ImplementsInterface(t *testing.T) {
	var a Adapter = &mockAdapter{}
	assert.NotNil(t, a)
}

func TestCommand_String(t *testing.T) {
	cmd := Command{
		Side:   model.Left,
		Action: ActionSetTemp,
		Value:  -20,
	}
	assert.Contains(t, cmd.String(), "left")
	assert.Contains(t, cmd.String(), "set_temp")
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/adapter/... -v`
Expected: FAIL - package does not exist

**Step 3: Write minimal implementation**

```go
// internal/adapter/adapter.go
package adapter

import (
	"context"
	"fmt"

	"github.com/steipete/eightctl/internal/model"
	"github.com/steipete/eightctl/internal/state"
)

// Adapter defines the interface for platform-specific integrations.
type Adapter interface {
	// Start begins the adapter's main loop.
	Start(ctx context.Context, provider state.StateProvider) error

	// HandleCommand processes an incoming command.
	HandleCommand(cmd Command) error

	// Stop gracefully shuts down the adapter.
	Stop() error
}

// Action represents a command action type.
type Action int

const (
	ActionSetTemp Action = iota
	ActionTurnOn
	ActionTurnOff
)

func (a Action) String() string {
	switch a {
	case ActionSetTemp:
		return "set_temp"
	case ActionTurnOn:
		return "turn_on"
	case ActionTurnOff:
		return "turn_off"
	default:
		return "unknown"
	}
}

// Command represents a control command from an adapter.
type Command struct {
	Side   model.Side
	Action Action
	Value  any
}

func (c Command) String() string {
	return fmt.Sprintf("Command{side=%s, action=%s, value=%v}", c.Side, c.Action, c.Value)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/adapter/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/adapter/adapter.go internal/adapter/adapter_test.go
git commit -m "feat(adapter): add Adapter interface and Command type"
```

---

## Phase 6: MQTT Adapter (Home Assistant)

### Task 6.1: Add MQTT Dependencies

**Files:**
- Modify: `go.mod`

**Step 1: Add Paho MQTT client**

```bash
go get github.com/eclipse/paho.mqtt.golang
```

**Step 2: Commit**

```bash
git add go.mod go.sum
git commit -m "chore: add paho mqtt dependency"
```

---

### Task 6.2: Create MQTT Discovery Config Generator

**Files:**
- Create: `internal/adapter/mqtt/discovery.go`
- Create: `internal/adapter/mqtt/discovery_test.go`

**Step 1: Write the failing test**

```go
// internal/adapter/mqtt/discovery_test.go
package mqtt

import (
	"encoding/json"
	"testing"

	"github.com/steipete/eightctl/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoveryConfig_BinarySensorPresence(t *testing.T) {
	cfg := NewDiscoveryConfig("homeassistant", "eight_sleep")

	topic, payload := cfg.BinarySensorPresence(model.Left)

	assert.Equal(t, "homeassistant/binary_sensor/eight_sleep/left_presence/config", topic)

	var config map[string]any
	err := json.Unmarshal(payload, &config)
	require.NoError(t, err)

	assert.Equal(t, "Eight Sleep Left Presence", config["name"])
	assert.Equal(t, "eight_sleep/left/presence/state", config["state_topic"])
	assert.Equal(t, "occupancy", config["device_class"])
}

func TestDiscoveryConfig_SensorTemperature(t *testing.T) {
	cfg := NewDiscoveryConfig("homeassistant", "eight_sleep")

	topic, payload := cfg.SensorBedTemperature(model.Right)

	assert.Equal(t, "homeassistant/sensor/eight_sleep/right_bed_temp/config", topic)

	var config map[string]any
	err := json.Unmarshal(payload, &config)
	require.NoError(t, err)

	assert.Equal(t, "Eight Sleep Right Bed Temperature", config["name"])
	assert.Equal(t, "temperature", config["device_class"])
	assert.Equal(t, "°F", config["unit_of_measurement"])
}

func TestDiscoveryConfig_Switch(t *testing.T) {
	cfg := NewDiscoveryConfig("homeassistant", "eight_sleep")

	topic, payload := cfg.Switch(model.Left)

	assert.Equal(t, "homeassistant/switch/eight_sleep/left_power/config", topic)

	var config map[string]any
	err := json.Unmarshal(payload, &config)
	require.NoError(t, err)

	assert.Equal(t, "eight_sleep/left/power/state", config["state_topic"])
	assert.Equal(t, "eight_sleep/left/power/set", config["command_topic"])
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/adapter/mqtt/... -v`
Expected: FAIL - package does not exist

**Step 3: Write minimal implementation**

```go
// internal/adapter/mqtt/discovery.go
package mqtt

import (
	"encoding/json"
	"fmt"

	"github.com/steipete/eightctl/internal/model"
)

// DiscoveryConfig generates Home Assistant MQTT Discovery payloads.
type DiscoveryConfig struct {
	DiscoveryPrefix string
	BaseTopic       string
}

// NewDiscoveryConfig creates a discovery config generator.
func NewDiscoveryConfig(discoveryPrefix, baseTopic string) *DiscoveryConfig {
	return &DiscoveryConfig{
		DiscoveryPrefix: discoveryPrefix,
		BaseTopic:       baseTopic,
	}
}

func (d *DiscoveryConfig) deviceInfo() map[string]any {
	return map[string]any{
		"identifiers":  []string{"eight_sleep_pod"},
		"name":         "Eight Sleep Pod",
		"manufacturer": "Eight Sleep",
		"model":        "Pod",
	}
}

// BinarySensorPresence returns topic and payload for presence sensor discovery.
func (d *DiscoveryConfig) BinarySensorPresence(side model.Side) (string, []byte) {
	sideName := side.String()
	topic := fmt.Sprintf("%s/binary_sensor/%s/%s_presence/config",
		d.DiscoveryPrefix, d.BaseTopic, sideName)

	config := map[string]any{
		"name":         fmt.Sprintf("Eight Sleep %s Presence", capitalize(sideName)),
		"unique_id":    fmt.Sprintf("eight_sleep_%s_presence", sideName),
		"state_topic":  fmt.Sprintf("%s/%s/presence/state", d.BaseTopic, sideName),
		"device_class": "occupancy",
		"payload_on":   "ON",
		"payload_off":  "OFF",
		"device":       d.deviceInfo(),
	}

	payload, _ := json.Marshal(config)
	return topic, payload
}

// SensorBedTemperature returns topic and payload for bed temperature sensor.
func (d *DiscoveryConfig) SensorBedTemperature(side model.Side) (string, []byte) {
	sideName := side.String()
	topic := fmt.Sprintf("%s/sensor/%s/%s_bed_temp/config",
		d.DiscoveryPrefix, d.BaseTopic, sideName)

	config := map[string]any{
		"name":                fmt.Sprintf("Eight Sleep %s Bed Temperature", capitalize(sideName)),
		"unique_id":           fmt.Sprintf("eight_sleep_%s_bed_temp", sideName),
		"state_topic":         fmt.Sprintf("%s/%s/bed_temp/state", d.BaseTopic, sideName),
		"device_class":        "temperature",
		"unit_of_measurement": "°F",
		"device":              d.deviceInfo(),
	}

	payload, _ := json.Marshal(config)
	return topic, payload
}

// SensorHeartRate returns topic and payload for heart rate sensor.
func (d *DiscoveryConfig) SensorHeartRate(side model.Side) (string, []byte) {
	sideName := side.String()
	topic := fmt.Sprintf("%s/sensor/%s/%s_heart_rate/config",
		d.DiscoveryPrefix, d.BaseTopic, sideName)

	config := map[string]any{
		"name":                fmt.Sprintf("Eight Sleep %s Heart Rate", capitalize(sideName)),
		"unique_id":           fmt.Sprintf("eight_sleep_%s_heart_rate", sideName),
		"state_topic":         fmt.Sprintf("%s/%s/heart_rate/state", d.BaseTopic, sideName),
		"unit_of_measurement": "bpm",
		"icon":                "mdi:heart-pulse",
		"device":              d.deviceInfo(),
	}

	payload, _ := json.Marshal(config)
	return topic, payload
}

// SensorSleepStage returns topic and payload for sleep stage sensor.
func (d *DiscoveryConfig) SensorSleepStage(side model.Side) (string, []byte) {
	sideName := side.String()
	topic := fmt.Sprintf("%s/sensor/%s/%s_sleep_stage/config",
		d.DiscoveryPrefix, d.BaseTopic, sideName)

	config := map[string]any{
		"name":        fmt.Sprintf("Eight Sleep %s Sleep Stage", capitalize(sideName)),
		"unique_id":   fmt.Sprintf("eight_sleep_%s_sleep_stage", sideName),
		"state_topic": fmt.Sprintf("%s/%s/sleep_stage/state", d.BaseTopic, sideName),
		"icon":        "mdi:sleep",
		"device":      d.deviceInfo(),
	}

	payload, _ := json.Marshal(config)
	return topic, payload
}

// Switch returns topic and payload for power switch discovery.
func (d *DiscoveryConfig) Switch(side model.Side) (string, []byte) {
	sideName := side.String()
	topic := fmt.Sprintf("%s/switch/%s/%s_power/config",
		d.DiscoveryPrefix, d.BaseTopic, sideName)

	config := map[string]any{
		"name":          fmt.Sprintf("Eight Sleep %s Power", capitalize(sideName)),
		"unique_id":     fmt.Sprintf("eight_sleep_%s_power", sideName),
		"state_topic":   fmt.Sprintf("%s/%s/power/state", d.BaseTopic, sideName),
		"command_topic": fmt.Sprintf("%s/%s/power/set", d.BaseTopic, sideName),
		"payload_on":    "ON",
		"payload_off":   "OFF",
		"device":        d.deviceInfo(),
	}

	payload, _ := json.Marshal(config)
	return topic, payload
}

// NumberTemperature returns topic and payload for temperature control number entity.
func (d *DiscoveryConfig) NumberTemperature(side model.Side) (string, []byte) {
	sideName := side.String()
	topic := fmt.Sprintf("%s/number/%s/%s_temperature/config",
		d.DiscoveryPrefix, d.BaseTopic, sideName)

	config := map[string]any{
		"name":          fmt.Sprintf("Eight Sleep %s Temperature", capitalize(sideName)),
		"unique_id":     fmt.Sprintf("eight_sleep_%s_temperature", sideName),
		"state_topic":   fmt.Sprintf("%s/%s/temperature/state", d.BaseTopic, sideName),
		"command_topic": fmt.Sprintf("%s/%s/temperature/set", d.BaseTopic, sideName),
		"min":           -100,
		"max":           100,
		"step":          5,
		"icon":          "mdi:thermometer",
		"device":        d.deviceInfo(),
	}

	payload, _ := json.Marshal(config)
	return topic, payload
}

// SensorRoomTemperature returns topic and payload for room temperature sensor.
func (d *DiscoveryConfig) SensorRoomTemperature() (string, []byte) {
	topic := fmt.Sprintf("%s/sensor/%s/room_temp/config",
		d.DiscoveryPrefix, d.BaseTopic)

	config := map[string]any{
		"name":                "Eight Sleep Room Temperature",
		"unique_id":           "eight_sleep_room_temp",
		"state_topic":         fmt.Sprintf("%s/room_temp/state", d.BaseTopic),
		"device_class":        "temperature",
		"unit_of_measurement": "°F",
		"device":              d.deviceInfo(),
	}

	payload, _ := json.Marshal(config)
	return topic, payload
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/adapter/mqtt/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/adapter/mqtt/discovery.go internal/adapter/mqtt/discovery_test.go
git commit -m "feat(mqtt): add Home Assistant discovery config generator"
```

---

### Task 6.3: Create MQTT Adapter

**Files:**
- Create: `internal/adapter/mqtt/adapter.go`
- Create: `internal/adapter/mqtt/adapter_test.go`

This task creates the main MQTT adapter that:
- Connects to MQTT broker
- Publishes discovery configs on startup
- Polls state and publishes updates
- Subscribes to command topics

**Implementation details in subsequent steps...**

---

### Task 6.4: Add `eightctl mqtt` Command

**Files:**
- Create: `internal/cmd/mqtt.go`

Creates the CLI command that starts the MQTT bridge.

---

## Phase 7: Hubitat Adapter

### Task 7.1: Create Hubitat HTTP Server

**Files:**
- Create: `internal/adapter/hubitat/server.go`
- Create: `internal/adapter/hubitat/server_test.go`

**Step 1: Write the failing test**

```go
// internal/adapter/hubitat/server_test.go
package hubitat

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/steipete/eightctl/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStateProvider struct {
	state *model.DeviceState
}

func (m *mockStateProvider) GetState(ctx context.Context) (*model.DeviceState, error) {
	return m.state, nil
}

func (m *mockStateProvider) SetTemperature(ctx context.Context, side model.Side, level int) error {
	return nil
}

func (m *mockStateProvider) TurnOn(ctx context.Context, side model.Side) error {
	return nil
}

func (m *mockStateProvider) TurnOff(ctx context.Context, side model.Side) error {
	return nil
}

func TestHubitatServer_GetStatus(t *testing.T) {
	provider := &mockStateProvider{
		state: &model.DeviceState{
			ID:              "device-1",
			RoomTemperature: 70.5,
			LeftUser: &model.UserState{
				Side:        model.Left,
				TargetLevel: -20,
				State:       model.PowerSmart,
			},
			RightUser: &model.UserState{
				Side:        model.Right,
				TargetLevel: 10,
				State:       model.PowerOff,
			},
		},
	}

	srv := NewServer(provider)

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, "device-1", resp["id"])
	assert.InDelta(t, 70.5, resp["room_temperature"], 0.01)
}

func TestHubitatServer_GetSideStatus(t *testing.T) {
	provider := &mockStateProvider{
		state: &model.DeviceState{
			LeftUser: &model.UserState{
				Side:        model.Left,
				TargetLevel: -20,
				State:       model.PowerSmart,
			},
		},
	}

	srv := NewServer(provider)

	req := httptest.NewRequest(http.MethodGet, "/left/status", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.UserState
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, model.Left, resp.Side)
	assert.Equal(t, -20, resp.TargetLevel)
}

func TestHubitatServer_SetTemperature(t *testing.T) {
	setCalls := []struct {
		side  model.Side
		level int
	}{}

	provider := &mockStateProvider{
		state: &model.DeviceState{
			LeftUser:  &model.UserState{Side: model.Left},
			RightUser: &model.UserState{Side: model.Right},
		},
	}

	srv := NewServer(provider)
	srv.OnSetTemperature = func(side model.Side, level int) error {
		setCalls = append(setCalls, struct {
			side  model.Side
			level int
		}{side, level})
		return nil
	}

	body := strings.NewReader(`{"level": -30}`)
	req := httptest.NewRequest(http.MethodPut, "/left/temperature", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	require.Len(t, setCalls, 1)
	assert.Equal(t, model.Left, setCalls[0].side)
	assert.Equal(t, -30, setCalls[0].level)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/adapter/hubitat/... -v`
Expected: FAIL - package does not exist

**Step 3: Write minimal implementation**

```go
// internal/adapter/hubitat/server.go
package hubitat

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/steipete/eightctl/internal/model"
	"github.com/steipete/eightctl/internal/state"
)

// Server provides HTTP endpoints for Hubitat drivers.
type Server struct {
	provider state.StateProvider
	mux      *http.ServeMux

	// Callbacks for control actions (set by adapter)
	OnSetTemperature func(side model.Side, level int) error
	OnTurnOn         func(side model.Side) error
	OnTurnOff        func(side model.Side) error
}

// NewServer creates a Hubitat HTTP server.
func NewServer(provider state.StateProvider) *Server {
	s := &Server{
		provider: provider,
		mux:      http.NewServeMux(),
	}
	s.mux.HandleFunc("GET /status", s.handleStatus)
	s.mux.HandleFunc("GET /{side}/status", s.handleSideStatus)
	s.mux.HandleFunc("PUT /{side}/temperature", s.handleSetTemperature)
	s.mux.HandleFunc("PUT /{side}/power", s.handleSetPower)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	state, err := s.provider.GetState(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

func (s *Server) handleSideStatus(w http.ResponseWriter, r *http.Request) {
	sideStr := r.PathValue("side")
	side, err := model.ParseSide(sideStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	state, err := s.provider.GetState(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := state.GetSide(side)
	if user == nil {
		http.Error(w, "side not configured", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (s *Server) handleSetTemperature(w http.ResponseWriter, r *http.Request) {
	sideStr := r.PathValue("side")
	side, err := model.ParseSide(sideStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req struct {
		Level int `json:"level"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if s.OnSetTemperature != nil {
		if err := s.OnSetTemperature(side, req.Level); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleSetPower(w http.ResponseWriter, r *http.Request) {
	sideStr := r.PathValue("side")
	side, err := model.ParseSide(sideStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req struct {
		On bool `json:"on"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.On && s.OnTurnOn != nil {
		if err := s.OnTurnOn(side); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if !req.On && s.OnTurnOff != nil {
		if err := s.OnTurnOff(side); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/adapter/hubitat/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/adapter/hubitat/server.go internal/adapter/hubitat/server_test.go
git commit -m "feat(hubitat): add HTTP server for Hubitat drivers"
```

---

### Task 7.2: Add `eightctl hubitat` Command

**Files:**
- Create: `internal/cmd/hubitat.go`

Creates the CLI command that starts the Hubitat HTTP bridge.

---

### Task 7.3: Create Hubitat Groovy Drivers

**Files:**
- Create: `drivers/hubitat/eight-sleep-pod.groovy`
- Create: `drivers/hubitat/eight-sleep-side.groovy`

Groovy drivers that Hubitat users install on their hubs.

---

## Phase 8: Documentation

### Task 8.1: Update CLI Reference

**Files:**
- Modify: `docs/cli-reference.md`

Add documentation for new commands: `sides`, `mqtt`, `hubitat`, and `--side` flag.

---

### Task 8.2: Create Home Assistant Setup Guide

**Files:**
- Create: `docs/home-assistant.md`

Document MQTT bridge setup and configuration.

---

### Task 8.3: Create Hubitat Setup Guide

**Files:**
- Create: `docs/hubitat.md`

Document driver installation and bridge setup.

---

## Summary

This plan has **35+ bite-sized tasks** across 8 phases:

1. **Domain Models** (4 tasks) - Types, UserState, DeviceState
2. **Client Enhancements** (4 tasks) - Dual-side methods
3. **State Manager** (3 tasks) - Interface, Observer, Manager
4. **CLI Commands** (4 tasks) - sides command, --side flags
5. **Adapter Interface** (1 task) - Common interface
6. **MQTT Adapter** (4 tasks) - Discovery, adapter, command
7. **Hubitat Adapter** (3 tasks) - Server, command, drivers
8. **Documentation** (3 tasks) - Docs updates

Each task is TDD: failing test → implementation → passing test → commit.
