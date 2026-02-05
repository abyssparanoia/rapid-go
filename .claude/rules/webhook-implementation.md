# Webhook Implementation Guidelines

## Overview

Webhooks are HTTP endpoints that receive notifications from external services. This guide provides a standardized pattern for implementing webhook endpoints that:
- Accept requests from external services (any HTTP method)
- Validate and process webhook payloads
- Return appropriate responses (JSON or plain text)
- Maintain consistency with the project's layered architecture

## Architecture Pattern

```
External Service → HTTP Server → Custom HTTP Handler → Internal gRPC → gRPC Handler → Usecase
                                 (Parse request)       (localhost:50051)  (Business logic)
```

**Key Principle**: HTTP layer should NOT call usecase directly. Always route through internal gRPC endpoint to maintain architectural consistency.

## Directory Structure

```
schema/proto/oshipo/
└── webhook_api/v1/              # Separate package for webhook definitions
    ├── api.proto                # Service definition
    └── api_{service}.proto      # Request/Response messages

internal/
├── usecase/
│   ├── webhook_{service}.go           # Interactor interface
│   ├── webhook_{service}_impl.go      # Implementation
│   └── input/
│       └── webhook_{service}.go       # Input DTOs
├── infrastructure/
    ├── grpc/
    │   └── internal/
    │       └── handler/
    │           └── webhook/
    │               ├── handler.go           # Handler struct
    │               └── {service}.go         # RPC implementations
    └── http/
        └── internal/
            └── handler/
                └── webhook_{service}.go     # Custom HTTP handler
```

## Implementation Steps

### Step 1: Proto Definitions

Create webhook-specific proto package in `schema/proto/oshipo/webhook_api/v1/`.

#### Service Definition (api.proto)

```protobuf
syntax = "proto3";

package oshipo.webhook_api.v1;

import "google/api/annotations.proto";
import "oshipo/webhook_api/v1/api_{service}.proto";

service WebhookService {
  rpc Receive{Service}Webhook(Receive{Service}WebhookRequest) returns (Receive{Service}WebhookResponse) {
    option (google.api.http) = {
      post: "/v1/webhook/{service}-internal"
      body: "*"
    };
  }
}
```

**Naming Convention:**
- Service: `WebhookService` (shared across all webhooks)
- RPC: `Receive{Service}Webhook` (e.g., `ReceiveAppDriverWebhook`)
- Internal path: `/v1/webhook/{service}-internal` (for gRPC, not exposed externally)

#### Request/Response Messages (api_{service}.proto)

```protobuf
syntax = "proto3";

package oshipo.webhook_api.v1;

import "protoc-gen-openapiv2/options/annotations.proto";

message Receive{Service}WebhookRequest {
  string field1 = 1;
  int64 field2 = 2;
  string field3 = 3;
  // Add all webhook parameters

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: {
      required: ["field1", "field2", "field3"]
    }
  };
}

message Receive{Service}WebhookResponse {
  bool success = 1;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: {
      required: ["success"]
    }
  };
}
```

**Key Points:**
- Request message contains all webhook parameters
- Response typically has `success` boolean
- Use OpenAPI annotations for required fields

#### Generate Proto

```bash
make generate.buf
```

### Step 2: Usecase Layer

If not already implemented, create webhook interactor.

#### Interface (webhook_{service}.go)

```go
package usecase

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_usecase
type Webhook{Service}Interactor interface {
    Receive(
        ctx context.Context,
        param *input.WebhookReceive{Service},
    ) error
}
```

#### Input DTO (input/webhook_{service}.go)

```go
package input

type WebhookReceive{Service} struct {
    Field1      string    `validate:"required"`
    Field2      int64     `validate:"required"`
    Field3      string    `validate:"required"`
    RequestTime time.Time `validate:"required"`
}

func NewWebhookReceive{Service}(
    field1 string,
    field2 int64,
    field3 string,
    requestTime time.Time,
) *WebhookReceive{Service} {
    return &WebhookReceive{Service}{
        Field1:      field1,
        Field2:      field2,
        Field3:      field3,
        RequestTime: requestTime,
    }
}

func (p *WebhookReceive{Service}) Validate() error {
    if err := validation.Validate(p); err != nil {
        return errors.RequestInvalidArgumentErr.Wrap(err)
    }
    return nil
}
```

### Step 3: gRPC Handler

#### Handler Struct (webhook/handler.go)

