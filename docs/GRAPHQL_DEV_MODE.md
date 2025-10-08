# GraphQL Playground with Dev Mode

## Overview

GraphQL Playground now supports dev mode, allowing you to test GraphQL queries with mock users without needing real JWT tokens.

---

## 🚀 Quick Start

### 1. Enable Dev Mode

```yaml
# internal/configs/app.yaml
auth:
  dev_mode: true
```

### 2. Start Your App

```bash
go run cmd/server/main.go
```

### 3. Open GraphQL Playground

Navigate to:
```
http://localhost:8080/playground
```

### 4. Run Queries (Default Admin User)

By default, queries run as the **admin** user:

```graphql
query {
  patients {
    id
    name
    email
  }
}
```

This will work because admin has `admin:all` permission!

---

## 🎯 Testing with Different Mock Users

### Option 1: Set HTTP Headers in Playground

In the GraphQL Playground interface:

1. Click **"HTTP HEADERS"** at the bottom left
2. Add the `X-Mock-User` header:

```json
{
  "X-Mock-User": "doctor"
}
```

3. Run your query - it will now execute as the doctor user

### Option 2: Use curl with Different Users

```bash
# Admin user (default)
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ patients { id name } }"
  }'

# Doctor user
curl -X POST http://localhost:8080/graphql \
  -H "X-Mock-User: doctor" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ patients { id name } }"
  }'

# Nurse user (limited permissions)
curl -X POST http://localhost:8080/graphql \
  -H "X-Mock-User: nurse" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ patients { id name } }"
  }'

# User with no permissions (should fail)
curl -X POST http://localhost:8080/graphql \
  -H "X-Mock-User: noperm" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ patients { id name } }"
  }'
```

---

## 🧪 Testing Permission Scenarios

### Test 1: Read Access (Everyone Can Read)

**Query:**
```graphql
query {
  patients {
    id
    name
  }
}
```

**Test with different users:**
```json
// HTTP HEADERS in Playground

// Admin - should work
{"X-Mock-User": "admin"}

// Doctor - should work (has patient:read)
{"X-Mock-User": "doctor"}

// Nurse - should work (has patient:read)
{"X-Mock-User": "nurse"}

// Readonly - should work (has patient:read)
{"X-Mock-User": "readonly"}
```

All should return **200 OK** with patient data.

---

### Test 2: Prescription Access (Healthcare Roles)

**Query:**
```graphql
query {
  prescriptions {
    id
    drug
    dose
    status
  }
}
```

**Expected Results:**

| User | Has Permission? | Result |
|------|----------------|--------|
| `admin` | ✅ `admin:all` | Success |
| `doctor` | ✅ `prescription:read` | Success |
| `pharmacist` | ✅ `prescription:read` | Success |
| `nurse` | ✅ `prescription:read` | Success |
| `readonly` | ✅ `prescription:read` | Success |

All should work!

---

### Test 3: Permission Failures

Create a custom user with no permissions to test authorization failures:

**Add to your code:**
```go
if auth.IsDevModeEnabled() {
    auth.AddMockUser("noperm", &auth.User{
        ID:          "test-noperm",
        Email:       "noperm@test.local",
        Name:        "No Permissions",
        Permissions: []string{},  // Empty!
    })
}
```

**Test in Playground:**
```graphql
query {
  patients {
    id
    name
  }
}
```

**HTTP Headers:**
```json
{"X-Mock-User": "noperm"}
```

**Expected Response:**
```json
{
  "errors": [
    {
      "message": "Forbidden: requires at least one of: [patient:read, admin:all]",
      "extensions": {
        "code": "FORBIDDEN",
        "status": 403,
        "required_permissions": ["patient:read", "admin:all"],
        "match": "any"
      }
    }
  ],
  "data": null
}
```

---

## 📋 Testing Checklist

Use this checklist to verify GraphQL permissions are working correctly:

### Patient Queries

```graphql
# Patient by ID
query GetPatient {
  patient(id: "1") {
    id
    name
    email
  }
}
```

Test with:
- [x] `admin` - Should work
- [x] `doctor` - Should work
- [x] `nurse` - Should work
- [x] `readonly` - Should work

