# Header Provider Examples - Practical Guide

## Overview

This guide shows **real working examples** from the codebase demonstrating how to add headers globally and per-endpoint.

---

## ✅ Example 1: Global Headers (All APIs)

### **Location:** `internal/integrations/integration_wire.go`

### **Code:**
```go
// Create global header provider for all API requests
// These headers will be added to ALL requests across all integrations
globalHeaderProvider := httpclient.NewStaticHeaderProvider(map[string]string{
    "X-IRIS-User-ID": "xyz", // ✅ Added to ALL API calls
    // Add more global headers here as needed:
    // "X-Client-Version": "1.0.0",
    // "X-Request-Source": "rxintake-app",
})

// Create shared HTTP client for all external API integrations
sharedHTTPClient := httpclient.NewClient(
    httpclient.Config{
        Timeout:        30 * time.Second,
        MaxIdleConns:   100,
        ServiceName:    "external_apis",
        HeaderProvider: globalHeaderProvider, // ✅ Global headers for ALL requests
    },
    logger,
    metricsInterceptor,
)
```

### **Result:**

**ALL API calls** (billing, pharmacy, etc.) will now include:
```
X-IRIS-User-ID: xyz
```

### **Example API Calls:**

```go
// Call 1: Get Invoice
invoice, err := billingClient.GetInvoice(ctx, "RX-123")

// HTTP Request includes:
// X-IRIS-User-ID: xyz               ← Global header
// X-IRIS-Env-Name: IRIS_stage       ← Endpoint-specific (see Example 2)
// Content-Type: application/json
// Accept: application/json

// Call 2: Get Prescription
prescription, err := pharmacyClient.GetPrescription(ctx, "RX-456")

// HTTP Request includes:
// X-IRIS-User-ID: xyz               ← Global header (automatically added)
// Content-Type: application/json
// Accept: application/json
```

---

## ✅ Example 2: Endpoint-Specific Headers (Some APIs)

### **Location:** `internal/integrations/iris_billing/http_client.go`

### **Code:**
```go
// GetInvoice retrieves an invoice for a given prescription ID
func (c *HTTPClient) GetInvoice(ctx context.Context, prescriptionID string) (*InvoiceResponse, error) {
    url := replacePathParams(c.endpoints.GetInvoiceEndpoint(), map[string]string{
        "prescriptionID": prescriptionID,
    })
    
    c.logger.Debug("fetching invoice",
        zap.String("prescription_id", prescriptionID),
        zap.String("url", url),
    )
    
    // ✅ Example: Add endpoint-specific header for THIS endpoint only
    resp, err := c.client.Get(ctx, url, map[string]string{
        "Content-Type":    "application/json",
        "Accept":          "application/json",
        "X-IRIS-Env-Name": "IRIS_stage", // ✅ Only for GetInvoice endpoint
    })
    
    // ... rest of implementation
}
```

### **Result:**

**Only GetInvoice** API calls will include:
```
X-IRIS-Env-Name: IRIS_stage
```

Other endpoints (CreateInvoice, AcknowledgeInvoice, etc.) won't have this header unless you add it to those specific methods.

---

## 📊 **Complete Example: Both Headers Combined**

### **GetInvoice Request:**
```
GET https://api.iris.example.com/billing/v1/invoices/RX-123

Headers:
├── X-IRIS-User-ID: xyz              ← From global HeaderProvider
├── X-IRIS-Env-Name: IRIS_stage      ← From endpoint-specific code
├── Content-Type: application/json   ← From endpoint-specific code
└── Accept: application/json         ← From endpoint-specific code
```

### **CreateInvoice Request:**
```
POST https://api.iris.example.com/billing/v1/invoices

Headers:
├── X-IRIS-User-ID: xyz              ← From global HeaderProvider
├── Content-Type: application/json   ← From endpoint code (uses GetJSON)
└── Accept: application/json         ← From endpoint code (uses GetJSON)

Note: NO X-IRIS-Env-Name (not added in CreateInvoice method)
```

### **GetPrescription Request:**
```
GET https://api.iris.example.com/pharmacy/v1/prescriptions/RX-456

Headers:
├── X-IRIS-User-ID: xyz              ← From global HeaderProvider
├── Content-Type: application/json   ← From GetJSON method
└── Accept: application/json         ← From GetJSON method

Note: NO X-IRIS-Env-Name (different service)
```

---

## 🎯 **Patterns**

### **Pattern 1: Add Header to ALL API Calls**

Update `integration_wire.go`:
```go
globalHeaderProvider := httpclient.NewStaticHeaderProvider(map[string]string{
    "X-IRIS-User-ID": "xyz",
    "X-New-Header": "value", // ✅ Add here for ALL APIs
})
```

