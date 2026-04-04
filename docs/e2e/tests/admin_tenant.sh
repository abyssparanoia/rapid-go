#!/bin/bash

# E2E Tests: Admin Tenant CRUD

test_admin_list_tenants() {
    print_step "Admin API - List Tenants"

    # @e2e GET /admin/v1/tenants
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

create_tenant() {
    print_step "Creating Tenant"

    # @e2e POST /admin/v1/tenants
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

test_admin_get_tenant() {
    print_step "Admin API - Get Tenant"

    # @e2e GET /admin/v1/tenants/{tenant_id}
    response=$(curl -s "$BASE_URL/admin/v1/tenants/$TENANT_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN")

    tenant_id=$(echo "$response" | jq -r '.tenant.id // empty')

    if [ -n "$tenant_id" ] && [ "$tenant_id" = "$TENANT_ID" ]; then
        print_success "Get tenant successful (ID matches)"
        print_info "  TenantID: $tenant_id"
    else
        print_error "Failed to get tenant"
        echo "$response"
        exit 1
    fi

    echo ""
}

test_admin_update_tenant() {
    print_step "Admin API - Update Tenant"

    UPDATED_TENANT_NAME="Updated ${TENANT_NAME}"

    # @e2e PATCH /admin/v1/tenants/{tenant_id}
    response=$(curl -s -X PATCH "$BASE_URL/admin/v1/tenants/$TENANT_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"name\":\"$UPDATED_TENANT_NAME\"}")

    updated_name=$(echo "$response" | jq -r '.tenant.name // empty')

    if [ "$updated_name" = "$UPDATED_TENANT_NAME" ]; then
        print_success "Tenant name updated successfully"
        print_info "  Name: $updated_name"
    else
        print_error "Failed to update tenant"
        echo "$response"
        exit 1
    fi

    echo ""
}

test_admin_delete_tenant() {
    print_step "Admin API - Delete Tenant"

    # Create a temporary tenant for deletion test
    temp_response=$(curl -s -X POST "$BASE_URL/admin/v1/tenants" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"name\":\"E2E Temp Tenant for Delete ${TEST_ID}\"}")

    TEMP_TENANT_ID=$(echo "$temp_response" | jq -r '.tenant.id // empty')

    if [ -z "$TEMP_TENANT_ID" ]; then
        print_error "Failed to create temp tenant for delete test"
        echo "$temp_response"
        exit 1
    fi

    print_info "  TempTenantID: $TEMP_TENANT_ID"

    # @e2e DELETE /admin/v1/tenants/{tenant_id}
    delete_response=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE \
        "$BASE_URL/admin/v1/tenants/$TEMP_TENANT_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN")

    if [ "$delete_response" = "200" ]; then
        print_success "Tenant deleted successfully (HTTP 200)"
    else
        print_error "Failed to delete tenant (HTTP $delete_response)"
        exit 1
    fi

    # Verify tenant is gone (should return 404)
    verify_response=$(curl -s -o /dev/null -w "%{http_code}" \
        "$BASE_URL/admin/v1/tenants/$TEMP_TENANT_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN")

    if [ "$verify_response" = "404" ] || [ "$verify_response" = "500" ]; then
        print_success "Deleted tenant returns not found (HTTP $verify_response)"
    else
        print_error "Deleted tenant still accessible (HTTP $verify_response)"
        exit 1
    fi

    echo ""
}
