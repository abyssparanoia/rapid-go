#!/bin/bash

# E2E Tests: Staff List & Get APIs

test_staff_list_staffs() {
    print_step "Staff API - List Staffs"

    # @e2e GET /staff/v1/staffs
    response=$(curl -s "$BASE_URL/staff/v1/staffs" \
        -H "Authorization: Bearer $STAFF_TOKEN")

    if echo "$response" | jq -e '.staffs' > /dev/null 2>&1; then
        staff_count=$(echo "$response" | jq '.staffs | length')
        print_success "Staff list staffs successful (found $staff_count staffs)"
    else
        print_error "Failed to list staffs via staff API"
        echo "$response"
        exit 1
    fi

    echo ""
}

test_staff_get_staff() {
    print_step "Staff API - Get Staff"

    # @e2e GET /staff/v1/staffs/{staff_id}
    response=$(curl -s "$BASE_URL/staff/v1/staffs/$STAFF_ID" \
        -H "Authorization: Bearer $STAFF_TOKEN")

    staff_id=$(echo "$response" | jq -r '.staff.id // empty')

    if [ -n "$staff_id" ] && [ "$staff_id" = "$STAFF_ID" ]; then
        print_success "Staff get staff successful (ID matches)"
        print_info "  StaffID: $staff_id"
    else
        print_error "Failed to get staff via staff API"
        echo "$response"
        exit 1
    fi

    echo ""
}
