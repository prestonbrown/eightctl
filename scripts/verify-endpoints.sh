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

preflight
