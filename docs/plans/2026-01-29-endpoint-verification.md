# Endpoint Verification Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Create a verification script to test all hidden API endpoints, then fix/re-enable working commands.

**Architecture:** Shell script tests endpoints in safety tiers (read-only → reversible writes), generates markdown report, then we modify command files based on results.

**Tech Stack:** Bash script, eightctl CLI, Go command modifications

---

## Task 1: Create Verification Script Skeleton

**Files:**
- Create: `scripts/verify-endpoints.sh`

**Step 1: Create the script with pre-flight checks**

```bash
#!/usr/bin/env bash
set -euo pipefail

# Endpoint Verification Script for eightctl
# Tests hidden commands and generates a report

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
REPORT_FILE="$PROJECT_ROOT/docs/plans/endpoint-verification-$(date +%Y-%m-%d-%H%M%S).md"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Results arrays
declare -a PASSED=()
declare -a FAILED=()
declare -a SKIPPED=()

log_pass() { echo -e "${GREEN}✓ PASS${NC}: $1"; PASSED+=("$1|$2"); }
log_fail() { echo -e "${RED}✗ FAIL${NC}: $1 - $2"; FAILED+=("$1|$2"); }
log_skip() { echo -e "${YELLOW}⊘ SKIP${NC}: $1 - $2"; SKIPPED+=("$1|$2"); }
log_info() { echo -e "  INFO: $1"; }

# Pre-flight checks
preflight() {
    echo "=== Pre-flight Checks ==="

    # Check eightctl exists
    if ! command -v eightctl &> /dev/null; then
        echo "ERROR: eightctl not found. Run 'go install ./cmd/eightctl' first."
        exit 1
    fi

    # Check authentication
    echo "Checking authentication..."
    if ! eightctl whoami &> /dev/null; then
        echo "ERROR: Authentication failed. Configure credentials first."
        exit 1
    fi

    echo "Authentication OK"
    echo ""
}

# Test a command, return 0 on success, 1 on failure
test_cmd() {
    local name="$1"
    local cmd="$2"
    local output
    local exit_code

    output=$(eval "$cmd" 2>&1) && exit_code=0 || exit_code=$?

    if [[ $exit_code -eq 0 ]]; then
        # Check for "Cannot GET" in output (API returns 200 but with error body)
        if echo "$output" | grep -qi "cannot get\|not found\|404"; then
            log_fail "$name" "API returned error in body"
            return 1
        fi
        log_pass "$name" "${output:0:50}..."
        return 0
    else
        log_fail "$name" "${output:0:80}"
        return 1
    fi
}

echo "Endpoint Verification Script"
echo "============================="
echo ""
```

**Step 2: Make executable and test skeleton runs**

Run:
```bash
chmod +x scripts/verify-endpoints.sh
./scripts/verify-endpoints.sh
```

Expected: Script runs preflight, confirms auth works, then exits (no tests yet)

**Step 3: Commit**

```bash
git add scripts/verify-endpoints.sh
git commit -m "feat(scripts): add endpoint verification script skeleton"
```

---

## Task 2: Add READ-ONLY Tests

**Files:**
- Modify: `scripts/verify-endpoints.sh`

**Step 1: Add read-only test functions**

Append after the `test_cmd` function:

