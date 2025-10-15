# API Integration Configuration Example

## YAML Configuration Structure

With the new architecture, endpoints are configured in YAML with **full URLs** instead of being hardcoded.

### Example Configuration (`internal/configs/app.yaml`)

```yaml
external:
  pharmacy:
    base_url: "https://api.iris.example.com"
    path: "/pharmacy/v1"
    use_mock: false
    timeout: "30s"
    api_key: "${PHARMACY_API_KEY}"

  billing:
    use_mock: false
    timeout: "30s"
    endpoints:
      # Full URLs with path parameters as {paramName}
      get_invoice: "https://api.iris.example.com/billing/v1/invoices/{prescriptionID}"
      create_invoice: "https://api.iris.example.com/billing/v1/invoices"
      acknowledge_invoice: "https://api.iris.example.com/billing/v1/invoices/{invoiceID}/acknowledge"
      get_invoice_payment: "https://api.iris.example.com/billing/v1/invoices/{invoiceID}/payment"
```

## Benefits of Config-Based Endpoints

### 1. **Flexibility**
- Change endpoints without code changes
- Different environments (dev, staging, prod) can have different URLs
- Easy to update API versions

### 2. **No Hardcoding**
- All URLs defined in one place (config file)
- Easy to see all available endpoints
- Environment-specific configurations

### 3. **Path Parameters**
Use `{paramName}` in URLs for dynamic values:
```yaml
get_invoice: "https://api.iris.com/v1/invoices/{prescriptionID}"
```

The client automatically replaces `{prescriptionID}` with the actual value:
```go
client.GetInvoice(ctx, "RX-12345")
// Calls: https://api.iris.com/v1/invoices/RX-12345
```

## Environment Variables

Use environment variables for sensitive data:
```yaml
billing:
  endpoints:
    get_invoice: "${BILLING_API_URL}/invoices/{prescriptionID}"
    create_invoice: "${BILLING_API_URL}/invoices"
```

## Adding New Endpoints

### 1. Update Config Structure

Add to `internal/platform/config/config.go`:
```go
type BillingEndpoints struct {
    GetInvoice         string `mapstructure:"get_invoice"`
    CreateInvoice      string `mapstructure:"create_invoice"`
    AcknowledgeInvoice string `mapstructure:"acknowledge_invoice"`
    GetInvoicePayment  string `mapstructure:"get_invoice_payment"`
    
    // NEW: Add your new endpoint
    RefundInvoice      string `mapstructure:"refund_invoice"`
}
```

### 2. Update Config Interface

Add to `internal/integrations/iris_billing/config.go`:
```go
type EndpointsConfig interface {
    GetInvoiceEndpoint() string
    CreateInvoiceEndpoint() string
    // ... existing ...
    
    // NEW: Add getter method
    RefundInvoiceEndpoint() string
}

// Implement the method
func (c *Config) RefundInvoiceEndpoint() string {
    return c.RefundInvoiceURL
}
```

### 3. Update YAML Config

Add to `internal/configs/app.yaml`:
```yaml
external:
  billing:
    endpoints:
      # ... existing endpoints ...
      refund_invoice: "https://api.iris.example.com/billing/v1/invoices/{invoiceID}/refund"
```

### 4. Update integration_wire.go

```go
billing := irisbilling.Module(irisbilling.ModuleDependencies{
    Config: irisbilling.Config{
        GetInvoiceURL:     deps.Config.External.Billing.Endpoints.GetInvoice,
        // ... existing ...
        RefundInvoiceURL:  deps.Config.External.Billing.Endpoints.RefundInvoice, // NEW
    },
    // ... rest ...
})
```

### 5. Use in Implementation

```go
func (c *HTTPClient) RefundInvoice(ctx, invoiceID, req) (*RefundInvoiceResponse, error) {
    url := replacePathParams(c.endpoints.RefundInvoiceEndpoint(), map[string]string{
        "invoiceID": invoiceID,
    })
    // ... make request ...
}
```

## Multiple Environments

### Development (`app.dev.yaml`)
```yaml
external:
  billing:
    use_mock: true  # Use mocks in dev
    endpoints:
      get_invoice: "http://localhost:8081/billing/invoices/{prescriptionID}"
```

### Staging (`app.staging.yaml`)
```yaml
external:
  billing:
    use_mock: false
    endpoints:
      get_invoice: "https://api-staging.iris.example.com/billing/v1/invoices/{prescriptionID}"
```

### Production (`app.prod.yaml`)
```yaml
external:
  billing:
    use_mock: false
    endpoints:
      get_invoice: "https://api.iris.example.com/billing/v1/invoices/{prescriptionID}"
```

## Query Parameters

For endpoints with query parameters, include them in the URL:
```yaml
endpoints:
  list_invoices: "https://api.iris.com/v1/invoices?status={status}&limit={limit}"
```

Use in code:
```go
url := replacePathParams(c.endpoints.ListInvoicesEndpoint(), map[string]string{
    "status": "paid",
    "limit":  "10",
})
// Result: https://api.iris.com/v1/invoices?status=paid&limit=10
```

## Best Practices

### 1. **Use Full URLs**
```yaml
# ✅ GOOD: Full URL
get_invoice: "https://api.iris.com/billing/v1/invoices/{prescriptionID}"

# ❌ BAD: Relative path
get_invoice: "/invoices/{prescriptionID}"
```

### 2. **Include API Version**
```yaml
# ✅ GOOD: Version in URL
get_invoice: "https://api.iris.com/billing/v1/invoices/{id}"

# Future: Easy to update
get_invoice: "https://api.iris.com/billing/v2/invoices/{id}"
```

### 3. **Use Descriptive Parameter Names**
```yaml
# ✅ GOOD: Clear parameter names
get_invoice: "https://api.iris.com/invoices/{prescriptionID}"

# ❌ BAD: Generic names
get_invoice: "https://api.iris.com/invoices/{id}"
```

### 4. **Group Related Endpoints**
```yaml
billing:
  endpoints:
    # Invoice operations
    get_invoice: "..."
    create_invoice: "..."
    update_invoice: "..."
    
    # Payment operations
    get_invoice_payment: "..."
    process_payment: "..."
```

## Migration from Hardcoded Endpoints

### Before (Hardcoded):
```go
// endpoints.go - Hardcoded
var GetInvoiceEndpoint = EndpointConfig{
    Path: "invoices/{prescriptionID}",
}

// http_client.go - Build URL manually
url := s.baseURL + "/" + s.basePath + "/invoices/" + prescriptionID
```

### After (Config-Based):
```yaml
# app.yaml - Configured
endpoints:
  get_invoice: "https://api.iris.com/billing/v1/invoices/{prescriptionID}"
```

```go
// http_client.go - Use config
url := replacePathParams(c.endpoints.GetInvoiceEndpoint(), map[string]string{
    "prescriptionID": prescriptionID,
})
```

## Summary

**Config-based endpoints provide:**
- ✅ Flexibility - change URLs without code changes
- ✅ Environment-specific configurations
- ✅ No hardcoded values
- ✅ Easy to maintain
- ✅ Clear visibility of all endpoints
- ✅ Support for different API versions

**Next Steps:**
1. Update your `app.yaml` with endpoint URLs
2. Set environment variables for sensitive data
3. Test with different environments

