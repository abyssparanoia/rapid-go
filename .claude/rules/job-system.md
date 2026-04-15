# Job System Guidelines

## Overview

The job system provides asynchronous task execution using a serverless architecture (AWS Batch or Cloud Run Jobs). This approach is chosen when:

- **Heavy processing**: Large record batch processing, file generation (CSV/PDF/ZIP)
- **Avoiding server load**: Moving intensive tasks off the main HTTP server
- **Scalability**: Serverless execution scales based on workload

## Architecture

```
HTTP Request
    ↓
Create Job Record (status=queued)
    ↓
Publisher.KickJob(jobID)
    ↓
SNS/Pub-Sub Topic
    ↓
AWS Batch / Cloud Run Jobs
    ↓
./app task process-job --job-id={jobID}
    ↓
CMD: JobInteractor.Start → Task{Type}Interactor.Process{JobType}Job → JobInteractor.Complete/Fail
    ↓
Job completed or failed
```

**Key principle**: The CMD layer orchestrates the lifecycle — `Start`, type-specific `Process{JobType}Job`, then `Complete` or `Fail` — using separate dedicated interactors.

## Table Design

### `jobs` - Main Job Table

```sql
CREATE TABLE `jobs` (
  `id`              VARCHAR(64)  NOT NULL COMMENT "ジョブID",
  `job_type`        VARCHAR(64)  NOT NULL COMMENT "ジョブタイプ",
  `status`          VARCHAR(64)  NOT NULL COMMENT "ジョブステータス",
  `auth_context`    VARCHAR(256) NOT NULL COMMENT "認可コンテキスト (type:identifier)",
  `idempotency_key` VARCHAR(256) NOT NULL COMMENT "冪等キー",
  `metadata`        JSON         NOT NULL COMMENT "ジョブパラメータ (JSON)",
  `error_code`      VARCHAR(64)  NULL COMMENT "エラーコード",
  `error_message`   TEXT         NULL COMMENT "エラーメッセージ",
  `created_at`      DATETIME     NOT NULL COMMENT "作成日時",
  `updated_at`      DATETIME     NOT NULL COMMENT "更新日時",
  CONSTRAINT `jobs_pkey` PRIMARY KEY (`id`),
  CONSTRAINT `jobs_fkey_job_type` FOREIGN KEY (`job_type`) REFERENCES `job_types` (`id`),
  CONSTRAINT `jobs_fkey_status` FOREIGN KEY (`status`) REFERENCES `job_statuses` (`id`),
  UNIQUE `jobs_unique_idempotency_key` (`idempotency_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT "ジョブの管理テーブル";
```

**Key design decisions:**

1. **`metadata` JSON NOT NULL** — Always stored, even if empty (`{}`). Job-type-specific parameters are stored as JSON instead of individual detail tables. Each job type has a dedicated typed metadata struct in the domain model.
2. **`idempotency_key` UNIQUE** — Prevents duplicate job creation. Caller sets an arbitrary key; the DB enforces uniqueness.
3. **`auth_context`** — Same `type:identifier` pattern as `AssetAuthContext`. Types: `staff`, `admin`, `command`. Handles authorization context without tenant dependency.
4. **No `tenant_id`** — Jobs are not always tenant-dependent (e.g., admin panel operations). Authorization is managed via `auth_context`.
5. **No individual parameter tables** — Unlike the original design, there are no `job_{type}` detail tables.

### Constant Tables

```yaml
- table: job_types
  values:
    - initialize_bots

- table: job_statuses
  values:
    - queued
    - started
    - completed
    - failed
```

## Domain Model

Location: `internal/domain/model/job.go`, `internal/domain/model/job_metadata.go`

### JobAuthContext

Same `type:identifier` pattern as `AssetAuthContext`:

```go
type JobAuthContext string

func NewStaffJobAuthContext(staffID string) JobAuthContext {
    return JobAuthContext("staff:" + staffID)
}
func NewAdminJobAuthContext(adminID string) JobAuthContext {
    return JobAuthContext("admin:" + adminID)
}
func NewCommandJobAuthContext(commandName string) JobAuthContext {
    return JobAuthContext("command:" + commandName)
}
```

### JobMetadata Interface

Location: `internal/domain/model/job_metadata.go`

Each job type has its own typed metadata struct implementing `JobMetadata`:

```go
// JobMetadata represents type-safe metadata for a job.
type JobMetadata interface {
    jobType() JobType
    ToMap() map[string]any
}

