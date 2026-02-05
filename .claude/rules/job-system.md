# Job System Guidelines

## Overview

The job system provides asynchronous task execution using a serverless architecture (AWS Batch or Cloud Run Jobs). This approach is chosen when:

- **Heavy processing**: Large record batch processing, file generation (CSV/PDF/ZIP)
- **Avoiding server load**: Moving intensive tasks off the main HTTP server
- **Scalability**: Serverless execution scales based on workload

## Table Design Philosophy

### Core Tables

#### `jobs` - Main Job Table

Tracks job execution state and errors.

```sql
CREATE TABLE `jobs` (
  `id`                       VARCHAR(64)    NOT NULL COMMENT "ジョブID",
  `tenant_id`                VARCHAR(64)    NOT NULL COMMENT "テナントID",
  `job_type`                 VARCHAR(64)    NOT NULL COMMENT "ジョブタイプ",
  `status`                   VARCHAR(64)    NOT NULL COMMENT "ジョブステータス",
  `error_code`               VARCHAR(64)    NULL COMMENT "エラーコード",
  `error_message`            TEXT           NULL COMMENT "エラーメッセージ",
  `created_at`               DATETIME       NOT NULL COMMENT "作成日時",
  `updated_at`               DATETIME       NOT NULL COMMENT "更新日時",
  CONSTRAINT `jobs_pkey` PRIMARY KEY (`id`),
  CONSTRAINT `jobs_fkey_tenant_id` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`),
  CONSTRAINT `jobs_fkey_job_type` FOREIGN KEY (`job_type`) REFERENCES `job_types` (`id`),
  CONSTRAINT `jobs_fkey_status` FOREIGN KEY (`status`) REFERENCES `job_statuses` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT "ジョブの管理テーブル";
```

#### `job_types` - Job Type Constants

```yaml
# db/main/constants/constants.yaml
- table: job_types
  values:
    - issue_gift_catalog
    - activate_gift_catalog
    - send_gift_catalog_codes
    - download_gift_catalog_codes
    - generate_gift_catalog_delivery_printing
```

#### `job_statuses` - Job Status Constants

```yaml
- table: job_statuses
  values:
    - queued      # Initial state - waiting to be kicked
    - started     # Execution in progress
    - completed   # Successfully finished
    - failed      # Execution failed with error
