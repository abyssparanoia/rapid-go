# step_11_test_staff_tenant_api.sh
# Designed to be sourced by run_e2e.sh
# Inputs: STAFF_TOKEN, TENANT_ID

step_11_test_staff_tenant_api() {
    print_step "Step 11: Testing Staff API - Get Tenant"

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
