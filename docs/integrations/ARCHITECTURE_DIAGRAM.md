# Integration Layer Architecture Diagram

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Application Layer                         │
│  (domain/prescription, domain/patient, etc.)                    │
└───────────────────────┬─────────────────────────────────────────┘
                        │
                        │ Uses domain-specific interfaces
                        │
┌───────────────────────▼─────────────────────────────────────────┐
│                    Integration Layer                             │
│           (internal/integrations/integration_wire.go)            │
│                                                                   │
│  ┌─────────────────────────┐  ┌─────────────────────────┐      │
│  │   BillingService        │  │   PharmacyService       │      │
│  │   (Interface)           │  │   (Interface)           │      │
│  └───────────┬─────────────┘  └───────────┬─────────────┘      │
│              │                             │                     │
│  ┌───────────▼─────────────┐  ┌───────────▼─────────────┐      │
│  │   HTTPService           │  │   HTTPService           │      │
│  │   or                    │  │   or                    │      │
│  │   MockService           │  │   MockService           │      │
│  └───────────┬─────────────┘  └───────────┬─────────────┘      │
└──────────────┼─────────────────────────────┼────────────────────┘
               │                             │
               │ Uses centralized client     │
               │                             │
┌──────────────▼─────────────────────────────▼────────────────────┐
│              Centralized HTTP Client                             │
│           (internal/platform/httpclient/)                        │
│                                                                   │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │                    Core Client                            │  │
│  │  • Automatic logging                                      │  │
│  │  • Request/response timing                                │  │
│  │  • Structured logging                                     │  │
│  │  • Context propagation                                    │  │
│  │  • Connection pooling                                     │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                   │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │                   Interceptors                            │  │
│  │  • AuthInterceptor (add auth headers)                     │  │
│  │  • MetricsInterceptor (collect metrics)                   │  │
│  │  • RetryInterceptor (retry logic)                         │  │
│  │  • Custom interceptors...                                 │  │
│  └──────────────────────────────────────────────────────────┘  │
└───────────────────────────┬───────────────────────────────────┘
                            │
                            │ HTTP/HTTPS
                            │
┌───────────────────────────▼───────────────────────────────────┐
│                    External Services                           │
│  • IRIS Billing API                                            │
│  • IRIS Pharmacy API                                           │
│  • Future integrations...                                      │
└────────────────────────────────────────────────────────────────┘
```

## Request Flow

```
1. Application calls BillingService.GetInvoice(ctx, id)
                    │
                    ▼
2. HTTPService receives the call
                    │
                    ▼
3. HTTPService uses centralized Client.Get()
                    │
                    ▼
4. Centralized client:
   a. Logs "request initiated"
   b. Creates HTTP request
   c. Runs "before" interceptors
   d. Executes HTTP request (starts timer)
   e. Reads response body (stops timer)
   f. Runs "after" interceptors
   g. Logs "request completed" with duration
                    │
                    ▼
5. HTTPService decodes response
                    │
                    ▼
6. HTTPService logs domain-specific details
                    │
                    ▼
7. Returns Invoice to application

Logs generated:
- INFO  http request initiated       (centralized client)
- INFO  http request completed       (centralized client)
- DEBUG fetching invoice              (service)
- DEBUG invoice retrieved successfully (service)
```

## Data Flow Diagram

```
┌────────────────────────────────────────────────────────────────┐
│                    HTTP Request                                 │
│                                                                  │
│  Context ──┐                                                    │
│            │                                                    │
│  URL ──────┼──► Centralized ──► Interceptors ──► HTTP ──► External │
│            │    Client                             Transport   API  │
│  Headers ──┤       │                                    │       │
│            │       │ Logs: initiated                    │       │
│  Body ─────┘       │                                    │       │
│                    │                                    │       │
└────────────────────┼────────────────────────────────────┼───────┘
                     │                                    │
