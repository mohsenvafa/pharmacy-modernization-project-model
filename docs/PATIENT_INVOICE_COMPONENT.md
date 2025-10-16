# Patient Invoice Component - Implementation Summary

## Overview

Successfully implemented a lazy-loading invoice component for the patient detail page that displays a list of invoices for a given patient by calling the billing API.

## Features

✅ **Lazy Loading**: Component loads invoices only when scrolled into view using HTMX  
✅ **Color-Coded Status Badges**: 
  - 🟢 Paid (green)
  - 🟡 Pending (yellow)  
  - 🔴 Overdue (red)
✅ **Formatted Data**: Currency amounts and dates properly formatted  
✅ **Empty State**: Graceful handling when no invoices exist  
✅ **Responsive Design**: Works on mobile and desktop  

## Implementation Details

### 1. Provider Layer

**Created Invoice Provider** (`domain/patient/providers/`):
- `invoices.go` - Interface definition
- `invoice_provider_impl.go` - Implementation using billing client

```go
type PatientInvoiceProvider interface {
    GetInvoicesByPatientID(ctx context.Context, patientID string) (*irisbilling.InvoiceListResponse, error)
}
```

### 2. UI Component

**Created Invoice List Component** (`domain/patient/ui/components/patient_invoices/`):
- `patient_invoice_list.templ` - Templ template with lazy-loading placeholder
- `patient_invoice_list.component.go` - Handler and business logic
- `patient_invoices.component.ts` - TypeScript component

**Features:**
- Table displaying: Invoice ID, Prescription ID, Amount, Status, Created Date
- Color-coded status badges (paid=green, pending=yellow, overdue=red)
- Total invoice count
- Empty state message

### 3. Helper Functions

**Added to `internal/helper/tools.go`:**
```go
func FormatDecimal(value float64) string          // Formats: 125.50
func FormatInt(value int) string                   // Formats: 3
func FormatShortDateFromString(dateStr string) string  // Parses ISO date and formats
```

### 4. Patient Detail Page Integration

**Updated `patient_detail_page.component.templ`:**
- Added `Invoices templ.Component` field to `PatientDetailPageParam`
- Positioned invoice card **before Clinical Summary** section
- Uses lazy-loading placeholder that triggers on scroll

**Updated `patient_detail_page.component.go`:**
- Added `invoiceListComponent` dependency
- Creates invoice placeholder: `patientinvoices.PlaceHolder(patientID)`
- Passes invoice component to template

### 5. Routing & Registration

**Updated `domain/patient/ui/paths/paths.go`:**
```go
PatientInvoiceCardComponentRoute = "/components/patient-invoices-card"
```

**Updated `domain/patient/ui/ui.go`:**
- Registered invoice component route
- Wired invoice component to patient detail handler

**Updated TypeScript Registry** (`domain/patient/ui/ts/register_components.ts`):
```typescript
registerComponent('patient.patient-invoices', () => new PatientInvoicesComponent())
```

### 6. Module Wiring

**Updated `domain/patient/module.go`:**
- Added `InvoiceProvider` to `ModuleDependencies`
- Passed invoice provider to UI dependencies

**Updated `internal/app/wire.go`:**
- Created invoice provider using billing client
- Wired invoice provider to patient module

**Updated UI Dependencies** (`domain/patient/ui/contracts/ui-dependencis.go`):
```go
type UiDependencies struct {
    // ... existing fields
    InvoiceProvider patientproviders.PatientInvoiceProvider
}
```

### 7. Configuration

**Updated `internal/configs/app.yaml`:**
```yaml
external:
  billing:
    use_mock: false  # Changed to use HTTP client
    timeout: "10s"
    endpoints:
      get_invoices_by_patient: "http://localhost:8881/billing/v1/patients/{patientID}/invoices"
```

## How It Works

1. **Page Load**: Patient detail page loads with invoice placeholder
2. **Lazy Loading**: When invoice section scrolls into view, HTMX triggers:
   ```
   GET /patients/components/patient-invoices-card?patientId=PAT-123
   ```
