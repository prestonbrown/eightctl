# Eight Sleep API Endpoint Audit

This document records the results of endpoint testing performed against the Eight Sleep API (January 2025).

## Summary

| Status | Count |
|--------|-------|
| Working | 6 |
| Broken | 28+ |

Most broken endpoints return "Cannot GET" errors, indicating the paths don't exist in Eight Sleep's current API. The API has likely changed since these endpoints were reverse-engineered.

## Working Endpoints

These endpoints have been verified to work:

| Command | HTTP Method | Endpoint | Notes |
|---------|-------------|----------|-------|
| `whoami` | GET | `/users/me` | Returns user profile, device list |
| `status` | GET | `/users/{userId}/temperature` | Returns current temperature state |
| `device info` | GET | `/devices/{deviceId}` | Returns device properties |
| `device peripherals` | GET | `/devices/{deviceId}/peripherals` | Returns peripheral info |
| `device online` | GET | `/devices/{deviceId}/online` | Returns online status |
| `sleep day` | GET | `/users/{userId}/trends` | Requires `tz` query param |
| `metrics trends` | GET | `/users/{userId}/trends` | Requires `tz` query param |

## Broken Endpoints

These endpoints return "Cannot GET" or similar errors:

### Presence

| Command | Endpoint | Error |
|---------|----------|-------|
| `presence` | `/users/{userId}/presence` | Cannot GET |

### Device

| Command | Endpoint | Error |
|---------|----------|-------|
| `device owner` | `/devices/{deviceId}/owner` | Cannot GET |
| `device warranty` | `/devices/{deviceId}/warranty` | Cannot GET |
| `device priming-tasks` | `/devices/{deviceId}/priming/tasks` | Cannot GET |
| `device priming-schedule` | `/devices/{deviceId}/priming/schedule` | Cannot GET |

### Alarms

| Command | Endpoint | Error |
|---------|----------|-------|
| `alarm list` | `/users/{userId}/alarms` | Cannot GET |
| `alarm create` | `/users/{userId}/alarms` | Cannot POST |
| `alarm update` | `/users/{userId}/alarms/{id}` | Cannot PUT |
| `alarm delete` | `/users/{userId}/alarms/{id}` | Cannot DELETE |
| `alarm snooze` | Various | Cannot access |
| `alarm dismiss` | Various | Cannot access |
| `alarm dismiss-all` | Various | Cannot access |
| `alarm vibration-test` | Various | Cannot access |

### Schedules

| Command | Endpoint | Error |
|---------|----------|-------|
| `schedule list` | `/users/{userId}/schedules` | Cannot GET |
| `schedule create` | `/users/{userId}/schedules` | Cannot POST |
| `schedule update` | `/users/{userId}/schedules/{id}` | Cannot PUT |
| `schedule delete` | `/users/{userId}/schedules/{id}` | Cannot DELETE |
| `schedule next` | (derived from list) | Depends on broken list |

### Temperature Modes

| Command | Endpoint | Error |
|---------|----------|-------|
| `tempmode nap on/off/extend/status` | `/users/{userId}/nap/*` | Cannot access |
| `tempmode hotflash on/off/status` | `/users/{userId}/hotflash/*` | Cannot access |
| `tempmode events` | `/users/{userId}/temp-events` | Cannot GET |

### Adjustable Base

| Command | Endpoint | Error |
|---------|----------|-------|
| `base info` | `/users/{userId}/base` | Cannot GET |
| `base angle` | `/users/{userId}/base/angle` | Cannot POST |
| `base presets` | `/users/{userId}/base/presets` | Cannot GET |
| `base preset-run` | `/users/{userId}/base/preset` | Cannot POST |
| `base vibration-test` | `/users/{userId}/base/vibration-test` | Cannot POST |

### Audio

