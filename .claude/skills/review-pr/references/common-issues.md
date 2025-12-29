# Common Issues by Layer

Frequently encountered issues during PR review. Check these patterns in your changes.

## Domain Layer

### Direct Field Assignment
```go
// Bad - Direct assignment in usecase
admin.Role = param.Role
admin.UpdatedAt = param.RequestTime

// Good - Use domain method
admin.UpdateRole(param.Role, param.RequestTime)
```

### Missing UpdatedAt
```go
// Bad - Forgot to update timestamp
func (m *Example) SetStatus(s Status) {
    m.Status = s
}

// Good - Always update timestamp
func (m *Example) SetStatus(s Status, t time.Time) {
    m.Status = s
    m.UpdatedAt = t
}
```

### Enum Without Unknown
```go
// Bad - No default value
type Status string
const (
    StatusActive Status = "active"
)

// Good - Start with Unknown
type Status string
const (
    StatusUnknown Status = "unknown"
    StatusActive  Status = "active"
)
```

## Repository Layer

### Wrong Type for Optional Enum Fields
```go
// Bad - Using pointer
type ListQuery struct {
    Status *model.Status
}

// Good - Use nullable.Type
type ListQuery struct {
    Status nullable.Type[model.Status]
}
```

### Missing OrFail Check
```go
// Bad - Always returns error on not found
if err == sql.ErrNoRows {
    return nil, errors.NotFoundErr.New()
}

// Good - Check OrFail option
if err == sql.ErrNoRows {
    if query.OrFail {
        return nil, errors.NotFoundErr.New()
    }
    return nil, nil
}
```

### Missing Context Executor
```go
// Bad - Direct database call
dbmodel.Examples().One(ctx, db)

// Good - Use context executor for transactions
dbmodel.Examples().One(ctx, transactable.GetContextExecutor(ctx))
```

## Usecase Layer

### Missing Validation Call
```go
// Bad - No validation
func (i *interactor) Create(ctx context.Context, param *input.Create) (*model.Example, error) {
    example := model.NewExample(...)
}

// Good - Validate first
func (i *interactor) Create(ctx context.Context, param *input.Create) (*model.Example, error) {
    if err := param.Validate(); err != nil {
        return nil, err
    }
    example := model.NewExample(...)
}
```

### IdP Sync Outside Transaction
```go
// Bad - IdP call outside transaction
if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
    return i.repo.Update(ctx, entity)
}); err != nil {
    return nil, err
}
i.authRepo.StoreClaims(ctx, ...)  // Outside transaction!

// Good - IdP call inside transaction
if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
    if err := i.repo.Update(ctx, entity); err != nil {
        return err
    }
    return i.authRepo.StoreClaims(ctx, ...)  // Inside transaction
}); err != nil {
    return nil, err
}
```

### Missing ForUpdate on Update Operations
```go
// Bad - No lock before update
entity, err := i.repo.Get(ctx, GetQuery{ID: id, OrFail: true})
entity.Update(...)
i.repo.Update(ctx, entity)

// Good - Lock with ForUpdate
entity, err := i.repo.Get(ctx, GetQuery{
    ID: id,
    BaseGetOptions: BaseGetOptions{
        OrFail:    true,
        ForUpdate: true,  // Lock row
    },
})
entity.Update(...)
i.repo.Update(ctx, entity)
```

### Missing Preload on Return
```go
// Bad - Return without relations
return i.repo.Get(ctx, GetQuery{ID: id, OrFail: true})

// Good - Always preload for returned entities
return i.repo.Get(ctx, GetQuery{
    ID: id,
    BaseGetOptions: BaseGetOptions{
        OrFail:  true,
        Preload: true,  // Load relations
    },
})
```

## Handler Layer

### Missing Nil Check for Optional Proto Fields
```go
// Bad - No nil check
status := marshaller.StatusToModel(req.Status)

// Good - Check before conversion
if req.Status != nil {
    status := marshaller.StatusToModel(*req.Status)
    param.Status = nullable.TypeFrom(status)
}
```

### Field Mapping Omission in Marshaller
```go
// Bad - Missing field
func ToPb(m *model.Example) *pb.Example {
    return &pb.Example{
        Id:   m.ID,
        Name: m.Name,
        // Missing: Status, CreatedAt, UpdatedAt
    }
}

// Good - All fields mapped
func ToPb(m *model.Example) *pb.Example {
    return &pb.Example{
        Id:        m.ID,
        Name:      m.Name,
        Status:    StatusToPb(m.Status),
        CreatedAt: timestamppb.New(m.CreatedAt),
        UpdatedAt: timestamppb.New(m.UpdatedAt),
    }
}
```

### Nullable Field Without Variable Declaration
```go
// Bad - Inline nullable handling
return &pb.Example{
    AcceptedAt: func() *timestamppb.Timestamp {
        if m.AcceptedAt.Valid {
            return timestamppb.New(m.AcceptedAt.Time)
        }
        return nil
    }(),
}

// Good - Variable declaration pattern
var acceptedAt *timestamppb.Timestamp
if m.AcceptedAt.Valid {
    acceptedAt = timestamppb.New(m.AcceptedAt.Time)
}
return &pb.Example{
    AcceptedAt: acceptedAt,
}
```

## Tests

### Missing Mock for New Dependency
```go
// Bad - Forgot to mock new dependency
interactor := NewInteractor(
    mockTransactable,
    mockRepo,
    // Missing: mockNewService
)

// Good - All dependencies mocked
interactor := NewInteractor(
    mockTransactable,
    mockRepo,
    mockNewService,
)
```

### Not Using Table-Driven Pattern
```go
// Bad - Separate test functions
func TestCreate_Success(t *testing.T) { ... }
func TestCreate_ValidationError(t *testing.T) { ... }

// Good - Table-driven
func TestCreate(t *testing.T) {
    tests := map[string]func(t *testing.T) (args, usecase, want){
        "success": func(t *testing.T) { ... },
        "validation error": func(t *testing.T) { ... },
    }
    for name, setup := range tests {
        t.Run(name, func(t *testing.T) { ... })
    }
}
```

## DI Registration

### Wrong Injection Order
```go
// Bad - Service created before its dependencies
staffService := service.NewStaff(staffRepo)  // staffRepo not created yet
staffRepo := repository.NewStaff()

// Good - Dependencies first
staffRepo := repository.NewStaff()
staffService := service.NewStaff(staffRepo)
```

### Missing Handler Update
```go
// Bad - New interactor not passed to handler
d.AdminHandler = admin.NewHandler(
    d.AdminTenantInteractor,
    // Missing: d.AdminNewEntityInteractor
)

// Good - All interactors passed
d.AdminHandler = admin.NewHandler(
    d.AdminTenantInteractor,
    d.AdminNewEntityInteractor,
)
```
