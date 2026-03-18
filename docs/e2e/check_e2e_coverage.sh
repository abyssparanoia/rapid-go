#!/bin/bash

# E2E Coverage Check Script
# Extracts all endpoints from OpenAPI swagger.json files and verifies
# that each endpoint has a corresponding # @e2e annotation in tests/*.sh.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
E2E_TESTS_DIR="$SCRIPT_DIR/tests"

# Swagger files to check
SWAGGER_FILES=(
    "$ROOT_DIR/schema/openapi/rapid/admin_api/v1/api.swagger.json"
    "$ROOT_DIR/schema/openapi/rapid/staff_api/v1/api.swagger.json"
    "$ROOT_DIR/schema/openapi/rapid/public_api/v1/api.swagger.json"
    "$ROOT_DIR/schema/openapi/rapid/debug_api/v1/api.swagger.json"
)

# Excluded endpoints (add paths here if an endpoint intentionally has no E2E test)
EXCLUDED_ENDPOINTS=(
)

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

is_excluded() {
    local endpoint="$1"
    for excluded in "${EXCLUDED_ENDPOINTS[@]+"${EXCLUDED_ENDPOINTS[@]}"}"; do
        if [ "$endpoint" = "$excluded" ]; then
            return 0
        fi
    done
    return 1
}

echo "=== E2E Coverage Check ==="
echo ""

# 1. Extract all endpoints from swagger files
all_endpoints=()
for swagger_file in "${SWAGGER_FILES[@]}"; do
    if [ ! -f "$swagger_file" ]; then
        echo -e "${YELLOW}WARNING: Swagger file not found: $swagger_file${NC}"
        continue
    fi

    # Extract "METHOD /path" from swagger paths
    endpoints=$(jq -r '.paths | to_entries[] | .key as $path | .value | to_entries[] | "\(.key | ascii_upcase) \($path)"' "$swagger_file" | sort)

    while IFS= read -r endpoint; do
        [ -z "$endpoint" ] && continue
        all_endpoints+=("$endpoint")
    done <<< "$endpoints"
done

# 2. Extract @e2e annotations from all test files in tests/
covered_endpoints=()
while IFS= read -r line; do
    [ -z "$line" ] && continue
    # Extract "METHOD /path" from "# @e2e METHOD /path"
    endpoint=$(echo "$line" | sed 's/^.*# @e2e //')
    covered_endpoints+=("$endpoint")
done < <(grep -rh '# @e2e ' "$E2E_TESTS_DIR" || true)

# 3. Compare and find uncovered endpoints
uncovered=()
for endpoint in "${all_endpoints[@]}"; do
    if is_excluded "$endpoint"; then
        continue
    fi

    found=false
    for covered in "${covered_endpoints[@]+"${covered_endpoints[@]}"}"; do
        if [ "$endpoint" = "$covered" ]; then
            found=true
            break
        fi
    done

    if [ "$found" = false ]; then
        uncovered+=("$endpoint")
    fi
done

# 4. Report results
echo "Total endpoints in OpenAPI: ${#all_endpoints[@]}"
echo "Covered by E2E tests:       ${#covered_endpoints[@]}"
echo "Excluded:                   ${#EXCLUDED_ENDPOINTS[@]}"
echo ""

if [ ${#uncovered[@]} -eq 0 ]; then
    echo -e "${GREEN}All endpoints are covered by E2E tests!${NC}"
    exit 0
else
    echo -e "${RED}Uncovered endpoints (${#uncovered[@]}):${NC}"
    for endpoint in "${uncovered[@]}"; do
        echo -e "  ${RED}- $endpoint${NC}"
    done
    echo ""
    echo -e "${YELLOW}To fix: Add '# @e2e METHOD /path' annotations in docs/e2e/tests/*.sh${NC}"
    echo -e "${YELLOW}Or add to EXCLUDED_ENDPOINTS in this script if intentionally uncovered.${NC}"
    exit 1
fi