┌────────────────────▼────────────────────────────────────▼───────┐
│                    HTTP Response                                 │
│                                                                  │
│  External ──► HTTP ──► Interceptors ──► Centralized ──► Service │
│  API        Transport                    Client          │      │
│             │                               │            │      │
│             │                               │ Logs:      │      │
│             │                               │ completed  │      │
│             │                               │ + timing   │      │
│             │                               │            │      │
│             └───────── Response ────────────┴────────────┴──► Application │
│                        • Status                                 │
│                        • Body                                   │
│                        • Duration                               │
│                        • Size                                   │
└────────────────────────────────────────────────────────────────┘
```

## Component Relationships

```
┌─────────────────────────────────────────────────────────────────┐
│                     Integration Components                       │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │ iris_billing/                                             │  │
│  │                                                           │  │
│  │  • service.go ──────► BillingService (interface)         │  │
│  │  • http_service.go ─► HTTPService (implementation)       │  │
│  │  • mock_service.go ─► MockService (implementation)       │  │
│  │  • model.go ────────► Invoice (domain model)             │  │
│  │  • config.go ───────► Config (configuration)             │  │
│  │  • module.go ───────► Module() (wire/DI)                 │  │
│  │                                                           │  │
│  │  HTTPService depends on:                                 │  │
│  │    • httpclient.Client (centralized)                     │  │
│  │    • zap.Logger                                           │  │
│  │    • Config                                               │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │ iris_pharmacy/                                            │  │
│  │                                                           │  │
│  │  • service.go ──────► PharmacyService (interface)        │  │
│  │  • http_service.go ─► HTTPService (implementation)       │  │
│  │  • mock_service.go ─► MockService (implementation)       │  │
│  │  • model.go ────────► Prescription (domain model)        │  │
│  │  • config.go ───────► Config (configuration)             │  │
│  │  • module.go ───────► Module() (wire/DI)                 │  │
│  │                                                           │  │
│  │  HTTPService depends on:                                 │  │
│  │    • httpclient.Client (centralized)                     │  │
│  │    • zap.Logger                                           │  │
│  │    • Config                                               │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │ integration_wire.go                                       │  │
│  │                                                           │  │
│  │  New() ──► Initializes all integrations                  │  │
│  │            • Creates/receives HTTP client                │  │
│  │            • Calls each module's Module()                │  │
│  │            • Returns Export with all services            │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## Interceptor Chain

```
Request Flow:
┌──────────┐
│  Client  │
│  Do()    │
└────┬─────┘
     │
     ▼
┌────────────────────┐
│ Interceptor 1      │
│ Before()           │
└────┬───────────────┘
     │
     ▼
┌────────────────────┐
│ Interceptor 2      │
│ Before()           │
└────┬───────────────┘
     │
     ▼
┌────────────────────┐
│ HTTP Request       │
│ (actual call)      │
└────┬───────────────┘
     │
     ▼
┌────────────────────┐
│ Interceptor 2      │
│ After()            │
└────┬───────────────┘
     │
     ▼
┌────────────────────┐
│ Interceptor 1      │
│ After()            │
└────┬───────────────┘
     │
     ▼
┌────────────────────┐
│ Return Response    │
└────────────────────┘

Example Interceptor Chain:
┌─────────────────────┐
│ AuthInterceptor     │ ──► Adds Authorization header
│ Before()            │
└─────────────────────┘
         │
         ▼
┌─────────────────────┐
│ MetricsInterceptor  │ ──► Records request start time
│ Before()            │
└─────────────────────┘
         │
         ▼
┌─────────────────────┐
│ HTTP Call           │ ──► Actual API request
└─────────────────────┘
         │
         ▼
┌─────────────────────┐
│ MetricsInterceptor  │ ──► Calculates duration, records metrics
│ After()             │
└─────────────────────┘
         │
         ▼
┌─────────────────────┐
│ AuthInterceptor     │ ──► No-op (could validate response headers)
│ After()             │
└─────────────────────┘
```

## Module Initialization Flow

```
main() ──► app.New()
              │
              ▼
         app.wire()
              │
              ├─► Create Logger
              │
              ├─► Create Centralized HTTP Client
              │      • Config: timeout, service name
              │      • Logger
              │      • Interceptors (optional)
              │
              ├─► integrations.New()
              │      │
              │      ├─► iris_pharmacy.Module()
              │      │      • Config (BaseURL, Path)
              │      │      • Logger
              │      │      • HTTP Client
              │      │      • UseMock flag
              │      │      │
              │      │      └─► Returns PharmacyService
              │      │
              │      └─► iris_billing.Module()
              │             • Config (BaseURL, Path)
              │             • Logger
              │             • HTTP Client
              │             • UseMock flag
              │             │
              │             └─► Returns BillingService
              │
              └─► Domain Modules (prescription, patient, etc.)
                     • Receive BillingService
                     • Receive PharmacyService
                     • Use interfaces (decoupled from implementation)
```