```go
package webhook

import (
    webhook_apiv1 "github.com/Mirrativ-com/livedx-oshipo-server/internal/infrastructure/grpc/pb/oshipo/webhook_api/v1"
    "github.com/Mirrativ-com/livedx-oshipo-server/internal/usecase"
)

type WebhookHandler struct {
    webhook{Service}Interactor usecase.Webhook{Service}Interactor
    // Add other webhook interactors as needed
}

func NewWebhookHandler(
    webhook{Service}Interactor usecase.Webhook{Service}Interactor,
) webhook_apiv1.WebhookServiceServer {
    return &WebhookHandler{
        webhook{Service}Interactor: webhook{Service}Interactor,
    }
}
```

#### RPC Implementation (webhook/{service}.go)

```go
package webhook

import (
    "context"
    "time"

    webhook_apiv1 "github.com/Mirrativ-com/livedx-oshipo-server/internal/infrastructure/grpc/pb/oshipo/webhook_api/v1"
    "github.com/Mirrativ-com/livedx-oshipo-server/internal/pkg/logger"
    "github.com/Mirrativ-com/livedx-oshipo-server/internal/pkg/logger/logger_field"
    "github.com/Mirrativ-com/livedx-oshipo-server/internal/usecase/input"
    "go.uber.org/zap"
)

func (h *WebhookHandler) Receive{Service}Webhook(
    ctx context.Context,
    req *webhook_apiv1.Receive{Service}WebhookRequest,
) (*webhook_apiv1.Receive{Service}WebhookResponse, error) {
    // 1. Parse/validate special fields (timestamps, etc.)
    specialField, err := parseSpecialField(req.Field3)
    if err != nil {
        logger.L(ctx).Warn("Invalid field format",
            zap.String("field3", req.Field3),
            logger_field.Error(err))
        return &webhook_apiv1.Receive{Service}WebhookResponse{
            Success: false,
        }, nil
    }

    // 2. Create input parameter
    param := input.NewWebhookReceive{Service}(
        req.Field1,
        req.Field2,
        specialField,
        time.Now(),
    )

    // 3. Call interactor
    err = h.webhook{Service}Interactor.Receive(ctx, param)
    if err != nil {
        logger.L(ctx).Error("Webhook processing failed",
            zap.String("identifier", req.Field1),
            logger_field.Error(err))
        return &webhook_apiv1.Receive{Service}WebhookResponse{
            Success: false,
        }, nil
    }

    logger.L(ctx).Info("Webhook processed successfully",
        zap.String("identifier", req.Field1))

    return &webhook_apiv1.Receive{Service}WebhookResponse{
        Success: true,
    }, nil
}
```

**Error Handling Pattern:**
- Parse/validation errors → `Success: false`, no gRPC error (return nil)
- Business logic errors → `Success: false`, no gRPC error (return nil)
- Always return `nil` error to prevent 500 responses
- Log all failures with structured context

### Step 4: Custom HTTP Handler

Create custom HTTP handler when:
- Non-standard HTTP method (e.g., GET with query params)
- Non-JSON response format (e.g., plain text "1"/"0")
- Special header requirements

#### HTTP Handler (http/internal/handler/webhook_{service}.go)

```go
package handler

import (
    "net/http"
    "strconv"

    webhook_apiv1 "github.com/Mirrativ-com/livedx-oshipo-server/internal/infrastructure/grpc/pb/oshipo/webhook_api/v1"
    "github.com/Mirrativ-com/livedx-oshipo-server/internal/pkg/logger"
    "go.uber.org/zap"
    "google.golang.org/grpc"
)

type Webhook{Service}Handler struct {
    grpcConn *grpc.ClientConn
}

func NewWebhook{Service}Handler(
    grpcConn *grpc.ClientConn,
) *Webhook{Service}Handler {
    return &Webhook{Service}Handler{
        grpcConn: grpcConn,
    }
}

func (h *Webhook{Service}Handler) Receive(
    w http.ResponseWriter,
    r *http.Request,
    pathParams map[string]string,
) {
    ctx := r.Context()

    // 1. Parse request (query params, body, headers, etc.)
    query := r.URL.Query()

    field2, _ := strconv.ParseInt(query.Get("field2"), 10, 64)

    // 2. Build gRPC request
    req := &webhook_apiv1.Receive{Service}WebhookRequest{
        Field1: query.Get("field1"),
        Field2: field2,
        Field3: query.Get("field3"),
    }

    // 3. Call internal gRPC endpoint
    client := webhook_apiv1.NewWebhookServiceClient(h.grpcConn)
    resp, err := client.Receive{Service}Webhook(ctx, req)

    // 4. Return response based on external service requirements
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)

    if err != nil {
        logger.L(ctx).Error("Webhook processing failed",
            zap.String("field1", req.Field1),
            zap.Error(err))
        _, _ = w.Write([]byte("0"))
        return
    }

    if resp.Success {
        logger.L(ctx).Info("Webhook processed successfully",
            zap.String("field1", req.Field1))
        _, _ = w.Write([]byte("1"))
    } else {
        logger.L(ctx).Warn("Webhook processing returned unsuccessful",
            zap.String("field1", req.Field1))
        _, _ = w.Write([]byte("0"))
    }
}
```

