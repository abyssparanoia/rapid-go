# RAPID GO

## Project Overview

Go-based gRPC API server with HTTP API via gRPC-Gateway.

**Architecture**: Layered Architecture + Domain-Driven Design (DDD)

## Tech Stack

| Category | Technology |
|----------|------------|
| Language | Go 1.25+ |
| API | gRPC + gRPC-Gateway |
| Database | MySQL / PostgreSQL (SQLBoiler ORM), Spanner (yo) |
| Auth | Firebase Auth, AWS Cognito |
| Storage | Google Cloud Storage, AWS S3 |

## Directory Structure

```
├── cmd/app/                    # Entry point
├── db/
│   ├── mysql/                  # MySQL migrations & SQLBoiler config
│   │   ├── migrations/         # Goose migrations
│   │   └── constants/          # Constant table definitions (YAML)
│   ├── postgresql/             # PostgreSQL migrations & SQLBoiler config
│   └── spanner/                # Spanner migrations (wrench)
├── schema/proto/rapid/         # Protocol Buffers definitions
│   ├── admin_api/v1/           # Admin API definitions
│   ├── debug_api/v1/           # Debug API definitions
│   └── public_api/v1/          # Public API definitions
├── internal/
│   ├── domain/                 # Domain layer (pure business logic)
│   │   ├── model/              # Domain models & entities
│   │   ├── repository/         # Repository interfaces
│   │   ├── service/            # Domain services
│   │   └── errors/             # Domain errors
│   ├── usecase/                # Application layer
│   │   ├── input/              # Input DTOs with validation
│   │   ├── output/             # Output DTOs
│   │   └── *_impl.go           # Interactor implementations
│   └── infrastructure/         # Infrastructure layer
│       ├── mysql/              # MySQL adapter (client, repositories, sqlboiler models)
│       ├── postgresql/         # PostgreSQL adapter
│       ├── spanner/            # Spanner adapter
│       ├── grpc/
│       │   └── internal/
│       │       └── handler/    # gRPC handlers by actor
│       │           ├── admin/  # Admin API handlers
│       │           ├── debug/  # Debug API handlers
│       │           └── public/ # Public API handlers
│       ├── http/               # gRPC-Gateway + HTTP server
│       ├── cognito/            # AWS Cognito integration
│       ├── firebase/           # Firebase Auth integration
│       ├── gcs/                # Google Cloud Storage integration
│       ├── s3/                 # AWS S3 integration
│       ├── redis/              # Redis cache integration
│       └── dependency/         # DI configuration
```

## Key Commands

| Command | Description |
|---------|-------------|
| `make http.dev` | Start server with hot reload |
| `make migrate.create` | Create new migration file |
| `make migrate.up` | Run migrations + generate SQLBoiler |
| `make generate.buf` | Generate Protocol Buffers code |
| `make generate.mock` | Generate mock files |
| `make test` | Run all tests |
| `make lint.go` | Lint Go code |

## Layer Dependencies

```
Infrastructure → Usecase → Domain
     ↓              ↓         ↓
  (implements)  (uses)    (no deps)
```

- **Domain**: No external dependencies - pure business logic only
- **Usecase**: Depends only on domain layer
- **Infrastructure**: Implements interfaces defined in domain/usecase

## Rules & Guidelines

Detailed coding rules are organized by theme in `.claude/rules/`:

