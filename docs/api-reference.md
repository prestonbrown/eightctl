# Eight Sleep API Reference

This document consolidates API information from reverse-engineering efforts, including APK decompilation of the Android app v7.41.66.

> **Note**: Eight Sleep does not publish a stable public API. These endpoints are undocumented and may change without notice.

## Base URLs

| Service | URL |
|---------|-----|
| Client API | `https://client-api.8slp.net/` |
| Auth API | `https://auth-api.8slp.net/` |

## Authentication

### OAuth2 Token Endpoint

**POST** `v1/tokens` (Auth API)

Default credentials (extracted from Android APK v7.41.66):
- `client_id`: `0894c7f33bb94800a03f1f4df13a4f38`
- `client_secret`: `f0954a3ed5763ba3d06834c73731a32f15f168f47d4f164751275def86db0c76`

**Login Request:**
```json
{
  "client_id": "...",
  "client_secret": "...",
  "grant_type": "password",
  "username": "user@example.com",
  "password": "..."
}
```

**Refresh Token Request:**
```json
{
  "client_id": "...",
  "client_secret": "...",
  "grant_type": "refresh_token",
  "refresh_token": "..."
}
```

**Response:**
```json
{
  "access_token": "...",
  "token_type": "Bearer",
  "expires_in": 86400,
  "refresh_token": "...",
  "userId": "..."
}
```

---

## User Management

### Get Current User

**GET** `v1/users/me`

Returns user profile, device list, and feature flags.

### Get User by ID

**GET** `v1/users/{userId}`

Returns user profile and `currentDevice` side assignment.

### Update User

**PUT** `v1/users/{userId}`

### Update Email

**POST** `v1/users/{userId}/email`

### Password Reset

**POST** `v1/users/password-reset`

---

## Device Management

### Get Device

**GET** `v1/devices/{deviceId}`

Query params: `filter=leftUserId,rightUserId,awaySides`

Returns device properties, sensor readings, priming status, water level.

### Update Device

**PUT** `v1/devices/{deviceId}`

### Get Last Heard Time

**GET** `v1/devices/{deviceId}/online`

Returns when device was last online.

### Set Device Owner

**PUT** `v1/devices/{deviceId}/owner`

### Manage Peripherals

**PUT** `v1/devices/{deviceId}/peripherals` - Set all peripherals
**PATCH** `v1/devices/{deviceId}/peripherals` - Add peripheral

### Get Device Warranty

**GET** `v1/devices/{deviceId}/warranty`

### Get BLE Device Key

**POST** `v1/devices/{deviceId}/security/key`

Returns encryption key for offline Bluetooth control.

---

## Temperature Control

### Get Temperature Settings (All)

**GET** `v1/users/{userId}/temperature/all`

Returns complete temperature settings for all modes.

### Get/Set Temperature

**GET** `v1/users/{userId}/temperature`
**PUT** `v1/users/{userId}/temperature/{deviceType}?ignoreDeviceErrors=false`

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

**PUT** `v1/users/{userId}/temperature/{deviceType}`

Turn on:
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

### Get Temperature Events

**GET** `v1/users/{userId}/temp-events`

---

## Nap Mode

### Get Nap Mode Settings

**GET** `v1/users/{userId}/temperature/nap-mode`

### Get Nap Mode Status

**GET** `v1/users/{userId}/temperature/nap-mode/status`

### Start Nap

**POST** `v1/users/{userId}/temperature/nap-mode/activate`

### Stop Nap

**PUT** `v1/users/{userId}/temperature/nap-mode/deactivate`

### Extend Nap

**POST** `v1/users/{userId}/temperature/nap-mode/extend`

### Nap Mode Alarm Settings

**GET/PUT** `v1/users/{userId}/temporary-mode/nap-mode`

---

## Hot Flash Mode

### Get Settings

**GET** `v1/users/{userId}/temperature/hot-flash-mode`

### Update Settings

**PUT** `v1/users/{userId}/temperature/hot-flash-mode`

### Delete Settings

**DELETE** `v1/users/{userId}/temperature/hot-flash-mode`

### Activate

**PUT** `v1/users/{userId}/temperature/hot-flash-mode/activate`

