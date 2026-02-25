# step_12_staff_signup.sh
# Designed to be sourced by run_e2e.sh
# Inputs: TEST_ID

step_12_staff_signup() {
    print_step "Step 12: Testing Staff Signup Flow"

    # Generate unique email for signup test
    SIGNUP_STAFF_EMAIL="e2e-signup-staff-${TEST_ID}@example.com"
    SIGNUP_STAFF_PASSWORD="TestPassword123!"

    # Step 1: Create auth credentials via debug API
    print_info "Creating staff auth credentials..."
    auth_response=$(curl -s -X POST "$BASE_URL/debug/v1/staffs/-/auth_uid" \
        -H "Content-Type: application/json" \
        -d "{
            \"email\":\"$SIGNUP_STAFF_EMAIL\",
            \"password\":\"$SIGNUP_STAFF_PASSWORD\"
        }")

    SIGNUP_AUTH_UID=$(echo "$auth_response" | jq -r '.auth_uid // empty')

    if [ -z "$SIGNUP_AUTH_UID" ]; then
        print_error "Failed to create staff auth credentials"
        echo "$auth_response"
        exit 1
    fi

    print_success "Staff auth credentials created"
    print_info "  AuthUID: $SIGNUP_AUTH_UID"
    print_info "  Email: $SIGNUP_STAFF_EMAIL"

    # Step 2: Get ID token for the new staff auth
    print_info "Getting ID token for signup staff..."
    token_response=$(curl -s -X POST "$BASE_URL/debug/v1/staffs/-/id_token" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$SIGNUP_STAFF_EMAIL\",\"password\":\"$SIGNUP_STAFF_PASSWORD\"}")

    SIGNUP_STAFF_TOKEN=$(echo "$token_response" | jq -r '.id_token // empty')

    if [ -n "$SIGNUP_STAFF_TOKEN" ]; then
        print_success "Signup staff ID token obtained"
        print_info "  Token: ${SIGNUP_STAFF_TOKEN:0:50}..."
    else
        print_error "Failed to get signup staff ID token"
        echo "$token_response"
        exit 1
    fi

    echo ""
}
