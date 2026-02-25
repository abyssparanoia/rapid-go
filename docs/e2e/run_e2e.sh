#!/bin/bash

# E2E Test Script
# This script performs end-to-end testing of the full authentication and API flow

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="${BASE_URL:-http://localhost:8080}"
CLI_PATH="${CLI_PATH:-.bin/app-cli}"

# Generate unique identifiers for this test run
TEST_ID="$(date +%s)-${RANDOM}"
ADMIN_EMAIL="${ADMIN_EMAIL:-e2e-admin-${TEST_ID}@example.com}"
ADMIN_DISPLAY_NAME="${ADMIN_DISPLAY_NAME:-E2E Admin ${TEST_ID}}"
STAFF_EMAIL="${STAFF_EMAIL:-e2e-staff-${TEST_ID}@example.com}"
STAFF_DISPLAY_NAME="${STAFF_DISPLAY_NAME:-E2E Staff ${TEST_ID}}"
TENANT_NAME="${TENANT_NAME:-E2E Test Tenant ${TEST_ID}}"

# Test results tracking
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Shared state (populated by steps)
ADMIN_ID=""
ADMIN_AUTH_UID=""
ADMIN_PASSWORD=""
ADMIN_TOKEN=""
TENANT_ID=""
STAFF_ID=""
STAFF_AUTH_UID=""
STAFF_PASSWORD=""
STAFF_TOKEN=""

# Helper functions
print_step() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
    PASSED_TESTS=$((PASSED_TESTS + 1))
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
    FAILED_TESTS=$((FAILED_TESTS + 1))
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

print_result() {
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}Test Results${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo -e "Total Tests: ${TOTAL_TESTS}"
    echo -e "${GREEN}Passed: ${PASSED_TESTS}${NC}"
    if [ $FAILED_TESTS -gt 0 ]; then
        echo -e "${RED}Failed: ${FAILED_TESTS}${NC}"
        exit 1
    else
        echo -e "${GREEN}All tests passed!${NC}"
    fi
}

# Source all step files
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
for step_file in "$SCRIPT_DIR"/steps/step_*.sh; do
    # shellcheck source=/dev/null
    source "$step_file"
done

# Main execution
main() {
    echo ""
    echo -e "${GREEN}=================================${NC}"
    echo -e "${GREEN}E2E Test${NC}"
    echo -e "${GREEN}=================================${NC}"
    echo ""

    step_01_check_prerequisites
    step_02_health_check
    step_03_create_admin
    step_04_get_admin_token
    step_05_test_admin_api
    step_06_create_tenant
    step_07_create_staff
    step_08_get_staff_token
    step_09_test_staff_api
    step_10_test_staff_update_me
    step_11_test_staff_tenant_api
    step_12_staff_signup

    print_result
}

# Run main function
main
