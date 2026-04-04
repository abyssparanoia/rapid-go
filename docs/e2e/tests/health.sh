#!/bin/bash

# E2E Tests: Health Check & Prerequisites

check_prerequisites() {
    print_step "Checking Prerequisites"

    # Check curl
    if command -v curl &> /dev/null; then
        print_success "curl is installed"
    else
        print_error "curl is not installed"
        exit 1
    fi

    # Check jq
    if command -v jq &> /dev/null; then
        print_success "jq is installed"
    else
        print_error "jq is not installed (brew install jq)"
        exit 1
    fi

    # Check CLI exists
    if [ -f "$CLI_PATH" ]; then
        print_success "CLI binary exists at $CLI_PATH"
    else
        print_error "CLI binary not found at $CLI_PATH (run 'make build')"
        exit 1
    fi

    echo ""
}

health_check() {
    print_step "Health Check"

    # Ping endpoint
    if curl -s -f "$BASE_URL/" > /dev/null; then
        print_success "Server is reachable (GET /)"
    else
        print_error "Server is not reachable at $BASE_URL"
        exit 1
    fi

    # @e2e GET /v1/deep_health_check
    response=$(curl -s "$BASE_URL/v1/deep_health_check")
    if echo "$response" | jq -e '.database_status == "up"' > /dev/null; then
        print_success "Deep health check passed (database connected)"
    else
        print_error "Deep health check failed (database not connected)"
        exit 1
    fi

    echo ""
}
