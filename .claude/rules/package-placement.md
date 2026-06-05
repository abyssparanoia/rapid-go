---
description: Where to put new Go packages — internal/pkg/ vs internal/domain/
globs:
  - "internal/pkg/**/*.go"
  - "internal/domain/**/*.go"
---

# Package Placement Guidelines

## Rule

`internal/pkg/` is reserved for **domain-agnostic** utilities. Anything that imports or returns a `domain/model` type, encodes a domain concept, or only makes sense within this product belongs under `internal/domain/`.

## Quick check

When deciding where to place a new package, ask:

- Does the package import from `internal/domain/`? → put it in `internal/domain/`.
- Could this package be lifted into another Go project unchanged? → `internal/pkg/` is OK.
- Is the package's name or types tied to a domain entity (`VehicleLocationLog`, `Tenant`, etc.)? → `internal/domain/`.

When in doubt, prefer `internal/domain/` — moving a misplaced helper out of `internal/pkg/` later is more painful than placing it correctly the first time.

## Examples

### Belongs in `internal/pkg/`

| Package | Why |
|---|---|
| `internal/pkg/id` | Generic ID generation, no domain coupling |
| `internal/pkg/now` | Time abstraction, reusable anywhere |
| `internal/pkg/geo` | Geometric calculations on raw lat/lng numbers |
| `internal/pkg/validation` | Generic validator wrapper |
| `internal/pkg/uuid` | UUID encoding, no domain coupling |

### Belongs in `internal/domain/`

| Package | Why |
|---|---|
| `internal/domain/cursor` | `VehicleLocationLog` cursor — encodes a product-specific pagination shape |
| `internal/domain/wifi` | WiFi positioning interface tied to `model.HWBotLocation` |
| `internal/domain/iot` | IoT device interface returning domain types |
| `internal/domain/geocode` | Geocoding interface returning domain types |

## Anti-pattern: Domain types under `internal/pkg/`

```go
// BAD - internal/pkg/cursor/cursor.go
package cursor

// VehicleLocationLog encodes a domain-specific cursor.
// This package only makes sense for the ListVehicleLocationLogs API,
// so it should live under internal/domain/, not internal/pkg/.
type VehicleLocationLog struct {
    LteTime time.Time `json:"t"`
    ID      string    `json:"i"`
}
```

```go
// GOOD - internal/domain/cursor/cursor.go
package cursor

type VehicleLocationLog struct { /* ... */ }
```

The same applies even when the package does not literally import a `domain/model` type, as long as its purpose is product-specific (e.g. a cursor for a specific API).