// Convert raw map to typed metadata for the given job type.
func JobMetadataFromMap(jt JobType, m map[string]any) (JobMetadata, error)
```

#### Typed Metadata Structs

```go
// InitializeBotsMetadata is the metadata for JobTypeInitializeBots jobs.
type InitializeBotsMetadata struct{}

func (m *InitializeBotsMetadata) ToMap() map[string]any {
    return map[string]any{}
}
```

Add fields to the struct as job-specific parameters grow.

### Job Entity

```go
type Job struct {
    ID             string
    JobType        JobType
    Status         JobStatus
    AuthContext    JobAuthContext
    IdempotencyKey string
    Metadata       JobMetadata   // Always non-nil; typed per job type
    ErrorCode      null.String
    ErrorMessage   null.String
    CreatedAt      time.Time
    UpdatedAt      time.Time
}
```

### Constructor

```go
func NewJob(
    jobType JobType,
    authContext JobAuthContext,
    idempotencyKey string,
    metadata JobMetadata,   // Always pass a typed struct, never nil
    t time.Time,
) *Job
```

### State Transition Methods

#### Start

```go
func (m *Job) Start(t time.Time) (*Job, error) {
    if m.Status != JobStatusQueued {
        return nil, errors.JobCanNotStartErr.New().
            WithDetail("job is not queued").
            WithValue("job_id", m.ID).
            WithValue("status", m.Status.String())
    }
    m.Status = JobStatusStarted
    m.UpdatedAt = t
    return m, nil
}
```

#### Complete

```go
func (m *Job) Complete(t time.Time) (*Job, error) {
    if m.Status != JobStatusStarted {
        return nil, errors.JobCanNotCompleteErr.New().
            WithDetail("job is not started").
            WithValue("job_id", m.ID).
            WithValue("status", m.Status.String())
    }
    m.Status = JobStatusCompleted
    m.UpdatedAt = t
    return m, nil
}
```

#### Fail

```go
func (m *Job) Fail(err error, t time.Time) (*Job, error) {
    if m.Status != JobStatusStarted {
        return nil, errors.JobCanNotFailErr.New().
            WithDetail("job is not started").
            WithValue("job_id", m.ID).
            WithValue("status", m.Status.String())
    }
    // Uses goerr.Unwrap to extract structured error code/message
    var errorCode string
    var errorMessage string
    if goErr := goerr.Unwrap(err); goErr != nil {
        errorCode = goErr.Code()
        errorMessage = goErr.Error()
    } else {
        errorCode = "unknown"
        errorMessage = err.Error()
    }
    m.Status = JobStatusFailed
    m.ErrorCode = null.StringFrom(errorCode)
    m.ErrorMessage = null.StringFrom(errorMessage)
    m.UpdatedAt = t
    return m, nil
}
```

## Domain Errors

```go
JobNotFoundErr       = NewNotFoundError("E201501", "Job not found")
JobCanNotStartErr    = NewConflictError("E201502", "Job cannot be started")
JobAlreadyExistsErr  = NewConflictError("E201503", "Job with this idempotency key already exists")
JobCanNotCompleteErr = NewConflictError("E201504", "Job cannot be completed")
JobCanNotFailErr     = NewConflictError("E201505", "Job cannot be failed")
```

## Repository

Location: `internal/domain/repository/job.go`

```go
type Job interface {
    Get(ctx context.Context, query GetJobQuery) (*model.Job, error)
    Create(ctx context.Context, job *model.Job) error
    Update(ctx context.Context, job *model.Job) error
}

