#!/bin/bash

# E2E Test Runner
# Sources test files from docs/e2e/tests/ and runs them in order.
# Each test file defines functions; this script calls them in the correct sequence.

set -e  # Exit on error

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Source common helpers & variables
source "$SCRIPT_DIR/common.sh"

# Source all test files
for f in "$SCRIPT_DIR"/tests/*.sh; do
    source "$f"
done

# PID of background HTTP server
HTTP_SERVER_PID=""

cleanup() {
    local exit_code=$?
    echo ""
    # Show server logs on failure (before killing processes)
    if [ $exit_code -ne 0 ]; then
        print_server_logs
    fi
    print_info "Cleaning up background processes..."
    if [ -n "$HTTP_SERVER_PID" ]; then
        kill "$HTTP_SERVER_PID" 2>/dev/null || true
        wait "$HTTP_SERVER_PID" 2>/dev/null || true
    fi
    # Remove temp log file
    [ -n "$SERVER_LOG_FILE" ] && rm -f "$SERVER_LOG_FILE"
    print_info "Cleanup complete."
}

trap cleanup EXIT

build_binary() {
    print_step "Building CLI Binary"

    print_info "Building CLI binary..."
    (cd "$REPO_ROOT" && /usr/bin/make build) || { echo "make build failed"; exit 1; }
    print_info "Build complete: $CLI_PATH"

    echo ""
}

start_server() {
    print_step "Starting Server"

    # Start HTTP server in background
    print_info "Starting HTTP server..."
    SERVER_LOG_FILE=$(mktemp /tmp/e2e-server-XXXXXX.log)
    (cd "$REPO_ROOT" && "$CLI_PATH" http-server run) > "$SERVER_LOG_FILE" 2>&1 &
    HTTP_SERVER_PID=$!
    print_info "HTTP server started (PID=$HTTP_SERVER_PID)"

    # Wait for server to become healthy
    wait_for_health

    echo ""
}

# Main execution
main() {
    echo ""
    echo -e "${GREEN}=================================${NC}"
    echo -e "${GREEN}E2E Test${NC}"
    echo -e "${GREEN}=================================${NC}"
    echo ""

    # Build binary first (required for migration CLI and server)
    build_binary

    # Reset DB and re-apply migrations for a clean state
    reset_database
    migrate_and_seed

    # Start server
    start_server

    # Phase 1: Prerequisites & Health
    check_prerequisites
    health_check

    # Phase 2: Admin Setup
    create_admin
    get_admin_token

    # Phase 3: Admin Tenant CRUD
    test_admin_list_tenants
    create_tenant
    test_admin_get_tenant
    test_admin_update_tenant
    test_admin_delete_tenant

    # Phase 4: Admin Staff CRUD
    create_staff
    test_admin_list_staffs
    test_admin_get_staff
    test_admin_update_staff

    # Phase 5: Staff APIs (requires staff token)
    get_staff_token
    test_staff_get_me
    test_staff_update_me
    test_staff_get_tenant
    test_staff_update_tenant
    test_staff_create_asset

    # Phase 6: Staff List/Get
    test_staff_list_staffs
    test_staff_get_staff

    # Phase 7: Staff Signup Flow
    staff_signup

    # Results
    print_result
}

# Run main function
main
