# Smart Home Integration Design

**Date:** 2025-01-31
**Status:** Approved
**Author:** Brainstorming session

## Overview

This document describes the design for adding smart home integrations to eightctl, including Home Assistant (via MQTT) and Hubitat Elevation support, along with the architectural refactoring needed to support them cleanly.

## Goals

1. **Personal use + community value** - Build integrations the maintainer will use and share
2. **Full bidirectional control** - Expose sensors AND control from smart home platforms
3. **Integration triggers** - Use bed presence to trigger automations (lights, fans, etc.)
4. **Modular architecture** - Avoid tech debt, enable future platform additions
5. **TDD throughout** - Tests before implementation

## Non-Goals (V1)

- Local server emulation (DNS spoofing for fully local control) - V2 stretch goal
- Pure Groovy Hubitat driver (no external dependencies) - V2
- Bluetooth local control
- HomeKit, Google Home integrations

---

## Feature Summary

### Foundation
- Dual-side support (`--side left|right`)
- `eightctl sides` command to show user/side mappings
- Real-time status with presence detection
- Strongly-typed domain models

### Home Assistant Integration
- MQTT bridge command (`eightctl mqtt`)
- Auto-discovery of entities (sensors, switches, climate)
- Presence, temperature, sleep stage, heart rate sensors
- Temperature and power control

### Hubitat Integration
- HTTP bridge command (`eightctl hubitat`)
- REST API for Groovy driver to call
- Parent/child device model (Pod + Left/Right sides)
- Switch, PresenceSensor, Thermostat capabilities

---

## Architecture

### Current Problems

- Client methods return `any` types (no structured data models)
- Commands directly call client and format output (tight coupling)
- No shared state representation between features
- Single-user/side only - no dual-side awareness

### Proposed Layer Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Commands (CLI)                          │
│         status, temp, on, off, realtime, mqtt, hubitat          │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
┌─────────────────────────────────▼───────────────────────────────┐
│                      Platform Adapters                          │
│                                                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │    CLI      │  │    MQTT     │  │   Hubitat   │   (future)  │
│  │   Adapter   │  │   Adapter   │  │   Adapter   │             │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘             │
│         │                │                │                     │
│         └────────────────┼────────────────┘                     │
│                          │                                      │
│              implements Adapter interface                       │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
┌─────────────────────────────────▼───────────────────────────────┐
│                      State Manager                              │
│                                                                 │
│  - Fetches unified device/user state                            │
│  - Caches with configurable TTL                                 │
│  - Notifies observers on state changes                          │
│  - Single source of truth for all adapters                      │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
┌─────────────────────────────────▼───────────────────────────────┐
│                      Domain Models                              │
│                                                                 │
│  DeviceState, UserState, Presence, SleepStage, Alarm, etc.      │
│  (Strongly typed, JSON-serializable, validation methods)        │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
┌─────────────────────────────────▼───────────────────────────────┐
│                      API Client                                 │
│                                                                 │
│  Low-level HTTP, auth, endpoints (existing, enhanced)           │
└─────────────────────────────────────────────────────────────────┘
```

### Package Structure

```
internal/
├── client/           # Low-level API (existing, keep lean)
├── model/            # NEW: Domain types
│   ├── device.go
│   ├── user.go
│   ├── presence.go
│   └── sleep.go
├── state/            # NEW: State management
│   ├── manager.go
│   └── cache.go
├── adapter/          # NEW: Platform adapters
│   ├── adapter.go    # Interface
│   ├── cli/          # CLI output adapter
│   ├── mqtt/         # MQTT bridge adapter
│   └── hubitat/      # Hubitat driver server
├── cmd/              # Commands (existing, refactored)
├── config/           # Config (existing)
├── output/           # Formatting (existing)
└── daemon/           # Scheduler (existing)
```

---

## Domain Models

### DeviceState

```go
type DeviceState struct {
    ID              string
    RoomTemperature float64
    HasWater        bool
    IsPriming       bool
    NeedsPriming    bool
    LeftUser        *UserState
    RightUser       *UserState
}
```

### UserState

```go
type UserState struct {
    ID                string
    Email             string
    Side              Side        // Left or Right
    Present           bool
    BedTemperature    float64
    TargetLevel       int
    State             PowerState  // Off, Smart, Manual
    SleepStage        SleepStage  // Awake, Light, Deep, REM
    HeartRate         float64
    HRV               float64
    BreathRate        float64
    LastHeartRateTime time.Time   // For presence calculation
}

const presenceTimeout = 10 * time.Minute

func (u *UserState) IsPresent() bool {
    return time.Since(u.LastHeartRateTime) < presenceTimeout
}
```

### Adapter Interface

```go
type Adapter interface {
    Start(ctx context.Context, state *StateManager) error
    HandleCommand(cmd Command) error
    Stop() error
}

type Command struct {
    Side   Side
    Action Action  // SetTemp, TurnOn, TurnOff
    Value  any
}
```

### State Manager

```go
type StateProvider interface {
    GetState(ctx context.Context) (*model.DeviceState, error)
    SetTemperature(ctx context.Context, side model.Side, level int) error
    TurnOn(ctx context.Context, side model.Side) error
    TurnOff(ctx context.Context, side model.Side) error
}

