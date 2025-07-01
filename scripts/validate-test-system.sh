#!/bin/bash

# Simple test validation: just run test-scenarios and verify working/answer match
# This confirms that our test system correctly validates results

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "=== Validating Test System ==="
echo "Running test-scenarios and verifying working/answer directories match..."
echo ""

cd "$PROJECT_ROOT"

# Run the actual test scenarios
if ! make test-scenarios; then
    echo "❌ test-scenarios failed!"
    exit 1
fi

echo ""
echo "✅ test-scenarios passed!"
echo ""

# Double-check by manually comparing working vs answer for each scenario
echo "Manual verification of working vs answer directories..."

success=true
for scenario in rules/testdata/templates/*; do
    if [ -d "$scenario" ]; then
        scenario_name=$(basename "$scenario")
        working_dir="rules/testdata/working/$scenario_name"
        answer_dir="rules/testdata/answer/$scenario_name"
        
        if [ -d "$answer_dir" ]; then
            echo "Checking $scenario_name..."
            if diff -r "$working_dir" "$answer_dir" > /dev/null 2>&1; then
                echo "  ✓ $scenario_name: working matches answer"
            else
                echo "  ✗ $scenario_name: working differs from answer"
                echo "  Differences:"
                diff -r "$working_dir" "$answer_dir" | head -10
                success=false
            fi
        else
            echo "  ⚠ $scenario_name: no answer directory found"
        fi
    fi
done

echo ""
if [ "$success" = true ]; then
    echo "🎉 Test system validation PASSED!"
    echo "All working directories match their corresponding answer directories."
else
    echo "❌ Test system validation FAILED!"
    echo "Some working directories differ from their answer directories."
    exit 1
fi