**HTTP Handler Patterns:**

| Request Type | Parse From | Example |
|--------------|------------|---------|
| GET with query params | `r.URL.Query()` | App Driver |
| POST with JSON body | `json.NewDecoder(r.Body).Decode()` | Most APIs |
| POST with form data | `r.FormValue()` | Legacy APIs |
| Custom headers | `r.Header.Get()` | Signature verification |

| Response Type | Content-Type | Write Method | Example |
|---------------|--------------|--------------|---------|
| Plain text | `text/plain` | `w.Write([]byte("1"))` | App Driver |
| JSON | `application/json` | `json.NewEncoder(w).Encode()` | Most APIs |
| XML | `application/xml` | `xml.NewEncoder(w).Encode()` | SOAP APIs |

### Step 5: Dependency Injection

#### Add to Dependency Struct

Location: `internal/infrastructure/dependency/dependency.go`

```go
type Dependency struct {
    // ... existing fields

    // Webhook interactors
    Webhook{Service}Interactor usecase.Webhook{Service}Interactor
}

func (d *Dependency) Inject(ctx context.Context, e *environment.Environment) {
    // ... existing initialization

    // Repositories (if needed)
    {entity}Repository := database_repository.New{Entity}()

    // Webhook interactors
    d.Webhook{Service}Interactor = usecase.NewWebhook{Service}Interactor(
        transactable,
        userRepository,
        {entity}Repository,
        // Add required dependencies
    )
}
```

### Step 6: Server Registration

#### gRPC Server Registration

Location: `internal/infrastructure/grpc/run.go`

```go
import (
    // ... existing imports
    webhook_apiv1 "github.com/Mirrativ-com/livedx-oshipo-server/internal/infrastructure/grpc/pb/oshipo/webhook_api/v1"
)

func NewServer(...) *grpc.Server {
    // ... existing code

    webhook_apiv1.RegisterWebhookServiceServer(
        server,
        webhook.NewWebhookHandler(
            dependency.Webhook{Service}Interactor,
        ),
    )

    return server
}
```

#### HTTP Server Registration

Location: `internal/infrastructure/http/run.go`

```go
// Webhook handler
webhook{Service}Handler := handler.NewWebhook{Service}Handler(conn)
if err = grpcGateway.HandlePath(
    http.MethodGet,  // or http.MethodPost
    "/v1/webhook/{service}",
    webhook{Service}Handler.Receive,
); err != nil {
    panic(err)
}
```

**Path Conventions:**
- External webhook path: `/v1/webhook/{service}` (e.g., `/v1/webhook/app-driver`)
- Internal gRPC path: `/v1/webhook/{service}-internal` (not exposed externally)

## Common Patterns

### Pattern 1: GET Request with Query Parameters

External service sends GET request with query params (e.g., App Driver).

```go
func (h *Handler) Receive(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
    query := r.URL.Query()

    // Parse integers
    intField, _ := strconv.ParseInt(query.Get("int_field"), 10, 64)

    // Parse booleans
    boolField, _ := strconv.ParseBool(query.Get("bool_field"))

    req := &webhook_apiv1.Request{
        StringField: query.Get("string_field"),
        IntField:    intField,
        BoolField:   boolField,
    }
    // ...
}
```

### Pattern 2: POST Request with JSON Body

Standard webhook pattern.

```go
func (h *Handler) Receive(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
    var payload struct {
        Field1 string `json:"field1"`
        Field2 int64  `json:"field2"`
    }

    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    req := &webhook_apiv1.Request{
        Field1: payload.Field1,
        Field2: payload.Field2,
    }
    // ...
}
```

