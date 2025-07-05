# Testing Guide

This document explains how to test the TFLint ruleset, especially the `--fix` functionality.

## Quick Testing

### Test All Scenarios (Recommended)
```bash
make test-all-scenarios
```
This will automatically:
1. Reset all scenarios to initial state
2. Run each scenario showing before/after states
3. Display results for easy verification

### Test Individual Scenarios
```bash
make test-scenario1  # Variable in main.tf → should move to variables.tf
make test-scenario2  # Invalid blocks in variables.tf → should show errors
make test-scenario3  # Append to existing variables.tf → should merge correctly
```

### Reset Scenarios Manually
```bash
make reset-scenarios
```
Restores all test scenarios to initial state for repeatable testing.

## Manual Testing

### Setup
1. Reset scenarios: `make reset-scenarios`
2. Navigate to working directory: `cd rules/testdata/working/scenario1_variable_in_main`
3. Build and install plugin: `make install`

### Test Workflow
1. Check initial state: `ls -la *.tf`
2. Run TFLint: `tflint`
3. Apply fixes: `tflint --fix`
4. Verify results: `ls -la *.tf` and `cat variables.tf`
5. Reset for next test: `make reset-scenarios`

## Test Scenarios

### Scenario 1: Variable in Wrong File
- **File**: `main.tf` contains variable blocks
- **Expected**: `tflint --fix` moves variables to `variables.tf`
- **Test**: `make test-scenario1`

### Scenario 2: Invalid Blocks in variables.tf
- **File**: `variables.tf` contains resource/output blocks
- **Expected**: TFLint reports errors (no fix available)
- **Test**: `make test-scenario2`

### Scenario 3: Append to Existing variables.tf
- **Files**: Existing `variables.tf` + `main.tf` with new variable
- **Expected**: New variable appended to existing `variables.tf`
- **Test**: `make test-scenario3`

## Directory Structure

```
rules/testdata/
├── templates/      # Original templates (never modified)
│   ├── scenario1_variable_in_main/
│   ├── scenario2_mixed_blocks_in_variables/
│   └── scenario3_existing_variables_tf/
├── working/        # Test execution area (gets modified by --fix)
│   ├── scenario1_variable_in_main/
│   ├── scenario2_mixed_blocks_in_variables/
│   └── scenario3_existing_variables_tf/
└── demo/          # File operation demos
    └── variables.tf
```

## Development Workflow

1. **Code Changes**: Modify rule implementation
2. **Quick Test**: `make test-all-scenarios`
3. **Manual Verification**: Navigate to `working/` directories
4. **Reset**: `make reset-scenarios` between tests
5. **Unit Tests**: `make test` for comprehensive testing

## Adding New Test Scenarios

1. Add new scenario to `templates/` directory
2. Update `scripts/reset-scenarios.sh` if needed
3. Add new Make target following existing pattern
4. Document the scenario in this file

## Troubleshooting

### Plugin Not Found
- Run `make install` to build and install plugin
- Check `~/.tflint.d/plugins/` for plugin binary

### Permission Errors
- Ensure scripts have execute permissions: `chmod +x scripts/*.sh`

### Files Not Resetting
- Manually clean: `make clean-testdata`
- Recreate templates: `make create-test-scenarios`

## File Operations Testing

For lower-level file manipulation testing:
```bash
make demo-file-ops  # Test BlockMover functionality directly
```