type GetJobQuery struct {
    BaseGetOptions
    ID             null.String
    IdempotencyKey null.String
}
```

## Marshaller

Location: `internal/infrastructure/mysql/internal/marshaller/job.go`

- `JobToModel`: returns `(*model.Job, error)` — `types.JSON` (NOT NULL) → `json.Unmarshal` to `map[string]any` (error returned on failure) → `model.JobMetadataFromMap(jobType, m)` for typed metadata (error returned on failure)
- `JobToDBModel`: `metadata.ToMap()` → `json.Marshal` → `sqltypes.JSON` (non-nullable `[]byte` alias)
- Uses `sqltypes "github.com/aarondl/sqlboiler/v4/types"` for DB types, `null/v9` for domain types

## Usecase Layer

### JobInteractor (lifecycle only)

Location: `internal/usecase/job.go`

```go
type JobInteractor interface {
    Start(ctx context.Context, param *input.JobStart) (*model.Job, error)
    Complete(ctx context.Context, param *input.JobComplete) error
    Fail(ctx context.Context, param *input.JobFail) error
}
```

Each method: `RWTx → Get(ForUpdate) → domain state change → Update`.

### Task{Type}Interactor (job creation + processing)

One interactor per job type with **two methods**:
- `Create{JobType}Job` — validates input, builds metadata, persists the job record, and kicks the publisher. Called from CLI commands or API handlers.
- `Process{JobType}Job` — fetches the job record, executes the actual business logic. Called exclusively from `process-job` CMD.

```go
// interface
type TaskBotInteractor interface {
    CreateInitializeBotsJob(ctx context.Context, param *input.TaskCreateInitializeBotsJob) (*model.Job, error)
    ProcessInitializeBotsJob(ctx context.Context, param *input.TaskProcessInitializeBotsJob) error
}

// Create input — holds the raw arguments (e.g. file bytes) needed to validate and enqueue the job
type TaskCreateInitializeBotsJob struct {
    CSVBytes    []byte    `validate:"required"`
    RequestTime time.Time `validate:"required"`
}

// Process input — holds only the job ID; all other parameters come from the persisted job record
type TaskProcessInitializeBotsJob struct {
    JobID       string    `validate:"required"`
    RequestTime time.Time `validate:"required"`
}
```

## Task Command Layer (CMD Orchestration)

Location: `internal/infrastructure/cmd/internal/task_cmd/process_job_cmd/`

The CMD layer orchestrates the three-phase execution:

```go
func (c *CMD) ProcessJob(cmd *cobra.Command) error {
    // 1. Start
    job, err := c.jobInteractor.Start(c.ctx, input.NewJobStart(jobID, requestTime))
    if err != nil {
        return err
    }

    // 2. Execute (type-specific)
    executeErr := c.execute(job)

    // 3. Complete or Fail
    if executeErr != nil {
        logger.L(c.ctx).Error("failed to execute job", ...)
        if err := c.jobInteractor.Fail(c.ctx, input.NewJobFail(jobID, executeErr, now.Now())); err != nil {
            logger.L(c.ctx).Error("failed to mark job as failed", ...)
            return executeErr  // Return executeErr so caller knows execution failed
        }
        return nil  // Return nil to prevent retry; job is already marked failed
    }
    return c.jobInteractor.Complete(c.ctx, input.NewJobComplete(jobID, now.Now()))
}

func (c *CMD) execute(job *model.Job) error {
    switch job.JobType {
    case model.JobTypeInitializeBots:
        return c.taskBotInteractor.ProcessInitializeBotsJob(
            c.ctx,
            input.NewTaskProcessInitializeBotsJob(job.ID, now.Now()),
        )
    case model.JobTypeUnknown:
        fallthrough
    default:
        return errors.InternalErr.Errorf("unknown job type: %s", job.JobType.String())
    }
}
```

**Why this separation:**

- `JobInteractor` is purely lifecycle management (Start/Complete/Fail) — no knowledge of job types
- `Task{Type}Interactor` owns both job creation (`Create{JobType}Job`) and job processing (`Process{JobType}Job`), keeping all job-type-specific logic in one place
- CMD layer wires them together, making each part independently testable

## Dependency Injection

```go
// Dependency struct
JobInteractor     usecase.JobInteractor
TaskBotInteractor usecase.TaskBotInteractor

