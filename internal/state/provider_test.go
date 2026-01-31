package state

import (
	"context"
	"testing"

	"github.com/steipete/eightctl/internal/model"
)

// Compile-time interface satisfaction checks.
var (
	_ StateProvider = (*Manager)(nil)
	_ Observer      = (*mockObserver)(nil)
)

// mockObserver is a test implementation of Observer.
type mockObserver struct {
	stateChanges    []StateChange
	presenceChanges []PresenceChange
}

func (m *mockObserver) OnStateChange(change StateChange) {
	m.stateChanges = append(m.stateChanges, change)
}

func (m *mockObserver) OnPresenceChange(change PresenceChange) {
	m.presenceChanges = append(m.presenceChanges, change)
}

// mockStateProvider is a test implementation of StateProvider.
type mockStateProvider struct {
	state        *model.DeviceState
	err          error
	setTempCalls []setTempCall
	turnOnCalls  []model.Side
	turnOffCalls []model.Side
}

type setTempCall struct {
	side  model.Side
	level int
}

func (m *mockStateProvider) GetState(ctx context.Context) (*model.DeviceState, error) {
	return m.state, m.err
}

func (m *mockStateProvider) SetTemperature(ctx context.Context, side model.Side, level int) error {
	m.setTempCalls = append(m.setTempCalls, setTempCall{side: side, level: level})
	return m.err
}

func (m *mockStateProvider) TurnOn(ctx context.Context, side model.Side) error {
	m.turnOnCalls = append(m.turnOnCalls, side)
	return m.err
}

func (m *mockStateProvider) TurnOff(ctx context.Context, side model.Side) error {
	m.turnOffCalls = append(m.turnOffCalls, side)
	return m.err
}

func TestStateProviderInterface(t *testing.T) {
	// Verify that mockStateProvider satisfies StateProvider interface.
	var _ StateProvider = (*mockStateProvider)(nil)

	provider := &mockStateProvider{
		state: &model.DeviceState{ID: "dev-1"},
	}

	ctx := context.Background()

	// Test GetState
	state, err := provider.GetState(ctx)
	if err != nil {
		t.Fatalf("GetState error: %v", err)
	}
	if state.ID != "dev-1" {
		t.Errorf("expected device ID 'dev-1', got '%s'", state.ID)
	}

	// Test SetTemperature
	if err := provider.SetTemperature(ctx, model.Left, 50); err != nil {
		t.Fatalf("SetTemperature error: %v", err)
	}
	if len(provider.setTempCalls) != 1 {
		t.Errorf("expected 1 SetTemperature call, got %d", len(provider.setTempCalls))
	}
	if provider.setTempCalls[0].side != model.Left || provider.setTempCalls[0].level != 50 {
		t.Errorf("unexpected SetTemperature args: %+v", provider.setTempCalls[0])
	}

	// Test TurnOn
	if err := provider.TurnOn(ctx, model.Right); err != nil {
		t.Fatalf("TurnOn error: %v", err)
	}
	if len(provider.turnOnCalls) != 1 || provider.turnOnCalls[0] != model.Right {
		t.Errorf("unexpected TurnOn calls: %v", provider.turnOnCalls)
	}

	// Test TurnOff
	if err := provider.TurnOff(ctx, model.Left); err != nil {
		t.Fatalf("TurnOff error: %v", err)
	}
	if len(provider.turnOffCalls) != 1 || provider.turnOffCalls[0] != model.Left {
		t.Errorf("unexpected TurnOff calls: %v", provider.turnOffCalls)
	}
}

func TestObserverInterface(t *testing.T) {
	observer := &mockObserver{}

	// Test OnStateChange
	oldState := &model.DeviceState{ID: "dev-1"}
	newState := &model.DeviceState{ID: "dev-1", RoomTemperature: 21.5}
	observer.OnStateChange(StateChange{Old: oldState, New: newState})

	if len(observer.stateChanges) != 1 {
		t.Fatalf("expected 1 state change, got %d", len(observer.stateChanges))
	}
	if observer.stateChanges[0].Old != oldState || observer.stateChanges[0].New != newState {
		t.Errorf("unexpected state change: %+v", observer.stateChanges[0])
	}

	// Test OnPresenceChange
	user := &model.UserState{ID: "user-1", Side: model.Left}
	observer.OnPresenceChange(PresenceChange{Side: model.Left, Present: true, User: user})

	if len(observer.presenceChanges) != 1 {
		t.Fatalf("expected 1 presence change, got %d", len(observer.presenceChanges))
	}
	if observer.presenceChanges[0].Side != model.Left || !observer.presenceChanges[0].Present {
		t.Errorf("unexpected presence change: %+v", observer.presenceChanges[0])
	}
	if observer.presenceChanges[0].User != user {
		t.Errorf("unexpected user in presence change: %+v", observer.presenceChanges[0].User)
	}
}

func TestStateChangeFields(t *testing.T) {
	change := StateChange{
		Old: &model.DeviceState{ID: "old"},
		New: &model.DeviceState{ID: "new"},
	}

	if change.Old.ID != "old" {
		t.Errorf("expected old ID 'old', got '%s'", change.Old.ID)
	}
	if change.New.ID != "new" {
		t.Errorf("expected new ID 'new', got '%s'", change.New.ID)
	}
}

func TestPresenceChangeFields(t *testing.T) {
	user := &model.UserState{ID: "user-1"}
	change := PresenceChange{
		Side:    model.Right,
		Present: true,
		User:    user,
	}

	if change.Side != model.Right {
		t.Errorf("expected side Right, got %v", change.Side)
	}
	if !change.Present {
		t.Error("expected Present to be true")
	}
	if change.User.ID != "user-1" {
		t.Errorf("expected user ID 'user-1', got '%s'", change.User.ID)
	}
}
