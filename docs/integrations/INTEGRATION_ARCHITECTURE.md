# Integration Layer Architecture

## Overview

The integration layer provides a centralized, well-architected approach to making external API calls with built-in observability, logging, and middleware support.

## Architecture Principles

### 1. **Centralized HTTP Client**
- Single source of truth for HTTP communication
- Located in `internal/platform/httpclient/`
- Eliminates code duplication across integrations
- Provides consistent behavior across all external API calls

### 2. **Service-Oriented Design**
- Each integration exposes a domain-specific service interface
- Example: `BillingService`, `PharmacyService` (not generic `Client`)
- Clear naming conventions that reflect business domains

### 3. **Observability by Default**
- Automatic request/response logging with structured fields
- Built-in timing metrics for every API call
- Configurable log levels based on response status
- Request/response size tracking

### 4. **Middleware/Interceptor Pattern**
- Extensible request/response interceptors
- Pre-built interceptors for common needs:
  - Authentication (`AuthInterceptor`)
  - Metrics collection (`MetricsInterceptor`)
  - Retry logic (`RetryInterceptor`)

### 5. **Mock-First Testing**
- Every service has a corresponding mock implementation
- Mock services follow the same interface
- Easy to switch between real and mock implementations via configuration

## Directory Structure

```
internal/
├── platform/
│   └── httpclient/              # Centralized HTTP client
│       ├── client.go            # Core HTTP client with observability
│       ├── interceptor.go       # Interceptor interface
│       └── interceptors/        # Built-in interceptors
│           ├── auth.go          # Authentication interceptor
│           ├── metrics.go       # Metrics collection
│           └── retry.go         # Retry logic
│
└── integrations/
    ├── integration_wire.go      # Integration layer assembly
    │
    ├── iris_billing/            # Billing service integration
    │   ├── service.go           # BillingService interface
    │   ├── http_service.go      # HTTP implementation
    │   ├── mock_service.go      # Mock implementation
    │   ├── model.go             # Domain models (Invoice)
    │   ├── config.go            # Configuration
    │   └── module.go            # Dependency injection module
    │
    └── iris_pharmacy/           # Pharmacy service integration
        ├── service.go           # PharmacyService interface
        ├── http_service.go      # HTTP implementation
        ├── mock_service.go      # Mock implementation
        ├── model.go             # Domain models (Prescription)
        ├── config.go            # Configuration
        └── module.go            # Dependency injection module
```

## Key Components

### Centralized HTTP Client (`internal/platform/httpclient/client.go`)

The `httpclient.Client` provides:

- **Automatic Logging**: Every request/response is logged with:
  - Service name
  - HTTP method
  - URL
  - Status code
  - Duration
  - Response size

- **Observability**: Built-in metrics tracking:
  ```go
  c.logger.Log(logLevel, "http request completed",
      zap.String("service", c.serviceName),
      zap.String("method", req.Method),
      zap.String("url", req.URL),
      zap.Int("status_code", httpResp.StatusCode),
      zap.Duration("duration", duration),
      zap.Int("response_size", len(body)),
  )
  ```

- **Convenience Methods**:
  - `Get(ctx, url, headers)`
  - `Post(ctx, url, body, headers)`
  - `Put(ctx, url, body, headers)`
  - `Delete(ctx, url, headers)`
  - `Do(ctx, request)` - for custom requests

### Service Interfaces

Each integration defines a clear, domain-specific interface:

```go
// BillingService - clear, business-focused name
type BillingService interface {
    GetInvoice(ctx context.Context, prescriptionID string) (*Invoice, error)
}

// PharmacyService - clear, business-focused name
type PharmacyService interface {
    GetPrescription(ctx context.Context, prescriptionID string) (*Prescription, error)
}
```

### HTTP Service Implementation

Each service implementation uses the centralized client:

