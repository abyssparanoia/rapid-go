# CLI Command Pattern Guidelines

## Overview

CLI commands are task-oriented commands executed via `./app task {command-name}`. These commands are used for:

- Initial setup operations (e.g., creating root admin/staff)
- One-time data migration tasks
- Administrative operations that don't fit into the API workflow

## Architecture

```
CLI Command Entry
    ↓
cmd/app/main.go → task_cmd → {entity}_cmd
    ↓
CMD struct → TaskInteractor → Repository
```

**Key Principle**: CLI commands use dedicated `Task{Entity}Interactor` implementations, separate from API interactors.

## Directory Structure

```
internal/
├── usecase/
│   ├── input/
│   │   └── task_{entity}.go          # Task input DTOs
│   ├── output/
│   │   └── task_{entity}.go          # Task output DTOs (if needed)
│   ├── task_{entity}.go               # TaskInteractor interface
│   └── task_{entity}_impl.go          # TaskInteractor implementation
├── infrastructure/
│   ├── dependency/
│   │   └── dependency.go              # DI registration
│   └── cmd/
│       └── internal/
│           └── task_cmd/
│               ├── cmd.go                          # Task command registration
│               └── {command_name}_cmd/
│                   ├── cmd.go                      # Cobra command definition
│                   └── {command_name}.go           # Command execution logic
```

## Implementation Components

### 1. Input DTO

**Location**: `internal/usecase/input/task_{entity}.go`

```go
package input

import (
    "time"
    "github.com/abyssparanoia/rapid-go/internal/domain/errors"
    "github.com/abyssparanoia/rapid-go/internal/pkg/validation"
)

type TaskCreateAdmin struct {
    Email       string    `validate:"required,email"`
    DisplayName string    `validate:"required"`
    Password    string    `validate:"required"`
    RequestTime time.Time `validate:"required"`
}

func NewTaskCreateAdmin(
    email string,
    displayName string,
    password string,
    t time.Time,
) *TaskCreateAdmin {
    return &TaskCreateAdmin{
        Email:       email,
        DisplayName: displayName,
        Password:    password,
        RequestTime: t,
    }
}

func (p *TaskCreateAdmin) Validate() error {
    if err := validation.Validate(p); err != nil {
        return errors.RequestInvalidArgumentErr.Wrap(err)
    }
    return nil
}
```

**Naming Pattern**: `Task{Action}{Entity}` (e.g., `TaskCreateAdmin`, `TaskMigrateUser`)

### 2. Output DTO (Optional)

**Location**: `internal/usecase/output/task_{entity}.go`

Only create when:

- Command needs to return multiple values to caller
- Values need to be formatted for display

```go
package output

type TaskCreateAdmin struct {
    AdminID  string
    AuthUID  string
    Password string
}

func NewTaskCreateAdmin(
    adminID string,
    authUID string,
    password string,
) *TaskCreateAdmin {
    return &TaskCreateAdmin{
        AdminID:  adminID,
        AuthUID:  authUID,
        Password: password,
    }
}
```

### 3. Interactor Interface

**Location**: `internal/usecase/task_{entity}.go`

```go
package usecase

import (
    "context"
    "github.com/abyssparanoia/rapid-go/internal/usecase/input"
    "github.com/abyssparanoia/rapid-go/internal/usecase/output"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_usecase
type TaskAdminInteractor interface {
    Create(
        ctx context.Context,
        param *input.TaskCreateAdmin,
    ) (*output.TaskCreateAdmin, error)
}
```

**Naming Pattern**:

- Interface: `Task{Entity}Interactor`
- Method: Descriptive action (e.g., `Create`, `Migrate`, `Cleanup`)

### 4. Interactor Implementation

**Location**: `internal/usecase/task_{entity}_impl.go`

