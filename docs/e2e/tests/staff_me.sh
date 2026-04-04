#!/bin/bash

# E2E Tests: Staff Me & Tenant APIs

get_staff_token() {
    print_step "Getting Staff ID Token"

    # @e2e POST /debug/v1/staffs/-/id_token
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

test_staff_get_me() {
    print_step "Staff API - Get Me"

    # @e2e GET /staff/v1/me
    response=$(curl -s "$BASE_URL/staff/v1/me" \
        -H "Authorization: Bearer $STAFF_TOKEN")

    staff_id=$(echo "$response" | jq -r '.staff.id // empty')

    if [ -n "$staff_id" ] && [ "$staff_id" = "$STAFF_ID" ]; then
        print_success "Staff API accessible (retrieved own profile)"
        print_info "  StaffID matches: $staff_id"
    else
        print_error "Failed to access staff API"
        echo "$response"
        exit 1
    fi

    echo ""
}

test_staff_update_me() {
    print_step "Staff API - Update Me"

    # Create new asset for updated image using Staff API
    print_info "Creating asset presigned URL for updated image..."
    # @e2e POST /staff/v1/assets/-/presigned_url
    asset_response=$(curl -s -X POST "$BASE_URL/staff/v1/assets/-/presigned_url" \
        -H "Authorization: Bearer $STAFF_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"asset_type\":\"ASSET_TYPE_USER_IMAGE\",
            \"content_type\":\"CONTENT_TYPE_IMAGE_PNG\"
        }")

    UPDATED_ASSET_ID=$(echo "$asset_response" | jq -r '.asset_id // empty')

    if [ -z "$UPDATED_ASSET_ID" ]; then
        print_error "Failed to create asset presigned URL for update"
        echo "$asset_response"
        exit 1
    fi

    print_info "  Updated AssetID: $UPDATED_ASSET_ID"

    # Update staff profile with new display name and image
    UPDATED_DISPLAY_NAME="${STAFF_DISPLAY_NAME} (Updated)"
    # @e2e PATCH /staff/v1/me
    response=$(curl -s -X PATCH "$BASE_URL/staff/v1/me" \
        -H "Authorization: Bearer $STAFF_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"display_name\":\"$UPDATED_DISPLAY_NAME\",
            \"image_asset_id\":\"$UPDATED_ASSET_ID\"
        }")

    updated_display_name=$(echo "$response" | jq -r '.staff.display_name // empty')
    staff_id=$(echo "$response" | jq -r '.staff.id // empty')

    if [ -n "$staff_id" ] && [ "$staff_id" = "$STAFF_ID" ] && [ "$updated_display_name" = "$UPDATED_DISPLAY_NAME" ]; then
        print_success "Staff profile updated successfully"
        print_info "  Display name updated: $updated_display_name"
        print_info "  Image asset updated: $UPDATED_ASSET_ID"
    else
        print_error "Failed to update staff profile"
        echo "$response"
        exit 1
    fi

    echo ""
}

test_staff_get_tenant() {
    print_step "Staff API - Get Tenant"

    # @e2e GET /staff/v1/me/tenant
    response=$(curl -s "$BASE_URL/staff/v1/me/tenant" \
        -H "Authorization: Bearer $STAFF_TOKEN")

    tenant_id=$(echo "$response" | jq -r '.tenant.id // empty')

    if [ -n "$tenant_id" ] && [ "$tenant_id" = "$TENANT_ID" ]; then
        print_success "Staff can access tenant (retrieved tenant info)"
        print_info "  TenantID matches: $tenant_id"
    else
        print_error "Failed to access tenant via staff API"
        echo "$response"
        exit 1
    fi

    echo ""
}

test_staff_update_tenant() {
    print_step "Staff API - Update Tenant"

    STAFF_UPDATED_TENANT_NAME="Staff Updated Tenant ${TEST_ID}"

    # @e2e PATCH /staff/v1/me/tenant
    response=$(curl -s -X PATCH "$BASE_URL/staff/v1/me/tenant" \
        -H "Authorization: Bearer $STAFF_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"name\":\"$STAFF_UPDATED_TENANT_NAME\"}")

    updated_name=$(echo "$response" | jq -r '.tenant.name // empty')

    if [ "$updated_name" = "$STAFF_UPDATED_TENANT_NAME" ]; then
        print_success "Staff updated tenant name successfully"
        print_info "  Name: $updated_name"
    else
        print_error "Failed to update tenant via staff API"
        echo "$response"
        exit 1
    fi

    echo ""
}

test_staff_create_asset() {
    print_step "Staff API - Create Asset Presigned URL"

    # @e2e POST /staff/v1/assets/-/presigned_url
    response=$(curl -s -X POST "$BASE_URL/staff/v1/assets/-/presigned_url" \
        -H "Authorization: Bearer $STAFF_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"asset_type\":\"ASSET_TYPE_USER_IMAGE\",
            \"content_type\":\"CONTENT_TYPE_IMAGE_PNG\"
        }")

    asset_id=$(echo "$response" | jq -r '.asset_id // empty')
    presigned_url=$(echo "$response" | jq -r '.presigned_url // empty')

    if [ -n "$asset_id" ] && [ -n "$presigned_url" ]; then
        print_success "Staff create asset presigned URL successful"
        print_info "  AssetID: $asset_id"
    else
        print_error "Failed to create asset presigned URL via staff API"
        echo "$response"
        exit 1
    fi

    echo ""
}
