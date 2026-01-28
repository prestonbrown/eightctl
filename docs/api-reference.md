# Eight Sleep API Reference

This document consolidates API information from reverse-engineering efforts across multiple projects, primarily [pyEight](https://github.com/mezz64/pyEight) and [lukas-clarke/pyEight](https://github.com/lukas-clarke/pyEight).

> **Note**: Eight Sleep does not publish a stable public API. These endpoints are undocumented and may change without notice.

## Base URLs

| Service | URL |
|---------|-----|
| Client API | `https://client-api.8slp.net/v1` |
| Auth API | `https://auth-api.8slp.net/v1` |
| App API | `https://app-api.8slp.net/` |

## Authentication

### OAuth2 Token Endpoint

**POST** `/tokens` (Auth API)

Default credentials (extracted from Android APK v7.39.17):
- `client_id`: `0894c7f33bb94800a03f1f4df13a4f38`
- `client_secret`: `f0954a3ed5763ba3d06834c73731a32f15f168f47d4f164751275def86db0c76`

**Request Body:**
```json
{
  "client_id": "...",
  "client_secret": "...",
  "grant_type": "password",
  "username": "user@example.com",
  "password": "..."
}
```

**Response:** Bearer token with expiration

### Legacy Session Auth

**POST** `/login` (Client API)

Returns session token in `Session-Token` header.

---

## User & Device Management

### Get Current User

**GET** `/users/me`

Returns user profile, device list, and feature flags (cooling, elevation, audio).

### Get User Profile

**GET** `/users/{userId}`

Returns user profile and `currentDevice` side assignment.

### Get/Set Current Device Side

**GET/PUT** `/users/{userId}/current-device`

**PUT Body:**
```json
{
  "side": "solo" | "left" | "right"
}
```

### Get Device Info

**GET** `/devices/{deviceId}`

Query params: `filter=leftUserId,rightUserId,awaySides`

Returns device properties, sensor readings, priming status, water level.

---

## Temperature Control

### Get Temperature State

**GET** `/v1/users/{userId}/temperature`

Returns current level, device level, smart schedule, bed state.

### Set Temperature Level

**PUT** `/v1/users/{userId}/temperature`

Set immediate level (-100 to 100):
```json
{
  "currentLevel": 20
}
```

Set level with duration:
```json
{
  "timeBased": {
    "level": 20,
    "durationSeconds": 3600
  }
}
```

Set smart/autopilot temperatures by sleep stage:
```json
{
  "smart": {
    "bedTimeLevel": 10,
    "initialSleepLevel": -20,
    "finalSleepLevel": 30
  }
}
```

### Turn On/Off

**PUT** `/v1/users/{userId}/temperature`

Turn on (activate side):
```json
{
  "currentState": {
    "type": "smart"
  }
}
```

Turn off:
```json
{
  "currentState": {
    "type": "off"
  }
}
```

---

## Alarms & Routines

### Get Routines (v2)

**GET** `/v2/users/{userId}/routines`

Returns all routines and next alarm info.

### Update Routine (v2)

**PUT** `/v2/users/{userId}/routines/{routineId}`

Update alarm times or bedtime settings.

### Create One-Time Alarm (v2)

**PUT** `/v2/users/{userId}/routines?ignoreDeviceErrors=false`

```json
{
  "oneOffAlarms": [{
    "time": "07:00",
    "enabled": true,
    "settings": {
      "vibration": true,
      "thermal": true
    }
  }]
}
```

### Get Alarms (v1)

**GET** `/v1/users/{userId}/alarms`

### Update Alarm (v1)

**PUT** `/v1/users/{userId}/alarms/{alarmId}`

Configure time, vibration, thermal, audio, smart features.

### Alarm Control (v1)

**PUT** `/v1/users/{userId}/routines`

Snooze:
```json
{
  "alarm": {
    "alarmId": "...",
    "snoozeForMinutes": 9
  }
}
```

Stop:
```json
{
  "alarm": {
    "alarmId": "...",
    "stopped": true
  }
}
```

Dismiss:
```json
{
  "alarm": {
    "alarmId": "...",
    "dismissed": true
  }
}
```

---

## Bedtime Schedule

### Set Bedtime

**PUT** `/v1/users/{userId}/bedtime`

Set bedtime schedule and temperature profile.

---

## Sleep Data & Trends

### Get Trends

**GET** `/users/{userId}/trends`

Query params:
- `tz`: timezone
- `from`: start date (YYYY-MM-DD)
- `to`: end date
- `include-main`: boolean
- `include-all-sessions`: boolean
- `model-version`: string

### Get Intervals

**GET** `/users/{userId}/intervals`

Returns sleep session scores and detailed metrics.

---

## Adjustable Base

### Get Base State

**GET** `/v1/users/{userId}/base`

Returns current angles and preset data.

### Set Base Angle

**POST** `/v1/users/{userId}/base/angle?ignoreDeviceErrors=false`

Set specific angles:
```json
{
  "deviceId": "...",
  "legAngle": 0,
  "torsoAngle": 15,
  "enableOfflineMode": false
}
```

Apply preset:
```json
{
  "deviceId": "...",
  "preset": "sleep" | "relaxing" | "reading",
  "enableOfflineMode": false
}
```

---

## Audio/Speaker

### Get Player State

**GET** `/v1/users/{userId}/audio/player`

Returns playback state and hardware info.

### Get Available Tracks

**GET** `/v1/users/{userId}/audio/tracks`

### Control Playback

**PUT** `/v1/users/{userId}/audio/player/state`

```json
{
  "state": "Playing" | "Paused"
}
```

### Set Volume

**PUT** `/v1/users/{userId}/audio/player/volume`

```json
{
  "volume": 50
}
```

### Select Track

**PUT** `/v1/users/{userId}/audio/player/currentTrack`

```json
{
  "id": "track-id",
  "stopCriteria": "ManualStop"
}
```

---

## Away Mode

**PUT** `/v1/users/{userId}/away-mode`

Start away mode:
```json
{
  "awayPeriod": {
    "start": "2024-01-15T00:00:00Z"
  }
}
```

End away mode:
```json
{
  "awayPeriod": {
    "end": "2024-01-20T00:00:00Z"
  }
}
```

---

## Device Maintenance

### Prime Pod

**POST** `/v1/devices/{deviceId}/priming/tasks`

```json
{
  "notifications": {
    "users": ["userId"],
    "meta": "rePriming"
  }
}
```

---

## Constants

### Temperature Range

- API value: -100 to 100 (unitless)
- Fahrenheit equivalent: 55°F to 110°F
- Celsius equivalent: 13°C to 44°C

### Sleep Stages

- `bedTimeLevel`
- `initialSleepLevel`
- `finalSleepLevel`

### Request Headers

```
Content-Type: application/json
User-Agent: okhttp/4.9.3
Accept-Encoding: gzip
```

### Timeouts

- Recommended: 20-240 seconds (API can be slow)
- Token refresh buffer: 120 seconds before expiry

---

## Rate Limiting

The API returns HTTP 429 when rate limited. Recommended handling:
- Retry with exponential backoff
- Re-authenticate on 401 responses

---

## Sources

- [mezz64/pyEight](https://github.com/mezz64/pyEight) - Original Python library
- [lukas-clarke/pyEight](https://github.com/lukas-clarke/pyEight) - Updated fork with OAuth2 and alarm fixes
- [lukas-clarke/eight_sleep](https://github.com/lukas-clarke/eight_sleep) - Home Assistant integration
