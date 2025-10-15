# Integrations Layer

## Quick Start

The integrations layer provides access to external services with built-in observability, logging, and best practices.

## Available Services

### BillingService (`iris_billing`)
Interface for IRIS billing system operations.

```go
type BillingService interface {
    GetInvoice(ctx context.Context, prescriptionID string) (*Invoice, error)
}
```

### PharmacyService (`iris_pharmacy`)
Interface for IRIS pharmacy system operations.

```go
type PharmacyService interface {
    GetPrescription(ctx context.Context, prescriptionID string) (*Prescription, error)
}
```

## Quick Examples

### Initialize Integration Layer

```go
import (
    "pharmacy-modernization-project-model/internal/integrations"
    "pharmacy-modernization-project-model/internal/platform/httpclient"
)

// Create centralized HTTP client
httpClient := httpclient.NewClient(
    httpclient.Config{
        Timeout:     30 * time.Second,
        ServiceName: "integrations",
    },
    logger,
)

// Initialize all integrations
services := integrations.New(integrations.Dependencies{
    Config:     config,
    Logger:     logger,
    HTTPClient: httpClient,
})

// Use services
invoice, err := services.BillingService.GetInvoice(ctx, "prescription-123")
prescription, err := services.PharmacyService.GetPrescription(ctx, "prescription-123")
```

### Use Mock Services (Testing/Development)

Update your config:
```yaml
external:
  billing:
    use_mock: true
  pharmacy:
    use_mock: true
```

Or programmatically:
```go
mockBilling := iris_billing.NewMockService(seedData, logger)
mockPharmacy := iris_pharmacy.NewMockService(seedData, logger)
```

## Observability

All API calls automatically log:
- Request details (method, URL, headers)
- Response details (status, size, duration)
- Errors with full context

Check your logs for entries like:
```
INFO  http request completed  service=iris_billing method=GET status_code=200 duration=245ms
```

## Documentation

- [Full Architecture Guide](./INTEGRATION_ARCHITECTURE.md) - Comprehensive architecture documentation
- [HTTP Client Package](../../internal/platform/httpclient/) - Centralized HTTP client implementation
- [Adding New Integrations](./INTEGRATION_ARCHITECTURE.md#adding-a-new-integration) - Step-by-step guide

## Best Practices

1. ✅ Always use context with appropriate timeouts
2. ✅ Handle errors appropriately with proper logging
3. ✅ Use mock services for testing
4. ✅ Configure timeouts per service based on SLAs
5. ✅ Use structured logging with zap fields

## Configuration

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

## Support

For questions or issues with the integration layer:
1. Check the [Architecture Guide](./INTEGRATION_ARCHITECTURE.md)
2. Review existing integration implementations as examples
3. Ensure proper configuration and logging setup