```go
package usecase

import (
    "context"
    "github.com/aarondl/null/v9"
    "github.com/abyssparanoia/rapid-go/internal/domain/model"
    "github.com/abyssparanoia/rapid-go/internal/domain/repository"
    "github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
    "github.com/abyssparanoia/rapid-go/internal/usecase/input"
    "github.com/abyssparanoia/rapid-go/internal/usecase/output"
)

type taskAdminInteractor struct {
    transactable                  repository.Transactable
    adminRepository               repository.Admin
    adminAuthenticationRepository repository.AdminAuthentication
}

func NewTaskAdminInteractor(
    transactable repository.Transactable,
    adminRepository repository.Admin,
    adminAuthenticationRepository repository.AdminAuthentication,
) TaskAdminInteractor {
    return &taskAdminInteractor{
        transactable:                  transactable,
        adminRepository:               adminRepository,
        adminAuthenticationRepository: adminAuthenticationRepository,
    }
}

func (i *taskAdminInteractor) Create(
    ctx context.Context,
    param *input.TaskCreateAdmin,
) (*output.TaskCreateAdmin, error) {
    // 1. Validate input
    if err := param.Validate(); err != nil {
        return nil, err
    }

    var admin *model.Admin
    var authUID string

    // 2. Execute in transaction
    if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
        // Business logic here
        // ...
        return nil
    }); err != nil {
        return nil, err
    }

    // 3. Return result
    return output.NewTaskCreateAdmin(admin.ID, authUID, param.Password), nil
}
```

**Key Points**:

- Use `RWTx` for write operations
- Validate input at the beginning
- Keep business logic in transaction scope
- Return output DTO or error

### 5. DI Registration

**Location**: `internal/infrastructure/dependency/dependency.go`

```go
type Dependency struct {
    // ... existing fields

    // task (separate section at the end)
    TaskAdminInteractor usecase.TaskAdminInteractor
}

func (d *Dependency) Inject(ctx context.Context, e *environment.Environment) {
    // ... existing initialization

    // Repositories
    adminRepository := database_repository.NewAdmin()

    // ... existing interactors

    // Task interactors (at the end of Inject)
    d.TaskAdminInteractor = usecase.NewTaskAdminInteractor(
        transactable,
        adminRepository,
        adminAuthenticationRepository,
    )
}
```

**Naming Convention**: Group all `Task*Interactor` fields under a `// task` comment section.

### 6. Cobra Command Definition

**Location**: `internal/infrastructure/cmd/internal/task_cmd/{command_name}_cmd/cmd.go`

```go
package create_root_admin_cmd

import (
    "context"
    "github.com/abyssparanoia/rapid-go/internal/infrastructure/dependency"
    "github.com/abyssparanoia/rapid-go/internal/infrastructure/environment"
    "github.com/abyssparanoia/rapid-go/internal/pkg/logger"
    "github.com/caarlos0/env/v11"
    "github.com/spf13/cobra"
)

func NewCreateRootAdminCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "create-root-admin",
        Short: "create root admin",
        Run: func(cmd *cobra.Command, args []string) {
            ctx := context.Background()

            // 1. Load environment
            e := &environment.Environment{}
            if err := env.Parse(e); err != nil {
                panic(err)
            }

            // 2. Initialize logger
            l := logger.New(e.MinLogLevel)
            ctx = logger.ToContext(ctx, l)

            // 3. Initialize dependencies
            d := &dependency.Dependency{}
            d.Inject(ctx, e)

            // 4. Execute command
            c := &CMD{
                ctx,
                d.TaskAdminInteractor,
            }
            if err := c.CreateRootAdmin(cmd); err != nil {
                panic(err)
            }
        },
    }

    // Define flags
    cmd.Flags().StringP("email", "e", "", "email address")
    cmd.Flags().StringP("display-name", "d", "", "display name")

    return cmd
}
```

**Command Naming**:

- Use: kebab-case (e.g., `create-root-admin`, `migrate-users`)
- Package: snake_case with `_cmd` suffix (e.g., `create_root_admin_cmd`)
- Function: PascalCase (e.g., `NewCreateRootAdminCmd`)

### 7. Command Execution Logic