// In Inject()
jobRepository := database_repository.NewJob()
d.JobInteractor = usecase.NewJobInteractor(transactable, jobRepository)
d.TaskBotInteractor = usecase.NewTaskBotInteractor(
    transactable,
    jobRepository,
    jobPublisher,
    // ... other dependencies required by Create and Process methods
)
```

## Adding a New Job Type

Follow these steps:

### 1. Add constant

```yaml
# db/mysql/constants/constants.yaml
- table: job_types
  values:
    - initialize_bots
    - your_new_job_type  # Add here
```

Run `make migrate.up`.

### 2. Add JobType constant

```go
// internal/domain/model/job.go
const (
    JobTypeUnknown        JobType = "unknown"
    JobTypeInitializeBots JobType = "initialize_bots"
    JobTypeYourNewJob     JobType = "your_new_job_type"  // Add here
)
```

### 3. Define typed metadata struct

```go
// internal/domain/model/job_metadata.go

// YourNewJobMetadata is the metadata for JobTypeYourNewJob jobs.
type YourNewJobMetadata struct {
    SomeParam string
    // Add job-specific fields here
}

func yourNewJobMetadataFromMap(m map[string]any) *YourNewJobMetadata {
    meta := &YourNewJobMetadata{}
    if v, ok := m["some_param"].(string); ok {
        meta.SomeParam = v
    }
    return meta
}

func (m *YourNewJobMetadata) jobType() JobType {
    return JobTypeYourNewJob
}

func (m *YourNewJobMetadata) ToMap() map[string]any {
    return map[string]any{
        "some_param": m.SomeParam,
    }
}
```

Also add the case to `JobMetadataFromMap`:

```go
func JobMetadataFromMap(jt JobType, m map[string]any) (JobMetadata, error) {
    switch jt {
    case JobTypeInitializeBots:
        return initializeBotsMetadataFromMap(m), nil
    case JobTypeYourNewJob:
        return yourNewJobMetadataFromMap(m), nil
    case JobTypeUnknown:
        fallthrough
    default:
        return nil, errors.InternalErr.Errorf("unknown job type for metadata: %s", jt.String())
    }
}
```

### 4. Create input DTOs

```go
// internal/usecase/input/task_your_new_job.go

// Create input — holds raw arguments for validation and job enqueuing
type TaskCreateYourNewJob struct {
    SomeParam   string    `validate:"required"`
    RequestTime time.Time `validate:"required"`
}

func NewTaskCreateYourNewJob(someParam string, requestTime time.Time) *TaskCreateYourNewJob {
    return &TaskCreateYourNewJob{SomeParam: someParam, RequestTime: requestTime}
}

func (p *TaskCreateYourNewJob) Validate() error {
    if err := validation.Validate(p); err != nil {
        return errors.RequestInvalidArgumentErr.Wrap(err)
    }
    return nil
}

// Process input — holds only the job ID; parameters come from the persisted job record
type TaskProcessYourNewJob struct {
    JobID       string    `validate:"required"`
    RequestTime time.Time `validate:"required"`
}

func NewTaskProcessYourNewJob(jobID string, requestTime time.Time) *TaskProcessYourNewJob {
    return &TaskProcessYourNewJob{JobID: jobID, RequestTime: requestTime}
}

func (p *TaskProcessYourNewJob) Validate() error {
    if err := validation.Validate(p); err != nil {
        return errors.RequestInvalidArgumentErr.Wrap(err)
    }
    return nil
}
```

### 5. Create interactor

```go
// internal/usecase/task_your_new_job.go
type TaskYourNewJobInteractor interface {
    CreateYourNewJob(ctx context.Context, param *input.TaskCreateYourNewJob) (*model.Job, error)
    ProcessYourNewJob(ctx context.Context, param *input.TaskProcessYourNewJob) error
}
```

### 6. Add case to CMD switch

```go
// process_job.go
case model.JobTypeYourNewJob:
    return c.taskYourNewJobInteractor.ProcessYourNewJob(c.ctx,
        input.NewTaskProcessYourNewJob(job.ID, now.Now()))
