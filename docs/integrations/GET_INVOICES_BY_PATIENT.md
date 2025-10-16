# Get Invoices by Patient ID - Implementation Guide

## Overview

Added a new API endpoint to retrieve all invoices for a specific patient ID in the `iris_billing` integration.

## What Was Added

### 1. **New Interface Method**

```go
// In internal/integrations/iris_billing/client.go
GetInvoicesByPatientID(ctx context.Context, patientID string) (*InvoiceListResponse, error)
```

### 2. **New Response Model**

```go
// In internal/integrations/iris_billing/models.go
type InvoiceListResponse struct {
    PatientID string            `json:"patient_id"`
    Invoices  []InvoiceResponse `json:"invoices"`
    Total     int               `json:"total"`
}
```

### 3. **Configuration**

Updated configuration files to include the new endpoint:

**`internal/configs/app.yaml`:**
```yaml
external:
  billing:
    use_mock: true
    timeout: "10s"
    endpoints:
      get_invoices_by_patient: "http://localhost:8881/billing/v1/patients/{patientID}/invoices"
      # ... other endpoints
```

**`internal/configs/app.prod.yaml`:**
```yaml
external:
  billing:
    endpoints:
      get_invoices_by_patient: "${PM_BILLING_BASE_URL}/billing/v1/patients/{patientID}/invoices"
```

### 4. **HTTP Implementation**

```go
// In internal/integrations/iris_billing/http_client.go
func (c *HTTPClient) GetInvoicesByPatientID(ctx context.Context, patientID string) (*InvoiceListResponse, error) {
    url := replacePathParams(c.endpoints.GetInvoicesByPatientEndpoint(), map[string]string{
        "patientID": patientID,
    })
    
    // Makes HTTP GET request and returns the invoice list
    // ...
}
```

### 5. **Mock Implementation**

```go
// In internal/integrations/iris_billing/mock_client.go
func (c *MockClient) GetInvoicesByPatientID(ctx context.Context, patientID string) (*InvoiceListResponse, error) {
    // Returns mock invoice data for testing
    // ...
}
```

### 6. **IRIS Mock Server**

Added handler in `cmd/iris_mock/main.go`:

```go
func handleGetInvoicesByPatient(w http.ResponseWriter, r *http.Request) {
    patientID := chi.URLParam(r, "patientID")
    
    // Returns mock invoice list with 3 sample invoices
    // - One paid invoice
    // - One pending invoice  
    // - One overdue invoice
}
```

## Usage Example

### Using the Billing Client

```go
import (
    "context"
    "fmt"
    irisbilling "pharmacy-modernization-project-model/internal/integrations/iris_billing"
)

func getPatientInvoices(billingClient irisbilling.BillingClient, patientID string) {
    ctx := context.Background()
    
    // Call the new method
    response, err := billingClient.GetInvoicesByPatientID(ctx, patientID)
    if err != nil {
        fmt.Printf("Error getting invoices: %v\n", err)
        return
    }
    
    fmt.Printf("Found %d invoices for patient %s:\n", response.Total, response.PatientID)
    
    for _, invoice := range response.Invoices {
        fmt.Printf("  - Invoice %s: $%.2f (%s)\n", 
            invoice.ID, 
            invoice.Amount, 
            invoice.Status,
        )
    }
}
```

## Testing with IRIS Mock Server

### 1. Start the Mock Server

```bash
cd /Users/mohsenvafa/code/sfd-root/rxintake_scaffold
go run ./cmd/iris_mock/main.go
```

### 2. Test the Endpoint

```bash
# Get invoices for patient PAT-123
curl http://localhost:8881/billing/v1/patients/PAT-123/invoices | jq .
```

### Expected Response

```json
{
  "patient_id": "PAT-123",
  "invoices": [
    {
      "id": "INV-001",
      "prescription_id": "RX-PAT-123-001",
      "amount": 125.50,
      "status": "paid",
      "created_at": "2025-10-01T10:00:00Z",
      "updated_at": "2025-10-02T14:30:00Z"
    },
    {
      "id": "INV-002",
      "prescription_id": "RX-PAT-123-002",
      "amount": 89.99,
      "status": "pending",
      "created_at": "2025-10-10T09:15:00Z",
      "updated_at": "2025-10-10T09:15:00Z"
    },
    {
      "id": "INV-003",
      "prescription_id": "RX-PAT-123-003",
      "amount": 250.00,
      "status": "overdue",
      "created_at": "2025-09-15T11:20:00Z",
      "updated_at": "2025-09-15T11:20:00Z"
    }
  ],
  "total": 3
}
```

## Files Modified

### Integration Layer
- ✅ `internal/integrations/iris_billing/client.go` - Added interface method
- ✅ `internal/integrations/iris_billing/models.go` - Added InvoiceListResponse model
- ✅ `internal/integrations/iris_billing/config.go` - Added endpoint configuration
- ✅ `internal/integrations/iris_billing/http_client.go` - Added HTTP implementation
- ✅ `internal/integrations/iris_billing/mock_client.go` - Added mock implementation
- ✅ `internal/integrations/integration_wire.go` - Updated wiring with new endpoint

### Configuration
- ✅ `internal/configs/app.yaml` - Added endpoint configuration for dev
- ✅ `internal/configs/app.prod.yaml` - Added endpoint configuration for prod
- ✅ `internal/platform/config/config.go` - Updated BillingEndpoints struct

### Mock Server
- ✅ `cmd/iris_mock/main.go` - Added handler for new endpoint

## Summary

The new `GetInvoicesByPatientID` API:
- ✅ Retrieves all invoices for a given patient ID
- ✅ Returns invoice details including amount, status, dates
- ✅ Supports both HTTP and mock implementations
- ✅ Fully configured for dev and prod environments
- ✅ Tested and working with IRIS mock server

## Next Steps

You can now use this API to:
1. Display patient invoice history in the UI
2. Calculate total outstanding balance for a patient
3. Filter invoices by status (paid, pending, overdue)
4. Generate patient billing reports

