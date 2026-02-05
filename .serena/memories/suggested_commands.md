## Setup
- `cp .envrc.tmpl .envrc && direnv allow` – load env vars; add `serviceAccount.json` for GCP.
- `docker-compose up -d` – start local databases/redis/etc.

## Run / develop
- `make http.dev` – start HTTP API with Air live reload (localhost:8080).
- `make build` – build CLI binary `./.bin/app-cli` used by tasks/migrations.
- `curl http://localhost:8080` – quick health check after start.

## Database & schema
- `make migrate.create` – scaffold new DB migration files (uses app CLI).
- `make migrate.up` – run migrations, sync constants, extract schema, regenerate diagrams & sqlboiler (postgres default).
- `make migrate.spanner.up` – apply Spanner migrations/loads via wrench & splanter.
- `make generate.sqlboiler.mysql|postgresql` – regenerate ORM models for respective DB.
- `make generate.yo` – generate Spanner dbmodel via yo.

## APIs / protobuf
- `make generate.buf` – regenerate protobuf + grpc-gateway + OpenAPI artifacts (cleans generated dirs).

## Linting & formatting
- `make lint.go` – run golangci-lint suite.
- `make lint.proto` – buf lint for protobuf.
- `make format` – go fmt + buf format + goimports + gofumpt + go mod tidy.

## Testing
- `make test` or `go test ./internal/...` – run unit tests.

## Local AWS emulation
- `make init.local.cognito` – initialize Localstack Cognito resources (scripts in `localstack/`).