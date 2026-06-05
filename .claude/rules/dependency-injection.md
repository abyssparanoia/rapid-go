---
description: Dependency injection configuration patterns
globs:
  - "internal/infrastructure/dependency/**/*.go"
---

# Dependency Injection Guidelines

## DI Configuration Location

Location: `internal/infrastructure/dependency/dependency.go`

## Dependency Struct

```go
package dependency

import (
    "github.com/abyssparanoia/rapid-go/internal/infrastructure/mysql"
    "github.com/abyssparanoia/rapid-go/internal/usecase"
)

type Dependency struct {
    DatabaseCli *database.Client

    // Admin interactors
    AdminTenantInteractor usecase.AdminTenantInteractor
    AdminStaffInteractor  usecase.AdminStaffInteractor
    AdminAssetInteractor  usecase.AdminAssetInteractor

    // Other interactors
    StaffInteractor          usecase.StaffInteractor
    AuthenticationInteractor usecase.AuthenticationInteractor
    DebugInteractor          usecase.DebugInteractor
}
```

## Inject Method

```go
func (d *Dependency) Inject(
    ctx context.Context,
    e *environment.Environment,
) {
    // 1. Database client
    d.DatabaseCli = database.NewClient(e.DBHost, e.DBUser, e.DBPassword, e.DBDatabase, e.DBLogEnable)

    // 2. External clients (Firebase, GCS, AWS, Cognito)
    firebaseCli := firebase.NewClient(e.GCPProjectID)
    gcsCli := gcs.NewClient(ctx)
    gcsBucketHandle := gcs.NewBucketHandle(gcsCli, e.GCPBucketName)
    awsSession := aws.NewConfig(ctx, e.AWSRegion)
    cognitoCli := cognito.NewClient(awsSession, e.AWSCognitoEmulatorHost)

    // 3. Transactable
    transactable := database_transactable.NewTransactable()

    // 4. Repositories
    staffAuthenticationRepository := cognito_repository.NewStaffAuthentication(
        ctx,
        cognitoCli,
        e.AWSCognitoStaffUserPoolID,
        e.AWSCognitoStaffClientID,
        e.AWSCognitoEmulatorHost,
        e.AWSRegion,
    )
    tenantRepository := database_repository.NewTenant()
    staffRepository := database_repository.NewStaff()
    assetRepository := gcs_repository.NewAsset(gcsBucketHandle)

    // 5. Caches
    assetPathCache := database_cache.NewAssetPath()

    // 6. Domain services
    assetService := service.NewAsset(
        assetRepository,
        assetPathCache,
    )
    staffService := service.NewStaff(
        staffRepository,
        staffAuthenticationRepository,
    )

    // 7. Interactors - Admin
    d.AdminTenantInteractor = usecase.NewAdminTenantInteractor(
        transactable,
        tenantRepository,
        assetService,
    )
    d.AdminStaffInteractor = usecase.NewAdminStaffInteractor(
        transactable,
        tenantRepository,
        staffRepository,
        staffService,
        assetService,
    )
    d.AdminAssetInteractor = usecase.NewAdminAssetInteractor(
        assetService,
    )

    // 8. Interactors - Other
    d.StaffInteractor = usecase.NewStaffInteractor(
        transactable,
        tenantRepository,
        staffService,
    )
    d.AuthenticationInteractor = usecase.NewAuthenticationInteractor(
        staffAuthenticationRepository,
    )
    d.DebugInteractor = usecase.NewDebugInteractor(
        staffAuthenticationRepository,
    )
}
```

## Injection Order

Follow this order to ensure dependencies are available:

```
1. Database client
2. External clients (Firebase, GCS, AWS, Cognito)
3. Transactable
4. Repositories
5. Caches
6. Domain services
7. Interactors (by actor: Admin → Other)
```

## Adding New Entity

When adding a new entity (e.g., `Staff`), update dependency.go:

### 1. Add Interactor Fields to Dependency Struct

```go
type Dependency struct {
    // ...existing fields...
    AdminStaffInteractor usecase.AdminStaffInteractor  // Add this
}
```

### 2. Initialize Repository, Service, and Interactor in Inject()

