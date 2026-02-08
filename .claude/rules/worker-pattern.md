# Worker Pattern Guidelines

## Overview

Workers are background processes that consume messages from messaging systems (SQS, Pub/Sub, etc.) and execute business logic asynchronously. This project uses a single binary with hierarchical Cobra commands for worker management.

## Supported Messaging Systems

- **AWS SQS** (Simple Queue Service) - `internal/infrastructure/sqs/`
- **GCP Pub/Sub** (Future) - `internal/infrastructure/pubsub/`
- **Other systems** can follow the same pattern

## Architecture

```
Client → API Server → Message Queue → Worker (Subscriber) → Handler → Interactor
                                         ↓
                                   Process messages
                                   (sequential or parallel)
```

## Directory Structure

```
internal/
├── usecase/
│   ├── task_{entity}.go              # Task interactor interface
│   ├── task_{entity}_impl.go         # Task interactor implementation
│   └── input/
│       └── task_{entity}.go          # Input DTOs for task
├── infrastructure/
│   ├── {messaging}/                  # Messaging system adapter (sqs, pubsub, etc.)
│   │   ├── client.go                 # Client wrapper
│   │   ├── subscriber/
│   │   │   └── subscriber.go         # Message polling/subscription logic
│   │   └── handler/
│   │       └── {entity}.go           # Message handler
│   └── cmd/
│       └── internal/
│           └── worker_cmd/
│               ├── cmd.go            # Worker root command
│               └── {entity}_cmd/
│                   └── cmd.go        # Entity-specific worker command
```

## Worker Command Structure

Workers use hierarchical Cobra commands:

```
app worker                           # Worker root command
    ├── project-key-creation        # Entity-specific subcommand
    ├── user-notification           # Another worker (example)
    └── data-export                 # Another worker (example)
```

**Naming Pattern:**

- Root command: `worker`
- Subcommand: `{entity-name}` (kebab-case, e.g., `project-key-creation`, `user-notification`)

## Implementation Components

### 1. Task Usecase Layer

#### Error Definitions

Location: `internal/domain/errors/errors.go`

Add task-specific errors (E300xxx series):

```go
// Task errors (E300xxx series)
TaskProjectNotFoundErr          = NewNotFoundError("E300001", "project not found for task processing")
TaskInvalidMessageErr           = NewBadRequestError("E300002", "invalid message format")
TaskProjectKeyCreationFailedErr = NewInternalError("E300003", "project key creation task failed")
```

**Error Code Pattern:**

- E3001xx - First worker entity
- E3002xx - Second worker entity
- E3003xx - Third worker entity (and so on)

#### Input DTO

Location: `internal/usecase/input/task_{entity}.go`

```go
package input

import (
    "time"

    "github.com/eaglys-platform/pandlock-api/internal/domain/errors"
    "github.com/eaglys-platform/pandlock-api/internal/pkg/validation"
)

type TaskProcessProjectKeyCreation struct {
    ProjectID   string    `validate:"required"`
    RequestTime time.Time `validate:"required"`
}

func NewTaskProcessProjectKeyCreation(
    projectID string,
    requestTime time.Time,
) *TaskProcessProjectKeyCreation {
    return &TaskProcessProjectKeyCreation{
        ProjectID:   projectID,
        RequestTime: requestTime,
    }
}

func (p *TaskProcessProjectKeyCreation) Validate() error {
    if err := validation.Validate(p); err != nil {
        return errors.RequestInvalidArgumentErr.Wrap(err)
    }
    return nil
}
```

**Naming Pattern:**

- Input struct: `TaskProcess{Entity}`
- Constructor: `NewTaskProcess{Entity}`

#### Interactor Interface

Location: `internal/usecase/task_{entity}.go`

```go
package usecase

import (
    "context"

    "github.com/eaglys-platform/pandlock-api/internal/domain/model"
    "github.com/eaglys-platform/pandlock-api/internal/usecase/input"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_usecase

type TaskProjectKeyCreationInteractor interface {
    Process(
        ctx context.Context,
        param *input.TaskProcessProjectKeyCreation,
    ) (*model.Project, error)
}
```

**Naming Pattern:**

- Interface: `Task{Entity}Interactor`
- Method: `Process`
- Return: Domain model or error

#### Interactor Implementation

Location: `internal/usecase/task_{entity}_impl.go`