**Location**: `internal/infrastructure/cmd/internal/task_cmd/{command_name}_cmd/{command_name}.go`

```go
package create_root_admin_cmd

import (
    "context"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "github.com/abyssparanoia/rapid-go/internal/domain/errors"
    "github.com/abyssparanoia/rapid-go/internal/pkg/now"
    "github.com/abyssparanoia/rapid-go/internal/usecase"
    "github.com/abyssparanoia/rapid-go/internal/usecase/input"
    "github.com/spf13/cobra"
)

type CMD struct {
    ctx                 context.Context
    taskAdminInteractor usecase.TaskAdminInteractor
}

func (c *CMD) CreateRootAdmin(cmd *cobra.Command) error {
    // 1. Parse flags
    email, err := cmd.Flags().GetString("email")
    if err != nil {
        return errors.InternalErr.Wrap(err)
    }
    if email == "" {
        return errors.InternalErr.WithDetail("email is required")
    }

    displayName, err := cmd.Flags().GetString("display-name")
    if err != nil {
        return errors.InternalErr.Wrap(err)
    }
    if displayName == "" {
        return errors.InternalErr.WithDetail("display-name is required")
    }

    // 2. Generate any required values (passwords, IDs, etc.)
    password, err := generatePassword(16)
    if err != nil {
        return errors.InternalErr.Wrap(err).WithDetail("failed to generate password")
    }

    // 3. Call interactor
    result, err := c.taskAdminInteractor.Create(
        c.ctx,
        input.NewTaskCreateAdmin(
            email,
            displayName,
            password,
            now.Now(),
        ),
    )
    if err != nil {
        return err
    }

    // 4. Output results
    fmt.Printf("AdminID: %s\n", result.AdminID)
    fmt.Printf("AuthUID: %s\n", result.AuthUID)
    fmt.Printf("Password: %s\n", result.Password)

    return nil
}

// Helper functions
func generatePassword(length int) (string, error) {
    byteLength := (length * 3) / 4
    if byteLength < 1 {
        byteLength = 1
    }

    bytes := make([]byte, byteLength)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }

    encoded := base64.URLEncoding.EncodeToString(bytes)
    if len(encoded) > length {
        encoded = encoded[:length]
    }

    return encoded, nil
}
```

**Key Points**:

- Validate all required flags
- Generate required values (passwords, codes) in CMD layer, not interactor
- Use `fmt.Printf` for standard output
- Return errors for proper error handling

### 8. Command Registration

**Location**: `internal/infrastructure/cmd/internal/task_cmd/cmd.go`

```go
package task_cmd

import (
    "github.com/abyssparanoia/rapid-go/internal/infrastructure/cmd/internal/task_cmd/create_root_admin_cmd"
    "github.com/abyssparanoia/rapid-go/internal/infrastructure/cmd/internal/task_cmd/create_root_staff_cmd"
    "github.com/spf13/cobra"
)

func NewTaskCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "task",
        Short: "cli task",
        Run: func(cmd *cobra.Command, args []string) {
            if len(args) == 0 {
                cmd.HelpFunc()(cmd, args)
            }
        },
    }

    // Register all task commands
    cmd.AddCommand(create_root_staff_cmd.NewCreateRootStaffCmd())
    cmd.AddCommand(create_root_admin_cmd.NewCreateRootAdminCmd())

    return cmd
}
```

## Implementation Workflow

When adding a new CLI command:

1. **Input DTO** - `internal/usecase/input/task_{entity}.go`
2. **Output DTO** - `internal/usecase/output/task_{entity}.go` (if needed)
3. **Interactor Interface** - `internal/usecase/task_{entity}.go`
4. **Interactor Implementation** - `internal/usecase/task_{entity}_impl.go`
5. **DI Registration** - `internal/infrastructure/dependency/dependency.go`
6. **Cobra Command** - `internal/infrastructure/cmd/internal/task_cmd/{command_name}_cmd/cmd.go`
7. **Command Logic** - `internal/infrastructure/cmd/internal/task_cmd/{command_name}_cmd/{command_name}.go`
8. **Command Registration** - `internal/infrastructure/cmd/internal/task_cmd/cmd.go`