```bash
# ==========================================
# READ-ONLY TESTS (Safe to run freely)
# ==========================================

test_readonly() {
    echo "=== READ-ONLY Endpoint Tests ==="
    echo ""

    # Alarm
    echo "-- Alarm --"
    test_cmd "alarm list" "eightctl alarm list --output json" || true

    # Audio
    echo "-- Audio --"
    test_cmd "audio tracks" "eightctl audio tracks --output json" || true
    test_cmd "audio categories" "eightctl audio categories --output json" || true
    test_cmd "audio state" "eightctl audio state --output json" || true
    test_cmd "audio favorites" "eightctl audio favorites --output json" || true

    # Autopilot
    echo "-- Autopilot --"
    test_cmd "autopilot details" "eightctl autopilot details --output json" || true
    test_cmd "autopilot history" "eightctl autopilot history --output json" || true
    test_cmd "autopilot recap" "eightctl autopilot recap --output json" || true

    # Base
    echo "-- Base --"
    test_cmd "base info" "eightctl base info --output json" || true
    test_cmd "base presets" "eightctl base presets --output json" || true

    # Temp Modes
    echo "-- Temp Modes --"
    test_cmd "tempmode nap status" "eightctl tempmode nap status --output json" || true
    test_cmd "tempmode hotflash status" "eightctl tempmode hotflash status --output json" || true
    test_cmd "tempmode events" "eightctl tempmode events --output json" || true

    # Household
    echo "-- Household --"
    test_cmd "household summary" "eightctl household summary --output json" || true
    test_cmd "household schedule" "eightctl household schedule --output json" || true
    test_cmd "household current-set" "eightctl household current-set --output json" || true
    test_cmd "household invitations" "eightctl household invitations --output json" || true
    test_cmd "household devices" "eightctl household devices --output json" || true
    test_cmd "household users" "eightctl household users --output json" || true
    test_cmd "household guests" "eightctl household guests --output json" || true

    # Metrics (hidden subcommands)
    echo "-- Metrics --"
    test_cmd "metrics intervals" "eightctl metrics intervals --output json" || true
    test_cmd "metrics summary" "eightctl metrics summary --output json" || true
    test_cmd "metrics aggregate" "eightctl metrics aggregate --output json" || true
    test_cmd "metrics insights" "eightctl metrics insights --output json" || true

    # Device (hidden subcommands)
    echo "-- Device --"
    test_cmd "device owner" "eightctl device owner --output json" || true
    test_cmd "device warranty" "eightctl device warranty --output json" || true
    test_cmd "device priming-tasks" "eightctl device priming-tasks --output json" || true
    test_cmd "device priming-schedule" "eightctl device priming-schedule --output json" || true

    # Travel
    echo "-- Travel --"
    test_cmd "travel trips" "eightctl travel trips --output json" || true
    test_cmd "travel airport-search JFK" "eightctl travel airport-search JFK --output json" || true

    # Standalone hidden commands
    echo "-- Other --"
    test_cmd "feats" "eightctl feats --output json" || true
    test_cmd "tracks" "eightctl tracks --output json" || true
    test_cmd "schedule list" "eightctl schedule list --output json" || true

    echo ""
}
```

**Step 2: Call the test function from main**

Add at the end of the script:

```bash
# Main execution
main() {
    preflight
    test_readonly
    echo "=== Read-only tests complete ==="
}

main "$@"
```

**Step 3: Run to verify structure (will fail since commands are hidden)**

Run:
```bash
./scripts/verify-endpoints.sh 2>&1 | head -50
```

