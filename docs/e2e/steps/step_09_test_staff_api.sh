# step_09_test_staff_api.sh
# Designed to be sourced by run_e2e.sh
# Inputs: STAFF_TOKEN, STAFF_ID

step_09_test_staff_api() {
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
