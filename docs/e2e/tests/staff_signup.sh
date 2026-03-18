#!/bin/bash

# E2E Tests: Staff Signup Flow

staff_signup() {
    print_step "Staff Signup Flow"

    # Generate unique email for signup test
    SIGNUP_STAFF_EMAIL="e2e-signup-staff-${TEST_ID}@example.com"
    SIGNUP_STAFF_PASSWORD="TestPassword123!"

    # Step 1: Create auth credentials via debug API
    print_info "Creating staff auth credentials..."
    # @e2e POST /debug/v1/staffs/-/auth_uid
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

    # Step 3: Create asset for signup via direct DB insert
    # The signup endpoint validates assets with auth_context="staff:<auth_uid>" (the Cognito sub).
    # The auth_uid returned by the debug API is the Cognito sub, so we use it directly.
    print_info "Creating asset for signup via DB insert..."

    SIGNUP_ASSET_ID="e2e-signup-asset-${TEST_ID}"
    SIGNUP_ASSET_AUTH_CONTEXT="staff:${SIGNUP_AUTH_UID}"
    SIGNUP_ASSET_PATH="private/user_images/e2e-signup-${TEST_ID}.png"
    NOW=$(date -u '+%Y-%m-%d %H:%M:%S+00')
    EXPIRES=$(date -u -v+15M '+%Y-%m-%d %H:%M:%S+00' 2>/dev/null || date -u -d '+15 minutes' '+%Y-%m-%d %H:%M:%S+00')

    run_psql "
        INSERT INTO assets (id, auth_context, content_type, type, path, expires_at, created_at, updated_at)
        VALUES ('${SIGNUP_ASSET_ID}', '${SIGNUP_ASSET_AUTH_CONTEXT}', 'image/png', 'private/user_images', '${SIGNUP_ASSET_PATH}', '${EXPIRES}', '${NOW}', '${NOW}');
    "

    if [ $? -ne 0 ]; then
        print_error "Failed to create asset for signup"
        exit 1
    fi

    print_success "Signup asset created"
    print_info "  AssetID: $SIGNUP_ASSET_ID"
    print_info "  AuthContext: $SIGNUP_ASSET_AUTH_CONTEXT"

    # Step 4: Call signup API
    print_info "Calling signup API..."
    # @e2e POST /staff/v1/me:signup
    signup_response=$(curl -s -X POST "$BASE_URL/staff/v1/me:signup" \
        -H "Authorization: Bearer $SIGNUP_STAFF_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"tenant_name\":\"E2E Signup Tenant ${TEST_ID}\",
            \"display_name\":\"E2E Signup Staff ${TEST_ID}\",
            \"image_asset_id\":\"$SIGNUP_ASSET_ID\"
        }")

    signup_staff_id=$(echo "$signup_response" | jq -r '.staff.id // empty')

    if [ -n "$signup_staff_id" ]; then
        print_success "Staff signup successful"
        print_info "  StaffID: $signup_staff_id"
    else
        print_error "Failed to signup staff"
        echo "$signup_response"
        exit 1
    fi

    echo ""
}
