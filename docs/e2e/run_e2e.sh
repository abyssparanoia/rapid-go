#!/bin/bash

# E2E Test Script
# This script performs end-to-end testing of the full authentication and API flow

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="${BASE_URL:-http://localhost:8080}"
CLI_PATH="${CLI_PATH:-.bin/app-cli}"

# Generate unique identifiers for this test run
TEST_ID="$(date +%s)-${RANDOM}"
ADMIN_EMAIL="${ADMIN_EMAIL:-e2e-admin-${TEST_ID}@example.com}"
ADMIN_DISPLAY_NAME="${ADMIN_DISPLAY_NAME:-E2E Admin ${TEST_ID}}"
STAFF_EMAIL="${STAFF_EMAIL:-e2e-staff-${TEST_ID}@example.com}"
STAFF_DISPLAY_NAME="${STAFF_DISPLAY_NAME:-E2E Staff ${TEST_ID}}"
TENANT_NAME="${TENANT_NAME:-E2E Test Tenant ${TEST_ID}}"

# Test results tracking
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Helper functions
print_step() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
    PASSED_TESTS=$((PASSED_TESTS + 1))
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
    FAILED_TESTS=$((FAILED_TESTS + 1))
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

print_result() {
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}Test Results${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo -e "Total Tests: ${TOTAL_TESTS}"
    echo -e "${GREEN}Passed: ${PASSED_TESTS}${NC}"
    if [ $FAILED_TESTS -gt 0 ]; then
        echo -e "${RED}Failed: ${FAILED_TESTS}${NC}"
        exit 1
    else
        echo -e "${GREEN}All tests passed!${NC}"
    fi
}

check_prerequisites() {
    print_step "Step 1: Checking Prerequisites"

    # Check curl
    if command -v curl &> /dev/null; then
        print_success "curl is installed"
    else
        print_error "curl is not installed"
        exit 1
    fi

    # Check jq
    if command -v jq &> /dev/null; then
        print_success "jq is installed"
    else
        print_error "jq is not installed (brew install jq)"
        exit 1
    fi

    # Check CLI exists
    if [ -f "$CLI_PATH" ]; then
        print_success "CLI binary exists at $CLI_PATH"
    else
        print_error "CLI binary not found at $CLI_PATH (run 'make build')"
        exit 1
    fi

    echo ""
}

health_check() {
    print_step "Step 2: Health Check"

    # Ping endpoint
    if curl -s -f "$BASE_URL/" > /dev/null; then
        print_success "Server is reachable (GET /)"
    else
        print_error "Server is not reachable at $BASE_URL"
        exit 1
    fi

    # Deep health check
    response=$(curl -s "$BASE_URL/v1/deep_health_check")
    if echo "$response" | jq -e '.database_status == "up"' > /dev/null; then
        print_success "Deep health check passed (database connected)"
    else
        print_error "Deep health check failed (database not connected)"
        exit 1
    fi

    echo ""
}

create_admin() {
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

get_admin_token() {
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

test_admin_api() {
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

create_tenant() {
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

create_staff() {
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

get_staff_token() {
    print_step "Step 8: Getting Staff ID Token"

    response=$(curl -s -X POST "$BASE_URL/debug/v1/staffs/-/id_token" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$STAFF_EMAIL\",\"password\":\"$STAFF_PASSWORD\"}")

    STAFF_TOKEN=$(echo "$response" | jq -r '.id_token // empty')

    if [ -n "$STAFF_TOKEN" ]; then
        print_success "Staff ID token obtained"
        print_info "  Token: ${STAFF_TOKEN:0:50}..."
    else
        print_error "Failed to get staff ID token"
        echo "$response"
        exit 1
    fi

    echo ""
}

test_staff_api() {
    print_step "Step 9: Testing Staff API - Get Me"

    response=$(curl -s "$BASE_URL/staff/v1/me" \
        -H "Authorization: Bearer $STAFF_TOKEN")

    staff_id=$(echo "$response" | jq -r '.staff.id // empty')

    if [ -n "$staff_id" ] && [ "$staff_id" = "$STAFF_ID" ]; then
        print_success "Staff API accessible (retrieved own profile)"
        print_info "  StaffID matches: $staff_id"
    else
        print_error "Failed to access staff API"
        echo "$response"
        exit 1
    fi

    echo ""
}

test_staff_update_me() {
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

test_staff_tenant_api() {
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

staff_signup() {
    print_step "Step 12: Testing Staff Signup Flow"

    # Generate unique email for signup test
    SIGNUP_STAFF_EMAIL="e2e-signup-staff-${TEST_ID}@example.com"
    SIGNUP_STAFF_PASSWORD="TestPassword123!"

    # Step 1: Create auth credentials via debug API
    print_info "Creating staff auth credentials..."
    auth_response=$(curl -s -X POST "$BASE_URL/debug/v1/staffs/-/auth_uid" \
        -H "Content-Type: application/json" \
        -d "{
            \"email\":\"$SIGNUP_STAFF_EMAIL\",
            \"password\":\"$SIGNUP_STAFF_PASSWORD\"
        }")

    SIGNUP_AUTH_UID=$(echo "$auth_response" | jq -r '.auth_uid // empty')

    if [ -z "$SIGNUP_AUTH_UID" ]; then
        print_error "Failed to create staff auth credentials"
        echo "$auth_response"
        exit 1
    fi

    print_success "Staff auth credentials created"
    print_info "  AuthUID: $SIGNUP_AUTH_UID"
    print_info "  Email: $SIGNUP_STAFF_EMAIL"

    # Step 2: Get ID token for the new staff auth
    print_info "Getting ID token for signup staff..."
    token_response=$(curl -s -X POST "$BASE_URL/debug/v1/staffs/-/id_token" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$SIGNUP_STAFF_EMAIL\",\"password\":\"$SIGNUP_STAFF_PASSWORD\"}")

    SIGNUP_STAFF_TOKEN=$(echo "$token_response" | jq -r '.id_token // empty')

    if [ -n "$SIGNUP_STAFF_TOKEN" ]; then
        print_success "Signup staff ID token obtained"
        print_info "  Token: ${SIGNUP_STAFF_TOKEN:0:50}..."
    else
        print_error "Failed to get signup staff ID token"
        echo "$token_response"
        exit 1
    fi

    echo ""
}

# Main execution
main() {
    echo ""
    echo -e "${GREEN}=================================${NC}"
    echo -e "${GREEN}E2E Test${NC}"
    echo -e "${GREEN}=================================${NC}"
    echo ""

    check_prerequisites
    health_check
    create_admin
    get_admin_token
    test_admin_api
    create_tenant
    create_staff
    get_staff_token
    test_staff_api
    test_staff_update_me
    test_staff_tenant_api
    staff_signup

    print_result
}

# Run main function
main