```go
func (d *Dependency) Inject(...) {
    // ...

    // Repository
    staffRepository := database_repository.NewStaff()

    // Domain service (if needed)
    staffService := service.NewStaff(
        staffRepository,
        staffAuthenticationRepository,
    )

    // Interactor
    d.AdminStaffInteractor = usecase.NewAdminStaffInteractor(
        transactable,
        tenantRepository,
        staffRepository,
        staffService,
        assetService,
    )
}
```

## Testing with DI

For integration tests, you can inject mock dependencies:

```go
func setupTestDependency(t *testing.T) *Dependency {
    ctrl := gomock.NewController(t)

    mockStaffRepo := mock_repository.NewMockStaff(ctrl)
    mockTransactable := mock_repository.NewMockTransactable(ctrl)

    return &Dependency{
        AdminStaffInteractor: usecase.NewAdminStaffInteractor(
            mockTransactable,
            // ... other mock dependencies
        ),
    }
}
```

## Environment Variables

### notEmpty by Default

All environment variables must use the `notEmpty` tag. Do not use `required` (allows empty string) or `envDefault:""` (masks misconfiguration).

```go
// Good - notEmpty (fails fast on missing OR empty config)
GoogleMapsAPIKey string `env:"GOOGLE_MAPS_API_KEY,notEmpty"`

// Bad - required only checks if env var is set; "" is allowed and slips into production
GoogleMapsAPIKey string `env:"GOOGLE_MAPS_API_KEY,required"`

// Bad - allows empty value (masks misconfiguration in production)
GoogleMapsAPIKey string `env:"GOOGLE_MAPS_API_KEY" envDefault:""`
```

**`required` vs `notEmpty`** (caarlos0/env v11):
- `required`: checks only that the env var is **set** — `FOO=""` passes
- `notEmpty`: checks that the env var is set **and** non-empty — `FOO=""` fails

For local development, set a dummy value in `.envrc` and switch to debug implementations via the `ApplicationEnvironmentLocal` branch in `Inject()`.

When adding a new env variable, also add it to:
- `.envrc.tmpl` (with a meaningful dummy/local value, not `""`)
- `.github/workflows/ci.yml` `env` section (for E2E tests)

**Exceptions** (do NOT use `notEmpty`):
- Fields with a meaningful `envDefault` value (e.g. `MIN_LOG_LEVEL` → `"info"`, `DB_LOG_ENABLE` → `"false"`)
- Local-only override fields that must be unset in production (e.g. `AWS_EMULATOR_HOST`, `AWS_COGNITO_EMULATOR_HOST`). Leave these untagged and document the intent in a comment.
- Fields where an empty value is a legitimate configured value in some tier (e.g. `DB_PASSWORD` — local TiDB / TiUP Playground uses no root password). Use `required` (not `notEmpty`) and document the intent in a comment.

### Local Debug Implementation Switching

Consolidate all Real/Debug switching in the single `ApplicationEnvironmentLocal` branch inside `Inject()`. Do not branch on whether an API key is empty.

```go
// Good - consolidated in the existing local branch
if e.Environment == environment.ApplicationEnvironmentLocal {
    thingClient = iot_core_iot.NewThingDebug()
    geocodeRepo = googlemaps_repository.NewGeocodeDebug()
} else {
    thingClient = iot_core_iot.NewThing(controlPlaneCli)
    geocodeRepo = googlemaps_repository.NewGeocode(e.GoogleMapsAPIKey)
}

// Bad - scattered API-key-based branching
if e.GoogleMapsAPIKey != "" {
    geocodeRepo = googlemaps_repository.NewGeocode(e.GoogleMapsAPIKey)
} else {
    geocodeRepo = googlemaps_repository.NewGeocodeDebug()
}
```

## Best Practices

1. **Group by layer** - Clients, then repositories, then caches, then services, then interactors
2. **Group interactors by actor** - Admin, then Other sections
3. **Single responsibility** - Each interactor should have focused dependencies
4. **Explicit dependencies** - Pass all dependencies through constructor
5. **No global state** - All dependencies should be in Dependency struct
6. **Domain services are optional** - Only create when business logic spans multiple entities
7. **Environment variables must be non-empty** - Use `notEmpty` tag (rejects unset AND empty string), not `required` or `envDefault:""`
