#!/bin/bash

# E2E Test Common Helpers
# This file is sourced by all test files and provides shared functions & variables.

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="${BASE_URL:-http://localhost:8080}"
CLI_PATH="${CLI_PATH:-.bin/app-cli}"

# Log temp file (created by run_e2e.sh via mktemp, cleaned up on exit)
SERVER_LOG_FILE=""

# Log offset: line count recorded at the start of each test step
SERVER_LOG_OFFSET=0

# Generate unique identifiers for this test run
TEST_ID="${TEST_ID:-$(date +%s)-${RANDOM}}"
ADMIN_EMAIL="${ADMIN_EMAIL:-e2e-admin-${TEST_ID}@example.com}"
ADMIN_DISPLAY_NAME="${ADMIN_DISPLAY_NAME:-E2E Admin ${TEST_ID}}"
STAFF_EMAIL="${STAFF_EMAIL:-e2e-staff-${TEST_ID}@example.com}"
STAFF_DISPLAY_NAME="${STAFF_DISPLAY_NAME:-E2E Staff ${TEST_ID}}"
TENANT_NAME="${TENANT_NAME:-E2E Test Tenant ${TEST_ID}}"

# Test results tracking
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Shared state (populated by test steps)
ADMIN_ID=""
ADMIN_AUTH_UID=""
ADMIN_PASSWORD=""
ADMIN_TOKEN=""
TENANT_ID=""
STAFF_ID=""
STAFF_AUTH_UID=""
STAFF_PASSWORD=""
STAFF_TOKEN=""

# PostgreSQL helper: uses local psql CLI if available, falls back to docker exec
POSTGRES_CONTAINER="${POSTGRES_CONTAINER:-rapid-go-postgresql-1}"

run_psql() {
    if command -v psql &> /dev/null; then
        PGPASSWORD=postgres psql -h 127.0.0.1 -U postgres -d maindb -c "$1" 2>/dev/null
    else
        docker exec "$POSTGRES_CONTAINER" psql -U postgres -d maindb -c "$1" 2>/dev/null
    fi
}

# MySQL helper: uses local mysql CLI if available, falls back to docker exec
MYSQL_CONTAINER="${MYSQL_CONTAINER:-rapid-go-mysql-1}"

run_mysql() {
    if command -v mysql &> /dev/null; then
        mysql -h 127.0.0.1 -P 3306 -u root -ppassword maindb -e "$1" 2>/dev/null
    else
        docker exec "$MYSQL_CONTAINER" mysql -u root -ppassword maindb -e "$1" 2>/dev/null
    fi
}

# Reset the database by dropping and recreating the public schema
reset_database() {
    print_step "Resetting Database"
    run_psql "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
    print_info "Database schema dropped and recreated"
}

# Run migrations and sync constant table data
migrate_and_seed() {
    print_step "Running Migrations & Seeding Constants"
    (cd "$REPO_ROOT" && "$CLI_PATH" schema-migration database up) || { echo "migrate up failed"; exit 1; }
    print_info "Migrations applied"
    (cd "$REPO_ROOT" && "$CLI_PATH" schema-migration database sync-constants) || { echo "sync-constants failed"; exit 1; }
    print_info "Constants synced"
}

# Helper functions
print_step() {
    # Record current log line count so failure output starts from this test
    if [ -f "$SERVER_LOG_FILE" ]; then
        SERVER_LOG_OFFSET=$(wc -l < "$SERVER_LOG_FILE")
    fi
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

wait_for_health() {
    local max_wait=30
    local interval=2
    local elapsed=0
    print_info "Waiting for server to become healthy (max ${max_wait}s)..."
    while [ $elapsed -lt $max_wait ]; do
        if curl -s -f "$BASE_URL/" > /dev/null 2>&1; then
            print_info "Server is healthy (${elapsed}s elapsed)"
            return 0
        fi
        sleep $interval
        elapsed=$((elapsed + interval))
    done
    echo -e "${RED}Server did not become healthy within ${max_wait}s${NC}"
    exit 1
}

print_server_logs() {
    local start=$((SERVER_LOG_OFFSET + 1))
    if [ -f "$SERVER_LOG_FILE" ] && [ -s "$SERVER_LOG_FILE" ]; then
        local relevant
        relevant=$(tail -n +"$start" "$SERVER_LOG_FILE")
        if [ -n "$relevant" ]; then
            echo ""
            echo -e "${RED}========================================${NC}"
            echo -e "${RED}Server Logs (from failing test)${NC}"
            echo -e "${RED}========================================${NC}"
            echo "$relevant"
        fi
    fi
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
