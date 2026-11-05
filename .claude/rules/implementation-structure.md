---
description: How implementation structs are shaped in every layer — public interface methods only, pointer receivers, gather inputs then build the result once
paths:
  - "internal/usecase/**/*.go"
  - "internal/domain/service/**/*.go"
  - "internal/infrastructure/**/*.go"
---

# Implementation Structure Guidelines (implementation layers)

Applies to every struct that implements an interface: usecase interactors, domain services,
repositories, external-service clients (`internal/infrastructure/<service>/<service>/`), gRPC
handlers, workers, CLI tasks. Layer rules (`usecase-interactor.md`, `domain-service.md`,
`repository.md`, …) add layer-specific detail; this file is the common baseline they share.

Out of scope (these packages are not interface implementations, so §1 does not apply to them):
`internal/domain/model/` (constructors, enums, slice types — see `domain-model.md`),
`internal/domain/errors/`, interface-only packages (`internal/domain/repository/`,
`internal/domain/<capability>/`), `internal/pkg/` (generic utilities are package-level functions
by design — see `package-placement.md`), and generated code (`internal/infrastructure/**/internal/generated/`,
SQLBoiler `dbmodel`).

## 1. Public interface methods only — no wrappers

An implementation struct defines **only the methods of the interface it implements**.
Do not add:

- package-level private (lowercase, receiver-less) functions, or
- private methods on the struct that merely wrap a step of a public method
  (`findXxx`, `resolveXxx`, `enrichXxx`, `fetchXxx`, `buildXxx`, …).

The public method does the work **inline, top to bottom**: call the repositories / external APIs
it needs in sequence (including pagination loops and fallback calls), then build the result once
at the end (§3).

Why:
- A wrapper exists for one caller's exact needs. The next slightly different case (another field,
  another filter, the same API for a different resource) does not fit it and spawns yet another
  wrapper, so the package fills with near-duplicate private functions nobody can tell apart.
- Straight-line code shows every I/O call a method performs in one place, which is what reviewers
  and operators need to see.
- Behavior is tested through the interface; private wrappers are neither testable on their own
  nor reusable across packages.

Duplicating a short loop between two public methods (or two capability structs) is **preferred**
over a shared private helper.

```go
// BAD - public method delegates to private wrappers
func (i *bundleTicket) GetOwnedByIssuerTicketID(ctx context.Context, param issuer.GetOwnedByIssuerTicketIDParam) (*model.IssuerCustomerBundleTicket, error) {
    owned, err := i.findOwnedBundleTicketByTicketID(ctx, param.IssuerTicketID)
    ...
    location, err := i.resolveEventLocation(ctx, owned.Event)
    ...
}
func (i *bundleTicket) findOwnedBundleTicketByTicketID(...) (*issuer_customer.V1BundleTicket, error) { ... }
func (i *bundleTicket) resolveEventLocation(...) (*issuer_customer.V1LocationPartial, error) { ... }

// GOOD - the public method calls what it needs, in order, and builds the model once
func (i *bundleTicket) GetOwnedByIssuerTicketID(ctx context.Context, param issuer.GetOwnedByIssuerTicketIDParam) (*model.IssuerCustomerBundleTicket, error) {
    ctx = withCustomerContext(ctx, param.IssuerCustomerID, param.IssuerCustomerAPIKey)

    var owned *issuer_customer.V1BundleTicket
    for page := int64(1); owned == nil && page <= maxPages; page++ {
        res, httpRes, err := i.cli.CustomerV1ServiceAPI.ListBundleTicketsByTicketId(ctx, param.IssuerTicketID).Page(page).Limit(limit).Execute()
        ...
        idx := slices.IndexFunc(res.BundleTickets, func(bt issuer_customer.V1BundleTicket) bool { ... })
        if idx >= 0 {
            owned = &res.BundleTickets[idx]
        }
        if len(res.BundleTickets) == 0 || !res.Pagination.HasNext {
            break
        }
    }
    if owned == nil {
        return nil, errors.IssuerCustomerBundleTicketNotFoundErr.New()...
    }

    location := owned.Event.Location
    if location == nil {
        eventRes, httpRes, err := i.cli.CustomerV1ServiceAPI.GetEvent(ctx, owned.Event.Id).Execute()
        ...
        location = eventRes.Event.Location
    }

    return marshaller.CustomerIssuerBundleTicketToBundleTicketModel(*owned, location), nil
}
```