### **Pattern 2: Add Header to ONE Endpoint**

Update specific method in `http_client.go`:
```go
func (c *HTTPClient) GetInvoice(ctx, prescriptionID) (*InvoiceResponse, error) {
    // ...
    resp, err := c.client.Get(ctx, url, map[string]string{
        "X-IRIS-Env-Name": "IRIS_stage", // ✅ Only this endpoint
        "X-Another-Header": "value",      // ✅ Add more here
    })
    // ...
}
```

### **Pattern 3: Add Header to SOME Endpoints (Same Service)**

Add to multiple methods in the same service:
```go
// GetInvoice - has the header
func (c *HTTPClient) GetInvoice(...) {
    resp, err := c.client.Get(ctx, url, map[string]string{
        "X-IRIS-Env-Name": "IRIS_stage", // ✅
    })
}

// CreateInvoice - also has the header
func (c *HTTPClient) CreateInvoice(...) {
    resp, err := c.client.Post(ctx, url, body, map[string]string{
        "X-IRIS-Env-Name": "IRIS_stage", // ✅
    })
}

// AcknowledgeInvoice - does NOT have the header
func (c *HTTPClient) AcknowledgeInvoice(...) {
    // Uses PostJSON which sets standard headers only
    err := c.client.PostJSON(ctx, url, req, &response)
}
```

---

## 🔄 **Header Merge Priority**

When headers are set at multiple levels:

```
1. Global HeaderProvider        ← Applied first
2. Endpoint-specific headers    ← Can override global

Example:
Global: {"Content-Type": "application/xml"}
Endpoint: {"Content-Type": "application/json"}
Result: Content-Type: application/json (endpoint wins)
```

---

## 💡 **Dynamic Headers (Advanced)**

For headers that change per request:

```go
// In integration_wire.go
dynamicHeaderProvider := httpclient.HeaderProviderFunc(func(ctx context.Context) (map[string]string, error) {
    // Get user ID from context
    userID := ctx.Value("user_id")
    if userID == nil {
        userID = "anonymous"
    }
    
    return map[string]string{
        "X-IRIS-User-ID": userID.(string),      // ✅ Dynamic per request
        "X-Request-Time": time.Now().String(),  // ✅ Timestamp per request
    }, nil
})

sharedHTTPClient := httpclient.NewClient(
    httpclient.Config{
        HeaderProvider: dynamicHeaderProvider,
    },
    logger,
)
```

---

## 📝 **Real-World Use Cases**

### **Use Case 1: User Tracking**

```go
// Global header with user ID from config/context
globalHeaderProvider := httpclient.NewStaticHeaderProvider(map[string]string{
    "X-IRIS-User-ID": getCurrentUserID(),
    "X-Client-ID":    "rxintake-app",
})
```

### **Use Case 2: Environment Identification**

```go
// Different endpoints need different environment headers
func (c *HTTPClient) GetInvoice(...) {
    resp, err := c.client.Get(ctx, url, map[string]string{
        "X-IRIS-Env-Name": getEnvironmentName(), // prod, staging, dev
    })
}
```

### **Use Case 3: API Versioning**

```go
// Global API version for all requests
globalHeaderProvider := httpclient.NewStaticHeaderProvider(map[string]string{
    "X-API-Version": "v1",
    "Accept":        "application/json",
})
```

### **Use Case 4: Correlation IDs**

```go
// Add correlation ID from request context
func (c *HTTPClient) GetInvoice(ctx, prescriptionID) (*InvoiceResponse, error) {
    correlationID := getCorrelationIDFromContext(ctx)
    
    resp, err := c.client.Get(ctx, url, map[string]string{
        "X-Correlation-ID": correlationID,  // ✅ Request tracing
    })
    // ...
}
```

---

## 🎉 **Summary**

### **Working Examples in Your Codebase:**

1. ✅ **Global Header** (`integration_wire.go`):
   ```go
   "X-IRIS-User-ID": "xyz"  // Added to ALL requests
   ```

2. ✅ **Endpoint-Specific Header** (`iris_billing/http_client.go`):
   ```go
   "X-IRIS-Env-Name": "IRIS_stage"  // Only for GetInvoice
   ```

### **To Use:**

**Global headers** (all APIs):
- Edit `integration_wire.go`
- Update `globalHeaderProvider` map

**Endpoint-specific** (some APIs):
- Edit `http_client.go` in the specific integration
- Add headers to specific method's `Get()`/`Post()` call

### **Result:**
- ✅ Clean separation of global vs endpoint-specific headers
- ✅ Easy to add/modify headers
- ✅ Clear, maintainable code
- ✅ Full flexibility

**Both patterns working in your codebase!** 🚀

