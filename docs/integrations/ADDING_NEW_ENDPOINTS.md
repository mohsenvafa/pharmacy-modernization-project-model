## Adding New API Endpoints - Simple Guide

The new architecture makes it **easy** to add new API endpoints with **minimal boilerplate**. Here's how:

## 🎯 What You Need to Do

To add a new endpoint, you only need to define:
1. **Endpoint configuration** (method, path, description)
2. **Request/Response models** (if needed)
3. **Interface method** (type signature)
4. **Implementation** (1-2 lines of code using base API service)

That's it! No repetitive HTTP code, no manual JSON parsing, no URL building.

---

## 📝 Example: Adding New Endpoints to Billing API

### Current State (after refactoring):
```
iris_billing/
├── client.go          # BillingClient interface
├── endpoints.go       # All endpoint definitions
├── models.go          # All models in one place
├── http_client.go     # HTTP implementation (clean, minimal)
├── mock_client.go     # Mock implementation
├── config.go          # Configuration
└── module.go          # Initialization
```

---

## Step 1: Define the Endpoint (`endpoints.go`)

Just add the endpoint configuration:

```go
// endpoints.go
var (
    // Existing endpoints...
    GetInvoiceEndpoint = httpclient.EndpointConfig{
        Method:      "GET",
        Path:        "invoices/{prescriptionID}",
        Description: "get invoice",
    }
    
    // ✅ NEW ENDPOINT - Just add this!
    UpdateInvoiceEndpoint = httpclient.EndpointConfig{
        Method:      "PUT",
        Path:        "invoices/{invoiceID}",
        Description: "update invoice",
    }
)
```

---

## Step 2: Define Models (`models.go`)

Add request/response models if needed:

```go
// models.go

// ✅ NEW MODEL - Add request/response types
type UpdateInvoiceRequest struct {
    Amount      float64 `json:"amount"`
    Status      string  `json:"status"`
    Description string  `json:"description,omitempty"`
}

// Response uses existing Invoice model
```

---

## Step 3: Update Interface (`client.go`)

Add method signature to the interface:

```go
// client.go
type BillingClient interface {
    // Existing methods...
    GetInvoice(ctx context.Context, prescriptionID string) (*Invoice, error)
    CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (*Invoice, error)
    
    // ✅ NEW METHOD - Add signature
    UpdateInvoice(ctx context.Context, invoiceID string, req UpdateInvoiceRequest) (*Invoice, error)
}
```

---

## Step 4: Implement HTTP Method (`http_client.go`)

Add the implementation - **only 2-3 lines**:

```go
// http_client.go

// ✅ NEW IMPLEMENTATION - Super simple!
func (c *HTTPClient) UpdateInvoice(ctx context.Context, invoiceID string, req UpdateInvoiceRequest) (*Invoice, error) {
    var invoice Invoice
    err := c.api.Put(ctx, UpdateInvoiceEndpoint, map[string]string{
        "invoiceID": invoiceID,
    }, req, &invoice)
    
    return &invoice, err
}
```

That's it! The base API service handles:
- ✅ URL building (replaces `{invoiceID}` with value)
- ✅ JSON marshaling/unmarshaling
- ✅ HTTP request execution
- ✅ Error handling
- ✅ Logging

---

## Step 5: Implement Mock Method (`mock_client.go`)

Add mock implementation for testing:

```go
// mock_client.go

// ✅ NEW MOCK - Simple test data
func (c *MockClient) UpdateInvoice(ctx context.Context, invoiceID string, req UpdateInvoiceRequest) (*Invoice, error) {
    // Find and update invoice in mock data
    for prescID, invoice := range c.invoices {
        if invoice.ID == invoiceID {
            invoice.Amount = req.Amount
            invoice.Status = req.Status
            c.invoices[prescID] = invoice
            
            c.logger.Debug("mock invoice updated",
                zap.String("invoice_id", invoiceID),
                zap.Float64("amount", req.Amount),
            )
            
            return &invoice, nil
        }
    }
    
    return nil, fmt.Errorf("invoice not found: %s", invoiceID)
}
```

---

## ✅ Done! Ready to Use

Now you can use your new endpoint:

```go
// In your business logic
invoice, err := billingClient.UpdateInvoice(ctx, "invoice-123", iris_billing.UpdateInvoiceRequest{
    Amount: 150.00,
    Status: "paid",
})
```

**Automatic features you get for free:**
- ✅ Structured logging
- ✅ Request/response timing
- ✅ Error handling
- ✅ URL parameter replacement
- ✅ JSON encoding/decoding
- ✅ Context propagation
- ✅ Timeout handling

---

## 🚀 Complete Example: Adding Multiple Endpoints

Let's add 3 new endpoints at once:

### 1. Define Endpoints (`endpoints.go`)
```go
var (
    // ... existing ...
    
    // Batch operations
    GetInvoicesByStatusEndpoint = httpclient.EndpointConfig{
        Method:      "GET",
        Path:        "invoices?status={status}",
        Description: "get invoices by status",
    }
    
    DeleteInvoiceEndpoint = httpclient.EndpointConfig{
        Method:      "DELETE",
        Path:        "invoices/{invoiceID}",
        Description: "delete invoice",
    }
    
    RefundInvoiceEndpoint = httpclient.EndpointConfig{
        Method:      "POST",
        Path:        "invoices/{invoiceID}/refund",
        Description: "refund invoice",
    }
)
```

