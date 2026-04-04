# E2E Test

This directory contains the E2E test suite for verifying the complete API flow.

## Directory Structure

```
docs/e2e/
├── run_e2e.sh              # Orchestrator: builds server, runs all tests
├── common.sh               # Shared helpers, config, state variables
├── check_e2e_coverage.sh   # CI coverage enforcement vs swagger.json
├── tests/
│   ├── health.sh           # Prerequisites & health check
│   ├── admin_setup.sh      # Admin creation + token retrieval
│   ├── admin_tenant.sh     # Admin tenant CRUD
│   ├── admin_staff.sh      # Admin staff CRUD
│   ├── staff_me.sh         # Staff me/tenant APIs
│   ├── staff_list.sh       # Staff list/get APIs
│   └── staff_signup.sh     # Staff signup flow
└── README.md
```

## Endpoint Coverage

All 22 endpoints are covered with `# @e2e` annotations:

| File | Endpoints |
|------|-----------|
| `health.sh` | `GET /v1/deep_health_check` |
| `admin_setup.sh` | `POST /debug/v1/admins/-/id_token` |
| `admin_tenant.sh` | `GET/POST/PATCH/DELETE /admin/v1/tenants[/{id}]` |
| `admin_staff.sh` | `POST /admin/v1/assets/-/presigned_url`, `GET/POST/PATCH /admin/v1/staffs[/{id}]` |
| `staff_me.sh` | `POST /debug/v1/staffs/-/id_token`, `GET/PATCH /staff/v1/me`, `GET/PATCH /staff/v1/me/tenant`, `POST /staff/v1/assets/-/presigned_url` |
| `staff_list.sh` | `GET /staff/v1/staffs`, `GET /staff/v1/staffs/{staff_id}` |
| `staff_signup.sh` | `POST /debug/v1/staffs/-/auth_uid`, `POST /staff/v1/me:signup` |

## Prerequisites

### 1. Start Infrastructure Services

```bash
docker compose up -d
```

### 2. Run Database Migrations

```bash
make migrate.up
```

### 3. Initialize Cognito User Pools

```bash
make init.local.cognito
```

## Running Tests

### Full E2E test (auto-builds and starts server)

```bash
make e2e
```

This automatically:
1. Builds the CLI binary (`make build`)
2. Starts the HTTP server in the background
3. Runs all test functions in dependency order
4. Stops the server and cleans up on exit (pass or fail)

### Check endpoint coverage

```bash
make e2e.cover
```

Reports any endpoints in `schema/openapi/rapid/*/api.swagger.json` that lack `# @e2e` annotations.

### Manual run (server already running)

```bash
bash docs/e2e/run_e2e.sh
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `BASE_URL` | `http://localhost:8080` | API server base URL |
| `CLI_PATH` | `.bin/app-cli` | Path to CLI binary |
| `ADMIN_EMAIL` | auto-generated | Admin user email |
| `ADMIN_DISPLAY_NAME` | auto-generated | Admin display name |
| `STAFF_EMAIL` | auto-generated | Staff user email |
| `STAFF_DISPLAY_NAME` | auto-generated | Staff display name |
| `TENANT_NAME` | auto-generated | Test tenant name |
| `POSTGRES_CONTAINER` | `rapid-go-postgresql-1` | PostgreSQL Docker container name |

## Adding New Tests

### 1. Add `# @e2e` annotation

Every `curl` call that tests an API endpoint must have a `# @e2e METHOD /path` comment directly above it:

```bash
# @e2e GET /admin/v1/examples
response=$(curl -s "$BASE_URL/admin/v1/examples" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
```

### 2. Place in the correct test file

Add the test function to the appropriate resource file in `tests/`. If the resource doesn't exist, create a new file.

### 3. Register in `run_e2e.sh`

Add the function call to the appropriate phase in `main()` in `run_e2e.sh`.

### 4. Verify coverage

```bash
make e2e.cover
```

## Troubleshooting

### Server not reachable

Ensure Docker Compose services are running: `docker compose up -d`

### Database not connected

Run migrations: `make migrate.up`

### Authentication failed

Initialize Cognito: `make init.local.cognito` and verify `ENV=local` in `.envrc`

### Admin already exists

Use a different email: `ADMIN_EMAIL=other@example.com make e2e`

### Reset everything

```bash
docker compose down -v
docker compose up -d
make migrate.up
make init.local.cognito
```
