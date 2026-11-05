---
description: How implementation structs are shaped in every layer — no package-level private functions (logic that needs the struct's dependencies is a method), pointer receivers, gather inputs then build the result once
paths:
  - "internal/usecase/**/*.go"
  - "internal/domain/service/**/*.go"
  - "internal/infrastructure/**/*.go"
---

# Implementation Structure Guidelines (implementation layers)

Applies to every struct that implements an interface: usecase interactors, domain services,
repositories, external-service clients (`internal/infrastructure/<service>/repository/`,
`internal/infrastructure/<service>/<service>/`), gRPC handlers, workers, CLI tasks.
Layer rules (`usecase-interactor.md`, `domain-service.md`, `repository.md`, …) add layer-specific
detail and may be **stricter** (e.g. `usecase-interactor.md` forbids private methods on interactors
altogether); this file is the common baseline they share.

Out of scope — not interface implementations, so §1 does not apply:

- `internal/domain/model/` (constructors, enums, slice types — see `domain-model.md`),
  `internal/domain/errors/`, interface-only packages (`internal/domain/repository/`, `internal/domain/cache/`, …)
- `internal/pkg/` (generic utilities are package-level functions by design — see `package-placement.md`)
- Generated code (`internal/infrastructure/**/internal/generated/`, SQLBoiler `dbmodel`, `*.pb.go`)
- Infrastructure plumbing that is not a domain interface implementation: `internal/infrastructure/http/`
  (server, JSON marshaler, middlewares), `internal/infrastructure/grpc/internal/interceptor/`,
  `*/transactable/`, `*/migration/`, `internal/infrastructure/cmd/` (cobra wiring),
  `internal/infrastructure/dependency/`, `internal/infrastructure/environment/`, `*/client.go`
  constructors, and the hand-written `internal/client/` of an external SDK wrapper

## 1. No package-level private functions — logic that needs the struct's dependencies is a method

An implementation package must **not** define package-level private (lowercase, receiver-less)
functions that do a step of a public method. The two shapes to reject on sight:

- a function that **takes the struct's dependency as a parameter** (`cli *billing.APIClient`,
  `repo repository.Staff`, `tx *sql.Tx`, …) — the struct already holds it, so this is a method
  written as a function;
- a function that **performs I/O** (external API call, repository / DB call) outside a method.

Private **methods on the struct** are allowed. Pick the smallest unit that reads well; when
`gocognit` / `funlen` fire, extract a private method on the receiver, never a package-level function.

Why:
- A receiver-less helper exists so two structs can share it. The sharing forces its dependencies
  into the parameter list, the next caller needs one more parameter, and the package fills with
  `findXxx` → `findXxxInPage` → `getXxx` chains nobody can tell apart or test in isolation.
- Behavior is tested through the interface. A method on the struct is covered by the struct's
  tests; a package-level function is not reachable through any interface.

If two structs really need the same lookup, **duplicate the loop in each struct's method**. If the
shared part is pure (no I/O), it belongs in `marshaller/` or on a `domain/model` type instead.

```go
// BAD - package-level functions that take the client and perform I/O, chained two levels deep,
//       shared by two capability structs in the same package
func findActiveSubscription(ctx context.Context, cli *billing.APIClient, tenantID string) (*billing.Subscription, error) {
    for page := int64(1); page <= maxPages; page++ {
        sub, hasNext, err := findActiveSubscriptionInPage(ctx, cli, tenantID, page)
        ...
    }
}
func findActiveSubscriptionInPage(ctx context.Context, cli *billing.APIClient, tenantID string, page int64) (*billing.Subscription, bool, error) {
    res, _, err := cli.SubscriptionAPI.List(ctx).TenantId(tenantID).Page(page).Execute()
    ...
}
func getPlanName(ctx context.Context, cli *billing.APIClient, planID string) (string, error) { ... }

func (r *tenantSubscription) GetByTenantID(ctx context.Context, tenantID string) (*model.TenantSubscription, error) {
    sub, err := findActiveSubscription(ctx, r.cli, tenantID)
    ...
}

// GOOD - the public method calls the client it holds; a private method on the receiver is fine
//        when it reads better than inlining
func (r *tenantSubscription) GetByTenantID(ctx context.Context, tenantID string) (*model.TenantSubscription, error) {
    var active *billing.Subscription
    for page := int64(1); active == nil && page <= maxPages; page++ {
        res, httpRes, err := r.cli.SubscriptionAPI.List(ctx).TenantId(tenantID).Page(page).Execute()
        if err != nil {
            return nil, errors.InternalErr.Wrap(err)
        }
        httpRes.Body.Close()
        if idx := slices.IndexFunc(res.Subscriptions, func(s billing.Subscription) bool { return s.Status == billing.ACTIVE }); idx >= 0 {
            active = &res.Subscriptions[idx]
        }
        if !res.Pagination.HasNext {
            break
        }
    }
    if active == nil {
        return nil, errors.TenantSubscriptionNotFoundErr.New().WithValue("tenant_id", tenantID)
    }

    plan, err := r.getPlan(ctx, active.PlanId) // private method on the receiver: allowed
    if err != nil {
        return nil, err
    }

    return marshaller.SubscriptionToModel(*active, plan), nil
}

func (r *tenantSubscription) getPlan(ctx context.Context, planID string) (*billing.Plan, error) { ... }
```

