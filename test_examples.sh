#!/bin/bash

# Script to test all .forky files in the fundamentals directory
# Runs the interpreter on each file and checks if it succeeds or fails
# Expected: all files except errors.forky should pass, errors.forky should fail

FUNDAMENTALS_DIR="examples/fundamentals"
FORKY_CMD="go run ."

# Find all .forky files in fundamentals
find "$FUNDAMENTALS_DIR" -name "*.forky" | while read -r file; do
    echo "Testing $file..."
    if $FORKY_CMD "$file" > /dev/null 2>&1; then
        result="PASSED"
        if [[ "$file" == *"errors.forky" ]]; then
            echo "✗ $file: PASSED (but expected to FAIL)"
        else
            echo "✓ $file: PASSED"
        fi
    else
        result="FAILED"
        if [[ "$file" == *"errors.forky" ]]; then
            echo "✓ $file: FAILED (as expected)"
        else
            echo "✗ $file: FAILED (unexpected)"
        fi
    fi
    echo
done