#!/bin/bash

# Reset test scenarios to initial state for repeatable testing

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
TESTDATA_DIR="$PROJECT_ROOT/rules/testdata"
TEMPLATES_DIR="$TESTDATA_DIR/templates"
WORKING_DIR="$TESTDATA_DIR/working"

echo "Resetting test scenarios..."

# Remove existing working directory
if [ -d "$WORKING_DIR" ]; then
    rm -rf "$WORKING_DIR"
    echo "Removed existing working directory"
fi

# Create fresh working directory
mkdir -p "$WORKING_DIR"

# Copy all templates to working directory
if [ -d "$TEMPLATES_DIR" ]; then
    cp -r "$TEMPLATES_DIR"/* "$WORKING_DIR"/
    echo "Copied templates to working directory"
    
    # List available scenarios
    echo ""
    echo "Available scenarios:"
    for scenario in "$WORKING_DIR"/*; do
        if [ -d "$scenario" ]; then
            scenario_name=$(basename "$scenario")
            echo "  - $scenario_name"
        fi
    done
    
    echo ""
    echo "Test scenarios reset successfully!"
    echo "You can now run 'tflint --fix' in any scenario under rules/testdata/working/"
else
    echo "Error: Templates directory not found at $TEMPLATES_DIR"
    exit 1
fi