### Pattern 3: Signature Verification

Verify webhook authenticity using HMAC signature.

```go
func (h *Handler) Receive(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
    signature := r.Header.Get("X-Webhook-Signature")

    // Read body for signature verification
    body, err := io.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    r.Body = io.NopCloser(bytes.NewBuffer(body))  // Restore body

    // Verify signature
    if !verifySignature(body, signature, h.webhookSecret) {
        logger.L(ctx).Warn("Invalid webhook signature")
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    // Continue processing...
}

func verifySignature(payload []byte, signature string, secret string) bool {
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write(payload)
    expected := hex.EncodeToString(mac.Sum(nil))
    return hmac.Equal([]byte(signature), []byte(expected))
}
```

### Pattern 4: Timestamp Parsing

Parse various timestamp formats.

```go
// ISO 8601 format: "2024-01-15T12:00:00"
acceptedTime, err := time.Parse("2006-01-02T15:04:05", req.AcceptedTime)

// RFC3339 format: "2024-01-15T12:00:00Z"
acceptedTime, err := time.Parse(time.RFC3339, req.AcceptedTime)

// Unix timestamp
acceptedTime := time.Unix(req.Timestamp, 0)
```

### Pattern 5: Idempotency Check

Prevent duplicate processing using unique identifiers.

```go
func (i *interactor) Receive(ctx context.Context, param *input.WebhookReceive) error {
    return i.transactable.RWTx(ctx, func(ctx context.Context) error {
        // Check if already processed
        existing, err := i.repository.Get(ctx, repository.GetQuery{
            UniqueID: null.StringFrom(param.UniqueID),
            BaseGetOptions: repository.BaseGetOptions{OrFail: false},
        })
        if err != nil {
            return err
        }
        if existing != nil {
            // Already processed - idempotent return
            logger.L(ctx).Info("Duplicate webhook request (idempotent)",
                zap.String("unique_id", param.UniqueID))
            return nil
        }

        // Process webhook...
    })
}
```

## Response Patterns

### Success/Failure Boolean Response

```go
return &webhook_apiv1.Response{
    Success: true,  // or false
}, nil
```

### HTTP Status Codes

| Status | Use Case |
|--------|----------|
| 200 OK | Success (most webhooks) |
| 400 Bad Request | Invalid request format |
| 401 Unauthorized | Invalid signature |
| 404 Not Found | Resource not found |
| 500 Internal Server Error | Server error (avoid if possible) |

**Important**: Most webhook services expect 2xx for success. Return 200 with `Success: false` instead of error status codes when possible.

### Plain Text Response

```go
w.Header().Set("Content-Type", "text/plain; charset=utf-8")
w.WriteHeader(http.StatusOK)
_, _ = w.Write([]byte("1"))  // Success
// or
_, _ = w.Write([]byte("0"))  // Failure
```

### JSON Response

```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(map[string]interface{}{
    "success": true,
    "message": "Webhook processed successfully",
})
```

## Security Considerations

### 1. Signature Verification

Always verify webhook signatures to prevent spoofing:

```go
// Store secret in environment variable
type Environment struct {
    Webhook{Service}Secret string `env:"WEBHOOK_{SERVICE}_SECRET,required"`
}

// Verify in HTTP handler
if !verifySignature(body, signature, h.webhookSecret) {
    w.WriteHeader(http.StatusUnauthorized)
    return
}
```

### 2. IP Whitelisting

**Recommended: Use Cloud Armor (GCP) or equivalent infrastructure layer**

For services with known IP ranges, implement IP whitelisting at the infrastructure layer using Cloud Armor. This provides better performance and security than application-level checks.

#### Terraform Implementation (GCP Cloud Armor)

**Step 1: Define IP whitelist in environment-specific locals**

```terraform
# deployments/environment/development/locals.tf
locals {
  webhook_app_driver_allowed_ips = ["27.110.48.28/32"]
  # Add more webhook services as needed
  # webhook_stripe_allowed_ips = ["54.187.174.169/32", "54.187.205.235/32"]
}
```

**Step 2: Create Cloud Armor security policy**

