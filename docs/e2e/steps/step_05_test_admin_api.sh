# step_05_test_admin_api.sh
# Designed to be sourced by run_e2e.sh
# Inputs: ADMIN_TOKEN

step_05_test_admin_api() {
    print_step "Step 5: Testing Admin API - List Tenants"

    response=$(curl -s "$BASE_URL/admin/v1/tenants" \
        -H "Authorization: Bearer $ADMIN_TOKEN")

    if echo "$response" | jq -e '.tenants' > /dev/null 2>&1; then
        tenant_count=$(echo "$response" | jq '.tenants | length')
        print_success "Admin API accessible (found $tenant_count tenants)"
    else
        print_error "Failed to access admin API"
        echo "$response"
        exit 1
    fi

    echo ""
}