## Benefits Visualization

### Before: Duplicated Code
```
iris_billing/          iris_pharmacy/
├── http_client.go     ├── http_client.go     ← DUPLICATE
│   • HTTP logic       │   • HTTP logic           DUPLICATE
│   • Logging          │   • Logging              DUPLICATE
│   • Error handling   │   • Error handling       DUPLICATE
│   • Timing           │   • Timing               DUPLICATE
```

### After: Centralized
```
platform/httpclient/
├── client.go          ← SINGLE SOURCE OF TRUTH
│   • HTTP logic
│   • Logging
│   • Error handling
│   • Timing
│   • Observability
       ▲      ▲
       │      │
       │      └─────────────┐
       │                    │
iris_billing/        iris_pharmacy/
└── http_service.go  └── http_service.go
    • Business logic     • Business logic
    • Domain models      • Domain models
```

## Observability Flow

```
┌──────────────────────────────────────────────────────────────────┐
│                         Observability                             │
│                                                                   │
│  Centralized Client                                              │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │ Request Start ──► Log: "request initiated"                 │ │
│  │                   • service name                            │ │
│  │                   • method                                  │ │
│  │                   • URL                                     │ │
│  └────────────────────────────────────────────────────────────┘ │
│                                                                   │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │ HTTP Call ──► Measure Duration                             │ │
│  │               • Start timer                                 │ │
│  │               • Execute request                             │ │
│  │               • Stop timer                                  │ │
│  └────────────────────────────────────────────────────────────┘ │
│                                                                   │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │ Request End ──► Log: "request completed"                   │ │
│  │                 • service name                              │ │
│  │                 • method                                    │ │
│  │                 • URL                                       │ │
│  │                 • status code                               │ │
│  │                 • duration ◄─── KEY METRIC                 │ │
│  │                 • response size                             │ │
│  └────────────────────────────────────────────────────────────┘ │
│                                                                   │
│  Service Implementation                                          │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │ Domain Logging ──► Log: "invoice retrieved"                │ │
│  │                    • prescription_id                        │ │
│  │                    • invoice_id                             │ │
│  │                    • amount                                 │ │
│  │                    • status                                 │ │
│  └────────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────┘

All logs are structured (zap) and can be:
• Filtered by service name
• Aggregated by duration
• Alerted on error patterns
• Analyzed for performance
```

## Testing Strategy

```
┌────────────────────────────────────────────────────────────┐
│                      Testing Layers                         │
│                                                             │
│  Unit Tests                                                 │
│  ┌─────────────────────────────────────────────────────┐  │
│  │ BillingService                                       │  │
│  │   • MockService with test data                      │  │
│  │   • No external dependencies                        │  │
│  │   • Fast execution                                   │  │
│  │   • Predictable results                             │  │
│  └─────────────────────────────────────────────────────┘  │
│                                                             │
│  Integration Tests                                         │
│  ┌─────────────────────────────────────────────────────┐  │
│  │ HTTPService                                          │  │
│  │   • Test server (httptest)                          │  │
│  │   • Real HTTP calls (mocked endpoint)               │  │
│  │   • Test client behavior                            │  │
│  │   • Test error handling                             │  │
│  └─────────────────────────────────────────────────────┘  │
│                                                             │
│  End-to-End Tests                                          │
│  ┌─────────────────────────────────────────────────────┐  │
│  │ Full Integration Flow                                │  │
│  │   • Real services (staging)                         │  │
│  │   • Full request lifecycle                          │  │
│  │   • Observability validation                        │  │
│  │   • Performance benchmarks                          │  │
│  └─────────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────┘
```

## Summary

This architecture provides:

1. **Single Responsibility**: Each component has one clear purpose
2. **Observability**: Built-in logging and metrics for all API calls
3. **Testability**: Easy to mock and test at any level
4. **Maintainability**: Clear structure, no duplication
5. **Extensibility**: Easy to add new integrations following the same pattern
6. **Performance**: Connection pooling, efficient resource usage
7. **Reliability**: Consistent error handling, timeout management

