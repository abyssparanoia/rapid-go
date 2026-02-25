# step_10_test_staff_update_me.sh
# Designed to be sourced by run_e2e.sh
# Inputs: STAFF_TOKEN, STAFF_ID, STAFF_DISPLAY_NAME

step_10_test_staff_update_me() {
    print_step "Step 10: Testing Staff API - Update Me"

    # Create new asset for updated image using Staff API
    print_info "Creating asset presigned URL for updated image..."
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