```go
type HTTPService struct {
    client   *httpclient.Client  // Centralized client
    endpoint string
    logger   *zap.Logger
}

func (s *HTTPService) GetInvoice(ctx context.Context, prescriptionID string) (*Invoice, error) {
    url := s.endpoint + prescriptionID
    
    // Use centralized client - gets automatic logging and observability
    resp, err := s.client.Get(ctx, url, map[string]string{
        "Content-Type": "application/json",
    })
    
    // Handle response...
}
```

## Usage Examples

### Basic Usage

```go
// Initialize centralized HTTP client
httpClient := httpclient.NewClient(
    httpclient.Config{
        Timeout:     30 * time.Second,
        ServiceName: "iris_billing",
    },
    logger,
)

// Create billing service
billingService := iris_billing.NewHTTPService(
    iris_billing.Config{
        BaseURL: "https://api.iris.example.com",
        Path:    "/billing/v1",
    },
    httpClient,
    logger,
)

// Use the service
invoice, err := billingService.GetInvoice(ctx, "prescription-123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Invoice: %+v\n", invoice)
```

### Using Interceptors

```go
// Create interceptors
authInterceptor := interceptors.NewAuthInterceptor("Bearer", "your-token")
metricsInterceptor := interceptors.NewMetricsInterceptor(logger)

// Create client with interceptors
httpClient := httpclient.NewClient(
    httpclient.Config{
        Timeout:     30 * time.Second,
        ServiceName: "iris_billing",
    },
    logger,
    authInterceptor,
    metricsInterceptor,
)
```

### Using Mock Services

```go
// Create mock service with seed data
mockService := iris_billing.NewMockService(
    map[string]iris_billing.Invoice{
        "prescription-123": {
            ID:             "invoice-456",
            PrescriptionID: "prescription-123",
            Amount:         125.50,
            Status:         "paid",
        },
    },
    logger,
)

// Use mock service (same interface as HTTP service)
invoice, err := mockService.GetInvoice(ctx, "prescription-123")
```

## Observability Features

### Automatic Logging

Every API call logs:

**Request Initiated:**
```
INFO    http request initiated
    service=iris_billing
    method=GET
    url=https://api.iris.example.com/billing/v1/prescription-123
```

**Request Completed:**
```
INFO    http request completed
    service=iris_billing
    method=GET
    url=https://api.iris.example.com/billing/v1/prescription-123
    status_code=200
    duration=245ms
    response_size=156
```

**Error Cases:**
```
ERROR   http request failed
    service=iris_billing
    method=GET
    url=https://api.iris.example.com/billing/v1/prescription-123
    duration=5.2s
    error=context deadline exceeded
```

### Performance Monitoring

The centralized client tracks:
- **Request Duration**: Time from request start to response received
- **Response Size**: Bytes received in response body
- **Success Rate**: Status codes 2xx vs 4xx/5xx
- **Service-level Metrics**: Per-service breakdown of all metrics

### Debug Logging

Service implementations provide additional debug logging:

```go
s.logger.Debug("fetching invoice",
    zap.String("prescription_id", prescriptionID),
    zap.String("url", url),
)

s.logger.Debug("invoice retrieved successfully",
    zap.String("prescription_id", prescriptionID),
    zap.String("invoice_id", invoice.ID),
    zap.Float64("amount", invoice.Amount),
    zap.String("status", invoice.Status),
)
```

## Configuration

### YAML Configuration Example

```yaml
external:
  billing:
    base_url: "https://api.iris.example.com"
    path: "/billing/v1"
    timeout: "30s"
    use_mock: false
  
  pharmacy:
    base_url: "https://api.iris.example.com"
    path: "/pharmacy/v1"
    timeout: "30s"
    use_mock: false
```

### Environment-Based Configuration

```go
// Development - use mocks
config.External.Billing.UseMock = true

// Production - use real services with longer timeouts
config.External.Billing.UseMock = false
config.External.Billing.Timeout = "60s"
```

## Best Practices

### 1. Always Use Context
```go
// Good: Context enables timeouts and cancellation
invoice, err := billingService.GetInvoice(ctx, id)

// Bad: Don't use background context in handlers
invoice, err := billingService.GetInvoice(context.Background(), id)
```

