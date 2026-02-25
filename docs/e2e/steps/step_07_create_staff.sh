# step_07_create_staff.sh
# Designed to be sourced by run_e2e.sh
# Inputs:  ADMIN_TOKEN, TENANT_ID, STAFF_EMAIL, STAFF_DISPLAY_NAME
# Outputs: STAFF_ID, STAFF_AUTH_UID, STAFF_PASSWORD

step_07_create_staff() {
    print_step "Step 7: Creating Staff"

    # Create asset presigned URL to get asset_id
    print_info "Creating asset presigned URL..."
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
