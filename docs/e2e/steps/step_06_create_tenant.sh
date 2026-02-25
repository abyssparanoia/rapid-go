# step_06_create_tenant.sh
# Designed to be sourced by run_e2e.sh
# Inputs:  ADMIN_TOKEN, TENANT_NAME
# Outputs: TENANT_ID

step_06_create_tenant() {
    print_step "Step 6: Creating Tenant"

    response=$(curl -s -X POST "$BASE_URL/admin/v1/tenants" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"name\":\"$TENANT_NAME\"}")

    TENANT_ID=$(echo "$response" | jq -r '.tenant.id // empty')

    if [ -n "$TENANT_ID" ]; then
        print_success "Tenant created successfully"
        print_info "  TenantID: $TENANT_ID"
        print_info "  Name: $TENANT_NAME"
    else
        print_error "Failed to create tenant"
        echo "$response"
        exit 1
    fi

    echo ""
}