```go
package usecase

import (
    "context"

    "github.com/eaglys-platform/pandlock-api/internal/domain/model"
    "github.com/eaglys-platform/pandlock-api/internal/domain/repository"
    "github.com/eaglys-platform/pandlock-api/internal/usecase/input"
    "github.com/aarondl/null/v9"
)

type taskProjectKeyCreationInteractor struct {
    transactable      repository.Transactable
    projectRepository repository.Project
}

func NewTaskProjectKeyCreationInteractor(
    transactable repository.Transactable,
    projectRepository repository.Project,
) TaskProjectKeyCreationInteractor {
    return &taskProjectKeyCreationInteractor{
        transactable:      transactable,
        projectRepository: projectRepository,
    }
}

func (i *taskProjectKeyCreationInteractor) Process(
    ctx context.Context,
    param *input.TaskProcessProjectKeyCreation,
) (*model.Project, error) {
    // 1. Validate input
    if err := param.Validate(); err != nil {
        return nil, err
    }

    var project *model.Project

    // 2. Execute in transaction
    if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
        // Get entity with lock
        p, err := i.projectRepository.Get(ctx, repository.GetProjectQuery{
            ID: null.StringFrom(param.ProjectID),
            BaseGetOptions: repository.BaseGetOptions{
                OrFail:    true,
                ForUpdate: true, // Lock for update
            },
        })
        if err != nil {
            return err
        }
        project = p

        // 3. Business logic implementation
        // TODO: Add actual processing logic here

        return nil
    }); err != nil {
        return nil, err
    }

    return project, nil
}
```

**Key Points:**

- Use transaction with `ForUpdate: true` for data consistency
- Validate input at the beginning
- Return domain model after processing

### 2. Messaging Infrastructure Layer

#### Client Wrapper

Location: `internal/infrastructure/{messaging}/client.go`

**AWS SQS Example:**

```go
package sqs

import (
    "context"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Client struct {
    Cli *sqs.Client
}

func NewClient(cfg aws.Config) *Client {
    return &Client{
        Cli: sqs.NewFromConfig(cfg),
    }
}

func (c *Client) ReceiveMessage(
    ctx context.Context,
    params *sqs.ReceiveMessageInput,
) (*sqs.ReceiveMessageOutput, error) {
    return c.Cli.ReceiveMessage(ctx, params)
}

func (c *Client) DeleteMessage(
    ctx context.Context,
    params *sqs.DeleteMessageInput,
) (*sqs.DeleteMessageOutput, error) {
    return c.Cli.DeleteMessage(ctx, params)
}
```

**GCP Pub/Sub Example (Future):**

```go
package pubsub

import (
    "context"

    "cloud.google.com/go/pubsub"
)

type Client struct {
    Cli *pubsub.Client
}

func NewClient(ctx context.Context, projectID string) (*Client, error) {
    cli, err := pubsub.NewClient(ctx, projectID)
    if err != nil {
        return nil, err
    }
    return &Client{Cli: cli}, nil
}
```

#### Subscriber

Location: `internal/infrastructure/{messaging}/subscriber/subscriber.go`

**AWS SQS Example:**

```go
type Subscriber struct {
    client            *sqs.Client
    queueURL          string
    maxMessages       int32
    waitTimeSeconds   int32
    visibilityTimeout int32
    handler           MessageHandler
    stopCh            chan struct{}
    stoppedCh         chan struct{}
}

type SubscriberConfig struct {
    QueueURL          string
    MaxMessages       int32  // Set to 1 for sequential processing
    WaitTimeSeconds   int32  // Long polling (e.g., 20 seconds)
    VisibilityTimeout int32  // Processing timeout (e.g., 300 seconds)
}

func NewSubscriber(
    client *sqs.Client,
    config SubscriberConfig,
    handler MessageHandler,
) *Subscriber {
    return &Subscriber{
        client:            client,
        queueURL:          config.QueueURL,
        maxMessages:       config.MaxMessages,
        waitTimeSeconds:   config.WaitTimeSeconds,
        visibilityTimeout: config.VisibilityTimeout,
        handler:           handler,
        stopCh:            make(chan struct{}),
        stoppedCh:         make(chan struct{}),
    }
}
```

**Key Configuration:**

- **Sequential Processing**: `MaxMessages: 1` (one message at a time)
- **Parallel Processing**: `MaxMessages: 10` (up to 10 concurrent messages)
- **Long Polling**: `WaitTimeSeconds: 20` (reduces empty responses)
- **Visibility Timeout**: `VisibilityTimeout: 300` (5 minutes for processing)

#### Message Handler

Location: `internal/infrastructure/{messaging}/handler/{entity}.go`