---

### Prescription Queries

```graphql
# All prescriptions
query GetPrescriptions {
  prescriptions {
    id
    drug
    patient {
      name
    }
  }
}
```

Test with:
- [x] `admin` - Should work
- [x] `doctor` - Should work
- [x] `pharmacist` - Should work
- [x] `nurse` - Should work

---

### Nested Fields with Different Permissions

```graphql
# Patient with prescriptions (requires both permissions)
query GetPatientWithPrescriptions {
  patient(id: "1") {
    id
    name
    prescriptions {
      id
      drug
    }
  }
}
```

**Expected Behavior:**
- Patient data requires: `patient:read` OR `admin:all`
- Prescriptions field requires: `prescription:read` OR healthcare roles

All mock users should be able to access this since they all have both permissions.

---

### Dashboard Query

```graphql
query GetDashboard {
  dashboardStats {
    totalPatients
    activePrescriptions
  }
}
```

Test with:
- [x] `admin` - Should work
- [x] `doctor` - Should work (has dashboard:view)
- [x] `pharmacist` - Should work (has dashboard:view)
- [x] `nurse` - Should FAIL (missing dashboard:view)
- [x] `readonly` - Should work (has dashboard:view)

---

## 🎨 GraphQL Playground Screenshots

### Setting Mock User Header

1. Open GraphQL Playground at `http://localhost:8080/playground`
2. At the bottom left, click **"HTTP HEADERS"**
3. Add header JSON:
```json
{
  "X-Mock-User": "doctor"
}
```
4. Run your query
5. Check response - should see successful data or permission error

### Default Behavior (No Header)

If you don't set the header:
- Uses **admin** user by default
- All queries should work (admin has full access)
- Useful for quick testing

---

## 🔧 Advanced Testing

### Test Permission Boundaries

```graphql
# This query requires prescription:read
query {
  prescriptions {
    id
    # This nested field requires patient:read
    patient {
      name
    }
  }
}
```

**Test Matrix:**

| User | Has prescription:read? | Has patient:read? | Result |
|------|----------------------|-------------------|--------|
| `admin` | ✅ (admin:all) | ✅ (admin:all) | Success |
| `doctor` | ✅ | ✅ | Success |
| `nurse` | ✅ | ✅ | Success |

---

### Test Different Query Combinations

**Query 1: Simple patient list**
```graphql
{ patients { id name } }
```
Header: `{"X-Mock-User": "nurse"}` → ✅ Success

**Query 2: Patient with prescriptions**
```graphql
{ patient(id: "1") { name prescriptions { drug } } }
```
Header: `{"X-Mock-User": "nurse"}` → ✅ Success (has both permissions)

**Query 3: Dashboard stats**
```graphql
{ dashboardStats { totalPatients activePrescriptions } }
```
Header: `{"X-Mock-User": "nurse"}` → ❌ 403 Forbidden (missing dashboard:view)

---

## 🐛 Troubleshooting

### Issue: Getting UNAUTHENTICATED Error

**Problem:**
```json
{
  "errors": [{
    "message": "Unauthenticated",
    "extensions": {"code": "UNAUTHENTICATED", "status": 401}
  }]
}
```

**Solution:**
1. Check dev mode is enabled in `app.yaml`
2. Check app startup logs for "AUTH DEV MODE ENABLED"
3. Try with explicit header: `{"X-Mock-User": "admin"}`
4. If still failing, dev mode might not be properly initialized

---

### Issue: Getting FORBIDDEN Error

**Problem:**
```json
{
  "errors": [{
    "message": "Forbidden: requires at least one of: [patient:read, admin:all]",
    "extensions": {"code": "FORBIDDEN", "status": 403}
  }]
}
```

**Solution:**
This is correct! The user doesn't have the required permission.
1. Check which user you're using: `{"X-Mock-User": "nurse"}`
2. Check what permissions that user has: `curl http://localhost:8080/__dev/auth`
3. Use a different user with the required permissions

---

### Issue: Playground Not Loading

**Problem:** Can't access `http://localhost:8080/playground`

**Solution:**
1. Check the app is running
2. Verify the path in `internal/platform/paths/registry.go`
3. Check logs for "GraphQL server mounted"

