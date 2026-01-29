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

echo "Endpoint Verification Script"
echo "============================="
echo ""

# Main execution
main() {
    local readonly_only=false
    if [[ "${1:-}" == "--readonly-only" ]]; then
        readonly_only=true
    fi

    preflight
    test_readonly

    if [[ "$readonly_only" == "false" ]]; then
        echo ""
        read -p "Continue to reversible write tests? (y/N) " -n 1 -r
        echo ""
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            test_reversible_writes
        else
            log_skip "reversible writes" "User skipped"
        fi
    else
        log_skip "reversible writes" "Skipped (--readonly-only)"
    fi

    echo ""
    echo "=== All tests complete ==="
    echo ""

    generate_report
}

main "$@"
