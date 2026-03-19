#!/bin/bash

# E2E Tests: Admin Staff CRUD

create_staff() {
    print_step "Creating Staff"

    # Create asset presigned URL to get asset_id
    print_info "Creating asset presigned URL..."
    # @e2e POST /admin/v1/assets/-/presigned_url
    asset_response=$(curl -s -X POST "$BASE_URL/admin/v1/assets/-/presigned_url" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"asset_type\":\"ASSET_TYPE_USER_IMAGE\",
            \"content_type\":\"CONTENT_TYPE_IMAGE_PNG\"
        }")

    ASSET_ID=$(echo "$asset_response" | jq -r '.asset_id // empty')

    if [ -z "$ASSET_ID" ]; then
        print_error "Failed to create asset presigned URL"
        echo "$asset_response"
        exit 1
    fi

    print_info "  AssetID: $ASSET_ID"

    # Create staff with asset_id
    # @e2e POST /admin/v1/staffs
    response=$(curl -s -X POST "$BASE_URL/admin/v1/staffs" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"tenant_id\":\"$TENANT_ID\",
            \"email\":\"$STAFF_EMAIL\",
            \"role\":\"STAFF_ROLE_ADMIN\",
            \"display_name\":\"$STAFF_DISPLAY_NAME\",
            \"image_asset_id\":\"$ASSET_ID\"
        }")

    STAFF_ID=$(echo "$response" | jq -r '.staff.id // empty')
    STAFF_AUTH_UID=$(echo "$response" | jq -r '.staff.auth_uid // empty')
    STAFF_PASSWORD=$(echo "$response" | jq -r '.password // empty')

    if [ -n "$STAFF_ID" ] && [ -n "$STAFF_AUTH_UID" ] && [ -n "$STAFF_PASSWORD" ]; then
        print_success "Staff created successfully"
        print_info "  StaffID: $STAFF_ID"
        print_info "  AuthUID: $STAFF_AUTH_UID"
        print_info "  Email: $STAFF_EMAIL"
        print_info "  Password: $STAFF_PASSWORD"
    else
        print_error "Failed to create staff"
        echo "$response"
        exit 1
    fi

    echo ""
}

test_admin_list_staffs() {
    print_step "Admin API - List Staffs"

    # @e2e GET /admin/v1/staffs
    response=$(curl -s "$BASE_URL/admin/v1/staffs?tenant_id=$TENANT_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN")

    if echo "$response" | jq -e '.staffs' > /dev/null 2>&1; then
        staff_count=$(echo "$response" | jq '.staffs | length')
        print_success "Admin list staffs successful (found $staff_count staffs)"
    else
        print_error "Failed to list staffs"
        echo "$response"
        exit 1
    fi

    echo ""
}

test_admin_get_staff() {
    print_step "Admin API - Get Staff"

    # @e2e GET /admin/v1/staffs/{staff_id}
    response=$(curl -s "$BASE_URL/admin/v1/staffs/$STAFF_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN")

    staff_id=$(echo "$response" | jq -r '.staff.id // empty')

    if [ -n "$staff_id" ] && [ "$staff_id" = "$STAFF_ID" ]; then
        print_success "Get staff successful (ID matches)"
        print_info "  StaffID: $staff_id"
    else
        print_error "Failed to get staff"
        echo "$response"
        exit 1
    fi

    echo ""
}

test_admin_update_staff() {
    print_step "Admin API - Update Staff"

    UPDATED_STAFF_NAME="Updated ${STAFF_DISPLAY_NAME}"

    # @e2e PATCH /admin/v1/staffs/{staff_id}
    response=$(curl -s -X PATCH "$BASE_URL/admin/v1/staffs/$STAFF_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"display_name\":\"$UPDATED_STAFF_NAME\"}")

    updated_name=$(echo "$response" | jq -r '.staff.display_name // empty')

    if [ "$updated_name" = "$UPDATED_STAFF_NAME" ]; then
        print_success "Staff display name updated successfully"
        print_info "  Display name: $updated_name"
    else
        print_error "Failed to update staff"
        echo "$response"
        exit 1
    fi

    echo ""
}