type Observer interface {
    OnStateChange(old, new *DeviceState)
    OnPresenceChange(side Side, present bool)
}
```

---

## Home Assistant Integration (MQTT)

### Command

```bash
eightctl mqtt --broker mqtt://localhost:1883 --interval 60
```

### Auto-Created Entities

Per side (left/right):
| Entity | Type | Purpose |
|--------|------|---------|
| `binary_sensor.eight_sleep_left_presence` | Binary Sensor | In bed / not in bed |
| `sensor.eight_sleep_left_bed_temp` | Sensor | Current bed temperature |
| `sensor.eight_sleep_left_heart_rate` | Sensor | Current heart rate |
| `sensor.eight_sleep_left_sleep_stage` | Sensor | awake/light/deep/rem |
| `climate.eight_sleep_left` | Climate | Temperature control |
| `switch.eight_sleep_left_power` | Switch | On/off control |

Device-level:
| Entity | Type | Purpose |
|--------|------|---------|
| `sensor.eight_sleep_room_temp` | Sensor | Room temperature |
| `binary_sensor.eight_sleep_water_low` | Binary Sensor | Water level warning |

### Topic Structure

```
homeassistant/binary_sensor/eight_sleep/left_presence/config   # Discovery
homeassistant/binary_sensor/eight_sleep/left_presence/state    # "ON"/"OFF"

eight_sleep/left/temperature/set   # Command topic
eight_sleep/left/power/set         # Command topic
```

### Configuration

```yaml
# ~/.config/eightctl/config.yaml
mqtt:
  broker: mqtt://localhost:1883
  username: homeassistant
  password: secret
  discovery_prefix: homeassistant
  poll_interval: 60
```

---

## Hubitat Integration

### Architecture

```
┌──────────────┐     HTTP      ┌──────────────┐     HTTPS    ┌──────────────┐
│   Hubitat    │ ◄──────────► │   eightctl   │ ◄──────────► │  Eight Sleep │
│   Driver     │  (local LAN) │   hubitat    │   (internet) │     API      │
│   (Groovy)   │              │   command    │              │              │
└──────────────┘              └──────────────┘              └──────────────┘
```

### Command

```bash
eightctl hubitat --port 8380 --poll-interval 60
```

### REST Endpoints

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/status` | GET | Full device state (both sides) |
| `/left/status` | GET | Left side state |
| `/right/status` | GET | Right side state |
| `/left/temperature` | PUT | Set left temp `{"level": -20}` |
| `/right/temperature` | PUT | Set right temp |
| `/left/power` | PUT | Turn on/off `{"on": true}` |
| `/right/power` | PUT | Turn on/off |

### Hubitat Devices Created

| Device | Capabilities | Use Case |
|--------|-------------|----------|
| Eight Sleep Pod | Refresh, Temp | Room temperature, refresh all |
| Eight Sleep Left | Switch, Presence, Temp, Thermostat | Left side control |
| Eight Sleep Right | Switch, Presence, Temp, Thermostat | Right side control |

### Groovy Drivers

Parent and child Groovy drivers will be provided in `drivers/hubitat/` for users to install on their hubs.

---

## Testing Strategy (TDD)

### Test Layers

1. **Model Tests** - Validation, JSON serialization, presence logic
2. **State Manager Tests** - Caching, observer dispatch, state diffing
3. **Adapter Tests** - MQTT topics/payloads, Hubitat HTTP handlers
4. **Integration Tests** - End-to-end with real MQTT broker

### Coverage Targets

| Package | Target |
|---------|--------|
| `model/` | 90%+ |
| `state/` | 85%+ |
| `adapter/mqtt/` | 80%+ |
| `adapter/hubitat/` | 80%+ |

### TDD Workflow

1. Write failing test for new behavior
2. Implement minimum code to pass
3. Refactor while keeping tests green

---

## Implementation Phases

### Phase 1: Foundation (No Breaking Changes)

- Add `internal/model/` with domain types
- Add `internal/state/` with StateManager
- Add dual-side client methods
- Add `eightctl sides` command

### Phase 2: Enhanced CLI

- Add `--side` flag to temp, on, off, status
- Add `eightctl realtime` command
- Add `eightctl presence` command
- Refactor existing commands to use StateManager

### Phase 3: MQTT Bridge

- Add `internal/adapter/mqtt/`
- Add HA discovery config generation
- Add `eightctl mqtt` command
- Observer integration for state changes

### Phase 4: Hubitat Bridge

- Add `internal/adapter/hubitat/`
- Add `eightctl hubitat` command
- Create Groovy parent driver
- Create Groovy child driver

### Phase 5: Polish & Documentation

- Update CLI reference docs
- Add Home Assistant setup guide
- Add Hubitat installation guide
- Ship Groovy drivers in repo

---

## Future Ideas (V2+)

| Idea | Description | Complexity |
|------|-------------|------------|
| Local Server Emulation | Emulate Eight Sleep API, DNS spoof for fully local control | High |
| Pure Groovy Driver | Full Hubitat driver with no external dependencies | Medium |
| WebSocket Events | Real-time push instead of polling | Medium |
| Bluetooth Local Control | Use BLE key for direct pod communication | High |
| Multi-pod Support | Households with multiple Eight Sleep devices | Low |
| Apple HomeKit | HomeKit bridge via HAP-go | Medium |
| Google Home | Local fulfillment integration | Medium |

---

## References

- [pyEight](https://github.com/lukas-clarke/pyEight) - Python library with working dual-side support
- [eight_sleep HA integration](https://github.com/lukas-clarke/eight_sleep) - Official HA integration
- [Eight Sleep API Reference](../api-reference.md) - Endpoint documentation