```go
package handler

import (
    "context"
    "encoding/json"
    "time"

    "github.com/aws/aws-sdk-go-v2/service/sqs/types"
    "github.com/eaglys-platform/pandlock-api/internal/domain/errors"
    "github.com/eaglys-platform/pandlock-api/internal/usecase"
    "github.com/eaglys-platform/pandlock-api/internal/usecase/input"
)

type ProjectKeyCreationMessage struct {
    ProjectID string `json:"project_id"`
}

type ProjectKeyCreationHandler struct {
    interactor usecase.TaskProjectKeyCreationInteractor
}

func NewProjectKeyCreationHandler(
    interactor usecase.TaskProjectKeyCreationInteractor,
) *ProjectKeyCreationHandler {
    return &ProjectKeyCreationHandler{
        interactor: interactor,
    }
}

func (h *ProjectKeyCreationHandler) Handle(ctx context.Context, message types.Message) error {
    // 1. Parse message body
    var msg ProjectKeyCreationMessage
    if err := json.Unmarshal([]byte(*message.Body), &msg); err != nil {
        return errors.TaskInvalidMessageErr.
            Wrap(err).
            WithDetail("failed to parse message body")
    }

    // 2. Validate required fields
    if msg.ProjectID == "" {
        return errors.TaskInvalidMessageErr.
            Errorf("project_id is required in message")
    }

    // 3. Call interactor
    param := input.NewTaskProcessProjectKeyCreation(
        msg.ProjectID,
        time.Now(),
    )

    _, err := h.interactor.Process(ctx, param)
    if err != nil {
        return err
    }

    return nil
}
```

**Message Format:**

```json
{
  "project_id": "xxx"
}
```

### 3. Worker Command

#### Root Worker Command

Location: `internal/infrastructure/cmd/internal/worker_cmd/cmd.go`

```go
package worker_cmd

import (
    "github.com/eaglys-platform/pandlock-api/internal/infrastructure/cmd/internal/worker_cmd/project_key_creation_cmd"
    "github.com/spf13/cobra"
)

func NewWorkerCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "worker",
        Short: "cli worker",
        Run: func(cmd *cobra.Command, args []string) {
            if len(args) == 0 {
                cmd.HelpFunc()(cmd, args)
            }
        },
    }
    cmd.AddCommand(project_key_creation_cmd.NewProjectKeyCreationCmd())
    return cmd
}
```

#### Entity-Specific Worker Command

Location: `internal/infrastructure/cmd/internal/worker_cmd/{entity}_cmd/cmd.go`

```go
package project_key_creation_cmd

import (
    "context"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/caarlos0/env/v11"
    "github.com/eaglys-platform/pandlock-api/internal/infrastructure/dependency"
    "github.com/eaglys-platform/pandlock-api/internal/infrastructure/environment"
    "github.com/eaglys-platform/pandlock-api/internal/infrastructure/sqs"
    "github.com/eaglys-platform/pandlock-api/internal/infrastructure/sqs/handler"
    "github.com/eaglys-platform/pandlock-api/internal/infrastructure/sqs/subscriber"
    "github.com/eaglys-platform/pandlock-api/internal/pkg/logger"
    "github.com/spf13/cobra"
    "go.uber.org/zap"
)

func NewProjectKeyCreationCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "project-key-creation",
        Short: "running project key creation worker",
        Run: func(cmd *cobra.Command, args []string) {
            run()
        },
    }
}

func run() {
    // 1. Load environment variables
    e := &environment.Environment{}
    if err := env.Parse(e); err != nil {
        panic(err)
    }

    // 2. Initialize logger
    l := logger.New(e.MinLogLevel)
    ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Second)
    defer cancel()
    ctx = logger.ToContext(ctx, l)

    // 3. Initialize dependencies
    dep := &dependency.Dependency{}
    dep.Inject(ctx, e)

    // 4. Wrap messaging client
    sqsClient := &sqs.Client{Cli: dep.SQSClient}

    // 5. Create handler
    h := handler.NewProjectKeyCreationHandler(dep.TaskProjectKeyCreationInteractor)

    // 6. Create subscriber
    sub := subscriber.NewSubscriber(
        sqsClient,
        subscriber.SubscriberConfig{
            QueueURL:          e.AWSSQSProjectKeyCreationQueueURL,
            MaxMessages:       1,   // Sequential processing
            WaitTimeSeconds:   20,  // Long polling
            VisibilityTimeout: 300, // 5 minutes
        },
        h,
    )

    logger.L(ctx).Info("[Worker] Starting project key creation worker",
        zap.String("queue_url", e.AWSSQSProjectKeyCreationQueueURL))

    // 7. Start subscriber in goroutine
    errCh := make(chan error, 1)
    go func() {
        if err := sub.Start(ctx); err != nil {
            errCh <- err
        }
    }()

    // 8. Setup signal handling
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGTERM, os.Interrupt)

    // 9. Wait for signal or error
    select {
    case sig := <-quit:
        logger.L(ctx).Info("[Worker] Received signal, shutting down...", zap.String("signal", sig.String()))
        sub.Stop()
        sub.WaitForShutdown()
        logger.L(ctx).Info("[Worker] Graceful shutdown complete")
    case err := <-errCh:
        logger.L(ctx).Error("[Worker] Subscriber error", zap.Error(err))
        os.Exit(1)
    }
}
```

