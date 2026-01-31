# Endpoint Verification & Command Restoration Design

**Date:** 2026-01-29
**Goal:** Verify hidden commands with valid endpoints, fix or re-enable them, and add commands for existing client methods.

## Background

The eightctl CLI has many commands marked as hidden due to reported "Cannot GET" errors. Additionally, the client library has methods with no corresponding CLI commands. This design covers systematically verifying endpoints and restoring functionality.

## Phase A: Verify Hidden Commands

### Approach

Create a verification script (`scripts/verify-endpoints.sh`) that tests all hidden endpoints against a live Eight Sleep device and generates a report.

### Safety Tiers

| Tier | Risk | Strategy | Examples |
|------|------|----------|----------|
| READ | None | Test freely | `alarm list`, `audio tracks` |
| REVERSIBLE | Low | Create → Verify → Delete | `alarm create/delete` pairs |
| IRREVERSIBLE | High | Skip with warning | `household remove-device` |

### Endpoints to Test

**READ-ONLY (safe to test):**
- `alarm list`
- `audio tracks`, `categories`, `state`, `favorites`
- `autopilot details`, `history`, `recap`
- `base info`, `presets`
- `temp-modes nap status`, `hotflash status`, `events`
- `household summary`, `schedule`, `devices`, `users`, `guests`, `invitations`, `current-set`
- `metrics intervals`, `summary`, `aggregate`, `insights`
- `device owner`, `warranty`, `priming-tasks`, `priming-schedule`
- `travel trips`, `plans`, `airport-search`, `flight-status`

**REVERSIBLE WRITES (test with cleanup):**
- `alarm create` → `alarm delete`
- `alarm dismiss-all` (user confirmed OK)
- `audio play` → `audio pause`
- `audio volume` (note current, change, restore)
- `temp-modes nap on` → `nap off`
- `temp-modes hotflash on` → `hotflash off`
- `travel create-trip` → `delete-trip`

**SKIP (irreversible):**
- `household remove-device`, `remove-guest`
- `device set-owner`, `set-peripherals`

### Script Flow

1. **Pre-flight checks**
   - Verify eightctl installed
   - Verify authentication (`eightctl whoami`)
   - Create timestamped report file

2. **Run READ-ONLY tests**
   - Execute each command with `--verbose`
   - Capture exit code and output
   - Log pass/fail to report

3. **Prompt for REVERSIBLE tests**
   - Show what will be tested
   - Require explicit confirmation
   - Run create/delete pairs atomically

4. **Generate summary report**

### Report Format

```markdown
# Endpoint Verification Report
Date: YYYY-MM-DD

## Summary
- Tested: X
- Passed: Y
- Failed: Z
- Skipped: N

## Results

### ✓ PASSED
| Command | Notes |
|---------|-------|
| alarm list | Returns empty array |

### ✗ FAILED
| Command | Error |
|---------|-------|
| autopilot details | 404 Not Found |

### ⊘ SKIPPED
| Command | Reason |
|---------|--------|
| household remove-device | Irreversible |
```

## Phase A Resolution: Handling Results

### For PASSED Endpoints

1. Remove `Hidden: true` from the command definition
2. Verify command works end-to-end
3. Update command-line help text if needed
4. Update `docs/` documentation
5. Commit: `fix(cmd): re-enable <command-group> commands`

### For FAILED Endpoints

1. Analyze error type (404, 401, 500, changed URL)
2. Compare against current API behavior if possible
3. Attempt to fix (update URL, method, payload)
4. If unfixable, document reason and keep hidden
5. Commit fixes: `fix(api): update <endpoint> for current API`

## Phase B: New Commands for Existing Client Methods

After Phase A, add CLI commands for client methods that lack them:

**Priority domains (based on Phase A results):**
- Settings: `tap`, `level-suggestions`, `blanket-recommendations`, `perks`, `referrals`
- Users: `get`, `update`, `update-email`, `password-reset`
- Insights: `llm-insights`, `settings`
- Subscriptions: `list`, `create-temporary`, `redeem`
- Challenges: `list`
- Health: `survey`, `checkpoints`, `upload`

**Implementation approach:**
- Follow existing command patterns
- One command group per commit
- Include docs and help text

## Commit Strategy

- One commit per command group fixed/added
- Example messages:
  - `fix(cmd): re-enable alarm commands after API verification`
  - `feat(cmd): add settings command group`
  - `docs: update command reference for restored commands`