```

### Job Detail Tables (One per Job Type)

Each `job_type` has a corresponding detail table storing type-specific parameters.

#### Pattern: One-to-One with `jobs` Table

```sql
CREATE TABLE `job_issue_gift_catalogs` (
  `id`                              VARCHAR(64)    NOT NULL COMMENT "ジョブ発行ギフトカタログID",
  `tenant_id`                       VARCHAR(64)    NOT NULL COMMENT "テナントID",
  `job_id`                          VARCHAR(64)    NOT NULL COMMENT "ジョブID",
  `gift_catalog_id`                 VARCHAR(64)    NOT NULL COMMENT "発行するギフトカタログID",
  `created_at`                      DATETIME       NOT NULL COMMENT "作成日時",
  `updated_at`                      DATETIME       NOT NULL COMMENT "更新日時",
  CONSTRAINT `job_issue_gcs_pkey` PRIMARY KEY (`id`),
  CONSTRAINT `job_issue_gcs_fkey_tenant_id` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`),
  CONSTRAINT `job_issue_gcs_fkey_job_id` FOREIGN KEY (`job_id`) REFERENCES `jobs` (`id`),
  CONSTRAINT `job_issue_gcs_fkey_gift_catalog_id` FOREIGN KEY (`gift_catalog_id`) REFERENCES `gift_catalogs` (`id`),
  UNIQUE `job_issue_gcs_unique_job_id` (`job_id`)  -- Enforces 1:1 relationship
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT "ギフトカタログの発行ジョブの詳細管理テーブル";
```

### Design Principles

1. **Separation of Concerns**: Main `jobs` table tracks execution state; detail tables store job-specific parameters
2. **1:1 Relationship**: UNIQUE constraint on `job_id` in detail tables ensures one-to-one mapping
3. **Tenant Scoping**: All tables include `tenant_id` for multi-tenant isolation
4. **Error Tracking**: `error_code` and `error_message` capture failure details
5. **Status Lifecycle**: `queued` → `started` → `completed` OR `failed`

## Domain Model

Location: `internal/domain/model/job.go`

### Job Entity

```go
type Job struct {
    ID                                  string
    TenantID                            string
    JobType                             JobType
    Status                              JobStatus
    ErrorCode                           null.String
    ErrorMessage                        null.String
    // Job-specific details (only one will be valid per job)
    IssueGiftCatalog                    nullable.Type[JobIssueGiftCatalog]
    ActivateGiftCatalog                 nullable.Type[JobActivateGiftCatalog]
    SendGiftCatalogCodes                nullable.Type[JobSendGiftCatalogCodes]
    DownloadGiftCatalogCodes            nullable.Type[JobDownloadGiftCatalogCodes]
    GenerateGiftCatalogDeliveryPrinting nullable.Type[JobGenerateGiftCatalogDeliveryPrinting]
    CreatedAt                           time.Time
    UpdatedAt                           time.Time
}
```

### State Transition Methods

#### Start

```go
func (m *Job) Start(requestTime time.Time) (*Job, error) {
    if m.Status != JobStatusQueued {
        return nil, errors.JobCanNotStartErr.New().
            WithDetail("job is not queued").
            WithValue("job_id", m.ID)
    }

    m.Status = JobStatusStarted
    m.UpdatedAt = requestTime

    return m, nil
}
```

#### Complete

```go
func (m *Job) Complete(requestTime time.Time) *Job {
    m.Status = JobStatusCompleted
    m.UpdatedAt = requestTime
    return m
}
```

#### Fail

```go
func (m *Job) Fail(err error, requestTime time.Time) *Job {
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
    m.UpdatedAt = requestTime
    return m
}
```

### Constructor Pattern (One per Job Type)

```go
func NewJobIssueGiftCatalog(
    tenantID string,
    giftCatalogID string,
    t time.Time,
) *Job {
    return &Job{
        ID:           id.New(),
        TenantID:     tenantID,
        JobType:      JobTypeIssueGiftCatalog,
        Status:       JobStatusQueued,
        ErrorCode:    null.String{},
        ErrorMessage: null.String{},
        IssueGiftCatalog: nullable.TypeFrom(
            JobIssueGiftCatalog{
                ID:            id.New(),
                GiftCatalogID: giftCatalogID,
                CreatedAt:     t,
                UpdatedAt:     t,
            },
        ),
        // Other detail fields initialized as empty nullable.Type
        ActivateGiftCatalog:                 nullable.Type[JobActivateGiftCatalog]{},
        SendGiftCatalogCodes:                nullable.Type[JobSendGiftCatalogCodes]{},
        DownloadGiftCatalogCodes:            nullable.Type[JobDownloadGiftCatalogCodes]{},
        GenerateGiftCatalogDeliveryPrinting: nullable.Type[JobGenerateGiftCatalogDeliveryPrinting]{},
        CreatedAt:                           t,
        UpdatedAt:                           t,
    }
}
```

## Publisher Integration (SNS → AWS Batch / Cloud Run Jobs)

### Interface Definition

Location: `internal/domain/message/publisher.go`

```go
package message

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_message
type Publisher interface {
    KickJob(ctx context.Context, jobId string) error
    // NOTE: 重い処理を行うジョブを起動する
    // 現時点では以下でのみ呼び出して良い
    // - CSVエクスポート
    // - PDF印刷
    KickHeavyJob(ctx context.Context, jobId string) error
}
```

### Two Job Queues

- **KickJob**: Light processing (fast, simple operations)
- **KickHeavyJob**: Heavy processing (file generation, large batch operations)

**Why separate queues?**
- Different resource allocation (CPU/Memory)
- Independent scaling policies
- Cost optimization (light jobs use smaller instances)

### Implementation (SNS → SQS → AWS Batch)

Location: `internal/infrastructure/sns/message/publisher.go`

```go
package message

type publisher struct {
    cli                    *sns.Client
    applicationEnvironment environment.ApplicationEnvironment
    batchJobLightTopicARN  string
    batchJobHeavyTopicARN  string
}

func NewPublisher(
    cli *sns.Client,
    applicationEnvironment environment.ApplicationEnvironment,
    batchJobLightTopicARN string,
    batchJobHeavyTopicARN string,
) message.Publisher {
    return &publisher{
        cli:                    cli,
        applicationEnvironment: applicationEnvironment,
        batchJobLightTopicARN:  batchJobLightTopicARN,
        batchJobHeavyTopicARN:  batchJobHeavyTopicARN,
    }
}

func (p *publisher) KickJob(ctx context.Context, jobId string) error {
    bytes, err := payload.EncodePublishMessage(
        payload.NewKickJob(jobId),
    )
    if err != nil {
        return err
    }

    // sqsの先では、AWS batchを実行するため、local環境では何もしない.
    if p.applicationEnvironment == environment.ApplicationEnvironmentLocal {
        return nil
    }

    pi := &sns.PublishInput{
        Message:          aws.String(string(bytes)),
        MessageStructure: aws.String("json"),
        TopicArn:         aws.String(p.batchJobLightTopicARN),
    }

    if _, err := p.cli.Publish(ctx, pi); err != nil {
        return errors.InternalErr.Wrap(err)
    }

    return nil
}

func (p *publisher) KickHeavyJob(ctx context.Context, jobId string) error {
    // Same as KickJob but uses batchJobHeavyTopicARN
    // ...
}
```

### Payload Structure

Location: `internal/infrastructure/sns/payload/job.go`

```go
type KickJob struct {
    JobID string `json:"job_id"`
}

func NewKickJob(jobID string) *KickJob {
    return &KickJob{
        JobID: jobID,
    }
}
```

## Task Command Pattern

### Command Structure

Location: `cmd/app/main.go` → `internal/infrastructure/cmd/`

```
cmd/app/
  └── main.go (entry point)
      └── internal/infrastructure/cmd/
          ├── root.go (cobra root command)
          └── internal/
              └── task_cmd/
                  ├── cmd.go (task subcommand registration)
                  └── process_job_cmd/
                      ├── cmd.go (cobra command definition)
                      └── process_job.go (business logic)
```

### Process Job Command (Generic Job Executor)

Location: `internal/infrastructure/cmd/internal/task_cmd/process_job_cmd/process_job.go`

```go
package process_job_cmd

type CMD struct {
    ctx           context.Context
    jobInteractor usecase.JobInteractor
}

func (c *CMD) Run(cmd *cobra.Command) error {
    jobID := cmd.Flag("job-id").Value.String()
    startTime := now.Now()
    logger.L(c.ctx).Info("start process job", zap.Time("start_time", startTime))

    if err := c.jobInteractor.Process(
        c.ctx,
        input.NewJobProcess(jobID, startTime),
    ); err != nil {
        endTime := now.Now()
        logger.L(c.ctx).Error(
            "failed to process job",
            logger_field.Error(err),
            zap.String("job_id", jobID),
            zap.Time("start_time", startTime),
            zap.Time("end_time", endTime),
            zap.Duration("duration", endTime.Sub(startTime)),
        )
        return err
    }

    endTime := now.Now()
    logger.L(c.ctx).Info(
        "completed to process job",
        zap.String("job_id", jobID),
        zap.Time("start_time", startTime),
        zap.Time("end_time", endTime),
        zap.Duration("duration", endTime.Sub(startTime)),
    )
    return nil
}
```

### Job Interactor (Generic Job Processor)

Location: `internal/usecase/job_impl.go`

```go
func (i *jobInteractor) Process(ctx context.Context, param *input.JobProcess) error {
    if err := param.Validate(); err != nil {
        return err
    }

    var job *model.Job

    // 1. Start the job (update status to "started")
    if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
        var err error
        job, err = i.jobRepository.Get(ctx, repository.GetJobQuery{
            ID: null.StringFrom(param.JobID),
            BaseGetOptions: repository.BaseGetOptions{
                OrFail:    true,
                ForUpdate: true,
            },
        })
        if err != nil {
            return err
        }

        job, err = job.Start(param.RequestTime)
        if err != nil {
            return err
        }

        if err := i.jobRepository.Update(ctx, job); err != nil {
            return err
        }
        return nil
    }); err != nil {
        return err
    }

    // 2. Execute the actual job based on job_type
    var jobErr error
    switch job.JobType {
    case model.JobTypeIssueGiftCatalog:
        jobErr = i.transactable.RWTx(ctx, func(ctx context.Context) error {
            _, err := i.giftCatalogService.Issue(ctx, service.GiftCatalogIssueParam{
                GiftCatalogID: job.IssueGiftCatalog.Ptr().GiftCatalogID,
                RequestTime:   param.RequestTime,
            })
            return err
        })
    case model.JobTypeActivateGiftCatalog:
        jobErr = i.transactable.RWTx(ctx, func(ctx context.Context) error {
            _, err := i.giftCatalogService.Activate(ctx, service.GiftCatalogActivateParam{
                GiftCatalogID: job.ActivateGiftCatalog.Ptr().GiftCatalogID,
                RequestTime:   param.RequestTime,
            })
            return err
        })
    // ... other job types
    }

    // 3. Handle job failure
    if jobErr != nil {
        if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
            job, err := i.jobRepository.Get(ctx, repository.GetJobQuery{
                ID: null.StringFrom(param.JobID),
                BaseGetOptions: repository.BaseGetOptions{
                    OrFail:    true,
                    ForUpdate: true,
                },
            })
            if err != nil {
                return err
            }
            logger.L(ctx).Error(
                fmt.Sprintf("failed to process %s job: %s", job.JobType.String(), job.ID),
                logger_field.Error(jobErr),
                zap.Reflect("job", job),
            )
            job = job.Fail(jobErr, param.RequestTime)
            if err := i.jobRepository.Update(ctx, job); err != nil {
                return err
            }
            return nil
        }); err != nil {
            return err
        }
        return nil  // Return nil to prevent retry - job status is already "failed"
    }

    // 4. Mark job as completed
    if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
        var err error
        job, err = i.jobRepository.Get(ctx, repository.GetJobQuery{
            ID: null.StringFrom(param.JobID),
            BaseGetOptions: repository.BaseGetOptions{
                OrFail:    true,
                ForUpdate: true,
            },
        })
        if err != nil {
            return err
        }
        job = job.Complete(param.RequestTime)
        if err := i.jobRepository.Update(ctx, job); err != nil {
            return err
        }
        return nil
    }); err != nil {
        return err
    }

    return nil
}
```

## Implementation Workflow

### Adding a New Job Type

Follow these steps when adding a new job type:

#### 1. Add Migration for Job Detail Table

Location: `db/main/migrations/XX_add_job_{job_name}.sql`

```sql
-- +goose Up
CREATE TABLE `job_{job_name}` (
  `id`          VARCHAR(64) NOT NULL COMMENT 'ジョブXXID',
  `tenant_id`   VARCHAR(64) NOT NULL COMMENT 'テナントID',
  `job_id`      VARCHAR(64) NOT NULL COMMENT 'ジョブID',
  -- Add job-specific parameter columns here
  `param_field` VARCHAR(64) NOT NULL COMMENT 'パラメータ',
  `created_at`  DATETIME    NOT NULL COMMENT '作成日時',
  `updated_at`  DATETIME    NOT NULL COMMENT '更新日時',
  CONSTRAINT `job_{job_name}_pkey` PRIMARY KEY (`id`),
  CONSTRAINT `job_{job_name}_fkey_tenant_id` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`),
  CONSTRAINT `job_{job_name}_fkey_job_id` FOREIGN KEY (`job_id`) REFERENCES `jobs` (`id`),
  UNIQUE `job_{job_name}_unique_job_id` (`job_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT 'XXジョブの詳細管理テーブル';

-- +goose Down
DROP TABLE `job_{job_name}`;
```

#### 2. Add Job Type to Constants

Location: `db/main/constants/constants.yaml`

```yaml
- table: job_types
  values:
    # ... existing values
    - {new_job_type}  # Add your new job type here
```

#### 3. Add Domain Model

Location: `internal/domain/model/job.go`

```go
// Add new job type constant
const (
    // ... existing constants
    JobType{NewJobName} JobType = "{new_job_type}"
)

// Add detail struct
type Job{NewJobName} struct {
    ID          string
    ParamField  string  // Job-specific parameters
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// Add field to Job entity
type Job struct {
    // ... existing fields
    {NewJobName} nullable.Type[Job{NewJobName}]
}

// Add constructor
func NewJob{NewJobName}(
    tenantID string,
    paramField string,  // Job-specific parameters
    t time.Time,
) *Job {
    return &Job{
        ID:           id.New(),
        TenantID:     tenantID,
        JobType:      JobType{NewJobName},
        Status:       JobStatusQueued,
        ErrorCode:    null.String{},
        ErrorMessage: null.String{},
        {NewJobName}: nullable.TypeFrom(
            Job{NewJobName}{
                ID:         id.New(),
                ParamField: paramField,
                CreatedAt:  t,
                UpdatedAt:  t,
            },
        ),
        // Initialize all other job detail fields as empty
        IssueGiftCatalog:                    nullable.Type[JobIssueGiftCatalog]{},
        ActivateGiftCatalog:                 nullable.Type[JobActivateGiftCatalog]{},
        // ... other fields
        CreatedAt: t,
        UpdatedAt: t,
    }
}
```

#### 4. Implement Job Logic in Domain Service

Location: `internal/domain/service/{resource}.go` or create new service

```go
type {Resource}Service interface {
    {Action}(ctx context.Context, param {Resource}{Action}Param) (*{Resource}{Action}Result, error)
}

type {Resource}{Action}Param struct {
    // Parameters from job detail table
    ParamField  string
    RequestTime time.Time
}

type {Resource}{Action}Result struct {
    // Return values if needed
}
```

#### 5. Add Case to JobInteractor.Process()

Location: `internal/usecase/job_impl.go`

```go
func (i *jobInteractor) Process(ctx context.Context, param *input.JobProcess) error {
    // ... existing code

    switch job.JobType {
    // ... existing cases
    case model.JobType{NewJobName}:
        jobErr = i.transactable.RWTx(ctx, func(ctx context.Context) error {
            _, err := i.{resource}Service.{Action}(ctx, service.{Resource}{Action}Param{
                ParamField:  job.{NewJobName}.Ptr().ParamField,
                RequestTime: param.RequestTime,
            })
            return err
        })
    }

    // ... rest of the code
}
```

#### 6. Create Job in Usecase

When you need to kick off a job:

```go
func (i *interactor) SomeOperation(ctx context.Context, param *input.SomeOperation) error {
    // Within transaction
    if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
        // Create job entity
        job := model.NewJob{NewJobName}(
            param.TenantID,
            param.ParamField,
            param.RequestTime,
        )

        // Persist job
        if err := i.jobRepository.Create(ctx, job); err != nil {
            return err
        }

        // Kick job via publisher
        // Use KickJob for light processing, KickHeavyJob for heavy processing
        if err := i.publisher.KickJob(ctx, job.ID); err != nil {
            return err
        }

        return nil
    }); err != nil {
        return err
    }

    return nil
}
```

## Best Practices

### 1. Transaction Boundaries

- **Create job within transaction**: Ensure job record is persisted before kicking
- **Publisher call within transaction**: If publisher fails, job record is rolled back
- **Job execution in separate transaction**: Each job type operation runs in its own transaction

### 2. Light vs Heavy Jobs

Use `KickJob` for:
- Database CRUD operations
- Simple API calls
- Fast operations (< 30 seconds)

Use `KickHeavyJob` for:
- File generation (CSV, PDF, ZIP)
- Large batch processing (1000+ records)
- Image/video processing
- Long-running operations (> 30 seconds)

### 3. Error Handling

- **Always call `job.Fail(err, t)`**: Captures error details for debugging
- **Log errors before failing**: Use structured logging with job context
- **Return nil after fail**: Prevent retry loops - job status is already "failed"

### 4. Job Parameters

- **Store in detail table**: Don't use JSON columns in main `jobs` table
- **Keep parameters minimal**: Only store IDs and essential config
- **Use domain service for logic**: Job interactor delegates to domain services

### 5. Local Development

- **Publisher no-op in local**: `if applicationEnvironment == Local { return nil }`
- **Manual execution**: Use task command directly: `./app task process-job --job-id={id}`
- **Testing**: Create job record, then call JobInteractor.Process() directly

## Infrastructure Setup

### AWS Batch Architecture

```
HTTP Request
    ↓
Create Job Record (status=queued)
    ↓
Publisher.KickJob(jobID)
    ↓
SNS Topic (batchJobLightTopicARN or batchJobHeavyTopicARN)
    ↓
SQS Queue
    ↓
AWS Batch (triggered by SQS)
    ↓
ECS Task runs: ./app task process-job --job-id={jobID}
    ↓
JobInteractor.Process() updates job status
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
JobInteractor.Process() updates job status
    ↓
Job completed or failed
```

### Key Differences Between AWS and GCP

| Aspect | AWS Batch | Cloud Run Jobs |
|--------|-----------|----------------|
| Message Queue | SNS → SQS | Pub/Sub |
| Compute | ECS Tasks | Cloud Run Containers |
| Publisher Implementation | `internal/infrastructure/sns/message/` | `internal/infrastructure/pubsub/message/` |
| Environment Variable | `BATCH_JOB_LIGHT_TOPIC_ARN` | `PUBSUB_LIGHT_TOPIC_NAME` |

## Why This Design?

### Problem: Heavy Processing on HTTP Server

- **Request timeout**: HTTP requests timeout after 30-60 seconds
- **Resource contention**: Heavy jobs block other requests
- **Scaling difficulty**: Can't scale HTTP and batch workloads independently

### Solution: Async Job Queue + Serverless Execution

- **Decoupled execution**: HTTP server just creates job record and returns immediately
- **Serverless scaling**: AWS Batch/Cloud Run Jobs scale based on queue depth
- **Cost efficient**: Pay only for execution time
- **Retry resilience**: Failed jobs can be retried without user intervention
- **Progress tracking**: Job status in database provides real-time feedback

### Alternative Approaches

| Approach | When to Use | Trade-offs |
|----------|-------------|------------|
| Synchronous | < 5 seconds, simple operations | Simple but blocks request |
| Background goroutine | < 30 seconds, fire-and-forget | No retry, no status tracking |
| Cron jobs | Scheduled batch operations | Fixed schedule, not event-driven |
| **Job queue + serverless** | Heavy, event-driven processing | Best for scale, requires infrastructure |

## Common Patterns

### Job with Multiple Steps

For jobs requiring multiple operations:

```go
case model.JobTypeGenerateGiftCatalogDeliveryPrinting:
    jobErr = func() error {
        // Step 1: Reserve resources (in transaction)
        var reserveResult *service.ReserveResult
        if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
            var err error
            reserveResult, err = i.service.Reserve(ctx, ...)
            return err
        }); err != nil {
            return err
        }

        // Step 2: Heavy processing (outside transaction)
        generateResult, err := i.service.Generate(ctx, ...)
        if err != nil {
            return err
        }

        // Step 3: Finalize (in transaction)
        if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
            _, err := i.service.Complete(ctx, ...)
            return err
        }); err != nil {
            return err
        }

        return nil
    }()
```

### Job Chaining (Kick Another Job After Completion)

```go
func (i *service) Complete(ctx context.Context, param Param) (*Result, error) {
    // Complete current operation
    // ...

    // Create next job
    nextJob := model.NewJob{NextJobName}(...)
    if err := i.jobRepository.Create(ctx, nextJob); err != nil {
        return nil, err
    }

    // Kick next job
    if err := i.publisher.KickJob(ctx, nextJob.ID); err != nil {
        return nil, err
    }

    return result, nil
}
```
