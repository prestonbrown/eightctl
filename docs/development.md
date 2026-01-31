# Development Guide

Guide for contributing to eightctl and reverse-engineering the Eight Sleep API.

## Building

```bash
# Format code
make fmt

# Run linter
make lint

# Run tests
go test ./...

# Build binary
go build -o bin/eightctl ./cmd/eightctl

# Install
go install ./cmd/eightctl
```

## Project Structure

```
eightctl/
├── cmd/eightctl/main.go     # Entry point
├── drivers/
│   └── hubitat/             # Groovy drivers for Hubitat
│       ├── eight-sleep-pod.groovy
│       └── eight-sleep-side.groovy
├── internal/
│   ├── adapter/             # Smart home adapter framework
│   │   ├── adapter.go       # Adapter interface and Command type
│   │   ├── mqtt/            # Home Assistant MQTT adapter
│   │   │   ├── adapter.go
│   │   │   └── discovery.go
│   │   └── hubitat/         # Hubitat HTTP adapter
│   │       └── server.go
│   ├── cmd/                 # Cobra commands
│   ├── client/              # Eight Sleep API client
│   ├── config/              # Viper configuration
│   ├── daemon/              # Schedule daemon
│   ├── model/               # Domain types (Side, PowerState, etc.)
│   ├── output/              # Table/JSON/CSV formatting
│   ├── state/               # State management with caching
│   └── tokencache/          # OS keyring token storage
└── docs/
    ├── api-reference.md
    ├── cli-reference.md
    ├── development.md       # This file
    ├── home-assistant.md    # Home Assistant setup guide
    └── hubitat.md           # Hubitat setup guide
```

## Adding a New Command

1. Create `internal/cmd/mycommand.go`:

```go
package cmd

import "github.com/spf13/cobra"

var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Brief description",
    RunE: func(cmd *cobra.Command, args []string) error {
        if err := requireAuthFields(); err != nil {
            return err
        }
        c, err := client.New(...)
        if err != nil {
            return err
        }
        result, err := c.MyEndpoint(...)
        if err != nil {
            return err
        }
        return output.Print(result, ...)
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

2. Add the API method to `internal/client/`.

## Adding an API Endpoint

1. Find the correct endpoint in [api-reference.md](./api-reference.md)
2. Add a method to `internal/client/` (use domain-specific file):

```go
func (c *Client) MyEndpoint(ctx context.Context, userID string) (*Response, error) {
    var resp Response
    err := c.do(ctx, "GET", "/v1/users/"+userID+"/endpoint", nil, &resp)
    return &resp, err
}
```

## Testing

- Use `httptest.NewServer` for API mocks
- Token cache tests use `SetOpenKeyringForTest()` for isolation
- See [testing.md](./testing.md) for comprehensive patterns and examples

```go
func TestMyEndpoint(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/v1/users/123/endpoint" {
            json.NewEncoder(w).Encode(Response{...})
            return
        }
        http.NotFound(w, r)
    }))
    defer srv.Close()

    c := &Client{baseURL: srv.URL, httpClient: srv.Client()}
    resp, err := c.MyEndpoint(context.Background(), "123")
    // assertions...
}
```

## Smart Home Adapters

eightctl includes adapters for smart home platforms.

### Architecture

```
┌──────────────┐     ┌───────────────┐     ┌─────────────┐
│ Smart Home   │────▶│ Adapter       │────▶│ State       │
│ Platform     │◀────│ (mqtt/hubitat)│◀────│ Manager     │
└──────────────┘     └───────────────┘     └─────────────┘
                            │                     │
                            │                     ▼
                            │              ┌─────────────┐
                            └─────────────▶│ API Client  │
                                          └─────────────┘
