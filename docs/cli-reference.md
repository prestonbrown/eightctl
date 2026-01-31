# eightctl CLI Reference

Command-line interface for controlling Eight Sleep Pods.

## Installation

```bash
go install github.com/pbrown/eightctl/cmd/eightctl@latest
```

Or build from source:

```bash
go build -o bin/eightctl ./cmd/eightctl
```

## Configuration

Configuration is loaded in this priority order (highest first):

1. Command-line flags (`--email`, `--password`)
2. Environment variables (`EIGHTCTL_EMAIL`, `EIGHTCTL_PASSWORD`)
3. Config file (`~/.config/eightctl/config.yaml`)

### Config File Example

```yaml
email: user@example.com
password: your-password
timezone: America/New_York
output: table
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `EIGHTCTL_EMAIL` | Eight Sleep account email |
| `EIGHTCTL_PASSWORD` | Eight Sleep account password |
| `EIGHTCTL_TIMEZONE` | Timezone for date/time operations |
| `EIGHTCTL_OUTPUT` | Default output format (table, json, csv) |

## Global Flags

| Flag | Description |
|------|-------------|
| `--email` | Eight Sleep account email |
| `--password` | Eight Sleep account password |
| `--output` | Output format: table (default), json, csv |
| `--fields` | Comma-separated list of fields to display |
| `--verbose` | Enable debug logging |
| `--quiet` | Suppress non-essential output |

## Working Commands

These commands have been verified to work with the current Eight Sleep API.

### Authentication & Status

| Command | Description |
|---------|-------------|
| `eightctl whoami` | Show current user profile and devices |
| `eightctl status` | Show current temperature state (`--side left\|right`) |
| `eightctl version` | Show eightctl version |

### Temperature Control

| Command | Description |
|---------|-------------|
| `eightctl on` | Turn on the pod (smart/autopilot mode) (`--side left\|right`) |
| `eightctl off` | Turn off the pod (`--side left\|right`) |
| `eightctl temp <level>` | Set temperature level (-100 to 100) (`--side left\|right`) |

### Device Information

| Command | Description |
|---------|-------------|
| `eightctl device info` | Show device properties |
| `eightctl device peripherals` | Show connected peripherals |
| `eightctl device online` | Show device online status |

### Sleep Data

| Command | Description |
|---------|-------------|
| `eightctl sleep day --date YYYY-MM-DD` | Get sleep data for a specific day |
| `eightctl sleep range --from DATE --to DATE` | Get sleep data for a date range |
| `eightctl metrics trends --from DATE --to DATE` | Get sleep trends |
| `eightctl metrics intervals` | Get sleep session intervals |

### Daemon

| Command | Description |
|---------|-------------|
| `eightctl daemon --schedule FILE` | Run scheduled automations |
| `eightctl daemon --dry-run` | Preview schedule without executing |

### Smart Home Integration

| Command | Description |
|---------|-------------|
| `eightctl mqtt` | Run MQTT bridge for Home Assistant |
| `eightctl hubitat` | Run HTTP server for Hubitat |

#### MQTT Flags

| Flag | Description |
|------|-------------|
| `--broker` | MQTT broker URL (default: tcp://localhost:1883) |
| `--topic-prefix` | Topic prefix for discovery (default: homeassistant) |
| `--device-name` | Device name in Home Assistant (default: Eight Sleep Pod) |
| `--client-id` | MQTT client ID (default: eightctl) |
| `--mqtt-username` | MQTT username (optional) |
| `--mqtt-password` | MQTT password (optional) |
| `--poll-interval` | State polling interval (default: 30s) |

#### Hubitat Flags

| Flag | Description |
|------|-------------|
| `--port` | HTTP server port (default: 8080) |
| `--poll-interval` | State polling interval (default: 30s) |

## Hidden Commands

These commands exist but are hidden because their API endpoints are currently broken. See [endpoint-audit.md](./endpoint-audit.md) for details.

### Alarms
- `alarm list`, `alarm create`, `alarm update`, `alarm delete`
- `alarm snooze`, `alarm dismiss`, `alarm dismiss-all`
- `alarm vibration-test`

### Schedules
- `schedule list`, `schedule create`, `schedule update`, `schedule delete`

### Temperature Modes
- `tempmode nap on|off|extend|status`
- `tempmode hotflash on|off|status`
- `tempmode events`

### Audio
- `audio tracks`, `audio categories`, `audio state`
- `audio play`, `audio pause`, `audio seek`, `audio volume`
- `audio pair`, `audio next`

### Adjustable Base
- `base info`, `base angle`, `base presets`, `base preset-run`
- `base vibration-test`

### Autopilot
- `autopilot details`, `autopilot history`, `autopilot recap`
- `autopilot set-level-suggestions`, `autopilot set-snore-mitigation`

### Household
- `household summary`, `household schedule`, `household current-set`
- `household invitations`

### Metrics (partial)
- `metrics summary`, `metrics aggregate`, `metrics insights`

### Travel
- `travel trips`, `travel create-trip`, `travel delete-trip`
- `travel plans`, `travel tasks`
- `travel airport-search`, `travel flight-status`

### Other
- `presence`
- `tracks`, `feats`

## Output Formats

### Table (default)

```
$ eightctl status
SIDE   STATE   LEVEL   TARGET
left   smart   -10     heating
```

### JSON

```
$ eightctl status --output json
[{"side":"left","state":"smart","level":-10,"target":"heating"}]
```

### CSV

```
$ eightctl status --output csv
side,state,level,target
left,smart,-10,heating
```

## Daemon Schedule Format

The daemon reads a YAML schedule file:

```yaml
schedules:
  - time: "22:00"
    action: "on"
  - time: "06:00"
    action: "temp"
    temperature: -20
  - time: "07:00"
    action: "off"
```

Run with:

```bash
eightctl daemon --schedule ~/.config/eightctl/schedule.yaml
```

## Smart Home Integration

eightctl provides built-in support for smart home platforms.

### MQTT Bridge (Home Assistant)

The MQTT command runs a bridge that publishes Eight Sleep state to an MQTT broker and subscribes to command topics. It supports Home Assistant MQTT Discovery for automatic device configuration.

```bash
eightctl mqtt --broker tcp://mqtt.local:1883
```

The bridge publishes:
- Temperature level and state for each side
- Online/offline status
- Supports commands: on, off, set temperature

See [Home Assistant Guide](./home-assistant.md) for complete setup instructions.

### Hubitat HTTP Server

The Hubitat command runs an HTTP server that exposes a REST API for Hubitat Maker API integration.

```bash
eightctl hubitat --port 8080
```

Endpoints:
- `GET /status` - Current state for both sides
- `POST /left/on`, `POST /right/on` - Turn on a side
- `POST /left/off`, `POST /right/off` - Turn off a side
- `POST /left/temp`, `POST /right/temp` - Set temperature (body: `{"level": -10}`)

See [Hubitat Guide](./hubitat.md) for complete setup instructions.

## See Also

- [API Reference](./api-reference.md) - Eight Sleep API endpoint documentation
- [Endpoint Audit](./endpoint-audit.md) - Status of endpoint testing
- [Development Guide](./development.md) - Contributing and reverse engineering
- [Home Assistant Guide](./home-assistant.md) - MQTT bridge setup for Home Assistant
- [Hubitat Guide](./hubitat.md) - HTTP server setup for Hubitat