| Rule File | Applies To | Description |
|-----------|------------|-------------|
| `domain-model.md` | `internal/domain/model/**` | Domain model conventions, state change methods |
| `domain-service.md` | `internal/domain/service/**` | Domain service patterns |
| `domain-errors.md` | `internal/domain/errors/**` | Error definition patterns |
| `repository.md` | `repository/** + marshaller/**` | Repository & marshaller patterns |
| `usecase-interactor.md` | `internal/usecase/**` | Interactor implementation, external service sync |
| `grpc-handler.md` | `internal/infrastructure/grpc/**` | gRPC handler patterns |
| `testing.md` | `**/*_test.go` | Testing conventions |
| `proto-definition.md` | `schema/proto/**` | Protocol Buffers style |
| `migration.md` | `db/{mysql,postgresql,spanner}/**` | Database migration patterns |
| `dependency-injection.md` | `internal/infrastructure/dependency/**` | DI configuration |
| `invitation-workflow.md` | `*invitation*` | Invitation/approval flow patterns |
| `external-service-integration.md` | `cognito/**`, `firebase/**`, `*authentication*` | Auth (Cognito/Firebase) integration patterns |
| `webhook-implementation.md` | `webhook/**`, `internal/infrastructure/http/internal/handler/webhook_*` | Webhook endpoint patterns (HTTP → gRPC routing) |
| `job-system.md` | `job/**`, `cmd/app/internal/task_cmd/process_job_cmd/` | Async job queue patterns (SNS/SQS → AWS Batch) |
| `worker-pattern.md` | `worker/**`, `cmd/app/internal/worker_cmd/` | Background worker patterns (SQS/Pub/Sub subscribers) |
| `cli-command-pattern.md` | `cmd/app/internal/task_cmd/**`, `internal/usecase/task_*` | CLI command implementation patterns (`./app task` commands) |

## Skills

Available automation skills in `.claude/skills/`:

| Skill | Description |
|-------|-------------|
| `code-investigation` | Efficient codebase investigation using available tools |
| `crud-implementation` | Complete workflow overview (references focused skills below) |
| `add-database-table` | Create migration file and constant tables |
| `add-domain-entity` | Create domain model, repository interface, and implementation |
| `add-api-endpoint` | Create usecase, proto definition, and gRPC handler |
| `review-diff` | Review & auto-fix current branch diff against main/master. Catches AI anti-patterns and rule violations, then fixes them automatically |
| `audit-rules` | Audit the **entire codebase** for rule-compliance. Partitions source by category, runs convention-reviewer + test-reviewer in parallel (audit mode), auto-fixes all violations, and opens a PR. Run `/audit-rules` (all) or `/audit-rules domain` (single partition) |
| `fix-review-comments` | Fetch unresolved GitHub PR review comments and auto-fix the code. Run `/fix-review-comments` (current branch PR) or `/fix-review-comments 123` (specific PR) |
| `create-pull-request` | PR creation guide with branch naming and body templates |
| `sync-claude-config` | Bidirectionally sync `.claude/` content with the rapid-go template (or a derived project added via `claude --add-dir`); opens a PR in each repo |

**Implementation Workflow**: `add-database-table` → `add-domain-entity` → `add-api-endpoint`

**Review Workflow**: Use `review-diff` to auto-fix issues on the current branch, then `create-pull-request` to create the PR

**Full Audit Workflow**: Use `audit-rules` for a periodic whole-codebase convention sweep (counterpart to `review-diff` but scoped to all files, not a diff)

**Post-Review Workflow**: Use `fix-review-comments` to address reviewer feedback automatically

**Investigation Workflow**: Use `code-investigation` before modifying existing code

## Common Errors & Solutions

| Error | Cause | Solution |
|-------|-------|----------|
| `undefined: dbmodel.ExampleWhere` | SQLBoiler not regenerated | Run `make migrate.up` |
| `cannot find package "github.com/.../mock_repository"` | Mocks not generated | Run `make generate.mock` |
| `undefined: pb.Example` | Proto not generated | Run `make generate.buf` |
| `foreign key constraint violation` | Missing constant table entry | Add record to `db/mysql/constants/*.yaml` then `make migrate.up` |
| `transaction has already been committed` | Transaction misuse | Ensure single `RWTx`/`ROTx` per operation flow |
| `make` command not found or not working | Shell function override in Claude Code | Use full path `/usr/bin/make` instead of `make` |

## PR Checklist

Before creating a PR, verify:

- [ ] `make lint.go` passes
- [ ] `make test` passes
- [ ] Migrations are idempotent (can run up/down)
- [ ] Proto changes are backward compatible
- [ ] New interactors registered in `dependency.go`
- [ ] New handlers added to gRPC server registration
- [ ] Mocks regenerated if repository interfaces changed