```

### Key Components

**adapter.Adapter interface:**
```go
type Adapter interface {
    Start(ctx context.Context) error
    HandleCommand(ctx context.Context, cmd Command) error
    Stop() error
}
```

**state.Manager:**
- Caches device state with configurable TTL
- Provides `GetState()`, `SetTemperature()`, `TurnOn()`, `TurnOff()`
- Supports observers for state change notifications

### Adding a New Adapter

1. Create package `internal/adapter/myplatform/`
2. Implement `adapter.Adapter` interface
3. Use `state.Manager` for state operations
4. Create `internal/cmd/myplatform.go` command
5. Add documentation in `docs/myplatform.md`

### Model Types

Domain types in `internal/model/`:
- `Side` - Left or Right (with JSON marshaling)
- `PowerState` - Off, Smart, Manual
- `UserState` - Per-side state (temperature, presence)
- `DeviceState` - Full device state with both sides

---

# Reverse Engineering the Eight Sleep API

Eight Sleep does not publish a public API. This section documents how to extract endpoint information from the mobile app.

## APK Sources

| Source | URL |
|--------|-----|
| APKPure | https://apkpure.com/eight-sleep/com.eightsleep.eight/versions |
| APKMirror | https://www.apkmirror.com/apk/eight-sleep-inc/eight-sleep/ |

Or extract from a device:

```bash
adb shell pm path com.eightsleep.eight
adb pull /data/app/.../base.apk
```

## Decompilation with JADX

JADX converts Android DEX bytecode to readable Java source.

**Install:**

```bash
# macOS
brew install jadx

# Linux
wget https://github.com/skylot/jadx/releases/latest/download/jadx-x.x.x.zip
unzip jadx-*.zip -d jadx
export PATH=$PATH:$(pwd)/jadx/bin
```

**Use:**

```bash
# GUI (recommended for exploration)
jadx-gui eightsleep.apk

# CLI (for scripted extraction)
jadx -d output_dir eightsleep.apk
```

## What to Search For

### API URLs

```
8slp.net
client-api.8slp.net
auth-api.8slp.net
/v1/
/v2/
```

### Retrofit Annotations

```java
@GET("/users/{userId}/temperature")
@POST("/v1/users/{userId}/alarms")
@PUT("/v2/users/{userId}/routines/{routineId}")
```

### Relevant Classes

```
*Api.java
*ApiService.java
*Repository.java
*Client.java
```

### BuildConfig

Check for API URLs, OAuth client IDs, and feature flags.

## Extraction Script

```bash
#!/bin/bash
APK_DIR="$1"

echo "=== API Base URLs ==="
grep -r "8slp.net" "$APK_DIR" --include="*.java"

echo ""
echo "=== HTTP Endpoints ==="
grep -rE '@(GET|POST|PUT|DELETE|PATCH)\(' "$APK_DIR" --include="*.java" -A1

echo ""
echo "=== URL Patterns ==="
grep -rE '"/v[12]/|"/users/|"/devices/' "$APK_DIR" --include="*.java"
```

## Known Findings

### APK v7.41.66 (January 2025)

OAuth credentials:
- `client_id`: `0894c7f33bb94800a03f1f4df13a4f38`
- `client_secret`: `f0954a3ed5763ba3d06834c73731a32f15f168f47d4f164751275def86db0c76`

See [api-reference.md](./api-reference.md) for complete endpoint documentation.

### API Versioning

Eight Sleep is migrating from v1 to v2 for some endpoints:
- Alarms: `v1/users/{userId}/alarms` and `v2/users/{userId}/alarms`
- Some features may only exist in v2

## Security Notes

- APK decompilation is for interoperability research
- Do not redistribute decompiled code
- OAuth credentials are semi-public (embedded in published apps)
- Respect Eight Sleep's terms of service

## Resources

- [JADX](https://github.com/skylot/jadx) - DEX to Java decompiler
- [APKtool](https://github.com/iBotPeaches/Apktool) - APK resource decoder
- [pyEight](https://github.com/lukas-clarke/pyEight) - Python reference implementation
- [eight_sleep](https://github.com/lukas-clarke/eight_sleep) - Home Assistant integration
