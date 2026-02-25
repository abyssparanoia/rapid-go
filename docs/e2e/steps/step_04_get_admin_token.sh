# step_04_get_admin_token.sh
# Designed to be sourced by run_e2e.sh
# Inputs:  ADMIN_EMAIL, ADMIN_PASSWORD
# Outputs: ADMIN_TOKEN

step_04_get_admin_token() {
    print_step "Step 4: Getting Admin ID Token"

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