Where logic should live instead of a wrapper:

| Situation | Solution |
|-----------|----------|
| A step of one public method | Inline it |
| Logic on a single entity (predicate, derived value, state change of a DB-backed entity) | Method on the `domain/model` type |
| Logic spanning several entities | `domain/service` (its public method, again inline) |
| Conversion between external/DB structs and domain models | Exported function in the package's `marshaller/` (takes every input it needs as a parameter) |
| Cross-cutting request setup every method of the package needs (auth header / API-key injection into `ctx`) | One package-level function in `<actor>_base.go`; it must not call any API |
| Generic, domain-agnostic utility | `internal/pkg/` (see `package-placement.md`) |

Sanctioned private methods named by a layer rule remain allowed and are the **only** ones:
repository `addPreload` / `buildPreload` and the query-mod builder `buildListQuery` (`repository.md`), and domain-service struct methods
that dispatch to model methods (`domain-service.md`).

### When a linter complains about complexity

If inlining pushes a method over `gocognit` (20) or `funlen`, **flatten, do not extract**:
replace an inner search loop with `slices.IndexFunc` / `slices.ContainsFunc`, return early, fold
conditions into the loop condition, or drop defensive branches that a hard cap already covers.
Extracting a private method to satisfy the linter is the anti-pattern this rule forbids.

## 2. Pointer receivers for implementation structs

Implementation structs use **pointer receivers** and their constructors return `&struct{...}`:

```go
type bundleTicket struct{ cli *issuer_customer.APIClient }

func NewCustomerBundleTicket(cli *issuer_customer.APIClient) issuer.CustomerBundleTicket {
    return &bundleTicket{cli: cli}
}

func (i *bundleTicket) Get(ctx context.Context, param issuer.GetParam) (*model.IssuerCustomerBundleTicket, error) { ... }
```

This also applies to the `*Debug` / E2E stand-ins of the same interface.
Value receivers are for value types only: enum/string types (`String()`, `Valid()`) and slice
types (`IDs()`, `MapByID()`, `FilterBy…()`).

## 3. Gather every input first, then build the result once

A public method **resolves everything the result needs, then constructs it in a single call**
(domain constructor or marshaller). It never builds a struct and then patches fields on it:

- If a field comes from a second query or API call (e.g. the list response omits the venue and a
  separate `GetEvent` call provides it), make that call **before** construction and pass the value
  into the constructor / marshaller as a parameter.
- Do **not** assign model fields after construction (`m.LocationName = …`) in usecase or
  infrastructure code.
- Do **not** add `Set{X}` / `Update{X}` methods to a read model fetched from an external system
  just to patch it after construction. Those methods are for state transitions of DB-backed
  entities, per `domain-model.md` "State Change Methods".

```go
// BAD - build, then patch
bt := marshaller.ToBundleTicketModel(res.BundleTicket)
if !bt.LocationName.Valid {
    ev, _, err := i.cli.GetEvent(ctx, bt.EventID).Execute()
    ...
    bt.LocationName = null.StringFrom(ev.Event.Location.Name)
}
return bt, nil

// GOOD - resolve, then build once
location := res.BundleTicket.Event.Location
if location == nil {
    ev, _, err := i.cli.GetEvent(ctx, res.BundleTicket.Event.Id).Execute()
    ...
    location = ev.Event.Location
}
return marshaller.ToBundleTicketModel(res.BundleTicket, location), nil
```

Marshaller functions therefore take **every input they need as parameters** and return a fully
populated model; they do not reach back into the client to fetch anything.
