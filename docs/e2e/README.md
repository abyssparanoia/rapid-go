# E2E Test

This directory contains an E2E test script to verify the complete environment setup and API flow.

## Overview

The `run_e2e.sh` script performs a comprehensive end-to-end test of the application, including:

1. Health checks
2. Admin user creation and authentication
3. Admin API operations (tenant creation)
4. Staff user creation and authentication
5. Staff API operations

## Prerequisites

Before running the E2E script, you must complete the following setup:

### 1. Start Infrastructure Services

```bash
docker compose up -d
```

This starts:
- MySQL/PostgreSQL database
- Redis
- AWS LocalStack (S3, SNS, SQS)
- AWS Cognito Local emulator
- Spanner emulator (if used)

### 2. Run Database Migrations

```bash
make migrate.up
```

This will:
- Run all pending migrations
- Sync constant tables
- Generate SQLBoiler models

### 3. Initialize Cognito User Pools

```bash
make init.local.cognito
```

This creates the required Cognito user pools for Admin and Staff authentication.

### 4. Build the CLI

```bash
make build
```

This creates the CLI binary at `.bin/app-cli`.

### 5. Start the HTTP Server

In a separate terminal:

```bash
make http.dev
```

Or manually:

```bash
.bin/app-cli http-server
```

The server should be running on `http://localhost:8080`.

## Running the E2E Test

### Basic Usage

```bash
bash docs/e2e/run_e2e.sh
```

### With Custom Configuration

```bash
# Custom base URL
BASE_URL=http://localhost:8080 bash docs/e2e/run_e2e.sh

# Custom CLI path
CLI_PATH=./app bash docs/e2e/run_e2e.sh

# Custom email addresses (to avoid conflicts)
ADMIN_EMAIL=test-admin@example.com \
STAFF_EMAIL=test-staff@example.com \
TENANT_NAME="Test Tenant" \
bash docs/e2e/run_e2e.sh
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `BASE_URL` | `http://localhost:8080` | Base URL for the API server |
| `CLI_PATH` | `.bin/app-cli` | Path to the CLI binary |
| `ADMIN_EMAIL` | `e2e-admin@example.com` | Email for the admin user |
| `ADMIN_DISPLAY_NAME` | `E2E Admin` | Display name for the admin user |
| `STAFF_EMAIL` | `e2e-staff@example.com` | Email for the staff user |
| `STAFF_DISPLAY_NAME` | `E2E Staff` | Display name for the staff user |
| `TENANT_NAME` | `E2E Test Tenant` | Name for the test tenant |

## Test Flow

The script performs the following steps:

### Step 1: Prerequisites Check
- Verifies `curl` is installed
- Verifies `jq` is installed
- Verifies CLI binary exists

### Step 2: Health Check
- Pings the server (`GET /`)
- Performs deep health check (`GET /v1/deep_health_check`)

### Step 3: Create Admin User
- Creates a root admin user via CLI
- Captures AdminID, AuthUID, and Password

### Step 4: Get Admin ID Token
- Authenticates admin user
- Obtains JWT token via debug endpoint

### Step 5: Test Admin API
- Lists tenants using admin token
- Verifies admin authorization works

### Step 6: Create Tenant
- Creates a new tenant via admin API
- Captures TenantID

### Step 7: Create Staff
- Creates an asset presigned URL to get asset_id
- Creates a staff user for the tenant with the asset_id
- Captures StaffID, AuthUID, and Password

### Step 8: Get Staff ID Token
- Authenticates staff user
- Obtains JWT token via debug endpoint

### Step 9: Test Staff API - Get Me
- Retrieves staff profile
- Verifies staff authorization works

### Step 10: Test Staff API - Get Tenant
- Retrieves tenant information
- Verifies staff can access tenant data

## Expected Output

