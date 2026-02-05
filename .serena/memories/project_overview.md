## Project purpose
- `rapid-go` is a Go 1.25 boilerplate/monorepo template for building HTTP-only API services following a layered architecture that clarifies responsibilities (domain/usecase/infrastructure) and is based on golang-standards/project-layout.

## Tech stack
- Language: Go 1.25
- Server runtime: Air live-reload (`make http.dev`), REST API in `cmd/app`
- Data: MySQL/PostgreSQL (sqlboiler), Spanner (wrench/splanter), Redis; migrations via app CLI and goose-compatible schema files
- API schema: Protocol Buffers + grpc-gateway/openapi via Buf
- Tooling: golangci-lint, buf, sqlboiler, goimports, gofumpt, air, mockgen, yo, docker-compose

## Repo structure (top-level)
- `cmd/app`: application entrypoint (HTTP server CLI)
- `internal/domain`: business models & domain logic
- `internal/usecase`: application use cases & DTOs
- `internal/infrastructure`: adapters (mysql/postgresql/spanner dbmodels, grpc pb, etc.)
- `internal/pkg`: shared/internal utilities
- `db`: database configs, sqlboiler templates, migrations, mermaid outputs
- `schema`: protobuf & openapi definitions
- `docker`: container assets; `docker-compose.yml` for local deps
- `localstack`: scripts/data for local AWS emulation (Cognito etc.)
- `docs`: developer docs (e.g., docs/golang.md)
- `deployments`: terraform and deployment assets
- `test-script`: sample/test scripts for spanner and others

## Entry points
- REST API server: run `make http.dev` then hit http://localhost:8080 (local)
- CLI: `./.bin/app-cli` produced by `make build` (used for migrations and tasks)

## Environment
- Uses direnv; copy `.envrc.tmpl` to `.envrc`, set `GCP_PROJECT_ID`, `FIREBASE_CLIENT_KEY`, and add GCP service account JSON as `serviceAccount.json` in repo root.