| Command | Endpoint | Error |
|---------|----------|-------|
| `audio tracks` | `/users/{userId}/audio/tracks` | Cannot GET |
| `audio categories` | `/users/{userId}/audio/categories` | Cannot GET |
| `audio state` | `/users/{userId}/audio/player` | Cannot GET |
| `audio play/pause/seek/volume` | Various player endpoints | Cannot access |
| `audio pair` | `/users/{userId}/audio/pair` | Cannot POST |
| `audio next` | `/users/{userId}/audio/recommended-next` | Cannot GET |
| `audio favorites *` | `/users/{userId}/audio/favorites` | Cannot access |

### Autopilot

| Command | Endpoint | Error |
|---------|----------|-------|
| `autopilot details` | `/users/{userId}/autopilot` | Cannot GET |
| `autopilot history` | `/users/{userId}/autopilot/history` | Cannot GET |
| `autopilot recap` | `/users/{userId}/autopilot/recap` | Cannot GET |
| `autopilot level-suggestions` | `/users/{userId}/autopilot/level-suggestions` | Cannot PUT |
| `autopilot snore-mitigation` | `/users/{userId}/autopilot/snore-mitigation` | Cannot PUT |

### Household

| Command | Endpoint | Error |
|---------|----------|-------|
| `household summary` | `/households/{id}/summary` | Cannot GET |
| `household schedule` | `/households/{id}/schedule` | Cannot GET |
| `household current-set` | `/households/{id}/current-set` | Cannot GET |
| `household invitations` | `/households/{id}/invitations` | Cannot GET |
| `household devices` | `/households/{id}/devices` | Cannot GET |
| `household users` | `/households/{id}/users` | Cannot GET |
| `household guests` | `/households/{id}/guests` | Cannot GET |

### Metrics (partial)

| Command | Endpoint | Error |
|---------|----------|-------|
| `metrics summary` | `/users/{userId}/metrics/summary` | Cannot GET |
| `metrics aggregate` | `/users/{userId}/metrics/aggregate` | Cannot GET |
| `metrics insights` | `/users/{userId}/insights` | Cannot GET |
| `metrics intervals` | `/users/{userId}/intervals/{id}` | Cannot GET |

### Travel

| Command | Endpoint | Error |
|---------|----------|-------|
| `travel trips` | `/users/{userId}/travel/trips` | Cannot GET |
| `travel create-trip` | `/users/{userId}/travel/trips` | Cannot POST |
| `travel delete-trip` | `/users/{userId}/travel/trips/{id}` | Cannot DELETE |
| `travel plans` | `/users/{userId}/travel/trips/{id}/plans` | Cannot GET |
| `travel create-plan` | `/users/{userId}/travel/trips/{id}/plans` | Cannot POST |
| `travel update-plan` | `/users/{userId}/travel/plans/{id}` | Cannot PUT |
| `travel tasks` | `/users/{userId}/travel/plans/{id}/tasks` | Cannot GET |
| `travel airport-search` | `/travel/airports` | Cannot GET |
| `travel flight-status` | `/travel/flights/{number}` | Cannot GET |

### Other

| Command | Endpoint | Error |
|---------|----------|-------|
| `tracks` | `/users/{userId}/audio/tracks` | Cannot GET |
| `feats` | `/release-features` | Cannot GET |

## Known Issues

### Missing `tz` Parameter

The `/users/{userId}/trends` endpoint requires a `tz` (timezone) query parameter. Without it, the endpoint may return errors or incorrect data.

**Fix applied:** Added `tz` parameter to `internal/client/metrics.go` Trends() function.

## Research Sources

For finding correct endpoints:

1. **pyEight (lukas-clarke fork)** - Most actively maintained
   - https://github.com/lukas-clarke/pyEight
   - Files: `pyeight/eight.py`, `pyeight/constants.py`

2. **eight_sleep Home Assistant integration**
   - https://github.com/lukas-clarke/eight_sleep
   - May have newer endpoint discoveries

3. **APK Decompilation** - See [apk-decompilation.md](./apk-decompilation.md)
   - Extract actual endpoints from Android app code

## Recommendations

1. Commands for broken endpoints have been hidden (not removed) using Cobra's `Hidden: true` field
2. APK decompilation should be performed to discover current API endpoints
3. Compare with pyEight's endpoint implementations for correct paths
4. Consider v2 API endpoints - Eight Sleep may have migrated many features
