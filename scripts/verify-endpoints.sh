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

echo "Endpoint Verification Script"
echo "============================="
echo ""

# Main execution
main() {
    preflight
    test_readonly
    echo "=== Read-only tests complete ==="
}

main "$@"
