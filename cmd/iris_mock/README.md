# IRIS Mock Server

## Overview

Mock server for all IRIS external APIs - pharmacy, billing, and Stargate authentication.

## Features

✅ **All API Endpoints Implemented**
- Pharmacy API (1 endpoint)
- Billing API (4 endpoints)
- Stargate OAuth (2 endpoints)

✅ **Header Logging**
- Logs all incoming headers
- Shows `X-IRIS-User-ID`, `X-IRIS-Env-Name`, `X-Idempotency-Key`
- Masks Authorization tokens for security

✅ **Request/Response Matching**
- Uses same models as actual integrations
- Response structures match production

---

## Quick Start

### **Run the Mock Server:**

```bash
# From project root
go run cmd/iris_mock/main.go

# Or build and run
go build -o bin/iris_mock cmd/iris_mock/main.go
./bin/iris_mock
```

### **Server Starts On:**
```
🚀 IRIS Mock Server starting on :8081
📍 Pharmacy API: http://localhost:8081/pharmacy/v1
📍 Billing API:  http://localhost:8081/billing/v1
📍 Stargate Auth: http://localhost:8081/oauth
```

---

## API Endpoints

### **Pharmacy API**

#### Get Prescription
```http
GET /pharmacy/v1/prescriptions/{prescriptionID}

Response:
{
  "id": "RX-123",
  "patient_id": "PAT-001",
  "drug": "Lisinopril",
  "dose": "10mg",
  "status": "active",
  "pharmacy_name": "CVS Pharmacy",
  "pharmacy_type": "Retail"
}
```

---

### **Billing API**

#### Get Invoice
```http
GET /billing/v1/invoices/{prescriptionID}

Headers:
  X-IRIS-Env-Name: IRIS_stage  (logged if present)

Response:
{
  "id": "INV-RX-123",
  "prescription_id": "RX-123",
  "amount": 125.50,
  "status": "pending",
  "created_at": "2025-10-14T10:00:00Z",
  "updated_at": "2025-10-14T10:00:00Z"
}
```

#### Create Invoice
```http
POST /billing/v1/invoices

Headers:
  X-Idempotency-Key: a1b2c3d4...  (logged if present)

Request:
{
  "prescription_id": "RX-123",
  "amount": 125.50,
  "description": "Prescription medication"
}

Response: 201 Created
{
  "id": "INV-NEW-RX-123",
  "prescription_id": "RX-123",
  "amount": 125.50,
  "status": "pending",
  "created_at": "2025-10-14T10:00:00Z"
}
```

#### Acknowledge Invoice
```http
POST /billing/v1/invoices/{invoiceID}/acknowledge

Request:
{
  "acknowledged_by": "user@example.com",
  "notes": "Invoice reviewed"
}

Response:
{
  "id": "INV-123",
  "prescription_id": "RX-123",
  "amount": 125.50,
  "status": "acknowledged",
  "updated_at": "2025-10-14T10:00:00Z"
}
```

#### Get Invoice Payment
```http
GET /billing/v1/invoices/{invoiceID}/payment

Response:
{
  "invoice_id": "INV-123",
  "payment_id": "PAY-INV-123",
  "amount": 125.50,
  "payment_method": "credit_card",
  "status": "completed",
  "paid_at": "2025-10-14T10:00:00Z"
}
```

---

### **Stargate OAuth API**

#### Get Access Token
```http
POST /oauth/token

Request:
{
  "grant_type": "client_credentials",
  "client_id": "rxintake-app",
  "client_secret": "your-secret",
  "scope": "api.read api.write"
}

Response:
{
  "access_token": "mock-access-token-rxintake-app",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "mock-refresh-token-rxintake-app",
  "scope": "api.read api.write"
}
```

#### Refresh Token
```http
POST /oauth/refresh

Request:
{
  "grant_type": "refresh_token",
  "refresh_token": "mock-refresh-token-...",
  "client_id": "rxintake-app",
  "client_secret": "your-secret"
}

Response:
{
  "access_token": "mock-refreshed-token-...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "mock-refresh-token-..."
}
```

---

## Header Logging

The mock server logs all important headers:

### **Example Logs:**

