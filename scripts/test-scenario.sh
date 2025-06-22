#!/bin/bash

# Test a specific scenario with automatic reset

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
TESTDATA_DIR="$PROJECT_ROOT/rules/testdata"
WORKING_DIR="$TESTDATA_DIR/working"

if [ $# -eq 0 ]; then
    echo "Usage: $0 <scenario_name>"
    echo ""
    echo "Available scenarios:"
    if [ -d "$TESTDATA_DIR/templates" ]; then
        for scenario in "$TESTDATA_DIR/templates"/*; do
            if [ -d "$scenario" ]; then
                echo "  - $(basename "$scenario")"
            fi
        done
    fi
    exit 1
fi

SCENARIO_NAME="$1"
SCENARIO_DIR="$WORKING_DIR/$SCENARIO_NAME"

# Reset scenarios first
echo "Resetting scenarios..."
"$SCRIPT_DIR/reset-scenarios.sh"

# Check if scenario exists
if [ ! -d "$SCENARIO_DIR" ]; then
    echo "Error: Scenario '$SCENARIO_NAME' not found"
    exit 1
fi

echo ""
echo "=== Testing scenario: $SCENARIO_NAME ==="
echo ""

# Show initial state
echo "Initial files:"
find "$SCENARIO_DIR" -name "*.tf" -exec echo "  {}" \; -exec cat {} \; -exec echo "" \;

echo ""
echo "=== Running tflint ==="
cd "$SCENARIO_DIR"

# Check if tflint config exists, if not create basic one
if [ ! -f ".tflint.hcl" ]; then
    cat > .tflint.hcl << EOF
plugin "file-name-is-resource-name" {
  enabled = true
}

rule "terraform_required_version" {
  enabled = false
}

rule "terraform_required_providers" {
  enabled = false
}
EOF
fi

# Run tflint to show issues
echo "Issues found:"
tflint || true

echo ""
echo "=== Running tflint --fix ==="
tflint --fix || true

echo ""
echo "=== Final state ==="
echo "Files after fix:"
find "$SCENARIO_DIR" -name "*.tf" -exec echo "  {}" \; -exec cat {} \; -exec echo "" \;

echo ""
echo "=== Test completed ==="
echo "Scenario directory: $SCENARIO_DIR"
echo "Run '$SCRIPT_DIR/reset-scenarios.sh' to reset for next test"