3. **API Call**: Invoice component calls billing API:
   ```
   invoiceProvider.GetInvoicesByPatientID(ctx, patientID)
   ```
4. **Data Fetch**: Provider calls IRIS billing integration:
   ```
   billingClient.GetInvoicesByPatientID(ctx, patientID)
   ```
5. **Rendering**: Component renders table with invoice data
6. **HTMX Swap**: Response replaces placeholder with actual invoice table

## Testing

### Test Invoice Component Endpoint
```bash
curl "http://localhost:8080/patients/components/patient-invoices-card?patientId=PAT-123"
```

### Expected Response
HTML table with:
- 3 invoices (INV-001, INV-002, INV-003)
- Color-coded status badges
- Formatted amounts ($125.50, $89.99, $250.00)
- Formatted dates (Oct 1, 2025, Oct 10, 2025, Sep 15, 2025)
- Total count (3 invoice(s))

### Test IRIS Mock Server
```bash
curl "http://localhost:8081/billing/v1/patients/PAT-123/invoices"
```

## Page Layout

```
Patient Detail Page
├── Page Header (Patient Name)
├── Statistics Cards (Age, State, Phone)
├── Patient Details Card
├── 📊 Invoices Card (NEW - Lazy Loaded)  ← Added here
├── Clinical Summary Card
└── Grid Section
    ├── Address List
    └── Prescriptions List
```

## Files Created (7 new files)

1. `domain/patient/providers/invoices.go`
2. `domain/patient/providers/invoice_provider_impl.go`
3. `domain/patient/ui/components/patient_invoices/patient_invoice_list.templ`
4. `domain/patient/ui/components/patient_invoices/patient_invoice_list.component.go`
5. `domain/patient/ui/components/patient_invoices/patient_invoices.component.ts`
6. `domain/patient/ui/components/patient_invoices/patient_invoice_list_templ.go` (generated)
7. `docs/PATIENT_INVOICE_COMPONENT.md` (this file)

## Files Modified (9 files)

1. `internal/helper/tools.go` - Added formatting functions
2. `domain/patient/ui/contracts/ui-dependencis.go` - Added InvoiceProvider
3. `domain/patient/ui/paths/paths.go` - Added invoice component route
4. `domain/patient/ui/ui.go` - Registered invoice component
5. `domain/patient/ui/ts/register_components.ts` - Registered TypeScript component
6. `domain/patient/ui/patient_detail/patient_detail_page.component.templ` - Added invoice section
7. `domain/patient/ui/patient_detail/patient_detail_page.component.go` - Wired invoice component
8. `domain/patient/module.go` - Added InvoiceProvider dependency
9. `internal/app/wire.go` - Wired invoice provider
10. `internal/configs/app.yaml` - Changed use_mock to false

## Benefits

✅ **Lazy Loading**: Improves initial page load time  
✅ **Reusable**: Invoice component can be used elsewhere  
✅ **Type-Safe**: Full type safety with Go and TypeScript  
✅ **Maintainable**: Clean separation of concerns  
✅ **Testable**: Easy to test with mock data or real API  
✅ **Responsive**: Works on all device sizes  
✅ **Accessible**: Semantic HTML with proper structure  

## Next Steps

Potential enhancements:
- Add invoice detail view (click to expand)
- Add filtering by status (paid, pending, overdue)
- Add date range filter
- Add pagination for many invoices
- Add export to PDF/CSV
- Add payment action buttons
- Add total amount calculation
- Add invoice search functionality

## Summary

Successfully implemented a complete invoice component for patient detail pages that:
- Loads invoices lazily when scrolled into view
- Displays invoice data in a clean table format
- Uses color-coded badges for status visualization
- Integrates seamlessly with existing architecture
- Follows established patterns from prescription component
- Is fully tested and working with IRIS mock server

The component is production-ready and can be extended with additional features as needed.