```terraform
# deployments/base/cloud_armor.tf
resource "google_compute_security_policy" "api_backend_waf" {
  name        = "api-backend-waf"
  description = "Protects the API backend with webhook-specific IP whitelists."

  # App Driver webhook: Deny non-whitelisted IPs
  rule {
    priority    = 1000
    action      = "deny(403)"
    description = "Block non-whitelisted IPs from accessing App Driver webhook."
    preview     = false

    match {
      expr {
        expression = "request.path == '/v1/webhook/app-driver' || request.path.startsWith('/v1/webhook/app-driver?')"
      }
    }
  }

  # App Driver webhook: Allow whitelisted IPs
  dynamic "rule" {
    for_each = length(local.webhook_app_driver_allowed_ips) > 0 ? [1] : []
    content {
      priority    = 900
      action      = "allow"
      description = "Allow App Driver whitelisted IPs."
      preview     = false

      match {
        expr {
          expression = "request.path == '/v1/webhook/app-driver' || request.path.startsWith('/v1/webhook/app-driver?')"
        }
        config {
          src_ip_ranges = local.webhook_app_driver_allowed_ips
        }
      }
    }
  }

  # Allow all other requests (including other webhook services)
  rule {
    priority    = 2000
    action      = "allow"
    description = "Allow all other requests."
    preview     = false

    match {
      versioned_expr = "SRC_IPS_V1"
      config { src_ip_ranges = ["*"] }
    }
  }
}
```

**Step 3: Attach to backend service**

```terraform
# deployments/base/compute_backend.tf
resource "google_compute_backend_service" "cloud_run_services" {
  # ...
  security_policy = (
    each.key == "api" && local.enable_api_backend_waf
  ) ? google_compute_security_policy.api_backend_waf.id : null
}
```

#### Application-Level Implementation (Fallback)

If infrastructure-level IP filtering is not available, implement in Go:

```go
var allowedIPs = []string{
    "203.0.113.0/24",
    "198.51.100.0/24",
}

func isAllowedIP(remoteAddr string) bool {
    ip := net.ParseIP(strings.Split(remoteAddr, ":")[0])
    for _, cidr := range allowedIPs {
        _, network, _ := net.ParseCIDR(cidr)
        if network.Contains(ip) {
            return true
        }
    }
    return false
}

// In HTTP handler
func (h *Handler) Receive(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
    if !isAllowedIP(r.RemoteAddr) {
        w.WriteHeader(http.StatusForbidden)
        return
    }
    // Continue processing...
}
```

**Note**: Infrastructure-level IP filtering is preferred because:
- Better performance (requests blocked before reaching application)
- Centralized management (environment-specific configuration)
- Enhanced security (protection against DDoS)
- Lower resource consumption

### 3. Rate Limiting

Prevent abuse with rate limiting:

```go
// Use existing rate limiter or implement webhook-specific limits
if !h.rateLimiter.Allow(ctx, req.Identifier) {
    w.WriteHeader(http.StatusTooManyRequests)
    return
}
```

## Testing

### Unit Test Pattern

```go
func TestWebhookHandler_Receive(t *testing.T) {
    tests := map[string]func(t *testing.T){
        "success": func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockInteractor := mock_usecase.NewMockWebhook{Service}Interactor(ctrl)
            mockInteractor.EXPECT().
                Receive(gomock.Any(), gomock.Any()).
                Return(nil)

            handler := webhook.NewWebhookHandler(mockInteractor)

            req := &webhook_apiv1.Receive{Service}WebhookRequest{
                Field1: "test",
                Field2: 123,
            }

            resp, err := handler.Receive{Service}Webhook(context.Background(), req)

            assert.NoError(t, err)
            assert.True(t, resp.Success)
        },
        "interactor error returns success false": func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockInteractor := mock_usecase.NewMockWebhook{Service}Interactor(ctrl)
            mockInteractor.EXPECT().
                Receive(gomock.Any(), gomock.Any()).
                Return(errors.InternalErr.New())

            handler := webhook.NewWebhookHandler(mockInteractor)

            req := &webhook_apiv1.Receive{Service}WebhookRequest{
                Field1: "test",
                Field2: 123,
            }

            resp, err := handler.Receive{Service}Webhook(context.Background(), req)

            assert.NoError(t, err)  // No gRPC error
            assert.False(t, resp.Success)
        },
    }

    for name, tc := range tests {
        t.Run(name, tc)
    }
}
```

### Integration Test Pattern

