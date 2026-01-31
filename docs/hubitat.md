# Hubitat Integration

Integrate your Eight Sleep Pod with Hubitat Elevation for smart home automation. This integration uses eightctl as a bridge between Hubitat and the Eight Sleep API.

## Prerequisites

- eightctl installed and configured with your Eight Sleep account credentials
- Hubitat Elevation hub on the same network as the machine running eightctl
- Eight Sleep Pod already set up and working with the official app

## Installation

### Step 1: Install the Drivers

1. Open your Hubitat web interface (typically `http://hubitat.local` or your hub's IP)
2. Navigate to **Drivers Code** in the sidebar
3. Click **New Driver**
4. Copy the entire contents of `drivers/hubitat/eight-sleep-pod.groovy` and paste it into the editor
5. Click **Save**
6. Click **New Driver** again
7. Copy the entire contents of `drivers/hubitat/eight-sleep-side.groovy` and paste it into the editor
8. Click **Save**

### Step 2: Start the eightctl Server

Run the hubitat server on a machine that will remain online:

```bash
eightctl hubitat --port 8080
```

The server will start and display:

```
Hubitat server listening on port 8080
```

For production use, run as a systemd service:

```ini
# /etc/systemd/system/eightctl-hubitat.service
[Unit]
Description=eightctl Hubitat Server
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/eightctl hubitat --port 8080
Restart=always
RestartSec=10
User=eightctl

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable eightctl-hubitat
sudo systemctl start eightctl-hubitat
```

### Step 3: Create the Device

1. In Hubitat, navigate to **Devices** > **Add Device** > **Virtual**
2. Enter a name (e.g., "Eight Sleep Pod")
3. Select **Eight Sleep Pod** as the driver
4. Click **Save Device**
5. Configure the device preferences:
   - **Server IP**: The IP address of the machine running eightctl (e.g., `192.168.1.100`)
   - **Server Port**: `8080` (or your custom port)
   - **Refresh Interval**: How often to poll for status updates (default: 60 seconds)
6. Click **Save Preferences**
7. Click **Create Child Devices** to create the left and right side devices

After creating child devices, you will have three devices:
- **Eight Sleep Pod** - Parent device that manages the connection
- **Eight Sleep Pod - Left** - Controls the left side of the bed
- **Eight Sleep Pod - Right** - Controls the right side of the bed

## Configuration Options

### Command-Line Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--port` | `8080` | HTTP server port |
| `--poll-interval` | `30s` | How often to refresh state from Eight Sleep API |

### Config File

Add hubitat settings to `~/.config/eightctl/config.yaml`:

```yaml
email: user@example.com
password: your-password

hubitat:
  port: 8080
  poll-interval: 30s
```

## Using the Devices

The left and right side devices expose the following capabilities:

### Switch Capability
- **on** - Turn on the side (activates smart/autopilot mode)
- **off** - Turn off the side

### Thermostat Capability
- **thermostatMode** - "auto" when active, "off" when inactive
- **heatingSetpoint** / **coolingSetpoint** - Temperature mapped from Eight Sleep level

### Temperature Level
- **level** - Current Eight Sleep level (-100 to +100)
- **targetLevel** - Target level
- **setLevel(n)** - Set temperature level directly
- **levelUp** / **levelDown** - Adjust level by step (default: 10)

### Level Reference
| Level | Effect |
|-------|--------|
| -100 | Maximum cooling (~55F) |
| -50 | Moderate cooling |
| 0 | Neutral (no heating or cooling) |
| +50 | Moderate heating |
| +100 | Maximum heating (~110F) |

## Example Rules

### Turn Off When Not in Bed

Using Rule Machine:

1. Create a new rule
2. Trigger: Eight Sleep Left switch changes
3. Condition: Eight Sleep Left switch is on
4. Action: If level > 50 for 30 minutes, turn off Eight Sleep Left

### Scheduled Heating Before Bedtime

Using Rule Machine:

1. Create a new rule
2. Trigger: Time 21:30
3. Actions:
   - Turn on Eight Sleep Left
   - Set Eight Sleep Left level to 20

### Integration with Motion Sensors

1. Create a new rule
2. Trigger: Bedroom motion sensor becomes inactive
3. Condition: Time between 22:00 and 06:00
4. Actions:
   - Delay 10 minutes (cancelable)
   - Turn on Eight Sleep Left
   - Set Eight Sleep Left level to -10

## Troubleshooting

### Connection Refused

If you see "connection refused" in Hubitat logs:

1. Verify eightctl hubitat server is running: `systemctl status eightctl-hubitat`
2. Check the server is accessible from Hubitat: `curl http://<server-ip>:8080/status`
3. Verify no firewall is blocking port 8080

### Device Shows "Not Configured"

The parent device will show this status if the Server IP is not set:

1. Open the Eight Sleep Pod device in Hubitat
2. Set the **Server IP** to your eightctl server's IP address
3. Click **Save Preferences**

### Authentication Errors

If the eightctl server fails to start with authentication errors:

1. Verify your credentials: `eightctl whoami`
2. Try re-authenticating: `eightctl login`
3. Check the config file: `cat ~/.config/eightctl/config.yaml`

### Enable Debug Logging

In Hubitat:

1. Open the Eight Sleep Pod device
2. Enable **Enable debug logging**
3. Click **Save Preferences**
4. View logs at **Logs** in the sidebar

On the eightctl server:

```bash
eightctl hubitat --port 8080 --verbose
```

### Child Devices Not Updating

1. Check the parent device's **connectionStatus** attribute
2. Click **Refresh** on the parent device
3. Verify the server is responding: `curl http://<server-ip>:8080/status`

## API Reference

The eightctl hubitat server exposes these HTTP endpoints:

### GET /status

Returns full device status for both sides.

**Request:**
```bash
curl http://localhost:8080/status
```

**Response:**
```json
{
  "left": {
    "isActive": true,
    "currentLevel": -10,
    "targetLevel": -10,
    "currentTemperature": 72
  },
  "right": {
    "isActive": false,
    "currentLevel": 0,
    "targetLevel": 0,
    "currentTemperature": 71
  }
}
```

### GET /{side}/status

Returns status for a specific side.

**Request:**
```bash
curl http://localhost:8080/left/status
```

**Response:**
```json
{
  "isActive": true,
  "currentLevel": -10,
  "targetLevel": -10,
  "currentTemperature": 72
}
```

### PUT /{side}/on

Turn on a side.

**Request:**
```bash
curl -X PUT http://localhost:8080/left/on
```

**Response:** `200 OK`
```json
{"status": "ok"}
```

### PUT /{side}/off

Turn off a side.

**Request:**
```bash
curl -X PUT http://localhost:8080/left/off
```

**Response:** `200 OK`
```json
{"status": "ok"}
```

### PUT /{side}/temperature?level=N

Set the temperature level for a side.

**Request:**
```bash
curl -X PUT "http://localhost:8080/left/temperature?level=-20"
```

**Response:** `200 OK`
```json
{"status": "ok"}
```

**Parameters:**
- `level` (required): Integer from -100 to 100

## See Also

- [CLI Reference](./cli-reference.md) - Full eightctl command documentation
- [API Reference](./api-reference.md) - Eight Sleep API endpoint documentation