Where logic should live instead of a package-level function:

| Situation | Solution |
|-----------|----------|
| A step of one public method | Inline it, or a private method on the receiver |
| Logic on a single entity (predicate, derived value, state change of a DB-backed entity) | Method on the `domain/model` type |
| Logic spanning several entities | `domain/service` (its public method, or a private method on the service struct) |
| Conversion between external/DB structs and domain models | Exported function in the package's `marshaller/` (takes every input it needs as a parameter) |
| Generic, domain-agnostic utility | `internal/pkg/` (see `package-placement.md`) |

### Allowed package-level private functions

Only cross-cutting **pure** functions that take none of the struct's dependencies, perform no I/O,
and carry no business logic:

- error mapping shared by every repository in the package (`repository/error.go` `handleError`)
- query-option builders (`repository/base.go` `addForUpdateFromBaseGetOptions` / `…ListOptions`)
- request setup every method of an external-API package needs, such as API-key / auth-header
  injection into `ctx` (`with<Actor>Context` in a `<actor>_base.go`); it must not call any API

Anything that would need the client, a repository, or a transaction to be passed in does not qualify.

## 2. Pointer receivers for implementation structs

Implementation structs use **pointer receivers** and their constructors return `&struct{...}`:

```go
type tenantSubscription struct{ cli *billing.APIClient }

func NewTenantSubscription(cli *billing.APIClient) repository.TenantSubscription {
    return &tenantSubscription{cli: cli}
}

func (r *tenantSubscription) GetByTenantID(ctx context.Context, tenantID string) (*model.TenantSubscription, error) { ... }
```

This also applies to the `*Debug` / E2E stand-ins of the same interface. Apply it to new and
modified structs; do not sweep existing value-receiver structs just to satisfy this rule.
Value receivers are for value types only: enum/string types (`String()`, `Valid()`) and slice
types (`IDs()`, `MapByID()`, `FilterBy…()`).

## 3. Gather every input first, then build the result once

A public method **resolves everything the result needs, then constructs it in a single call**
(domain constructor or marshaller). It never builds a struct and then patches fields on it:

- If a field comes from a second query or API call (e.g. the list response omits the plan name and
  a separate `GetPlan` call provides it), make that call **before** construction and pass the value
  into the constructor / marshaller as a parameter.
- Do **not** assign model fields after construction (`m.PlanName = …`) in usecase or
  infrastructure code.
- Do **not** add `Set{X}` / `Update{X}` methods to a read model fetched from an external system
  just to patch it after construction. Those methods are for state transitions of DB-backed
  entities, per `domain-model.md` "State Change Methods".

```go
// BAD - build, then patch
sub := marshaller.SubscriptionToModel(res.Subscription)
if !sub.PlanName.Valid {
    plan, _, err := r.cli.PlanAPI.Get(ctx, sub.PlanID).Execute()
    ...
    sub.PlanName = null.StringFrom(plan.Name)
}
return sub, nil

// GOOD - resolve, then build once
plan := res.Subscription.Plan
if plan == nil {
    planRes, _, err := r.cli.PlanAPI.Get(ctx, res.Subscription.PlanId).Execute()
    ...
    plan = &planRes.Plan
}
return marshaller.SubscriptionToModel(res.Subscription, plan), nil
```

Marshaller functions therefore take **every input they need as parameters** and return a fully
populated model; they do not reach back into the client to fetch anything.