---

## 💡 Tips for GraphQL Testing

### 1. Start with Admin
When exploring the schema, use admin first:
```json
{"X-Mock-User": "admin"}
```
This ensures all queries work so you can focus on schema exploration.

### 2. Test Permission Boundaries
Switch to limited users to verify permissions:
```json
{"X-Mock-User": "nurse"}
```
Try queries that should fail - verify you get 403 errors.

### 3. Use Introspection Query
Explore the schema:
```graphql
query IntrospectSchema {
  __schema {
    types {
      name
      fields {
        name
      }
    }
  }
}
```

### 4. Check Applied Directives
See which fields have auth directives:
```graphql
query {
  __type(name: "Query") {
    fields {
      name
      # Directives are visible in the schema
    }
  }
}
```

---

## 📝 Example Testing Session

### Session 1: Verify All Users Can Read Patients

**Query:**
```graphql
query GetPatients {
  patients {
    id
    name
    phone
  }
}
```

**Test each user:**
1. `{"X-Mock-User": "admin"}` → ✅ Success
2. `{"X-Mock-User": "doctor"}` → ✅ Success
3. `{"X-Mock-User": "pharmacist"}` → ✅ Success
4. `{"X-Mock-User": "nurse"}` → ✅ Success
5. `{"X-Mock-User": "readonly"}` → ✅ Success

---

### Session 2: Verify Dashboard Access

**Query:**
```graphql
query GetDashboard {
  dashboardStats {
    totalPatients
    activePrescriptions
  }
}
```

**Expected Results:**
1. `{"X-Mock-User": "admin"}` → ✅ Success
2. `{"X-Mock-User": "doctor"}` → ✅ Success (has dashboard:view)
3. `{"X-Mock-User": "pharmacist"}` → ✅ Success (has dashboard:view)
4. `{"X-Mock-User": "nurse"}` → ❌ 403 Forbidden (no dashboard:view)
5. `{"X-Mock-User": "readonly"}` → ✅ Success (has dashboard:view)

---

### Session 3: Nested Field Permissions

**Query:**
```graphql
query GetPatientWithPrescriptions {
  patient(id: "1") {
    name
    prescriptions {
      drug
      dose
    }
  }
}
```

This tests:
- `patient(id)` requires: `patient:read` OR `admin:all`
- `prescriptions` field requires: `prescription:read` OR healthcare roles

All mock users should succeed (they all have both permissions).

---

## 🎯 GraphQL Playground Features with Dev Mode

### ✅ What Works
- Execute queries with mock users
- Test permission scenarios
- Switch users via HTTP headers
- See proper error messages (401, 403)
- Introspection and schema exploration
- Default to admin user (no header needed)

### 🎨 Visual Flow

```
GraphQL Playground
    ↓
User sets header: {"X-Mock-User": "doctor"}
    ↓
Execute query
    ↓
POST /graphql
    ↓
RequireAuthWithDevMode() middleware
    ↓
Dev mode enabled?
    ├─ YES → Use mock user "doctor"
    └─ NO → Validate JWT token
    ↓
User in context
    ↓
GraphQL directives check permissions
    ├─ @auth → User authenticated?
    ├─ @permissionAny → Has any of the permissions?
    └─ @permissionAll → Has all of the permissions?
    ↓
Execute resolver
    ↓
Return response
```

---

## 📚 Example Queries for Each User

### Admin (Full Access)

```graphql
# Can query everything
query AdminTest {
  patients { id name }
  prescriptions { id drug }
  dashboardStats { totalPatients }
}
```

**Header:** None needed (defaults to admin) or `{"X-Mock-User": "admin"}`

---

### Doctor (Clinical Management)

```graphql
# Can manage patients and prescriptions
query DoctorTest {
  patients {
    id
    name
    prescriptions {
      id
      drug
      status
    }
  }
  dashboardStats {
    totalPatients
    activePrescriptions
  }
}
```

**Header:** `{"X-Mock-User": "doctor"}`

---

### Pharmacist (Dispensing Focus)