### Deactivate

**PUT** `v1/users/{userId}/temperature/hot-flash-mode/deactivate`

---

## Autopilot

### Get Autopilot Settings

**GET** `v1/users/{userId}/level-suggestions-mode`

### Update Autopilot Settings

**PUT** `v1/users/{userId}/level-suggestions-mode`

### Get Autopilot Details

**GET** `v1/users/{userId}/autopilotDetails`

### Get Autopilot History

**GET** `v1/users/{userId}/autopilot-history`

### Get Autopilot Recap

**GET** `v1/users/{userId}/autopilotDetails/autopilotRecap`

Query params:
- `day`: LocalDate (YYYY-MM-DD)
- `tz`: Timezone ID

### Update Snoring Mitigation

**PUT** `v1/users/{userId}/autopilotDetails/snoringMitigation`

---

## Bedtime Schedule

### Get Bedtime Schedule

**GET** `v1/users/{userId}/temperature`

### Update Bedtime Schedule

**PUT** `v1/users/{userId}/bedtime`

---

## Alarms

### Get Alarms (v2)

**GET** `v2/users/{userId}/alarms`

### Get Alarms (v1)

**GET** `v1/users/{userId}/alarms`

### Create Alarm

**POST** `v1/users/{userId}/alarms`

### Update Alarm

**PUT** `v1/users/{userId}/alarms/{alarmId}`

### Delete Alarm

**DELETE** `v1/users/{userId}/alarms/{alarmId}`

### Dismiss Alarm

**PUT** `v1/users/{userId}/alarms/{alarmId}/dismiss`

### Dismiss All Alarms

**PUT** `v1/users/{userId}/alarms/active/dismiss-all`

### Snooze Alarm

**PUT** `v1/users/{userId}/alarms/{alarmId}/snooze`

### Vibration Test

**POST** `v1/users/{userId}/vibration-test`

---

## Sleep Metrics & Trends

### Get Trends

**GET** `v1/users/{userId}/trends`

Query params:
- `tz`: timezone
- `from`: start date (YYYY-MM-DD)
- `to`: end date
- `include-main`: boolean
- `include-all-sessions`: boolean
- `model-version`: string

### Get Metrics Aggregate

**GET** `v1/users/{userId}/metrics/aggregate?v2=true`

Query params:
- `to`: end date
- `tz`: timezone
- `metrics`: comma-separated metric names
- `periods`: comma-separated periods
- `refreshCache`: boolean

### Get Metrics Summary

**GET** `v1/users/{userId}/metrics/summary`

Query params:
- `from`: start date
- `to`: end date
- `tz`: timezone
- `metrics`: comma-separated metric names

### Update Sleep Session

**PUT** `v1/users/{userId}/intervals/{sessionId}`

### Delete Sleep Session

**DELETE** `v1/users/{userId}/intervals/{sessionId}`

### Send Feedback

**POST** `/v1/users/{userId}/feedback`

### Get Insights

**GET** `v1/users/{userId}/insights`

Query params:
- `date`: LocalDate

---

## AI/LLM Insights

### Get AI Insights

**GET** `v1/users/{userId}/llm-insights`

Query params:
- `from`: start date
- `to`: end date

### Create AI Insights Batch

**POST** `v1/users/{userId}/llm-insights/batch`

### Get AI Insights Settings

**GET** `v1/users/{userId}/llm-insights/settings`

### Update AI Insights Settings

**PUT** `v1/users/{userId}/llm-insights/settings`

### Submit AI Insight Feedback

**POST** `v1/users/{userId}/llm-insights/{insightId}/feedback`

---

## Audio/Speaker

### Get Categories

**GET** `v1/audio/categories`

### Get Available Tracks

**GET** `v1/users/{userId}/audio/tracks`

Query params:
- `category`: category ID (optional)

### Get Recommended Track

**GET** `v1/users/{userId}/audio/tracks/recommended-next-track`

### Get Player State

**GET** `v1/users/{userId}/audio/player`

### Get Player State by Device

**GET** `v1/devices/{deviceId}/audio/player`

### Delete Player (Stop)

**DELETE** `v1/devices/{deviceId}/audio/player`

