## Coding style & linting
- Enforced by `.golangci.yml` with many strict linters: errcheck, exhaustive, exhaustruct (for domain/usecase/dbmodel structs), forbidigo (disallows fmt.Print*, zap.Error*), funlen/cyclop/gocognit limits, nilerr/nakedret/noctx/wrapcheck, sloglint, govet (enable-all except fieldalignment), gosec (except G115), testifylint/tparallel, etc.
- Formatting uses `go fmt`, `buf format`, `goimports`, `gofumpt` followed by `go mod tidy` (invoked via `make format` or appended to other make targets).
- Prefer structured logging (zapdriver) instead of fmt.Print; avoid long functions; ensure exhaustive switch/map handling; ensure close checks on SQL resources.
- Generated code kept under `internal/infrastructure/grpc/pb`, `internal/infrastructure/*/internal/dbmodel`, schema outputs, etc.

## Architectural conventions
- Layered layout: domain (entities/rules) ← usecase (application logic/DTOs) ← infrastructure (DB, GRPC/OpenAPI, external services) with `cmd/app` wiring.
- DTOs and models commonly validated via `go-playground/validator`.
- Migrations and schema sync handled through the CLI tasks (`schema-migration`), with sqlboiler/yo for ORM-style models.
- Prefer dependency injection through constructors; avoid globals (gochecknoglobals lint).

## Naming/documentation
- Standard Go naming; doc comments expected on exported symbols (godoc-lint enabled implicitly via linters). Keep functions small; avoid magic numbers (go-mnd) and unused params (unparam/unused).

## Protobuf/OpenAPI
- Buf is canonical; definitions live in `schema/proto`; generated Go + OpenAPI artifacts placed under `schema/openapi/...` and `internal/infrastructure/grpc/pb`.