### 4. Dependency Injection

Add to existing `Dependency` struct:

Location: `internal/infrastructure/dependency/dependency.go`

```go
import (
    awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Dependency struct {
    DatabaseCli *database.Client
    SQSClient   *awssqs.Client  // Add messaging client

    // ... existing fields

    // task (add new section)
    TaskProjectKeyCreationInteractor usecase.TaskProjectKeyCreationInteractor
}

func (d *Dependency) Inject(ctx context.Context, e *environment.Environment) {
    // ... existing initialization

    // Initialize messaging client
    awsSession := aws.NewConfig(ctx, e.AWSRegion)
    d.SQSClient = awssqs.NewFromConfig(awsSession)

    // ... existing repository initialization

    // Initialize task interactors
    d.TaskProjectKeyCreationInteractor = usecase.NewTaskProjectKeyCreationInteractor(
        transactable,
        projectRepository,
    )
}
```

### 5. Environment Variables

Location: `internal/infrastructure/environment/env.go`

```go
type AWSEnvironment struct {
    AWSRegion                             string `env:"AWS_REGION,required"`
    AWSEmulatorHost                       string `env:"AWS_EMULATOR_HOST"`
    // ... existing fields
    AWSSQSProjectKeyCreationQueueURL      string `env:"AWS_SQS_PROJECT_KEY_CREATION_QUEUE_URL,required"`
}
```

**Naming Convention:**

- Prefix: `AWS_SQS_` for SQS queues
- Prefix: `GCP_PUBSUB_` for Pub/Sub topics (future)
- Pattern: `{PROVIDER}_{SERVICE}_{ENTITY}_{RESOURCE_TYPE}`

### 6. Command Registration

Location: `internal/infrastructure/cmd/root.go`

```go
import (
    "github.com/eaglys-platform/pandlock-api/internal/infrastructure/cmd/internal/worker_cmd"
)

func NewCmdRoot() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "app",
        Short: "cli tool for app",
        // ...
    }
    cmd.AddCommand(http_server_cmd.NewHTTPServerCmd())
    cmd.AddCommand(task_cmd.NewTaskCmd())
    cmd.AddCommand(schema_migration_cmd.NewSchemaMigrationCmd())
    cmd.AddCommand(worker_cmd.NewWorkerCmd())  // Add worker command
    return cmd
}
```

## Error Handling Strategy

### Message Acknowledgment

| Error Type       | Action         | Reason                   |
| ---------------- | -------------- | ------------------------ |
| Success          | Delete message | Completed successfully   |
| Validation error | Delete message | Won't succeed on retry   |
| Not found error  | Delete message | Entity doesn't exist     |
| Transient error  | Keep visible   | Will retry after timeout |

### Implementation

```go
func (s *Subscriber) processMessage(ctx context.Context, message types.Message) error {
    err := s.handler.Handle(ctx, message)
    if err != nil {
        // Check if error is retryable
        if isRetryableError(err) {
            logger.L(ctx).Warn("[SQS Subscriber] Retryable error, keeping message visible",
                zap.Error(err))
            return nil // Don't delete, let visibility timeout expire
        }
        // Non-retryable error, delete message
        logger.L(ctx).Error("[SQS Subscriber] Non-retryable error, deleting message",
            zap.Error(err))
    }

    // Delete message on success or non-retryable error
    if err := s.deleteMessage(ctx, message); err != nil {
        logger.L(ctx).Error("[SQS Subscriber] Failed to delete message", zap.Error(err))
        return err
    }

    return nil
}
```

## LocalStack Setup (Development)

### Docker Compose Configuration

Location: `docker-compose.yml`

```yaml
services:
  aws:
    build: ./localstack
    environment:
      SERVICES: s3,sns,sqs # Enable SQS service
      AWS_DEFAULT_REGION: ap-northeast-1
```

### Initialization Script

