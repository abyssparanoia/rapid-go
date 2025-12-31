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

## Best Practices

1. **Group by layer** - Clients, then repositories, then caches, then services, then interactors
2. **Group interactors by actor** - Admin, then Other sections
3. **Single responsibility** - Each interactor should have focused dependencies
4. **Explicit dependencies** - Pass all dependencies through constructor
5. **No global state** - All dependencies should be in Dependency struct
6. **Domain services are optional** - Only create when business logic spans multiple entities