```go
func TestWebhook_Integration(t *testing.T) {
    // Start test server
    server := httptest.NewServer(http.HandlerFunc(handler.Receive))
    defer server.Close()

    // Send webhook request
    resp, err := http.Get(server.URL + "?field1=test&field2=123")
    require.NoError(t, err)
    defer resp.Body.Close()

    // Verify response
    body, _ := io.ReadAll(resp.Body)
    assert.Equal(t, "1", string(body))
}
```

## Best Practices

1. **Always Route Through gRPC**: HTTP handler → internal gRPC → usecase
   - Never call usecase directly from HTTP layer
   - Maintains architectural consistency

2. **Return Success/Failure, Not Errors**: Return `Success: false` instead of gRPC errors
   - External services expect 2xx responses
   - Errors should be logged, not returned

3. **Log All Webhook Requests**: Include structured context
   - Identifier/unique ID
   - Processing result
   - Error details if failed

4. **Implement Idempotency**: Use unique identifiers to prevent duplicate processing
   - Database unique constraint on webhook ID
   - Check before processing

5. **Validate Early**: Parse and validate in gRPC handler
   - Timestamp formats
   - Required fields
   - Value ranges

6. **Secure by Default**: Implement signature verification
   - Store secrets in environment variables
   - Use HMAC-SHA256 or similar
   - Reject invalid signatures

7. **Separate Proto Package**: Use `webhook_api/v1` for webhook definitions
   - Keeps webhook logic isolated
   - Easier to manage multiple webhook services

8. **Handle Async Processing**: For long-running operations
   - Return success immediately
   - Process in background job
   - Use job queue pattern (see job-system.md)

## Common Pitfalls

### Pitfall 1: Calling Usecase from HTTP Layer

```go
// Bad - HTTP handler calls usecase directly
func (h *Handler) Receive(...) {
    err := h.interactor.Receive(...)  // Don't do this!
}

// Good - HTTP handler calls internal gRPC
func (h *Handler) Receive(...) {
    client := webhook_apiv1.NewWebhookServiceClient(h.grpcConn)
    resp, err := client.ReceiveWebhook(ctx, req)
}
```

### Pitfall 2: Returning gRPC Errors

```go
// Bad - Returns gRPC error
func (h *Handler) ReceiveWebhook(...) (*Response, error) {
    if err := h.interactor.Receive(...); err != nil {
        return nil, err  // External service sees 500 error
    }
}

// Good - Returns Success: false
func (h *Handler) ReceiveWebhook(...) (*Response, error) {
    if err := h.interactor.Receive(...); err != nil {
        return &Response{Success: false}, nil  // External service sees 200 OK
    }
}
```

### Pitfall 3: Missing Idempotency Check

```go
// Bad - Processes duplicate requests
func (i *interactor) Receive(...) error {
    // Directly create record
    return i.repository.Create(...)  // May fail on duplicate
}

// Good - Checks for duplicates first
func (i *interactor) Receive(...) error {
    existing, _ := i.repository.Get(ctx, ...)
    if existing != nil {
        return nil  // Already processed
    }
    return i.repository.Create(...)
}
```

### Pitfall 4: Using Public API Proto Package

```go
// Bad - Webhook in public_api package
package oshipo.public_api.v1;
service PublicV1Service {
    rpc ReceiveWebhook(...) returns (...);
}

// Good - Separate webhook_api package
package oshipo.webhook_api.v1;
service WebhookService {
    rpc ReceiveWebhook(...) returns (...);
}
```

## Multiple Webhook Services

When adding multiple webhooks, reuse the same `WebhookService` and handler:

```go
// Proto definition
service WebhookService {
    rpc ReceiveAppDriverWebhook(...) returns (...);
    rpc ReceiveStripeWebhook(...) returns (...);
    rpc ReceiveTwilioWebhook(...) returns (...);
}

// Handler struct
type WebhookHandler struct {
    webhookAppDriverInteractor usecase.WebhookAppDriverInteractor
    webhookStripeInteractor    usecase.WebhookStripeInteractor
    webhookTwilioInteractor    usecase.WebhookTwilioInteractor
}

// Each webhook gets its own RPC method
func (h *WebhookHandler) ReceiveAppDriverWebhook(...) {...}
func (h *WebhookHandler) ReceiveStripeWebhook(...) {...}
func (h *WebhookHandler) ReceiveTwilioWebhook(...) {...}
```

## References

- See `usecase-interactor.md` for usecase layer patterns
- See `grpc-handler.md` for gRPC handler patterns
- See `job-system.md` for async processing patterns
- See `domain-errors.md` for error handling
