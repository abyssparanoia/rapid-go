---
name: code-investigation
description: Efficient codebase investigation in Cursor. ALWAYS use this skill BEFORE modifying existing code, fixing bugs, or adding features to existing modules. Required for understanding code structure, tracing data flow, and impact analysis.
---

# Code Investigation Guide

This skill describes a fast, low-noise way to understand the codebase using tools available in Cursor:
- `grep` (exact text / regex)
- `codebase_search` (semantic search)
- `glob_file_search` / `list_dir` (find files and explore structure)
- `read_file` (read only what you need)

## Tool Selection Strategy

- If you **know the exact string/symbol** → use `grep`
- If you **know what it does but not where** → use `codebase_search`
- If you **only know the filename pattern** → use `glob_file_search`
- Once you have the file → use `read_file` (prefer partial reads for big files)

## Rapid-go: Key Places to Look

### API contract (proto)

`schema/proto/rapid/{admin_api|public_api|debug_api}/v1/*.proto`

### gRPC handlers & interceptors

`internal/infrastructure/grpc/internal/handler/{admin|public|debug}/`

`internal/infrastructure/grpc/internal/interceptor/`
- `session_interceptor`: staff session context
- `authorization_interceptor`: access control for admin API
- `request_interceptor`: request time, error mapping helpers

### Usecases

`internal/usecase/*_impl.go`, `internal/usecase/input/*`

### DB implementations (multiple backends)

`internal/infrastructure/{mysql|postgresql|spanner}/`
- `internal/dbmodel/` (generated for mysql/postgresql; yo for spanner)
- `internal/marshaller/`
- `repository/`
- `transactable/`

## Investigation Workflows

### Workflow A: “How does an HTTP request become a DB write?”

1. **Find the RPC / HTTP path** in proto (`schema/proto/rapid/**`):
   - `grep` for `option (google.api.http)` or the specific path (e.g. `/admin/v1/tenants`)

2. **Find the handler method**:
   - `grep` for the RPC name in `internal/infrastructure/grpc/internal/handler/`

3. **Trace into the usecase interactor**:
   - the handler will call a usecase interactor (e.g. `AdminTenantInteractor`)

4. **Trace into repositories**:
   - usecase calls domain repositories (`internal/domain/repository/**`)
   - implementations live under `internal/infrastructure/{mysql|postgresql|spanner}/repository/`

### Workflow B: “Where is auth/authorization checked?”

1. Search for interceptors:
   - `list_dir internal/infrastructure/grpc/internal/interceptor/`

2. Identify session creation and authorization:
   - `session_interceptor/*` (extracts and stores staff claims)
   - `authorization_interceptor/*` (e.g. Admin API gating)

### Workflow C: Impact analysis before changing a type

1. `grep` the type name (e.g. `StaffClaims`, `BaseGetOptions`) in `internal/`
2. Check callers first (usecases/handlers), then implementations (repositories)
3. For interface changes: also check mock generation targets and tests

## Tips

- Prefer starting with **proto → handler → usecase → repository** for request flows.
- When you see a misleading doc path, verify it exists with `list_dir` before trusting it.

