# Home Assistant Integration

Integrate your Eight Sleep Pod with Home Assistant using MQTT Discovery. The eightctl MQTT bridge automatically publishes device configurations and state updates, allowing seamless control through the Home Assistant interface.

## Prerequisites

- eightctl installed and configured with your Eight Sleep account credentials
- MQTT broker (e.g., Mosquitto) accessible to both eightctl and Home Assistant
- Home Assistant with MQTT integration enabled
- Eight Sleep Pod already set up and working with the official app

## Installation

### Step 1: Configure MQTT in Home Assistant

1. Open Home Assistant
2. Navigate to **Settings** > **Devices & Services** > **Integrations**
3. Click **Add Integration** and search for **MQTT**
4. Configure your MQTT broker connection (typically `localhost:1883` if running Mosquitto on the same machine)
5. Ensure **Enable discovery** is checked (this is the default)

### Step 2: Start the eightctl MQTT Bridge

Run the MQTT bridge to connect your Eight Sleep Pod to Home Assistant:

```bash
eightctl mqtt --broker tcp://localhost:1883
```

You should see output like:

```
MQTT bridge connected to tcp://localhost:1883
Publishing to homeassistant discovery prefix
```

For production use, run as a systemd service:

```ini
# /etc/systemd/system/eightctl-mqtt.service
[Unit]
Description=eightctl MQTT Bridge for Home Assistant
After=network.target mosquitto.service

[Service]
Type=simple
ExecStart=/usr/local/bin/eightctl mqtt --broker tcp://localhost:1883
Restart=always
RestartSec=10
User=eightctl

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable eightctl-mqtt
sudo systemctl start eightctl-mqtt
```

### Step 3: Verify Discovery

1. In Home Assistant, go to **Settings** > **Devices & Services** > **MQTT**
2. Click on **Devices** to see the newly discovered Eight Sleep Pod
3. Climate entities should appear automatically for both sides of the bed

## Configuration Options

### Command-Line Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--broker` | `tcp://localhost:1883` | MQTT broker URL |
| `--topic-prefix` | `homeassistant` | Topic prefix for MQTT discovery |
| `--device-name` | `Eight Sleep Pod` | Device name in Home Assistant |
| `--client-id` | `eightctl` | MQTT client ID |
| `--mqtt-username` | | MQTT broker username (optional) |
| `--mqtt-password` | | MQTT broker password (optional) |
| `--poll-interval` | `30s` | How often to poll state from Eight Sleep API |

### Config File

Add MQTT settings to `~/.config/eightctl/config.yaml`:

```yaml
email: user@example.com
password: your-password

mqtt:
  broker: tcp://localhost:1883
  topic-prefix: homeassistant
  device-name: Bedroom Pod
  client-id: eightctl
  mqtt-username: homeassistant
  mqtt-password: your-mqtt-password
  poll-interval: 30s
```

### Environment Variables

All settings can be configured via environment variables:

```bash
export EIGHTCTL_EMAIL=user@example.com
export EIGHTCTL_PASSWORD=your-password
export EIGHTCTL_MQTT_BROKER=tcp://localhost:1883
export EIGHTCTL_MQTT_TOPIC_PREFIX=homeassistant
export EIGHTCTL_MQTT_DEVICE_NAME="Bedroom Pod"
export EIGHTCTL_MQTT_POLL_INTERVAL=30s
```

## Entities Created

The MQTT bridge creates climate entities for each side of the bed:

| Entity ID | Description |
|-----------|-------------|
| `climate.eight_sleep_pod_left` | Left side climate control |
| `climate.eight_sleep_pod_right` | Right side climate control |

Each climate entity supports:

- **Temperature**: Set the target temperature level (-100 to +100)
- **Modes**: `off`, `heat`, `cool`
- **Current Temperature**: Read the current bed temperature

## Understanding Temperature Levels

Eight Sleep uses a level system from -100 to +100, not traditional temperature units:

| Level | Effect |
|-------|--------|
| -100 | Maximum cooling |
| -50 | Moderate cooling |
| 0 | Neutral (no active heating or cooling) |
| +50 | Moderate heating |
| +100 | Maximum heating |

In the Home Assistant interface, the temperature slider displays these levels. While Home Assistant requires a temperature unit, the values represent Eight Sleep's proprietary level system, not degrees Fahrenheit or Celsius.

**Mode Behavior:**
- `heat` mode: Positive levels (warming the bed)
- `cool` mode: Negative levels (cooling the bed)
- `off` mode: System is disabled

## Example Automations

### Pre-heat Bed Before Bedtime

Warm up the bed 30 minutes before your typical bedtime:

```yaml
automation:
  - alias: "Pre-heat bed before bedtime"
    trigger:
      - platform: time
        at: "21:30:00"
    condition:
      - condition: state
        entity_id: person.your_name
        state: "home"
    action:
      - service: climate.set_hvac_mode
        target:
          entity_id: climate.eight_sleep_pod_left
        data:
          hvac_mode: heat
      - service: climate.set_temperature
        target:
          entity_id: climate.eight_sleep_pod_left
        data:
          temperature: 30
```

### Turn Off When Leaving Home

Automatically turn off the Pod when everyone leaves:

```yaml
automation:
  - alias: "Turn off Eight Sleep when away"
    trigger:
      - platform: state
        entity_id: zone.home
        to: "0"
    action:
      - service: climate.set_hvac_mode
        target:
          entity_id:
            - climate.eight_sleep_pod_left
            - climate.eight_sleep_pod_right
        data:
          hvac_mode: "off"
```

### Cool Down After Falling Asleep

Gradually cool the bed after you've likely fallen asleep:

```yaml
automation:
  - alias: "Cool bed after falling asleep"
    trigger:
      - platform: time
        at: "23:30:00"
    condition:
      - condition: state
        entity_id: climate.eight_sleep_pod_left
        attribute: hvac_mode
        state: heat
    action:
      - service: climate.set_hvac_mode
        target:
          entity_id: climate.eight_sleep_pod_left
        data:
          hvac_mode: cool
      - service: climate.set_temperature
        target:
          entity_id: climate.eight_sleep_pod_left
        data:
          temperature: -20
```

### Wake-up Warming Routine

Gently warm the bed before your alarm:

```yaml
automation:
  - alias: "Morning warmup routine"
    trigger:
      - platform: time
        at: "06:30:00"
    condition:
      - condition: time
        weekday:
          - mon
          - tue
          - wed
          - thu
          - fri
    action:
      - service: climate.set_hvac_mode
        target:
          entity_id: climate.eight_sleep_pod_left
        data:
          hvac_mode: heat
      - service: climate.set_temperature
        target:
          entity_id: climate.eight_sleep_pod_left
        data:
          temperature: 40
```

## Lovelace Dashboard

### Thermostat Card

Add a simple thermostat card for controlling the Pod:

```yaml
type: thermostat
entity: climate.eight_sleep_pod_left
name: My Side
```

### Custom Button Card

Create dedicated buttons for quick temperature presets:

```yaml
type: horizontal-stack
cards:
  - type: button
    name: Cool
    icon: mdi:snowflake
    tap_action:
      action: call-service
      service: climate.set_temperature
      target:
        entity_id: climate.eight_sleep_pod_left
      data:
        temperature: -30
        hvac_mode: cool

  - type: button
    name: Neutral
    icon: mdi:bed
    tap_action:
      action: call-service
      service: climate.set_hvac_mode
      target:
        entity_id: climate.eight_sleep_pod_left
      data:
        hvac_mode: "off"

  - type: button
    name: Warm
    icon: mdi:fire
    tap_action:
      action: call-service
      service: climate.set_temperature
      target:
        entity_id: climate.eight_sleep_pod_left
      data:
        temperature: 30
        hvac_mode: heat
```

### Entities Card with Both Sides

Display both sides of the bed together:

```yaml
type: entities
title: Eight Sleep Pod
entities:
  - entity: climate.eight_sleep_pod_left
    name: Left Side
  - entity: climate.eight_sleep_pod_right
    name: Right Side
```

## Troubleshooting

### No Entities Appearing

If climate entities don't appear in Home Assistant:

1. **Check MQTT broker connection**: Verify the broker is running and accessible
   ```bash
   mosquitto_pub -h localhost -t test -m "hello"
   ```

2. **Verify eightctl is connected**: Check that the MQTT bridge is running
   ```bash
   systemctl status eightctl-mqtt
   ```

3. **Check discovery prefix**: Ensure `--topic-prefix` matches Home Assistant's discovery prefix (default: `homeassistant`)

4. **View MQTT messages**: Use an MQTT client to verify discovery messages are published
   ```bash
   mosquitto_sub -h localhost -t "homeassistant/climate/#" -v
   ```

### Entities Show as Unavailable

If entities appear but show as "unavailable":

1. **Check eightctl is running**: The bridge must be running to publish availability
   ```bash
   systemctl status eightctl-mqtt
   ```

2. **Check availability topic**: Verify the availability message is published
   ```bash
   mosquitto_sub -h localhost -t "eightsleep/+/availability" -v
   ```

3. **Restart the bridge**: Sometimes a restart helps reconnect
   ```bash
   sudo systemctl restart eightctl-mqtt
   ```

### Commands Not Working

If temperature or mode changes don't take effect:

1. **Check eightctl logs**: Look for authentication or API errors
   ```bash
   journalctl -u eightctl-mqtt -f
   ```

2. **Verify Eight Sleep credentials**: Test with the CLI directly
   ```bash
   eightctl status
   ```

3. **Check command topics**: Verify commands are received
   ```bash
   mosquitto_sub -h localhost -t "eightsleep/+/+/set_#" -v
   ```

### Authentication Errors

If the MQTT bridge fails to start:

1. **Verify your Eight Sleep credentials**:
   ```bash
   eightctl whoami
   ```

2. **Re-authenticate if needed**:
   ```bash
   eightctl login
   ```

3. **Check config file permissions**:
   ```bash
   ls -la ~/.config/eightctl/config.yaml
   ```

## MQTT Topics Reference

### Discovery Topics

Discovery configurations are published to:

```
homeassistant/climate/{device_id}_left/config
homeassistant/climate/{device_id}_right/config
```

### State Topics

Current state is published to:

| Topic | Description | Example Value |
|-------|-------------|---------------|
| `eightsleep/{device_id}/{side}/temperature` | Target temperature level | `25` |
| `eightsleep/{device_id}/{side}/mode` | Current mode | `heat`, `cool`, `off` |
| `eightsleep/{device_id}/{side}/current_temperature` | Bed temperature reading | `72.5` |

### Command Topics

Commands are received on:

| Topic | Description | Expected Payload |
|-------|-------------|------------------|
| `eightsleep/{device_id}/{side}/set_temperature` | Set target level | `-100` to `100` |
| `eightsleep/{device_id}/{side}/set_mode` | Set mode | `heat`, `cool`, `off` |

### Availability Topic

Bridge availability is published to:

```
eightsleep/{device_id}/availability
```

Payload: `online` or `offline`

## See Also

- [CLI Reference](./cli-reference.md) - Full eightctl command documentation
- [Hubitat Integration](./hubitat.md) - Alternative smart home integration
- [API Reference](./api-reference.md) - Eight Sleep API endpoint documentation
