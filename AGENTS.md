# AGENTS Guide for livedx-oshipo-server

## High-Level Overview

- Monorepo Go service providing gRPC APIs (admin, user, public, debug) that are proxied to HTTP via gRPC-Gateway; exposed through `cmd/app/main.go` and the Cobra CLI under `app`.
- Layered architecture: **domain** (pure business logic and contracts), **usecase** (application services/interactors), **infrastructure** (framework, transport, persistence, third-party integrations).
- Primary dependencies: PostgreSQL (SQLBoiler models), Firebase Authentication, Google Cloud Storage, X (Twitter) API, gRPC + HTTP gateway, Zap-based structured logging.

## Entry Points & Execution Flow

- `cmd/app/main.go` → `internal/infrastructure/cmd/root.go` registers subcommands:
  - `http-server run` → starts `internal/infrastructure/http.Run()` which spins up the gRPC server (`internal/infrastructure/grpc`) and HTTP gateway with CORS and ping handler.
  - `task ...` → maintenance tasks (e.g., create root staff/admin) defined in `internal/infrastructure/cmd/internal/task_cmd` and backed by usecase interactors.
  - `schema-migration database ...` → wraps migration helpers in `internal/infrastructure/postgresql/migration` for Goose migrations, schema extraction, and constant sync.
- HTTP requests hit the gRPC-Gateway; gRPC calls flow through interceptors (logging, auth, authorization) to handlers that delegate into usecase interactors and domain services, which rely on repository interfaces implemented in infrastructure packages.

## Configuration & Environment

- Runtime configuration populated from environment variables via `internal/infrastructure/environment/env.go` using `caarlos0/env`.
- Key groups: application (`PORT`, `ENV`, `MIN_LOG_LEVEL`), database (`DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_DATABASE`, `DB_LOG_ENABLE`), GCP (`GCP_PROJECT_ID`, bucket name, Firebase client key, Identity Platform tenant IDs), X API credentials.
- `internal/pkg/logger` builds a Zap logger whose minimum level derives from `MIN_LOG_LEVEL`; logger is attached to contexts for interceptors and usecases.

## Layered Architecture Details

### Domain Layer (`internal/domain`)

- **Model**: immutable-ish structs (e.g., `model.Asset`, `model.User`, `model.Staff`) with factory helpers and behavior (status transitions, derived IDs).
- **Model factory**: `internal/domain/model/factory.NewFactory` seeds model instances via `faker.FakeData` (mirroring the admin setup) and then overrides key fields manually—prefer tweaking the overrides instead of reintroducing constructor helpers.
- **Repository interfaces**: define persistence contracts (`internal/domain/repository`). Concrete implementations live under infrastructure (`postgresql/repository`, `firebase/repository`, `gcs/repository`). Mock interfaces generated under `internal/domain/repository/mock` for testing.
- **Services**: orchestrate complex domain operations (e.g., `service.Asset`, `service.Staff`, `service.Question`) often combining multiple repositories/caches.
- **Errors**: centralized typed errors in `internal/domain/errors` returning goerr instances with codes/messages used by interceptors for gRPC status mapping.
- **Cache contracts**: e.g., `cache.AssetPath` executed by Postgres-backed cache storing presigned asset metadata.
- **Common utilities**: `internal/pkg` packages for IDs, UUID, time (`now`), nullable helpers, validation.

### Usecase Layer (`internal/usecase`)

- Exposes interactors grouped by audience (admin, staff, user, task, debug). Each interactor validates inputs (`internal/usecase/input/*` structs using `validation.Validate`) and coordinates domain services/repositories inside transactions.
- Transaction boundary defined by `repository.Transactable`; Postgres implementation wraps functions with retry-aware `RunTx` (`internal/infrastructure/postgresql/transactable`).
- Within write transactions, keep `Preload` disabled by default; enable it only when the usecase needs to inspect domain `ReadonlyReference` data inside the transaction. After mutations, perform a fresh load outside the transaction (with `Preload: true`) so readonly references are hydrated for response shaping.
- Outputs under `internal/usecase/output` when custom view models are needed.
- Authentication interactor mediates Firebase ID token verification for staff/admin/user contexts.
- Example: `adminStaffInteractor.Create` fetches tenant, validates asset via `service.Asset`, creates staff via `service.Staff`, and reloads the persisted entity for response.

### Infrastructure Layer