```
=================================
E2E Test
=================================

========================================
Step 1: Checking Prerequisites
========================================
✓ curl is installed
✓ jq is installed
✓ CLI binary exists at .bin/app-cli

========================================
Step 2: Health Check
========================================
✓ Server is reachable (GET /)
✓ Deep health check passed (database connected)

========================================
Step 3: Creating Admin User
========================================
✓ Admin created successfully
ℹ   AdminID: xxx
ℹ   AuthUID: xxx
ℹ   Password: xxx

========================================
Step 4: Getting Admin ID Token
========================================
✓ Admin ID token obtained
ℹ   Token: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...

========================================
Step 5: Testing Admin API - List Tenants
========================================
✓ Admin API accessible (found 0 tenants)

========================================
Step 6: Creating Tenant
========================================
✓ Tenant created successfully
ℹ   TenantID: xxx
ℹ   Name: E2E Test Tenant

========================================
Step 7: Creating Staff
========================================
ℹ Creating asset presigned URL...
ℹ   AssetID: xxx
✓ Staff created successfully
ℹ   StaffID: xxx
ℹ   AuthUID: xxx
ℹ   Email: e2e-staff@example.com
ℹ   Password: xxx

========================================
Step 8: Getting Staff ID Token
========================================
✓ Staff ID token obtained
ℹ   Token: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...

========================================
Step 9: Testing Staff API - Get Me
========================================
✓ Staff API accessible (retrieved own profile)
ℹ   StaffID matches: xxx

========================================
Step 10: Testing Staff API - Get Tenant
========================================
✓ Staff can access tenant (retrieved tenant info)
ℹ   TenantID matches: xxx

========================================
Test Results
========================================
Total Tests: 10
Passed: 10
All tests passed!
```

## Troubleshooting

### Server Not Reachable

**Error**: `Server is not reachable at http://localhost:8080`

**Solution**: Make sure the server is running:
```bash
make http.dev
```

### Database Not Connected

**Error**: `Deep health check failed (database not connected)`

**Solution**:
1. Check if MySQL/PostgreSQL is running: `docker compose ps`
2. Verify migrations: `make migrate.up`
3. Check database connection in `.envrc`

### CLI Binary Not Found

**Error**: `CLI binary not found at .bin/app-cli`

**Solution**: Build the CLI:
```bash
make build
```

### Admin Already Exists

**Error**: `Admin already exists, attempting to use existing credentials`

**Solution**: Use a different email address:
```bash
ADMIN_EMAIL=new-admin@example.com bash docs/e2e/run_e2e.sh
```

### Authentication Failed

**Error**: `Failed to get admin/staff ID token`

**Possible Causes**:
1. Cognito user pools not initialized
2. Wrong credentials
3. Debug endpoints not enabled (only work in local/development)

**Solution**:
1. Run: `make init.local.cognito`
2. Verify `ENV=local` in `.envrc`
3. Restart the server

### jq Not Installed

**Error**: `jq is not installed`

**Solution**:
```bash
# macOS
brew install jq

# Ubuntu/Debian
apt-get install jq
```

## Cleanup

The script does **not** automatically clean up created resources. To clean up:

### Delete Created Users

```bash
# Via Cognito (if using AWS)
aws cognito-idp admin-delete-user \
  --user-pool-id <USER_POOL_ID> \
  --username <AUTH_UID>

# Or reset the database
docker compose down -v
docker compose up -d
make migrate.up
```

### Reset Database

```bash
# Stop containers and remove volumes
docker compose down -v

# Start fresh
docker compose up -d
make migrate.up
make init.local.cognito
```

## Advanced Usage

### Running Specific Tests

You can modify the script to skip certain steps by commenting out function calls in the `main()` function.

### Integration with CI/CD

The script returns exit code 1 on failure, making it suitable for CI/CD pipelines:

```bash
#!/bin/bash
# In CI pipeline
docker compose up -d
make migrate.up
make init.local.cognito
make build

# Start server in background
make http.dev &
SERVER_PID=$!

# Wait for server to be ready
sleep 5

# Run E2E test
bash docs/e2e/run_e2e.sh

# Cleanup
kill $SERVER_PID
docker compose down
```

## Next Steps

After successful E2E verification, you can:

1. Use the created admin/staff credentials for manual testing
2. Explore other API endpoints with Postman or curl
3. Run unit tests: `make test`
4. Run linter: `make lint.go`

## Support

For issues or questions:
- Check the main project README
- Review the project rules in `.claude/rules/`
- Check server logs for detailed error messages
