# Endpoint Verification Report

**Date:** 2026-01-29 17:39:55
**Total Tested:** 33
**Passed:** 2
**Failed:** 31
**Skipped:** 1

---

## ✓ PASSED (2)

| Command | Notes |
|---------|-------|
| `audio favorites` | Favorite tracks |
| `metrics intervals` | Using config file: /Users/pbrown/.config/eightctl/... |

## ✗ FAILED (31)

| Command | Error |
|---------|-------|
| `alarm list` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `audio tracks` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `audio categories` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `audio state` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `autopilot details` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `autopilot history` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `autopilot recap` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `base info` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `base presets` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `tempmode nap status` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `tempmode hotflash status` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `tempmode events` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `household summary` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `household schedule` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `household current-set` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `household invitations` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `household devices` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `household users` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `household guests` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `metrics summary` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `metrics aggregate` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `metrics insights` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `device owner` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `device warranty` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `device priming-tasks` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `device priming-schedule` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `travel trips` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `travel airport-search JFK` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `feats` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `tracks` | Using config file: /Users/pbrown/.config/eightctl/config.yam |
| `schedule list` | Using config file: /Users/pbrown/.config/eightctl/config.yam |

## ⊘ SKIPPED (1)

| Command | Reason |
|---------|--------|
| `reversible writes` | Skipped (--readonly-only) |

---

## Next Steps

For PASSED commands:
1. Remove `Hidden: true` from command definition
2. Update command help text if needed
3. Update docs/

For FAILED commands:
1. Analyze error (404? Changed URL? Auth issue?)
2. Check API reference for endpoint changes
3. Fix client method or document as permanently broken
