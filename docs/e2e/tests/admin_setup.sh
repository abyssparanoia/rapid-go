#!/bin/bash

# E2E Tests: Admin Setup (create admin user & get token)

create_admin() {
    print_step "Creating Admin User"

    # Create admin via CLI
    print_info "Creating admin with email: $ADMIN_EMAIL"
    output=$("$CLI_PATH" task create-root-admin --email "$ADMIN_EMAIL" --display-name "$ADMIN_DISPLAY_NAME" 2>&1 || true)

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
        if echo "$output" | grep -q "already exists"; then
            print_info "Admin already exists, attempting to use existing credentials"
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

get_admin_token() {
    print_step "Getting Admin ID Token"

    # @e2e POST /debug/v1/admins/-/id_token
    response=$(curl -s -X POST "$BASE_URL/debug/v1/admins/-/id_token" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}")

    ADMIN_TOKEN=$(echo "$response" | jq -r '.id_token // empty')

    if [ -n "$ADMIN_TOKEN" ]; then
        print_success "Admin ID token obtained"
        print_info "  Token: ${ADMIN_TOKEN:0:50}..."
    else
        print_error "Failed to get admin ID token"
        echo "$response"
        exit 1
    fi

    echo ""
}
