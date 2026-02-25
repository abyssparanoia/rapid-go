# step_03_create_admin.sh
# Designed to be sourced by run_e2e.sh
# Outputs: ADMIN_ID, ADMIN_AUTH_UID, ADMIN_PASSWORD

step_03_create_admin() {
    print_step "Step 3: Creating Admin User"

    # Create admin via CLI
    print_info "Creating admin with email: $ADMIN_EMAIL"
    output=$($CLI_PATH task create-root-admin --email "$ADMIN_EMAIL" --display-name "$ADMIN_DISPLAY_NAME" 2>&1 || true)

    # Parse output
    ADMIN_ID=$(echo "$output" | grep "AdminID:" | awk '{print $2}')
    ADMIN_AUTH_UID=$(echo "$output" | grep "AuthUID:" | awk '{print $2}')
    ADMIN_PASSWORD=$(echo "$output" | grep "Password:" | awk '{print $2}')

    if [ -n "$ADMIN_ID" ] && [ -n "$ADMIN_AUTH_UID" ] && [ -n "$ADMIN_PASSWORD" ]; then
        print_success "Admin created successfully"
        print_info "  AdminID: $ADMIN_ID"
        print_info "  AuthUID: $ADMIN_AUTH_UID"
        print_info "  Password: $ADMIN_PASSWORD"
    else
        # Admin might already exist, check error message
        if echo "$output" | grep -q "already exists"; then
            print_info "Admin already exists, attempting to use existing credentials"
            # For existing admin, we need to get the AuthUID manually or use known values
            print_error "Cannot retrieve existing admin credentials automatically"
            print_info "Please delete existing admin or use different ADMIN_EMAIL"
            exit 1
        else
            print_error "Failed to create admin"
            echo "$output"
            exit 1
        fi
    fi

    echo ""
}