- **Dependency wiring**: `internal/infrastructure/dependency.Dependency.Inject` constructs all repositories, services, and interactors, wiring Firebase, GCS, X API clients, and Postgres client.
- **Transport**:
  - gRPC server (`internal/infrastructure/grpc/run.go`) registers handlers for admin/public/debug APIs, sets keepalive configs, and applies interceptors: request logging (structured zap logs + goerr mapping to gRPC status), session authentication (extracts Bearer tokens, verifies ID tokens via usecases, saves claims onto context), and simple authorization checks.
  - HTTP Gateway (`internal/infrastructure/http/run.go`) hosts gRPC-Gateway mux with custom JSON marshaler, registers gRPC handlers, adds `/` ping endpoint, and wraps mux with permissive CORS middleware.
- **Persistence**:
  - Postgres client configured in `internal/infrastructure/postgresql/client.go`; SQLBoiler generated models under `internal/infrastructure/postgresql/internal/dbmodel` power repository implementations.
  - Repository packages translate between DB models and domain models (marshal/unmarshal helpers under `internal/infrastructure/postgresql/internal/marshaller`).
- **External services**:
  - Firebase auth client wrappers provide staff/admin/user authentication repositories (verify tokens, CRUD users, custom claims).
  - GCS repository handles signed URL generation and asset uploads.
  - X integration (`internal/infrastructure/x/x`) fetches profile data for onboarding.
- **CLI Tasks**: `task create-root-staff` and `task create-admin` parse flags, load dependencies, and invoke corresponding interactors with current time (`internal/pkg/now`).
- **Migrations**: Goose-based runner manages migrations embedded via `db/postgresql/migrations`; constants sync reads YAML and upserts enumerations; schema extraction outputs `db/postgresql/schema.sql` and Mermaid diagrams.

## Protocol & Schema Assets

- Protobuf definitions under `schema/proto`; Buf configuration (`buf.gen.yaml`, `buf.work.yaml`) generates gRPC stubs into `internal/infrastructure/grpc/pb` and OpenAPI specs into `schema/openapi/oshipo/...`.
- Regenerate stubs with `make generate.buf`. SQLBoiler models regenerated via `make generate.sqlboiler.postgresql` after migrations.

## Database & Data Modeling

- PostgreSQL schemas managed via Goose migrations in `db/postgresql/migrations` and constants in `db/postgresql/constants` (embedded YAML).
- SQLBoiler view/table files in `internal/infrastructure/postgresql/internal/dbmodel` are generated; avoid manual edits.
- Transaction helper retries on deadlocks (`pq` error codes `40P01`, `55P03`, `40001`).

## CLI & Development Recipes

- Build CLI: `make build` (binary at `./.bin/app-cli`).
- Run API locally: `docker-compose up -d` (db), `make migrate.up`, then `make http.dev` (Air hot-reload) or `go run ./cmd/app http-server run`.
- Health check: `curl http://localhost:8080/` returns `pongpong`.
- Linting: `make lint.go`, `make lint.proto`.
- Tests: `make test` (runs `go test ./internal/...`).
- Mock generation: `make generate.mock` (runs `go generate` across repository).
- Admin bootstrap: `go run ./cmd/app task create-admin -e <email> -p <pwd> -d <display> -r <role>`.

## Testing & Observability

- Unit tests present in domain/usecase packages (e.g., `internal/domain/service/*_test.go`, `internal/usecase/*_test.go`). Prefer using mocks from `internal/domain/.../mock` for isolation.
- Domain service tests (e.g., wallet service) must stay table-driven (`tests := map[string]testcaseFunc`) with `testcase` structs containing only `args`, `service`, and `want`; avoid introducing ad-hoc matcher helpers and rely on existing gomock matchers such as `gomock.Eq` alongside the generated mocks. When deterministic IDs are required, call `id.Mock()` directly without overriding or restoring `id.New` manually.
- Logging uses Zap with structured fields including operation IDs; interceptors map domain errors to gRPC codes and attach request IDs to `errdetails.RequestInfo` for clients.
- When adding handlers, ensure gRPC methods propagate context with logger via `request_interceptor` for consistent tracing.

## Contribution Guidelines for Agents

