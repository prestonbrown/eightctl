# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

eightctl is a Go CLI for controlling Eight Sleep Pods. It communicates with the undocumented Eight Sleep API using OAuth credentials extracted from the mobile app, and includes a daemon for scheduled automations.

## Common Commands

```bash
# Build, test, lint
make fmt                    # Format with gofumpt
make lint                   # Run golangci-lint
make test                   # Run go test ./...

# Build locally
go build -o bin/eightctl ./cmd/eightctl

# Install
go install ./cmd/eightctl

# Run with verbose debugging
eightctl --verbose <command>
```

## Architecture

### Package Structure

- `cmd/eightctl/main.go` - Entry point, calls `cmd.Execute()`
- `internal/cmd/` - Cobra commands; each command follows pattern: `requireAuthFields() → client.New() → API call → output.Print()`
- `internal/client/` - Eight Sleep API client; domain-specific files (alarms.go, autopilot.go, etc.) contain related endpoints
- `internal/config/` - Viper-based config (YAML/env/flags)
- `internal/daemon/` - YAML schedule parser and minute-tick execution loop
- `internal/output/` - Table/JSON/CSV output formatting
- `internal/tokencache/` - OS keyring token persistence (macOS Keychain, Linux SecretService, Windows Credential Manager)

### Configuration Priority (highest first)

1. Command-line flags (`--email`, `--password`)
2. Environment variables (`EIGHTCTL_EMAIL`, `EIGHTCTL_PASSWORD`)
3. Config file (`~/.config/eightctl/config.yaml`)

### API Client Details

- Base URL: `https://client-api.8slp.net/v1`
- Auth URL: `https://auth-api.8slp.net/v1/tokens`
- Default OAuth credentials are built-in (from Android APK v7.41.66)
- HTTP/2 disabled (workaround for Eight Sleep frontend)
- 20-second timeout, auto-retry on 429
- See `docs/api-reference.md` for complete endpoint documentation

### Adding a New Command

1. Create `internal/cmd/mycommand.go` with `var myCmd = &cobra.Command{}`
2. Add `rootCmd.AddCommand(myCmd)` in that file's `init()`
3. Implement `RunE` using existing patterns

### Adding an API Endpoint

1. Add method to `Client` in `internal/client/` (use domain-specific file or create new one)
2. Use `c.do()` for HTTP requests (handles auth headers, JSON decoding)

### Testing

- Standard Go testing with `httptest.NewServer` for API mocks
- Token cache tests use `SetOpenKeyringForTest()` for isolated keyring
- See `internal/client/eightsleep_test.go` for mock server pattern

## Key Patterns

- Commands return `[]map[string]any` rows; `output.Print()` handles all format rendering
- Config values via `viper.GetString()`, `viper.GetBool()`
- Token caching is transparent: first auth caches token, subsequent calls reuse it
- Daemon uses simple minute-tick loop with execute-once-per-day semantics

## Git Worktrees

Use `.worktrees/` directory for isolated workspaces. This directory is gitignored.