### Update Player State

**PUT** `v1/users/{userId}/audio/player`

### Control Playback

**PUT** `v1/users/{userId}/audio/player/state`

```json
{
  "state": "Playing" | "Paused"
}
```

### Set Volume

**PUT** `v1/users/{userId}/audio/player/volume`

```json
{
  "volume": 50
}
```

### Seek Position

**PUT** `v1/users/{userId}/audio/player/seek`

### Preview Track

**PUT** `v1/users/{userId}/audio/player/preview-track`

### Favorite Tracks

**PUT** `v1/users/{userId}/audio/tracks/{trackId}/favorites` - Add favorite
**DELETE** `v1/users/{userId}/audio/tracks/{trackId}/favorites` - Remove favorite

### Pair Device Speaker

**PUT** `v1/devices/{deviceId}/audio/player/pair`

---

## Adjustable Base

### Get Base State

**GET** `v1/users/{userId}/base`

Returns current angles and preset data.

### Set Base Angle

**POST** `v1/users/{userId}/base/angle?ignoreDeviceErrors=false`

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

## Household Management

### Get Household Summary

**GET** `v1/household/users/{userId}/summary`

### Get Invitations

**GET** `v1/household/users/{userId}/invitations`

### Set Current Device Side

**PUT** `v1/household/users/{userId}/current-set`

### Clear Current Device (Away)

**DELETE** `v1/household/users/{userId}/current-set`

### Set Return Date

**POST** `v1/household/users/{userId}/schedule`

### Remove Return Date

**DELETE** `v1/household/users/{userId}/schedule/{setId}`

### Add Device to Household

**POST** `v1/household/households/{householdId}/devices`

### Update Household Device

**PUT** `v1/household/devices/{deviceId}`

### Remove Household Device

**DELETE** `v1/household/devices/{deviceId}`

### Invite User

**POST** `v1/household/households/{householdId}/users`

### Respond to Invitation

**POST** `v1/household/households/{householdId}/users/{userId}`

### Remove Guest

**DELETE** `v1/household/households/{householdId}/users/{userId}`

### Add Guests

**POST** `v1/household/households/{householdId}/devices/{deviceId}/guests`

### Update Device Set

**PUT** `v1/household/households/{householdId}/sets/{setId}`

### Remove Device Set

**DELETE** `v1/household/households/{householdId}/sets/{setId}`

### Remove Device Assignment

**DELETE** `v1/household/devices/{deviceId}/assignment/users/{userId}`

---

## Travel / Jet Lag

### Get Trips

**GET** `v1/users/{userId}/travel/trips`

### Get Trip

**GET** `v1/users/{userId}/travel/trips/{tripId}`

### Create Trip

**POST** `v1/users/{userId}/travel/trips`

### Update Trip

**PUT** `v1/users/{userId}/travel/trips/{tripId}`

### Delete Trip

**DELETE** `v1/users/{userId}/travel/trips/{tripId}`

### Search Airports

**GET** `v1/travel/airport-search`

Query params:
- `maxResults`: int
- `query`: search string

### Find Flight

**GET** `v1/travel/flight-status`

Query params:
- `flightNumber`: string
- `date`: YYYY-MM-DD

### Get Jet Lag Plans

**GET** `v1/users/{userId}/travel/trips/{tripId}/plans`

### Create Jet Lag Plan

**POST** `v1/users/{userId}/travel/trips/{tripId}/plans`

### Bulk Update Plan Tasks

**PATCH** `v1/users/{userId}/travel/plans/{planId}/tasks`

---

## Subscriptions

### Get Subscriptions

**GET** `v3/users/{userId}/subscriptions`

### Create Temporary Subscription

**POST** `v3/users/{userId}/subscriptions/temporary`

### Redeem Subscription

**POST** `v3/users/{userId}/subscriptions/redeem`

---

## Device Priming

### Get Priming Schedule

**GET** `v1/devices/{deviceId}/priming/schedule`

### Update Priming Schedule

**PUT** `v1/devices/{deviceId}/priming/schedule`

### Get Priming Tasks

**GET** `v1/devices/{deviceId}/priming/tasks`

### Create Priming Task