- Favor the existing layer boundaries: domain models stay pure (no infrastructure imports), usecases orchestrate business rules, infrastructure wires dependencies and IO.
- Always validate incoming params using the input structs and `validation.Validate`; surface errors via domain error types so interceptors format responses correctly.
- When defining domain service APIs (packages under `internal/domain/service`), name request/response structs using the `XxxParam` / `XxxResult` pattern (e.g., `WalletGrantPointsParam`), not `Input` / `Output`.
- Wrap DB mutations inside `repository.Transactable.RWTx` to get retry and context-aware executors; use `transactable.GetContextExecutor` within repositories to access the active transaction.
- Keep transaction boundaries in the usecase layer: the interactor should invoke `repository.Transactable.RWTx/ROTx`, while downstream domain services must assume a transaction context is already in place and must not start their own transactions.
- Leverage services (asset, staff, question, etc.) instead of duplicating logic; they encapsulate cross-repository coordination.
- For usecases, always set `Preload` to `true` on repository `Get`/`List` options so readonly relations are hydrated, and after `Create`/`Update` flows, issue a final `Get` to build the response payload.
- Even when no URLs need to be populated today, implement and invoke the appropriate `assetService.BatchXXX` helper from usecases to keep response shaping consistent.
- Respect generated code ownership: Protobuf/SQLBoiler outputs and embedded assets should be regenerated via Make targets, not edited manually.
- Attach new dependencies in `internal/infrastructure/dependency.Inject` and update relevant handlers/interactors; keep constructor signature order consistent for readability.
- For new gRPC endpoints, extend protobufs under `schema/proto`, regenerate stubs, implement handler functions in infrastructure, and delegate to usecase interactors for business logic.
- Route all state mutations and validation logic through domain models/services; usecases must call domain helpers rather than mutating `model` fields or reimplementing validation directly.
- Do not introduce package-level helper functions/methods in usecase or domain/service packages for business logic; extend the relevant domain model (e.g., `model.Wallet`) instead so that state mutations live alongside their data.
- When adding domain-specific helper logic, attach it as a method on the relevant domain model (e.g., make expiration helpers receiver methods on `model.Wallet`) rather than free functions, even within `domain/model`.
- When adding gomock expectations, mirror existing patterns: only use `gomock.Any()` for `context.Context`, avoid `DoAndReturn`, prefer exact argument matching with `Return`, and skip helpers like `gomock.AssignableToTypeOf`.
- When writing rate limit tests, build the domain objects via helpers in `internal/domain/model/factory` and pass those concrete values directly into gomock expectations instead of introducing custom matchers so comparisons stay consistent across test suites.
- Keep test `want` structs minimal: include only the fields asserted in the test flow (typically `got`/`expectedResult`) and move custom checks into small helper closures or inline assertions.
- Write table-driven tests using `tests := map[string]testcaseFunc` where `type testcaseFunc func(ctx context.Context, ctrl *gomock.Controller) testcase`.
- Define `testcase` structs with exactly `args`, `usecase` (or `service` for domain-layer tests), and `want` fields; construct mocks and test data inside each `testcaseFunc` before returning the struct.
- Instantiate all test fixtures inside the returned `testcase` (e.g., factories, mocked repositories/services) so individual cases stay isolated and self-contained.
- Do not embed assertion functions inside `want`; compare the actual result with the expected object directly using equality assertions.
- Use `internal/pkg/logger` helpers to add contextual fields instead of ad-hoc logging; propagate contexts down call stacks.
- Within the UserService code paths (inputs, handlers, usecases), name the current authenticated user's identifier `UserID` to keep terminology consistent.
- UserService handlers must obtain the current time from the request context (e.g., `request_interceptor.GetRequestTime`) and pass it through inputs; do not call `now.Now`, `time.Now`, or similar clock helpers inside usecase or domain layers.

## Development Philosophy

- Keep persistence logic cohesive: when an aggregate owns child entities, write repositories so create/update operate transactionally on the whole aggregate; avoid “readonly reference” patterns unless mandated by performance.
- Preload semantics should be explicit and consistent: use a boolean preload flag on repository queries, but allow aggregates that must always hydrate critical relations to override flags for correctness; document any such exceptions near the helper.
- Marshallers/DTO mappers are responsible for validation and normalization when crossing boundaries (DB ↔ domain ↔ transport); keep them thin, deterministic, and co-located with generated models to reduce drift.
- Prefer declarative guides (like this file) that capture principles instead of one-off case notes; when adding new rules, generalize them to future work rather than binding them to a single feature.
