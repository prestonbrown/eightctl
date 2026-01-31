# eightctl Project Specification

## Purpose

Eight Sleep Pod CLI for macOS/Linux users who want terminal-based pod control, scheduled automations, and metrics export using the same API as the mobile app.

## API Reality

Eight Sleep does not publish a public API. This project reverse-engineers endpoints from the mobile app. API changes can break functionality without warning.

- **Current source**: Android APK v7.41.66 (January 2025)
- **Auth**: OAuth2 password grant with embedded client credentials
- **Rate limiting**: 429 responses with retry; re-auth on 401

See [api-reference.md](./api-reference.md) for complete endpoint documentation.

## Architecture

```
User → CLI (Cobra) → Client → Eight Sleep API
                  ↓
              Config (Viper)
                  ↓
            Token Cache (OS Keyring)
```

**Key packages:**
- `internal/cmd/` - Cobra commands
- `internal/client/` - API client with domain-specific files
- `internal/config/` - Viper configuration (YAML/env/flags)
- `internal/daemon/` - YAML schedule parser and executor
- `internal/output/` - Table/JSON/CSV formatting
- `internal/tokencache/` - OS keyring persistence

## Current State

Many API endpoints discovered from older APK versions no longer work. Commands for broken endpoints are hidden (not removed) pending re-implementation with correct endpoints.

| Category | Status |
|----------|--------|
| Auth, status, on/off/temp | Working |
| Device info, online | Working |
| Sleep trends | Working |
| Alarms, schedules | Hidden (broken) |
| Audio, base, autopilot | Hidden (broken) |
| Household, travel | Hidden (broken) |

See [endpoint-audit.md](./endpoint-audit.md) for detailed status.

## Documentation

| Document | Purpose |
|----------|---------|
| [cli-reference.md](./cli-reference.md) | User-facing command documentation |
| [api-reference.md](./api-reference.md) | Eight Sleep API endpoints |
| [endpoint-audit.md](./endpoint-audit.md) | Broken endpoint tracking |
| [development.md](./development.md) | Contributing and reverse engineering |
| [testing.md](./testing.md) | Testing infrastructure and patterns |

## Prior Work

- [clim8](https://github.com/blacktop/clim8) - Go CLI
- [pyEight](https://github.com/lukas-clarke/pyEight) - Python library (OAuth2 fork)
- [eight_sleep](https://github.com/lukas-clarke/eight_sleep) - Home Assistant integration
- [8sleep-mcp](https://github.com/elizabethtrykin/8sleep-mcp) - MCP server (Node/TS)
