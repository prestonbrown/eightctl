package state

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/steipete/eightctl/internal/client"
	"github.com/steipete/eightctl/internal/model"
)

// DefaultCacheTTL is the default cache time-to-live.
const DefaultCacheTTL = 30 * time.Second

// Compile-time check that Manager implements StateProvider.
var _ StateProvider = (*Manager)(nil)

// Manager implements StateProvider with caching and observer notifications.
type Manager struct {
	client   *client.Client
	deviceID string
	cacheTTL time.Duration

	mu          sync.RWMutex
	cachedState *model.DeviceState
	cacheExpiry time.Time
	observers   []Observer
}

// Option configures the Manager.
type Option func(*Manager)

// WithCacheTTL sets the cache TTL.
func WithCacheTTL(ttl time.Duration) Option {
	return func(m *Manager) {
		m.cacheTTL = ttl
	}
}

// NewManager creates a new state manager.
func NewManager(c *client.Client, deviceID string, opts ...Option) *Manager {
	m := &Manager{
		client:   c,
		deviceID: deviceID,
		cacheTTL: DefaultCacheTTL,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// AddObserver registers an observer for state changes.
func (m *Manager) AddObserver(o Observer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.observers = append(m.observers, o)
}

// RemoveObserver unregisters an observer.
func (m *Manager) RemoveObserver(o Observer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, obs := range m.observers {
		if obs == o {
			m.observers = append(m.observers[:i], m.observers[i+1:]...)
			return
		}
	}
}

// GetState returns the current device state, using cache if valid.
func (m *Manager) GetState(ctx context.Context) (*model.DeviceState, error) {
	m.mu.RLock()
	if m.cachedState != nil && time.Now().Before(m.cacheExpiry) {
		state := m.cachedState
		m.mu.RUnlock()
		return state, nil
	}
	m.mu.RUnlock()

	return m.refreshState(ctx)
}

// refreshState fetches fresh state from the API and updates cache.
func (m *Manager) refreshState(ctx context.Context) (*model.DeviceState, error) {
	// Fetch device info with user assignments
	device, err := m.client.Device().GetWithUsers(ctx)
	if err != nil {
		return nil, err
	}

	state := &model.DeviceState{
		ID:              device.ID,
		RoomTemperature: device.RoomTemperature,
		HasWater:        device.WaterLevel > 0,
		IsPriming:       device.IsPriming,
		NeedsPriming:    device.NeedsPriming,
	}

	// Fetch left user state if assigned
	if device.LeftUserID != "" {
		leftState, err := m.fetchUserState(ctx, device.LeftUserID, model.Left)
		if err == nil {
			state.LeftUser = leftState
		}
	}

	// Fetch right user state if assigned
	if device.RightUserID != "" {
		rightState, err := m.fetchUserState(ctx, device.RightUserID, model.Right)
		if err == nil {
			state.RightUser = rightState
		}
	}

	m.mu.Lock()
	oldState := m.cachedState
	m.cachedState = state
	m.cacheExpiry = time.Now().Add(m.cacheTTL)
	observers := make([]Observer, len(m.observers))
	copy(observers, m.observers)
	m.mu.Unlock()

	// Notify observers of state change
	if oldState != nil {
		m.notifyStateChange(observers, oldState, state)
	}

	return state, nil
}

// fetchUserState fetches temperature status for a user.
func (m *Manager) fetchUserState(ctx context.Context, userID string, side model.Side) (*model.UserState, error) {
	temp, err := m.client.GetUserTemperature(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &model.UserState{
		ID:          userID,
		Side:        side,
		TargetLevel: temp.CurrentLevel,
		State:       model.ParsePowerState(temp.CurrentState.Type),
	}, nil
}

// notifyStateChange notifies observers of state and presence changes.
func (m *Manager) notifyStateChange(observers []Observer, old, new *model.DeviceState) {
	change := StateChange{Old: old, New: new}
	for _, o := range observers {
		o.OnStateChange(change)
	}

	// Check for presence changes
	m.checkPresenceChange(observers, model.Left, old.LeftUser, new.LeftUser)
	m.checkPresenceChange(observers, model.Right, old.RightUser, new.RightUser)
}

// checkPresenceChange detects and notifies presence changes for a side.
func (m *Manager) checkPresenceChange(observers []Observer, side model.Side, old, new *model.UserState) {
	oldPresent := old != nil && old.IsPresent()
	newPresent := new != nil && new.IsPresent()

	if oldPresent != newPresent {
		change := PresenceChange{
			Side:    side,
			Present: newPresent,
			User:    new,
		}
		for _, o := range observers {
			o.OnPresenceChange(change)
		}
	}
}

// InvalidateCache forces the next GetState call to refresh from API.
func (m *Manager) InvalidateCache() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cacheExpiry = time.Time{}
}

// SetTemperature sets the target temperature for a side and invalidates cache.
func (m *Manager) SetTemperature(ctx context.Context, side model.Side, level int) error {
	userID, err := m.getUserID(ctx, side)
	if err != nil {
		return err
	}

	if err := m.client.SetUserTemperature(ctx, userID, level); err != nil {
		return err
	}

	m.InvalidateCache()
	return nil
}

// TurnOn powers on a side and invalidates cache.
func (m *Manager) TurnOn(ctx context.Context, side model.Side) error {
	userID, err := m.getUserID(ctx, side)
	if err != nil {
		return err
	}

	if err := m.client.TurnOnUser(ctx, userID); err != nil {
		return err
	}

	m.InvalidateCache()
	return nil
}

// TurnOff powers off a side and invalidates cache.
func (m *Manager) TurnOff(ctx context.Context, side model.Side) error {
	userID, err := m.getUserID(ctx, side)
	if err != nil {
		return err
	}

	if err := m.client.TurnOffUser(ctx, userID); err != nil {
		return err
	}

	m.InvalidateCache()
	return nil
}

// getUserID returns the user ID for a side from cached state.
func (m *Manager) getUserID(ctx context.Context, side model.Side) (string, error) {
	state, err := m.GetState(ctx)
	if err != nil {
		return "", err
	}

	user := state.GetSide(side)
	if user == nil {
		return "", fmt.Errorf("no user assigned to %s side", side)
	}

	return user.ID, nil
}