```graphql
# Can view patients and prescriptions
query PharmacistTest {
  patient(id: "1") {
    name
    phone
    prescriptions {
      id
      drug
      dose
      status
    }
  }
  dashboardStats {
    totalPatients
  }
}
```

**Header:** `{"X-Mock-User": "pharmacist"}`

---

### Nurse (Read-Only Clinical)

```graphql
# Can view patients and prescriptions, but NOT dashboard
query NurseTest {
  patients {
    id
    name
  }
  prescriptions {
    id
    drug
  }
  # This will fail - nurse doesn't have dashboard:view
  # dashboardStats { totalPatients }
}
```

**Header:** `{"X-Mock-User": "nurse"}`

**Dashboard query will return:**
```json
{
  "errors": [{
    "message": "Forbidden: requires at least one of: [dashboard:view, admin:all]",
    "extensions": {"code": "FORBIDDEN", "status": 403}
  }]
}
```

---

### Read-Only (Minimal Access)

```graphql
# Can view data and dashboard
query ReadonlyTest {
  patients { id name }
  prescriptions { id drug }
  dashboardStats { totalPatients }
}
```

**Header:** `{"X-Mock-User": "readonly"}`

---

## 🔒 Permission Testing Matrix

| Query | Admin | Doctor | Pharmacist | Nurse | Readonly |
|-------|-------|--------|------------|-------|----------|
| `patients` | ✅ | ✅ | ✅ | ✅ | ✅ |
| `patient(id)` | ✅ | ✅ | ✅ | ✅ | ✅ |
| `prescriptions` | ✅ | ✅ | ✅ | ✅ | ✅ |
| `prescription(id)` | ✅ | ✅ | ✅ | ✅ | ✅ |
| `dashboardStats` | ✅ | ✅ | ✅ | ❌ | ✅ |
| `Patient.prescriptions` | ✅ | ✅ | ✅ | ✅ | ✅ |
| `Prescription.patient` | ✅ | ✅ | ✅ | ✅ | ✅ |

---

## 🎨 GraphQL Playground UI Tips

### Setting Headers in Playground

**Bottom left panel:**
```
┌─────────────────────────────────────┐
│  HTTP HEADERS                       │
├─────────────────────────────────────┤
│  {                                  │
│    "X-Mock-User": "doctor"         │
│  }                                  │
└─────────────────────────────────────┘
```

### Variables Panel

You can also use variables:
```graphql
query GetPatient($id: ID!) {
  patient(id: $id) {
    name
  }
}
```

**Query Variables:**
```json
{
  "id": "1"
}
```

**HTTP Headers:**
```json
{
  "X-Mock-User": "nurse"
}
```

---

## 🚀 Production Mode

When dev mode is disabled (`auth.dev_mode: false`):

### Use Real JWT Tokens

**HTTP Headers in Playground:**
```json
{
  "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

Or set JWT in cookie and it will automatically be sent.

### No Mock Users

The `X-Mock-User` header will be ignored - only real JWT tokens work.

---

## 📝 Quick Reference

### Default User (No Header)
```graphql
query { patients { id name } }
```
→ Runs as **admin** user

### Doctor User
**Headers:** `{"X-Mock-User": "doctor"}`
```graphql
query { patients { id name prescriptions { drug } } }
```
→ ✅ Success (has both permissions)

### Nurse User (Test Failure)
**Headers:** `{"X-Mock-User": "nurse"}`
```graphql
query { dashboardStats { totalPatients } }
```
→ ❌ 403 Forbidden (missing dashboard:view)

### Check All Users
Visit: `http://localhost:8080/__dev/auth`

---

## 🎉 Summary

GraphQL Playground now supports dev mode:
- ✅ **Default admin user** - No headers needed for quick testing
- ✅ **Mock user selection** - Use `X-Mock-User` header
- ✅ **All 5 mock users** - Test different permission scenarios
- ✅ **Proper error codes** - 401 for auth, 403 for permissions
- ✅ **Directive enforcement** - `@auth`, `@permissionAny`, `@permissionAll` work
- ✅ **Easy testing** - No JWT tokens needed locally

**Open Playground and start testing your GraphQL API!** 🚀

URL: `http://localhost:8080/playground`

