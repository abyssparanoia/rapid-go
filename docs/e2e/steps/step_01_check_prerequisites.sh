# step_01_check_prerequisites.sh
# Designed to be sourced by run_e2e.sh

step_01_check_prerequisites() {
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