## Running Commands

```bash
# Build
go build -o app cmd/app/main.go

# Execute
./app task {command-name} [flags]

# Example
./app task create-root-admin --email admin@example.com --display-name "Root Admin"
```

## Best Practices

1. **Separate Task Interactors** - Don't reuse API interactors for CLI commands
   - CLI commands often need different validation rules
   - CLI commands may bypass authorization checks
   - CLI commands may have different transaction scopes

2. **Generate Values in CMD Layer** - Generate passwords, codes, etc. in CMD layer, not interactor
   - Keeps interactor focused on business logic
   - Makes testing easier
   - Allows CMD to control output format

3. **Validate Required Flags** - Always validate required flags in CMD layer
   - Return clear error messages
   - Use `errors.InternalErr.WithDetail()` for context

4. **Use `now.Now()` for Timestamps** - Consistent with codebase patterns
   - Allows for time mocking in tests
   - Ensures consistent timezone handling

5. **Output to stdout** - Use `fmt.Printf` for command output
   - Easy to capture in scripts
   - Clear separation from logs

6. **Keep Commands Focused** - One command = one operation
   - Avoid combining multiple operations in single command
   - Use separate commands for related but distinct operations

7. **Document Command Purpose** - Add clear `Short` description in cobra command
   - Helps users understand command purpose via `--help`

## Common Patterns

### Pattern 1: Generate and Output Secrets

When generating passwords/tokens that need to be displayed:

```go
// Generate in CMD layer
password, err := generatePassword(16)
if err != nil {
    return errors.InternalErr.Wrap(err)
}

// Pass to interactor
result, err := c.interactor.Create(ctx, input.New(..., password, ...))

// Output to user
fmt.Printf("Password: %s\n", result.Password)
```

### Pattern 2: Dry Run Flag

For commands that modify data, support dry-run:

```go
cmd.Flags().Bool("dry-run", false, "preview changes without applying")

// In execution
dryRun, _ := cmd.Flags().GetBool("dry-run")
if dryRun {
    fmt.Println("DRY RUN: Would create admin with email:", email)
    return nil
}
```

### Pattern 3: Confirmation Prompt

For destructive operations, add confirmation:

```go
fmt.Printf("This will delete all users. Continue? (y/n): ")
var response string
fmt.Scanln(&response)
if response != "y" {
    fmt.Println("Operation cancelled")
    return nil
}
```

### Pattern 4: Progress Output

For long-running operations:

```go
fmt.Println("Processing users...")
for i, user := range users {
    if err := processUser(user); err != nil {
        return err
    }
    if (i+1)%100 == 0 {
        fmt.Printf("Processed %d/%d users\n", i+1, total)
    }
}
fmt.Println("Complete!")
```

## Error Handling

CLI commands should:

- Return errors from CMD execution methods
- Use `panic(err)` in cobra `Run` function for fatal errors
- Log detailed context with structured logging if needed

```go
Run: func(cmd *cobra.Command, args []string) {
    // ... setup code ...

    if err := c.ExecuteCommand(cmd); err != nil {
        logger.L(ctx).Error("Command failed", zap.Error(err))
        panic(err)  // Exit with error
    }
}
```

## Testing CLI Commands

### Unit Test Interactor

Test interactor logic separately from CMD:

```go
func TestTaskAdminInteractor_Create(t *testing.T) {
    // Test interactor logic with mocked repositories
}
```

### Integration Test CMD

Test full command execution:

```go
func TestCreateRootAdminCmd(t *testing.T) {
    // Setup test environment
    // Execute command
    // Verify output
    // Verify database state
}
```

## References

- See `usecase-interactor.md` for interactor patterns
- See `testing.md` for testing conventions
- See `dependency-injection.md` for DI patterns