Expected: Commands fail with "unknown command" (since they're hidden) - this confirms our test targets are correct

**Step 4: Commit**

```bash
git add scripts/verify-endpoints.sh
git commit -m "feat(scripts): add read-only endpoint tests"
```

---

## Task 3: Add Reversible Write Tests

**Files:**
- Modify: `scripts/verify-endpoints.sh`

**Step 1: Add reversible write tests**

Add after `test_readonly` function:

```bash
# ==========================================
# REVERSIBLE WRITE TESTS (Require confirmation)
# ==========================================

test_reversible_writes() {
    echo "=== REVERSIBLE Write Tests ==="
    echo ""
    echo "These tests will:"
    echo "  - Create and delete test alarms"
    echo "  - Dismiss all alarms (user confirmed OK)"
    echo "  - Play/pause audio"
    echo "  - Toggle nap/hotflash modes"
    echo "  - Create and delete test travel trip"
    echo ""

    read -p "Proceed with write tests? (y/N) " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_skip "reversible writes" "User declined"
        return
    fi

    # Alarm create/delete pair
    echo "-- Alarm Create/Delete --"
    local alarm_id
    alarm_id=$(eightctl alarm create --time "23:59" --days "mon" --output json 2>/dev/null | jq -r '.id // empty')
    if [[ -n "$alarm_id" ]]; then
        log_pass "alarm create" "Created alarm $alarm_id"
        if eightctl alarm delete "$alarm_id" 2>/dev/null; then
            log_pass "alarm delete" "Deleted alarm $alarm_id"
        else
            log_fail "alarm delete" "Failed to delete $alarm_id"
        fi
    else
        log_fail "alarm create" "Could not create alarm"
    fi

    # Alarm dismiss-all (user said OK)
    echo "-- Alarm Dismiss All --"
    test_cmd "alarm dismiss-all" "eightctl alarm dismiss-all" || true

    # Audio play/pause
    echo "-- Audio Play/Pause --"
    # Get first track ID
    local track_id
    track_id=$(eightctl audio tracks --output json 2>/dev/null | jq -r '.[0].id // empty')
    if [[ -n "$track_id" ]]; then
        if eightctl audio play "$track_id" 2>/dev/null; then
            log_pass "audio play" "Started track $track_id"
            sleep 2
            if eightctl audio pause 2>/dev/null; then
                log_pass "audio pause" "Paused playback"
            else
                log_fail "audio pause" "Failed to pause"
            fi
        else
            log_fail "audio play" "Failed to play track"
        fi
    else
        log_skip "audio play/pause" "No tracks available"
    fi

    # Nap mode on/off
    echo "-- Nap Mode Toggle --"
    if eightctl tempmode nap on 2>/dev/null; then
        log_pass "tempmode nap on" "Activated nap mode"
        sleep 2
        if eightctl tempmode nap off 2>/dev/null; then
            log_pass "tempmode nap off" "Deactivated nap mode"
        else
            log_fail "tempmode nap off" "Failed to deactivate"
        fi
    else
        log_fail "tempmode nap on" "Failed to activate"
    fi

    # Hot flash mode on/off
    echo "-- Hot Flash Mode Toggle --"
    if eightctl tempmode hotflash on 2>/dev/null; then
        log_pass "tempmode hotflash on" "Activated hot flash mode"
        sleep 2
        if eightctl tempmode hotflash off 2>/dev/null; then
            log_pass "tempmode hotflash off" "Deactivated hot flash mode"
        else
            log_fail "tempmode hotflash off" "Failed to deactivate"
        fi
    else
        log_fail "tempmode hotflash on" "Failed to activate"
    fi

    # Travel trip create/delete
    echo "-- Travel Trip Create/Delete --"
    local trip_id
    trip_id=$(eightctl travel create-trip --name "Test Trip" --output json 2>/dev/null | jq -r '.id // empty')
    if [[ -n "$trip_id" ]]; then
        log_pass "travel create-trip" "Created trip $trip_id"
        if eightctl travel delete-trip "$trip_id" 2>/dev/null; then
            log_pass "travel delete-trip" "Deleted trip $trip_id"
        else
            log_fail "travel delete-trip" "Failed to delete $trip_id"
        fi
    else
        log_fail "travel create-trip" "Could not create trip"
    fi

    echo ""
}
```

**Step 2: Update main to call write tests**

Replace main function:

```bash
# Main execution
main() {
    preflight
    test_readonly

    echo ""
    read -p "Continue to reversible write tests? (y/N) " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        test_reversible_writes
    else
        log_skip "reversible writes" "User skipped"
    fi

    echo "=== All tests complete ==="
}

main "$@"
```

**Step 3: Commit**

```bash
git add scripts/verify-endpoints.sh
git commit -m "feat(scripts): add reversible write tests with safety prompts"
```

---

## Task 4: Add Report Generation

**Files:**
- Modify: `scripts/verify-endpoints.sh`

**Step 1: Add report generation function**

Add before `main`:

```bash
# ==========================================
# REPORT GENERATION
# ==========================================

generate_report() {
    echo "Generating report: $REPORT_FILE"

    cat > "$REPORT_FILE" << EOF
# Endpoint Verification Report

**Date:** $(date +"%Y-%m-%d %H:%M:%S")
**Total Tested:** $((${#PASSED[@]} + ${#FAILED[@]}))
**Passed:** ${#PASSED[@]}
**Failed:** ${#FAILED[@]}
**Skipped:** ${#SKIPPED[@]}

---

## ✓ PASSED (${#PASSED[@]})

| Command | Notes |
|---------|-------|
EOF

    for entry in "${PASSED[@]}"; do
        IFS='|' read -r cmd notes <<< "$entry"
        echo "| \`$cmd\` | ${notes:0:60} |" >> "$REPORT_FILE"
    done

    cat >> "$REPORT_FILE" << EOF

## ✗ FAILED (${#FAILED[@]})

| Command | Error |
|---------|-------|
EOF

    for entry in "${FAILED[@]}"; do
        IFS='|' read -r cmd error <<< "$entry"
        # Escape pipe characters in error
        error="${error//|/\\|}"
        echo "| \`$cmd\` | ${error:0:60} |" >> "$REPORT_FILE"
    done

    cat >> "$REPORT_FILE" << EOF

## ⊘ SKIPPED (${#SKIPPED[@]})

| Command | Reason |
|---------|--------|
EOF

    for entry in "${SKIPPED[@]}"; do
        IFS='|' read -r cmd reason <<< "$entry"
        echo "| \`$cmd\` | $reason |" >> "$REPORT_FILE"
    done

    cat >> "$REPORT_FILE" << EOF

---

## Next Steps

For PASSED commands:
1. Remove \`Hidden: true\` from command definition
2. Update command help text if needed
3. Update docs/

For FAILED commands:
1. Analyze error (404? Changed URL? Auth issue?)
2. Check API reference for endpoint changes
3. Fix client method or document as permanently broken
EOF

    echo ""
    echo "Report saved to: $REPORT_FILE"
}
```

**Step 2: Update main to generate report**

```bash
# Main execution
main() {
    preflight
    test_readonly

    echo ""
    read -p "Continue to reversible write tests? (y/N) " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        test_reversible_writes
    else
        log_skip "reversible writes" "User skipped"
    fi

    echo ""
    echo "=== All tests complete ==="
    echo ""

    generate_report
}

main "$@"
```

**Step 3: Test report generation with dummy data**

Run:
```bash
./scripts/verify-endpoints.sh
```

Expected: Script runs all tests, generates markdown report

**Step 4: Commit**

```bash
git add scripts/verify-endpoints.sh
git commit -m "feat(scripts): add markdown report generation"
```

---

## Task 5: Temporarily Unhide Commands for Testing

**Files:**
- Modify: `internal/cmd/alarm.go`
- Modify: `internal/cmd/audio.go`
- Modify: `internal/cmd/autopilot.go`
- Modify: `internal/cmd/base.go`
- Modify: `internal/cmd/temp_modes.go`
- Modify: `internal/cmd/household.go`
- Modify: `internal/cmd/travel.go`
- Modify: `internal/cmd/metrics.go`
- Modify: `internal/cmd/device.go`
- Modify: `internal/cmd/schedule.go`
- Modify: `internal/cmd/feats.go`
- Modify: `internal/cmd/tracks.go`

**Step 1: Comment out Hidden: true in all files**

For each file, change:
```go
Hidden: true, // reason
```
To:
```go
// Hidden: true, // TESTING - reason
```

**Step 2: Rebuild eightctl**

```bash
go build -o bin/eightctl ./cmd/eightctl
```

**Step 3: Verify commands are now visible**

```bash
./bin/eightctl --help | grep -E "alarm|audio|autopilot"
```

Expected: alarm, audio, autopilot appear in help

**Step 4: Commit (temporary)**

```bash
git add internal/cmd/
git commit -m "test: temporarily unhide all commands for endpoint verification"
```

---

## Task 6: Run Full Verification

**Step 1: Run the verification script**

```bash
PATH="./bin:$PATH" ./scripts/verify-endpoints.sh
```

**Step 2: Review results interactively**

- Answer prompts for write tests
- Watch for pass/fail output
- Note any unexpected errors

**Step 3: Review generated report**

```bash
cat docs/plans/endpoint-verification-*.md
```

---

## Task 7: Restore Hidden Status Based on Results

**Files:**
- Modify: All command files based on report

**Step 1: For each PASSED command group**

Remove the `Hidden: true` line entirely (don't just comment it out).

Update the comment to indicate it's working:
```go
var alarmCmd = &cobra.Command{
    Use:   "alarm",
    Short: "Manage alarms",
    // Verified working 2026-01-29
}
```

**Step 2: For each FAILED command group**

Restore `Hidden: true` with updated comment:
```go
var autopilotCmd = &cobra.Command{
    Use:    "autopilot",
    Short:  "Autopilot settings",
    Hidden: true, // Verified broken 2026-01-29: 404 Not Found
}
```

**Step 3: Rebuild and verify**

```bash
go build -o bin/eightctl ./cmd/eightctl
./bin/eightctl --help
```

**Step 4: Run tests**

```bash
go test ./...
```

**Step 5: Commit per command group**

```bash
git add internal/cmd/alarm.go
git commit -m "fix(cmd): re-enable alarm commands after API verification"
```

---

## Task 8: Update Documentation

**Files:**
- Modify: `docs/api-reference.md` (update endpoint status)
- Modify: `README.md` if command list changed

**Step 1: Update API reference with verification results**

Add status indicators to documented endpoints.

**Step 2: Update README command list**

If commands were re-enabled, add them to the command list in README.

**Step 3: Commit**

```bash
git add docs/ README.md
git commit -m "docs: update command status after endpoint verification"
```

---

## Task 9: Final Verification and Cleanup

**Step 1: Run full test suite**

```bash
make test
make lint
```

**Step 2: Verify all re-enabled commands work**

```bash
./bin/eightctl alarm list
./bin/eightctl audio tracks
# etc for each re-enabled command
```

**Step 3: Clean up verification report**

Move report to permanent location or delete if no longer needed.

**Step 4: Final commit**

```bash
git add .
git commit -m "chore: endpoint verification complete"
```

---

## Summary

| Task | Description | Estimated Steps |
|------|-------------|-----------------|
| 1 | Script skeleton | 3 |
| 2 | Read-only tests | 4 |
| 3 | Write tests | 3 |
| 4 | Report generation | 4 |
| 5 | Unhide commands | 4 |
| 6 | Run verification | 3 |
| 7 | Restore based on results | 5 |
| 8 | Update docs | 3 |
| 9 | Final cleanup | 4 |