### 2. Handle Errors Appropriately
```go
invoice, err := billingService.GetInvoice(ctx, id)
if err != nil {
    // Log the error with context
    logger.Error("failed to get invoice",
        zap.String("prescription_id", id),
        zap.Error(err),
    )
    // Return appropriate error to caller
    return nil, fmt.Errorf("billing service error: %w", err)
}
```

### 3. Use Structured Logging
```go
// Good: Structured fields
logger.Info("invoice retrieved",
    zap.String("invoice_id", invoice.ID),
    zap.Float64("amount", invoice.Amount),
)

// Bad: String concatenation
logger.Info("invoice retrieved: " + invoice.ID + " amount: " + fmt.Sprintf("%.2f", invoice.Amount))
```

### 4. Leverage Mock Services for Testing
```go
func TestBillingFlow(t *testing.T) {
    mockBilling := iris_billing.NewMockService(
        map[string]iris_billing.Invoice{
            "test-id": {ID: "invoice-1", Amount: 100.0},
        },
        zap.NewNop(),
    )
    
    // Test your business logic with predictable mock data
    result, err := yourBusinessLogic(mockBilling)
    // assertions...
}
```

## Benefits of This Architecture

### 1. **Single Source of Truth**
- One HTTP client implementation instead of duplicated code in each integration
- Consistent behavior across all API calls
- Easy to update/enhance HTTP logic in one place

### 2. **Observability**
- Built-in logging and metrics for all API calls
- Easy to identify slow or failing external services
- Request/response tracing for debugging

### 3. **Maintainability**
- Clear separation of concerns
- Easy to add new integrations following the same pattern
- Reduced code duplication

### 4. **Testability**
- Mock services for testing business logic
- No external dependencies required for tests
- Predictable test data

### 5. **Scalability**
- Easy to add new interceptors (circuit breakers, rate limiting, etc.)
- Centralized configuration management
- Service-specific timeout and retry policies

### 6. **Developer Experience**
- Clear, domain-focused interface names
- Consistent patterns across all integrations
- Easy to understand and use

## Adding a New Integration

To add a new integration, follow this pattern:

1. **Create directory**: `internal/integrations/new_service/`

2. **Define service interface** (`service.go`):
   ```go
   type NewService interface {
       SomeOperation(ctx context.Context, params) (*Result, error)
   }
   ```

3. **Create HTTP implementation** (`http_service.go`):
   ```go
   type HTTPService struct {
       client   *httpclient.Client
       endpoint string
       logger   *zap.Logger
   }
   ```

4. **Create mock implementation** (`mock_service.go`):
   ```go
   type MockService struct {
       data   map[string]Result
       logger *zap.Logger
   }
   ```

5. **Define models** (`model.go`):
   ```go
   type Result struct {
       // fields...
   }
   ```

6. **Create module** (`module.go`):
   ```go
   func Module(deps ModuleDependencies) ModuleExport {
       // initialization logic
   }
   ```

7. **Update integration wire** (`integration_wire.go`):
   ```go
   export.NewService = newservice.Module(...).NewService
   ```

## Migration Guide

If you have existing integrations using the old pattern:

1. **Rename interface**: `Client` → `[Domain]Service`
2. **Rename implementation**: `HTTPClient` → `HTTPService`
3. **Replace http.Client**: Use `httpclient.Client` instead
4. **Update imports**: Use centralized client
5. **Remove duplicated http_client.go**: Delete old HTTP client files
6. **Update callers**: Use new service interface name

## Performance Considerations

- **Connection Pooling**: Centralized client uses connection pooling (MaxIdleConns: 100)
- **Timeout Configuration**: Per-service timeout configuration
- **Context Propagation**: Always pass context for proper cancellation
- **Resource Cleanup**: Response bodies are properly closed

## Security Considerations

- **Authentication**: Use `AuthInterceptor` for API key/token management
- **TLS**: Default http.Transport uses secure TLS settings
- **Secrets**: Never log sensitive data (tokens, API keys)
- **Context**: Respect context cancellation to prevent resource leaks