```
📥 GET /billing/v1/invoices/RX-123
   └─ X-IRIS-User-ID: xyz
   └─ X-IRIS-Env-Name: IRIS_stage
✅ Returned invoice for prescription: RX-123

📥 POST /billing/v1/invoices
   └─ X-IRIS-User-ID: xyz
   └─ X-Idempotency-Key: a1b2c3d4e5f67890
💡 Idempotency key: a1b2c3d4e5f67890
✅ Created invoice: INV-NEW-RX-123 (Amount: 125.50)

📥 POST /oauth/token
🔐 Token request from client: rxintake-app
✅ Issued token for: rxintake-app (expires in 3600 seconds)
```

---

## Testing with Your App

### **1. Configure Your App to Use Mock:**

```yaml
# internal/configs/app.yaml
external:
  pharmacy:
    use_mock: false  # Use HTTP
    endpoints:
      get_prescription: "http://localhost:8081/pharmacy/v1/prescriptions/{prescriptionID}"
  
  billing:
    use_mock: false  # Use HTTP
    endpoints:
      get_invoice: "http://localhost:8081/billing/v1/invoices/{prescriptionID}"
      create_invoice: "http://localhost:8081/billing/v1/invoices"
      acknowledge_invoice: "http://localhost:8081/billing/v1/invoices/{invoiceID}/acknowledge"
      get_invoice_payment: "http://localhost:8081/billing/v1/invoices/{invoiceID}/payment"
  
  stargate:
    use_mock: false  # Use HTTP
    endpoints:
      token: "http://localhost:8081/oauth/token"
      refresh_token: "http://localhost:8081/oauth/refresh"
```

### **2. Start Mock Server:**
```bash
go run cmd/iris_mock/main.go
```

### **3. Start Your App:**
```bash
# In another terminal
make run
# or
go run cmd/server/main.go
```

### **4. Make API Calls:**

Your app will now call the mock server, and you'll see:

**Mock Server Logs:**
```
📥 GET /billing/v1/invoices/RX-123
   └─ X-IRIS-User-ID: xyz
   └─ X-IRIS-Env-Name: IRIS_stage
✅ Returned invoice for prescription: RX-123
```

**Your App Logs:**
```
INFO  http request completed
      service=external_apis method=GET status_code=200 duration=5ms
INFO  http metrics
      method=GET duration=5ms response_bytes=156
DEBUG invoice retrieved successfully
      prescription_id=RX-123 invoice_id=INV-RX-123
```

---

## Testing Features

### **Test Global Headers:**
Look for log entry: `└─ X-IRIS-User-ID: xyz`

### **Test Endpoint-Specific Headers:**
Look for log entry: `└─ X-IRIS-Env-Name: IRIS_stage` (only on GetInvoice)

### **Test Idempotency:**
Look for log entry: `└─ X-Idempotency-Key: a1b2c3d4...` (only on CreateInvoice)

### **Test Authentication:**
1. Call any billing/pharmacy endpoint
2. Should see auth token request first
3. Then see cached token usage

---

## cURL Examples

### **Test Pharmacy:**
```bash
curl -X GET http://localhost:8081/pharmacy/v1/prescriptions/RX-123 \
  -H "X-IRIS-User-ID: xyz"
```

### **Test Get Invoice:**
```bash
curl -X GET http://localhost:8081/billing/v1/invoices/RX-123 \
  -H "X-IRIS-User-ID: xyz" \
  -H "X-IRIS-Env-Name: IRIS_stage"
```

### **Test Create Invoice:**
```bash
curl -X POST http://localhost:8081/billing/v1/invoices \
  -H "Content-Type: application/json" \
  -H "X-IRIS-User-ID: xyz" \
  -H "X-Idempotency-Key: test-key-123" \
  -d '{
    "prescription_id": "RX-123",
    "amount": 125.50,
    "description": "Test invoice"
  }'
```

### **Test OAuth Token:**
```bash
curl -X POST http://localhost:8081/oauth/token \
  -H "Content-Type: application/json" \
  -d '{
    "grant_type": "client_credentials",
    "client_id": "test-client",
    "client_secret": "test-secret",
    "scope": "api.read api.write"
  }'
```

---

## Summary

**IRIS Mock Server provides:**
- ✅ All 7 API endpoints working
- ✅ Header logging for verification
- ✅ Response structures match production
- ✅ OAuth token simulation
- ✅ Idempotency key tracking
- ✅ Easy local testing

**Perfect for development and testing!** 🚀

