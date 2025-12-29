# Review Checklists by Category

Detailed checklists for each file category. Apply the relevant sections based on changed files.

## Domain Model (`internal/domain/model/**`)

- [ ] Entity has `ReadonlyReference` struct pointer for relations
- [ ] Constructor uses `id.New()` for ID generation
- [ ] Constructor sets both `CreatedAt` and `UpdatedAt` to same time
- [ ] Update methods use `null.String`, `null.Int64` for optional fields
- [ ] Update methods always update `UpdatedAt`
- [ ] Status/Enum types have `Unknown` as first constant
- [ ] Status types have `String()` and `Valid()` methods
- [ ] State changes done via domain methods (not direct field assignment)
- [ ] Role types have helper methods like `IsRoot()`, `IsNormal()`
- [ ] Slice types have `IDs()` and `MapByID()` helpers
- [ ] Type aliases defined: `{Entity}MapByID`, `{Entity}s`

## Repository Interface (`internal/domain/repository/**`)

- [ ] Has `//go:generate` directive for mockgen
- [ ] Query structs use `null.String`, `null.Uint64` for optional string/numeric fields
- [ ] Query structs use `nullable.Type[T]` for optional enum/custom type fields
- [ ] `BaseGetOptions`, `BaseBatchGetOptions`, `BaseListOptions` properly embedded

## Repository Implementation (`internal/infrastructure/**/repository/**`)

- [ ] Uses `transactable.GetContextExecutor(ctx)` for all queries
- [ ] `Get` method handles `OrFail` correctly (nil vs error on not found)
- [ ] `Get` method handles `ForUpdate` option
- [ ] `List` method applies pagination with `Page` and `Limit`
- [ ] `List` method validates sort key with `query.SortKey.Valid && query.SortKey.Ptr().Valid()`
- [ ] Preload helper defined if relations exist

## Marshaller (`internal/infrastructure/**/internal/marshaller/**`)

- [ ] `ToModel` handles relations via `R != nil` check
- [ ] `ToDBModel` sets `R: nil, L: struct{}{}`
- [ ] Enum conversions handle all cases including Unknown/default
- [ ] Slice conversion function defined (`{Entity}sToModel`)
- [ ] Related entity's `ReadonlyReference` always nil (no recursive loading)

## Usecase Interactor (`internal/usecase/**`)

- [ ] Interface has `//go:generate` directive
- [ ] Implementation uses dependency injection via constructor
- [ ] All methods start with `param.Validate()` check
- [ ] Write operations wrapped in `transactable.RWTx`
- [ ] Get before update uses `ForUpdate: true` for locking
- [ ] State changes use domain methods (not direct field assignment)
- [ ] IdP sync (StoreClaims/DeleteUser) happens within transaction
- [ ] On delete: IdP deletion before database deletion
- [ ] Final return fetches entity with `Preload: true` for fresh data
- [ ] Asset service called even if no assets currently exist

## Input Struct (`internal/usecase/input/**`)

- [ ] Named as `{Actor}{Action}{Resource}`
- [ ] Has `RequestTime` field with `validate:"required"`
- [ ] Has `Validate()` method that uses `validation.Validate()`
- [ ] Optional update fields use `nullable.Type[T]` (not pointers)
- [ ] Validation includes business rule checks for optional fields

## gRPC Handler (`internal/infrastructure/grpc/internal/handler/**`)

- [ ] Gets claims via `session_interceptor.Require{Actor}SessionContext(ctx)`
- [ ] Gets request time via `request_interceptor.GetRequestTime(ctx)`
- [ ] Converts proto to input struct correctly
- [ ] Handles optional proto fields with `if req.Field != nil` pattern
- [ ] Returns error directly (interceptor handles conversion)
- [ ] Uses marshaller for domain-to-proto conversion

## Handler Marshaller (`internal/infrastructure/grpc/internal/handler/**/marshaller/**`)

- [ ] Each resource has its own file (not combined)
- [ ] `ToPb` handles nil input
- [ ] `ToPb` uses variable declaration pattern for optional/nullable fields
- [ ] Enum conversions have both `ToPb` and `ToModel` directions
- [ ] Slice conversion function defined (`{Entity}sToPb`)
- [ ] All proto fields are explicitly mapped (check for omissions)

## Proto Definition (`schema/proto/**`)

- [ ] Enum values start with `{ENUM_NAME}_UNSPECIFIED = 0`
- [ ] Field names use snake_case
- [ ] Request/Response named as `{Action}{Resource}Request/Response`
- [ ] HTTP annotations follow REST patterns
- [ ] Optional fields marked with `optional` keyword
- [ ] List requests have `page` and `limit` fields
- [ ] Required fields annotated with `openapiv2_schema`

## Migration (`db/**/migrations/**`)

- [ ] Has both `+goose Up` and `+goose Down` sections
- [ ] Column types match Go types (see migration.md for mapping)
- [ ] Foreign keys named as `{table}_fkey_{column}`
- [ ] Indexes named as `{table}_idx_{column}`
- [ ] Unique constraints named as `{table}_uq_{columns}`
- [ ] `TIMESTAMPTZ` used for all timestamps (not TIMESTAMP)
- [ ] Constant tables have corresponding YAML in `db/**/constants/`

## Tests (`**/*_test.go`)

- [ ] Uses table-driven tests with `map[string]testcaseFunc`
- [ ] Test function returns `(args, usecase/service, want)` tuple
- [ ] Mock setup uses closure pattern with `func(ctrl *gomock.Controller)`
- [ ] Transactable mock uses `DoAndReturn` to execute function
- [ ] Error assertions use `assert.ErrorIs(t, err, want.err)`

## Dependency Injection (`internal/infrastructure/dependency/**`)

- [ ] New repository registered in Dependency struct
- [ ] New interactor registered in Dependency struct
- [ ] Constructor call added in `Inject()` method
- [ ] Handler updated to include new interactor
- [ ] Injection order: clients -> transactable -> repos -> services -> interactors -> handlers

## External Service Integration

- [ ] Claims model uses `null.String` and `nullable.Type[T]` for optional fields
- [ ] `StoreClaims` called after entity creation/update
- [ ] `DeleteUser` called before database deletion
- [ ] All IdP operations within transaction boundary

## Invitation Workflow (`*invitation*`)

- [ ] Status enum includes: Pending, Accepted, Rejected, Invalidated
- [ ] `ExpiresAt` field set in constructor
- [ ] State transition methods validate current state
- [ ] `IsExpired()` helper method exists
- [ ] Email sent within transaction