Location: `localstack/script/init.sh`

```bash
#!/bin/sh

echo "SQS setup start!"
echo "Creating SQS queues..."

# Create queue for project key creation
aws --endpoint-url=http://aws:4566 sqs create-queue --queue-name project-key-creation
echo 'project-key-creation queue created!'

echo "SQS setup Done!"
```

### Environment Variables

Location: `.envrc` / `.envrc.tmpl`

```bash
export AWS_REGION="ap-northeast-1"
export AWS_EMULATOR_HOST="localhost:4566"
export AWS_SQS_PROJECT_KEY_CREATION_QUEUE_URL="http://localhost:4566/000000000000/project-key-creation"
```

## Processing Modes

### Sequential Processing

For tasks that **cannot run in parallel** (e.g., key generation, financial transactions):

```go
subscriber.SubscriberConfig{
    QueueURL:          queueURL,
    MaxMessages:       1,   // One message at a time
    WaitTimeSeconds:   20,
    VisibilityTimeout: 300,
}
```

### Parallel Processing

For tasks that **can run in parallel** (e.g., email sending, data export):

```go
subscriber.SubscriberConfig{
    QueueURL:          queueURL,
    MaxMessages:       10,  // Up to 10 concurrent messages
    WaitTimeSeconds:   20,
    VisibilityTimeout: 300,
}
```

## Running Workers

```bash
# View available workers
./app worker --help

# Run specific worker
./app worker project-key-creation

# Run with environment variables
source .envrc
./app worker project-key-creation
```

## Testing Workers

### Manual Testing

```bash
# Send test message (LocalStack)
aws --endpoint-url=http://localhost:4566 sqs send-message \
  --queue-url http://localhost:4566/000000000000/project-key-creation \
  --message-body '{"project_id": "test-project-123"}'

# Start worker
./app worker project-key-creation
```

### Verify Logs

```
[Worker] Starting project key creation worker queue_url=http://localhost:4566/000000000000/project-key-creation
[SQS Subscriber] Starting subscriber queue_url=http://localhost:4566/000000000000/project-key-creation max_messages=1
[SQS Subscriber] Processing message: msg-abc123
[SQS Subscriber] Successfully processed message: msg-abc123
[SQS Subscriber] Deleted message: msg-abc123
```

## Best Practices

1. **Single Binary** - Use `cmd/app` with Cobra commands, not separate binaries
2. **Hierarchical Commands** - Group workers under `worker` root command
3. **Reuse Dependencies** - Use existing `Dependency` struct, don't create new DI files
4. **Environment Naming** - Follow existing conventions with proper prefixes (AWS*, GCP*, etc.)
5. **Structured Logging** - Always use `logger.L(ctx)` with structured fields
6. **Graceful Shutdown** - Handle SIGINT/SIGTERM signals properly
7. **Transaction Safety** - Use `ForUpdate: true` when modifying entities
8. **Error Handling** - Distinguish retryable vs non-retryable errors
9. **Message Validation** - Validate message format before processing
10. **Idempotency** - Design handlers to be idempotent (safe to retry)

## Common Patterns

### Idempotent Processing

```go
func (i *taskInteractor) Process(ctx context.Context, param *input.TaskProcess) error {
    return i.transactable.RWTx(ctx, func(ctx context.Context) error {
        entity, err := i.repository.Get(ctx, repository.GetQuery{
            ID:        null.StringFrom(param.EntityID),
            ForUpdate: true,
        })
        if err != nil {
            return err
        }

        // Check if already processed (idempotency)
        if entity.ProcessedAt.Valid {
            logger.L(ctx).Info("Already processed, skipping", zap.String("entity_id", entity.ID))
            return nil
        }

        // Process entity
        entity.MarkAsProcessed(time.Now())

        return i.repository.Update(ctx, entity)
    })
}
```

### Batch Processing

```go
type BatchMessage struct {
    EntityIDs []string `json:"entity_ids"`
}

func (h *Handler) Handle(ctx context.Context, message types.Message) error {
    var msg BatchMessage
    if err := json.Unmarshal([]byte(*message.Body), &msg); err != nil {
        return errors.TaskInvalidMessageErr.Wrap(err)
    }

    // Process in batches
    for _, id := range msg.EntityIDs {
        param := input.NewTaskProcess(id, time.Now())
        if _, err := h.interactor.Process(ctx, param); err != nil {
            logger.L(ctx).Error("Failed to process entity",
                zap.String("entity_id", id),
                zap.Error(err))
            // Continue processing other entities
        }
    }

    return nil
}
```