### 2. Define Models (`models.go`)
```go
type RefundInvoiceRequest struct {
    Reason string  `json:"reason"`
    Amount float64 `json:"amount"`
}

type InvoiceList struct {
    Invoices []Invoice `json:"invoices"`
    Total    int       `json:"total"`
}
```

### 3. Update Interface (`client.go`)
```go
type BillingClient interface {
    // ... existing ...
    
    GetInvoicesByStatus(ctx context.Context, status string) (*InvoiceList, error)
    DeleteInvoice(ctx context.Context, invoiceID string) error
    RefundInvoice(ctx context.Context, invoiceID string, req RefundInvoiceRequest) (*Invoice, error)
}
```

### 4. Implement (`http_client.go`)
```go
func (c *HTTPClient) GetInvoicesByStatus(ctx context.Context, status string) (*InvoiceList, error) {
    var list InvoiceList
    err := c.api.Get(ctx, GetInvoicesByStatusEndpoint, map[string]string{
        "status": status,
    }, &list)
    return &list, err
}

func (c *HTTPClient) DeleteInvoice(ctx context.Context, invoiceID string) error {
    return c.api.Delete(ctx, DeleteInvoiceEndpoint, map[string]string{
        "invoiceID": invoiceID,
    })
}

func (c *HTTPClient) RefundInvoice(ctx context.Context, invoiceID string, req RefundInvoiceRequest) (*Invoice, error) {
    var invoice Invoice
    err := c.api.Post(ctx, RefundInvoiceEndpoint, map[string]string{
        "invoiceID": invoiceID,
    }, req, &invoice)
    return &invoice, err
}
```

**Done!** Three new endpoints added with minimal code.

---

## 📊 Code Reduction Comparison

### Before (Old Architecture):
```go
// Each endpoint needed ~50 lines
func (s *HTTPService) GetInvoice(ctx context.Context, prescriptionID string) (*Invoice, error) {
    url := s.endpoint + prescriptionID  // Manual URL building
    
    s.logger.Debug("fetching invoice", ...)  // Manual logging
    
    // Manual HTTP request
    resp, err := s.client.Get(ctx, url, map[string]string{
        "Content-Type": "application/json",
        "Accept":       "application/json",
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get invoice: %w", err)
    }
    
    // Manual status check
    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("iris billing service returned error status %d", resp.StatusCode)
    }
    
    // Manual JSON decoding
    var invoice Invoice
    if err := json.Unmarshal(resp.Body, &invoice); err != nil {
        s.logger.Error("failed to decode invoice response", ...)
        return nil, fmt.Errorf("failed to decode invoice response: %w", err)
    }
    
    s.logger.Debug("invoice retrieved successfully", ...)
    
    return &invoice, nil
}
```

**50 lines** per endpoint × 4 endpoints = **200 lines**

### After (New Architecture):
```go
// Each endpoint needs ~3-5 lines
func (c *HTTPClient) GetInvoice(ctx context.Context, prescriptionID string) (*Invoice, error) {
    var invoice Invoice
    err := c.api.Get(ctx, GetInvoiceEndpoint, map[string]string{
        "prescriptionID": prescriptionID,
    }, &invoice)
    return &invoice, err
}
```

**5 lines** per endpoint × 4 endpoints = **20 lines**

**90% code reduction!** 🎉

---

## 🎯 Benefits

### 1. **Less Boilerplate**
- No manual URL building
- No manual JSON parsing
- No repetitive error handling
- No HTTP client management

### 2. **Consistent Behavior**
- All endpoints use same base service
- Consistent logging format
- Consistent error handling
- Consistent observability

### 3. **Easy to Maintain**
- Endpoints defined in one place
- Models organized together
- Clear naming conventions
- Self-documenting code

### 4. **Easy to Test**
- Mock implementations are simple
- Clear interface contracts
- Type-safe

### 5. **Developer Friendly**
- Just define what, not how
- Clear patterns to follow
- Minimal cognitive load
- Fast to add new endpoints

---

## 📝 Naming Conventions

### Files:
- `client.go` - Interface definition
- `http_client.go` - HTTP implementation
- `mock_client.go` - Mock implementation
- `endpoints.go` - Endpoint definitions
- `models.go` - All request/response models
- `config.go` - Configuration
- `module.go` - Initialization

### Types:
- `BillingClient` - Main interface (not generic "Client")
- `HTTPClient` - HTTP implementation
- `MockClient` - Mock implementation
- Models: `Invoice`, `CreateInvoiceRequest`, `InvoicePayment`

### Methods:
- Verb + Noun pattern: `GetInvoice`, `CreateInvoice`, `UpdateInvoice`
- Clear, self-documenting names

---

## 🚀 Summary

**Adding a new endpoint now takes ~5 minutes:**

1. Add endpoint config to `endpoints.go` (1 min)
2. Add models to `models.go` if needed (1 min)
3. Add method to interface in `client.go` (30 sec)
4. Implement in `http_client.go` with 3-5 lines (2 min)
5. Add mock implementation if needed (30 sec)

**vs. Old way: ~30 minutes** of writing repetitive boilerplate code!

**Result**: Clean, organized, easy to maintain, and developer-friendly! 🎉

