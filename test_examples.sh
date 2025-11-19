#!/bin/bash

# Script to test all .forky files in the fundamentals directory
# Runs the interpreter on each file and checks if it succeeds or fails
# Expected: all files except errors.forky should pass, errors.forky should fail

FUNDAMENTALS_DIR="examples/fundamentals"
FORKY_CMD="go run ."

# ANSI color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counters
total=0
passed=0
failed_expected=0
failed_unexpected=0
passed_unexpected=0

echo "Running tests for .forky files in $FUNDAMENTALS_DIR"
echo "=================================================="
echo

# Find all .forky files in fundamentals
find "$FUNDAMENTALS_DIR" -name "*.forky" | sort | while read -r file; do
    total=$((total + 1))
    filename=$(basename "$file")
    echo -n "Testing $filename... "

    if $FORKY_CMD "$file" > /dev/null 2>&1; then
        if [[ "$file" == *"errors.forky" ]]; then
            echo -e "${YELLOW}✗ PASSED (but expected to FAIL)${NC}"
            passed_unexpected=$((passed_unexpected + 1))
        else
            echo -e "${GREEN}✓ PASSED${NC}"
            passed=$((passed + 1))
        fi
    else
        if [[ "$file" == *"errors.forky" ]]; then
            echo -e "${GREEN}✓ FAILED (as expected)${NC}"
            failed_expected=$((failed_expected + 1))
        else
            echo -e "${RED}✗ FAILED (unexpected)${NC}"
            failed_unexpected=$((failed_unexpected + 1))
        fi
    fi
done

echo
echo "=================================================="
echo "Test Summary:"
echo "Total files tested: $total"
expected=$((passed + failed_expected))
unexpected=$((passed_unexpected + failed_unexpected))
echo -e "Expected results: ${GREEN}$expected${NC}"
echo -e "Unexpected results: ${RED}$unexpected${NC}"

if [ $failed_unexpected -gt 0 ] || [ $passed_unexpected -gt 0 ]; then
    echo -e "${RED}Some tests did not behave as expected.${NC}"
else
    echo -e "${GREEN}All tests behaved as expected!${NC}"
fi