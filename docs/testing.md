# Testing Guide

Guide to testing infrastructure and patterns in eightctl.

## Running Tests

```bash
# Run all tests
go test ./...

# Run all tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./internal/client/...

# Run a specific test by name
go test -run TestAlarm ./internal/client/...

# Run tests matching a pattern
go test -run "TestAlarm.*" ./internal/client/...

# Run tests in parallel (default is GOMAXPROCS)
go test -parallel 4 ./...

# Run with race detector (slower, catches data races)
go test -race ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Test File Organization

Tests live alongside source files with `_test.go` suffix:

```
internal/client/
├── eightsleep.go      → eightsleep_test.go  (core client, auth)
├── alarms.go          → alarms_test.go
├── audio.go           → audio_test.go
├── autopilot.go       → autopilot_test.go
├── base.go            → base_test.go
├── device.go          → device_test.go
├── household.go       → household_test.go
├── metrics.go         → metrics_test.go
├── presence.go        → presence_test.go
├── schedules.go       → schedules_test.go
├── tempmodes.go       → tempmodes_test.go
└── travel.go          → travel_test.go
```

## Test Patterns

### HTTP API Mocking

Use `httptest.NewServer` to mock Eight Sleep API responses:

```go
func TestMyEndpoint(t *testing.T) {
    mux := http.NewServeMux()
    mux.HandleFunc("/users/uid-123/endpoint", func(w http.ResponseWriter, r *http.Request) {
        // Verify request
        if r.Method != http.MethodGet {
            t.Errorf("expected GET, got %s", r.Method)
        }

        // Return mock response
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"data": "value"}`))
    })
    srv := httptest.NewServer(mux)
    defer srv.Close()

    // Create client pointing to test server
    c := New("email", "pass", "uid-123", "", "")
    c.BaseURL = srv.URL
    c.token = "t"
    c.tokenExp = time.Now().Add(time.Hour)
    c.HTTP = srv.Client()

    // Call method and verify
    result, err := c.MyMethod(context.Background())
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    // assertions...
}
```

### Capturing Request Bodies

Verify what the client sends:

```go
func TestCreateSomething(t *testing.T) {
    var capturedBody map[string]any

    mux := http.NewServeMux()
    mux.HandleFunc("/endpoint", func(w http.ResponseWriter, r *http.Request) {
        if err := json.NewDecoder(r.Body).Decode(&capturedBody); err != nil {
            t.Fatal(err)
        }
        w.WriteHeader(http.StatusCreated)
    })
    // ... setup client ...

    c.Create(context.Background(), input)

    if capturedBody["field"] != "expected" {
        t.Errorf("expected field=expected, got %v", capturedBody["field"])
    }
}
```

### Testing Error Conditions

```go
func TestNotFound(t *testing.T) {
    mux := http.NewServeMux()
    mux.HandleFunc("/endpoint", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"error": "not found"}`))
    })
    // ... setup ...

    _, err := c.Get(context.Background(), "missing-id")
    if err == nil {
        t.Fatal("expected error for 404")
    }
}
```

### Table-Driven Tests

Test multiple cases efficiently:

```go
func TestTemperatureRange(t *testing.T) {
    tests := []struct {
        name    string
        level   int
        wantErr bool
    }{
        {"min valid", -100, false},
        {"max valid", 100, false},
        {"below min", -101, true},
        {"above max", 101, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := c.SetTemperature(context.Background(), tt.level)
            if (err != nil) != tt.wantErr {
                t.Errorf("SetTemperature(%d) error = %v, wantErr %v",
                    tt.level, err, tt.wantErr)
            }
        })
    }
}
```

### Token Cache Testing

Use isolated keyring for token cache tests:

```go
func TestTokenCache(t *testing.T) {
    // Use test keyring instead of real OS keyring
    tokencache.SetOpenKeyringForTest()
    defer tokencache.ResetOpenKeyring()

    // Test token operations...
}
```

## Assertions

Go's standard library doesn't include assertions. Use explicit checks:

```go
// Check equality
if got != want {
    t.Errorf("got %v, want %v", got, want)
}

// Check error occurred
if err == nil {
    t.Fatal("expected error")
}

// Check error message
if !strings.Contains(err.Error(), "expected text") {
    t.Errorf("error should contain 'expected text', got: %v", err)
}

// Check nil
if result == nil {
    t.Fatal("result should not be nil")
}

// Check slice length
if len(items) != 3 {
    t.Errorf("expected 3 items, got %d", len(items))
}
```

## Test Helpers

Mark functions as helpers so failures report the caller's line:

```go
func setupTestClient(t *testing.T) (*httptest.Server, *Client) {
    t.Helper()  // Makes failures point to caller

    mux := http.NewServeMux()
    // setup...
    srv := httptest.NewServer(mux)

    c := New("email", "pass", "", "", "")
    c.BaseURL = srv.URL
    c.token = "t"
    c.tokenExp = time.Now().Add(time.Hour)
    c.HTTP = srv.Client()

    return srv, c
}

func TestSomething(t *testing.T) {
    srv, c := setupTestClient(t)
    defer srv.Close()
    // test...
}
```

## Skipping Tests

Skip tests conditionally:

```go
func TestRequiresNetwork(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping network test in short mode")
    }
    // test that needs network...
}
```

Run without skipped tests: `go test -short ./...`

## Common Issues

### Race Conditions

If tests fail intermittently, run with race detector:

```bash
go test -race ./...
```

### Parallel Test Isolation

Tests run in parallel by default within a package. Ensure tests don't share state:

```go
func TestA(t *testing.T) {
    t.Parallel()  // Explicitly allow parallel execution
    // Use local variables, not package globals
}
```

### Cleanup

Always close test servers:

```go
srv := httptest.NewServer(mux)
defer srv.Close()  // Don't forget this
```

## IDE Integration

Most Go IDEs (VS Code with Go extension, GoLand) support:
- Running individual tests with a click
- Debugging tests with breakpoints
- Viewing coverage inline

## See Also

- [Go testing package docs](https://pkg.go.dev/testing)
- [httptest package docs](https://pkg.go.dev/net/http/httptest)
- [Development Guide](./development.md)
