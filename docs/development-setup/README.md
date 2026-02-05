# Development Setup

## Overview

This guide covers setting up the local development environment for rapid-go, including database setup, environment configuration, and common development tasks.

## Prerequisites

- Go 1.25 or higher
- Docker and Docker Compose
- direnv (recommended for environment variable management)
- make

## Environment Configuration

### Using direnv (Recommended)

1. Copy the template environment file:

```bash
cp .envrc.tmpl .envrc
direnv allow
```

2. Configure required variables in `.envrc`:

- **GCP_PROJECT_ID**: Your GCP project ID
- **FIREBASE_CLIENT_KEY**: Firebase client API key
- **Database credentials**: Update DB_HOST, DB_USER, DB_PASSWORD, DB_DATABASE as needed
- **AWS Cognito**: Configure AWS_COGNITO_* variables for authentication

3. Add GCP service account (if using GCP services):

Save your service account JSON as `serviceAccount.json` in the project root.

## Running the Application

### 1. Start Database

```bash
docker-compose up -d
```

This starts the MySQL database and other required services.

### 2. Run Database Migrations

```bash
make migrate.up
```

This applies all pending migrations and regenerates SQLBoiler models.

### 3. Start Development Server

```bash
make http.dev
```

The server will start with hot reload enabled at `http://localhost:8080`.

### 4. Verify Server

```bash
curl http://localhost:8080
```

## Development Tasks

### Code Generation

#### Generate Protocol Buffers + OpenAPI

```bash
make generate.buf
```

Generates gRPC service definitions and OpenAPI (v2) specifications from proto files.

#### Generate SQLBoiler Models

```bash
make generate.sqlboiler
```

Regenerates ORM models from database schema. Usually run automatically by `make migrate.up`.

#### Generate Mocks

```bash
make generate.mock
```

Generates mock implementations for testing using mockgen.

### Database Migrations

#### Create New Migration

```bash
make migrate.create
```

Creates a new timestamped migration file in `db/{mysql,postgresql,spanner}/migrations/`.

#### Apply Migrations

```bash
make migrate.up
```

Applies all pending migrations and regenerates SQLBoiler models.

#### Check Migration Status

```bash
make migrate.status
```

Shows which migrations have been applied.

### Linting

#### Lint Go Code

```bash
make lint.go
```

Runs golangci-lint on all Go code.

#### Lint Proto Files

```bash
make lint.proto
```

Lints Protocol Buffer definitions using buf.

### Testing

```bash
make test
```

Runs all unit tests.

## Initial Admin Setup

After setting up the environment, create an initial root admin:

```bash
make build
./app task create-root-admin \
  --email admin@example.com \
  --display-name "Root Admin"
```

See [create-root-admin CLI documentation](../create-root-admin-cli/README.md) for detailed usage.

## Common Makefile Commands

| Command | Description |
|---------|-------------|
| `make http.dev` | Start server with hot reload |
| `make build` | Build application binary |
| `make migrate.create` | Create new migration file |
| `make migrate.up` | Run migrations + generate SQLBoiler |
| `make migrate.status` | Check migration status |
| `make generate.buf` | Generate Protocol Buffers code |
| `make generate.sqlboiler` | Generate SQLBoiler models |
| `make generate.mock` | Generate mock files |
| `make test` | Run all tests |
| `make lint.go` | Lint Go code |
| `make lint.proto` | Lint proto files |

## Docker Services

The `docker-compose.yml` includes:

- **MySQL**: Primary database (port 3306)
- **Redis**: Cache layer (port 6379)
- **AWS LocalStack**: Local AWS services emulation (port 4566)
  - S3, SNS, SQS, Cognito emulation

## Troubleshooting

### Database Connection Issues

If you encounter database connection errors:

1. Check if Docker containers are running:
   ```bash
   docker-compose ps
   ```

2. Verify database credentials in `.envrc`:
   ```bash
   echo $DB_HOST
   echo $DB_DATABASE
   ```

3. Test database connection:
   ```bash
   mysql -h $DB_HOST -u $DB_USER -p$DB_PASSWORD $DB_DATABASE -e "SELECT 1"
   ```

### Make Command Not Found

If `make` commands fail, use the full path:

```bash
/usr/bin/make migrate.up
```

This is a known issue with shell function overrides in Claude Code environments.

### Proto Generation Fails

Ensure buf is installed and up to date:

```bash
go install github.com/bufbuild/buf/cmd/buf@latest
```

## Next Steps

- Read [Project Architecture Overview](../../.claude/CLAUDE.md) for understanding the codebase structure
- Review [CLI Command Patterns](../../.claude/rules/cli-command-pattern.md) for implementing new commands
- Check [Testing Guidelines](../../.claude/rules/testing.md) for writing unit tests

## Additional Resources

- [Admin API Documentation](../admin-api/)
- [Migration Guidelines](../../.claude/rules/migration.md)
- [Domain Model Patterns](../../.claude/rules/domain-model.md)