**POST** `v1/devices/{deviceId}/priming/tasks`

```json
{
  "notifications": {
    "users": ["userId"],
    "meta": "rePriming"
  }
}
```

### Cancel Priming Task

**DELETE** `v1/devices/{deviceId}/priming/tasks`

---

## Vibration Test

### Start Vibration Test (v2)

**POST** `v2/devices/{deviceId}/vibration-test`

### Stop Vibration Test (v2)

**PUT** `v2/devices/{deviceId}/vibration-test/stop`

---

## Hub Auto-Pairing

### Start Auto-Pairing

**POST** `v1/devices/{deviceId}/auto-pairing/start`

### Get Pairing Status

**GET** `v1/devices/{deviceId}/auto-pairing/status/{pairingId}`

---

## Challenges

### Get Challenges

**GET** `v1/users/{userId}/challenges`

Query params:
- `state`: challenge state filter

---

## Feature Flags

### Get Release Features

**GET** `v1/users/{userId}/release-features`

---

## User Settings

### Get Tap Settings

**GET** `v1/users/{userId}/devices/{deviceId}/tap-settings`

### Update Tap Settings

**PUT** `v1/users/{userId}/devices/{deviceId}/tap-settings`

### Get Tap History

**GET** `/v1/users/{userId}/tap-history`

Query params:
- `from`: timestamp

### Get Level Suggestions

**GET** `v1/users/{userId}/level-suggestions`

### Get Blanket Temperature Recommendations

**GET** `v1/users/{userId}/recommendations/blanket`

### Get Member Perks

**GET** `v1/users/{userId}/perks`

### Get Referral Link

**PUT** `v2/users/{userId}/referral/personal-referral-link`

### Get Referral Campaigns

**GET** `v2/users/{userId}/referral/campaigns`

### Get Purchases

**GET** `v1/purchase-tracker`

### Get Maintenance Insert Status

**GET** `v1/user/{userId}/device_maintenance/maintenance_insert?v=2`

---

## App State

### Get Messages State

**GET** `v1/users/{userId}/app-state/messages`

### Update Messages State

**PUT** `v1/users/{userId}/app-state/messages`
**PATCH** `v1/users/{userId}/app-state/messages`

---

## Health Integrations

### Get Health Survey

**GET** `v1/health-survey/test-drive`

### Update Health Survey

**PATCH** `v1/health-survey/test-drive?enableValidation=true`

### Get Health Integration Checkpoints

**GET** `v1/users/{userId}/health-integrations/sources/{sourceId}/checkpoints`

### Upload Health Integration Data

**POST** `v1/users/{userId}/health-integrations/sources/{sourceId}`

---

## Push Notifications

### Update Push Token

**PUT** `v1/users/me/push-targets/{deviceId}`

### Delete Push Token

**DELETE** `v1/users/me/push-targets/token/{token}`

---

## Extracted Keys & IDs

These values were extracted from Android APK v7.41.66 decompilation.

### OAuth Credentials

Used for authentication with the Auth API:
- `client_id`: `0894c7f33bb94800a03f1f4df13a4f38`
- `client_secret`: `f0954a3ed5763ba3d06834c73731a32f15f168f47d4f164751275def86db0c76`

### Third-Party Service Keys

These are used by the app for analytics/notifications (not needed for API access):

| Service | Key/DSN |
|---------|---------|
| Iterable (Push) | `7992202542544950be7bd5727747d990` |
| Sentry (Errors) | `https://a95f740f7fdc6735289abafdb6fee00e@o4507766157017088.ingest.us.sentry.io/4508844791365632` |

---

## Constants

### Temperature Range

- API value: -100 to 100 (unitless)
- Fahrenheit equivalent: 55°F to 110°F
- Celsius equivalent: 13°C to 44°C

### Request Headers

```
Content-Type: application/json
User-Agent: okhttp/4.9.3
Accept-Encoding: gzip
Authorization: Bearer {access_token}
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

- Android APK v7.41.66 decompilation (January 2025)
- [mezz64/pyEight](https://github.com/mezz64/pyEight) - Original Python library
- [lukas-clarke/pyEight](https://github.com/lukas-clarke/pyEight) - Updated fork with OAuth2
