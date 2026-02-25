# step_02_health_check.sh
# Designed to be sourced by run_e2e.sh

step_02_health_check() {
    print_step "Step 2: Health Check"

    # Ping endpoint
    if curl -s -f "$BASE_URL/" > /dev/null; then
        print_success "Server is reachable (GET /)"
    else
        print_error "Server is not reachable at $BASE_URL"
        exit 1
    fi

    # Deep health check
    response=$(curl -s "$BASE_URL/v1/deep_health_check")
    if echo "$response" | jq -e '.database_status == "up"' > /dev/null; then
        print_success "Deep health check passed (database connected)"
    else
        print_error "Deep health check failed (database not connected)"
        exit 1
    fi

    echo ""
}