```

### 7. Register in DI

```go
d.TaskYourNewJobInteractor = usecase.NewTaskYourNewJobInteractor(...)
```

## Creating a Job (in Task{Type}Interactor.Create)

Job creation lives inside `Task{Type}Interactor.Create{JobType}Job`. It validates input, assembles metadata, persists the job record, and kicks the publisher — all within a transaction.

```go
func (i *taskBotInteractor) CreateInitializeBotsJob(
    ctx context.Context,
    param *input.TaskCreateInitializeBotsJob,
) (*model.Job, error) {
    if err := param.Validate(); err != nil {
        return nil, err
    }

    // Validate and build typed metadata from raw input
    // (e.g., parse CSV, upload to S3, store object key in metadata)
    metadata := &model.InitializeBotsMetadata{
        // SomeParam: derivedFromInput,
    }

    job := model.NewJob(
        model.JobTypeInitializeBots,
        model.NewCommandJobAuthContext("initialize-bots"),
        buildIdempotencyKey(param),        // deterministic key derived from input content
        metadata,
        param.RequestTime,
    )

    if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
        if err := i.jobRepository.Create(ctx, job); err != nil {
            return err
        }
        return i.jobPublisher.KickJob(ctx, job.ID)
    }); err != nil {
        return nil, err
    }

    return job, nil
}
```

**Idempotency key**: Must be **deterministic** — derived from the content/intent of the job, not from time. This allows the unique constraint to prevent duplicate jobs for the same input.

```go
// Good — content-based hash
func buildIdempotencyKey(param *input.TaskCreateInitializeBotsJob) string {
    sum := sha256.Sum256(param.CSVBytes)
    return fmt.Sprintf("initialize-bots:%x", sum)
}

// Bad — time-based (non-deterministic, defeats idempotency)
func buildIdempotencyKey(t time.Time) string {
    return fmt.Sprintf("initialize-bots:%d", t.UnixNano())
}
```

## Infrastructure Setup

### AWS Batch Architecture

```
HTTP Request
    ↓
Create Job Record (status=queued)
    ↓
Publisher.KickJob(jobID)
    ↓
SNS Topic
    ↓
SQS Queue
    ↓
AWS Batch (triggered by SQS)
    ↓
ECS Task runs: ./app task process-job --job-id={jobID}
    ↓
CMD: Start → Process → Complete/Fail
    ↓
Job completed or failed
```

### Cloud Run Jobs Architecture (GCP Alternative)

```
HTTP Request
    ↓
Create Job Record (status=queued)
    ↓
Publisher.KickJob(jobID)  [GCP Pub/Sub implementation]
    ↓
Pub/Sub Topic
    ↓
Cloud Run Jobs (triggered by Pub/Sub)
    ↓
Container runs: ./app task process-job --job-id={jobID}
    ↓
CMD: Start → Process → Complete/Fail
    ↓
Job completed or failed
```

## Best Practices

1. **Idempotency key design** — Use a deterministic key derived from input content (e.g., SHA-256 of file bytes, or a stable resource identifier). Never use timestamps — they produce a new key every call, defeating the purpose.
2. **Typed metadata** — Always use the concrete metadata struct for each job type; never use `map[string]any` directly
3. **metadata NOT NULL** — Always pass a typed metadata struct to `NewJob`; never pass nil
4. **Return nil after Fail** — Prevents retry loops; job is already marked failed. If `Fail` itself fails, log the failure and return `executeErr` to surface the original error.
5. **Status guards in Complete/Fail** — Both methods validate the job is in `started` status before transitioning; return `JobCanNotCompleteErr`/`JobCanNotFailErr` otherwise
6. **Two-method Task interactor** — `Create{JobType}Job` owns job creation (validate → build metadata → persist → kick); `Process{JobType}Job` owns job execution (fetch job → run business logic). Keep lifecycle (Start/Complete/Fail) in `JobInteractor` only.
7. **Type switch is exhaustive** — Always handle `JobTypeUnknown` with an error case
8. **Local development** — Publisher is no-op in local env; run manually with `./app task process-job --job-id={id}`
9. **Process input carries only job ID** — `TaskProcess{Type}` input structs hold only the job ID and request time. All job parameters come from the persisted job record (metadata). Never pass raw input data to the process method.
