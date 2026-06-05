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
| `make e2e` | Run E2E tests (requires server running) |
| `make ao-start` | Start Agent Orchestrator |
| `make ao-stop` | Stop Agent Orchestrator |
| `make ao-spawn ISSUE=123` | Spawn agent on specific issue |
| `make ao-watch` | Open ao-tui dashboard |

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
| `device-group-authorization.md` | `session_interceptor/**`, `handler/**`, `device_group_*` | Device group 3-layer authorization patterns |
| `package-placement.md` | `internal/pkg/**`, `internal/domain/**` | Where to put new packages — pkg is domain-agnostic only |
| `object-storage-paths.md` | `internal/domain/model/asset.go`, `s3/**`, `gcs/**`, `job_*` | S3/GCS path prefix 集約と private/ prefix 規約 |

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
| `fix-review-comments` | Fetch unresolved GitHub PR review comments and auto-fix the code. Run `/fix-review-comments` (current branch PR) or `/fix-review-comments 123` (specific PR) |
| `create-pull-request` | PR creation guide with branch naming and body templates |
| `sync-claude-config` | Bidirectionally sync `.claude/` content with the rapid-go template (or a derived project added via `claude --add-dir`); opens a PR in each repo |
| `create-ao-issue` | Agent Orchestrator 用の GitHub Issue を正しいフォーマットで作成 |

**Implementation Workflow**: `add-database-table` → `add-domain-entity` → `add-api-endpoint`

**Review Workflow**: Use `review-diff` to auto-fix issues, then `create-pull-request` to create the PR

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

## Implementation Preflight

新規ファイル作成・大きな変更の前に必ず実施:

1. **同カテゴリの既存ファイルを Read する** — 記憶ではなく実物を見てパターンを確認
2. **該当する `.claude/rules/*.md` を参照する** — 対象レイヤーのルールを確認してから実装
3. **`internal/pkg/` を検索する** — 新しいユーティリティ定義前に既存実装を確認
4. **新規パッケージを置く位置に迷ったら `package-placement.md` を確認** — `internal/pkg/` はドメイン非依存のみ、ドメイン依存は `internal/domain/` 配下
5. **テスト作成時は `testing.md` と `ai-antipatterns.md` を参照する**

## PR Checklist

Before creating a PR, verify:

- [ ] `make lint.go` passes
- [ ] `make test` passes
- [ ] Migrations are idempotent (can run up/down)
- [ ] Proto changes are backward compatible
- [ ] New interactors registered in `dependency.go`
- [ ] New handlers added to gRPC server registration
- [ ] Mocks regenerated if repository interfaces changed
