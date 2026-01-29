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
| `eightctl status` | Show current temperature state |
| `eightctl version` | Show eightctl version |

### Temperature Control

| Command | Description |
|---------|-------------|
| `eightctl on` | Turn on the pod (smart/autopilot mode) |
| `eightctl off` | Turn off the pod |
| `eightctl temp <level>` | Set temperature level (-100 to 100) |

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

### Daemon

| Command | Description |
|---------|-------------|
| `eightctl daemon --schedule FILE` | Run scheduled automations |
| `eightctl daemon --dry-run` | Preview schedule without executing |

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
- `metrics intervals`

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

## See Also

- [API Reference](./api-reference.md) - Eight Sleep API endpoint documentation
- [Endpoint Audit](./endpoint-audit.md) - Status of endpoint testing
- [Development Guide](./development.md) - Contributing and reverse engineering
