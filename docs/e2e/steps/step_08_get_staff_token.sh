# step_08_get_staff_token.sh
# Designed to be sourced by run_e2e.sh
# Inputs:  STAFF_EMAIL, STAFF_PASSWORD
# Outputs: STAFF_TOKEN

step_08_get_staff_token() {
    print_step "Step 8: Getting Staff ID Token"

    response=$(curl -s -X POST "$BASE_URL/debug/v1/staffs/-/id_token" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$STAFF_EMAIL\",\"password\":\"$STAFF_PASSWORD\"}")

    STAFF_TOKEN=$(echo "$response" | jq -r '.id_token // empty')

    if [ -n "$STAFF_TOKEN" ]; then
        print_success "Staff ID token obtained"
        print_info "  Token: ${STAFF_TOKEN:0:50}..."
    else
        print_error "Failed to get staff ID token"
        echo "$response"
        exit 1
    fi

    echo ""
}
