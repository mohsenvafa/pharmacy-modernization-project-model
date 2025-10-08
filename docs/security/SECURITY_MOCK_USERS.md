# Mock Users Reference Guide

## Overview

When development mode is enabled, RxIntake provides 5 pre-configured mock users for testing different permission scenarios. This guide shows you how to check their permissions and use them effectively.

---

## 🔍 How to Check Mock User Permissions

### Option 1: Dev Info Endpoint (Recommended)

**When dev mode is enabled**, access the dev info endpoint:

```bash
# Get all mock users and their permissions
curl http://localhost:8080/__dev/auth

# Pretty print with jq
curl http://localhost:8080/__dev/auth | jq '.'

# Just list the users
curl http://localhost:8080/__dev/auth | jq '.mock_users[] | {key, name, permissions}'
```

**Response Example:**
```json
{
  "dev_mode": true,
  "message": "Development mode is ACTIVE - Authentication bypassed",
  "mock_users": [
    {
      "key": "admin",
      "id": "mock-admin-001",
      "email": "admin@dev.local",
      "name": "Dev Admin",
      "permissions": ["admin:all"]
    },
    {
      "key": "doctor",
      "id": "mock-doctor-001",
      "email": "doctor@dev.local",
      "name": "Dr. Dev",
      "permissions": [
        "patient:read",
        "patient:write",
        "prescription:read",
        "prescription:write",
        "prescription:approve",
        "doctor:role",
        "dashboard:view"
      ]
    }
    // ... more users
  ],
  "usage": {
    "header": "X-Mock-User: <key>",
    "default": "admin",
    "example_curl": "curl -H 'X-Mock-User: doctor' http://localhost:8080/patients"
  }
}
```

### Option 2: Source Code

Check the mock user definitions directly in code:

**File**: `internal/platform/auth/dev_mode.go`

**Function**: `initializeMockUsers()` (lines 26-85)

### Option 3: This Documentation

See the reference table below for a quick lookup.

---

## 👥 Mock Users Reference

### 1. Admin User

**Key**: `admin`

**Purpose**: Super administrator with full access to everything

**Permissions**:
- `admin:all` - Grants access to all endpoints and operations

**Usage**:
```bash
curl -H "X-Mock-User: admin" http://localhost:8080/patients
curl -H "X-Mock-User: admin" http://localhost:8080/prescriptions
curl -H "X-Mock-User: admin" http://localhost:8080/dashboard
```

**Can Access**:
- ✅ All patient operations (read, write, export, delete)
- ✅ All prescription operations (read, write, approve, dispense, cancel)
- ✅ Dashboard and analytics
- ✅ Everything (bypasses all permission checks via `admin:all`)

**Cannot Access**:
- Nothing - has complete access

---

### 2. Doctor User

**Key**: `doctor`

**Purpose**: Healthcare provider who manages patients and prescriptions

**Permissions**:
- `patient:read` - View patient data
- `patient:write` - Create/update patients
- `prescription:read` - View prescriptions
- `prescription:write` - Create/update prescriptions
- `prescription:approve` - Approve prescriptions
- `doctor:role` - Doctor role identifier
- `dashboard:view` - View dashboard

**Usage**:
```bash
# Can read and write patients
curl -H "X-Mock-User: doctor" http://localhost:8080/patients
curl -X POST -H "X-Mock-User: doctor" http://localhost:8080/api/v1/patients

# Can read and write prescriptions
curl -H "X-Mock-User: doctor" http://localhost:8080/prescriptions
curl -X POST -H "X-Mock-User: doctor" http://localhost:8080/api/v1/prescriptions

# Can view dashboard
curl -H "X-Mock-User: doctor" http://localhost:8080/
```

**Can Access**:
- ✅ Read patient data
- ✅ Create/update patients
- ✅ Read prescriptions
- ✅ Create/update prescriptions
- ✅ Approve prescriptions
- ✅ View dashboard

**Cannot Access**:
- ❌ Export patient data (missing `patient:export`)
- ❌ Delete patients (missing `patient:delete`)
- ❌ Dispense prescriptions (missing `prescription:dispense`)

---

### 3. Pharmacist User

**Key**: `pharmacist`

**Purpose**: Pharmacy staff who dispense prescriptions

**Permissions**:
- `patient:read` - View patient data (read-only)
- `prescription:read` - View prescriptions
- `prescription:dispense` - Dispense prescriptions
- `pharmacist:role` - Pharmacist role identifier
- `dashboard:view` - View dashboard

**Usage**:
```bash
# Can read patients (but not modify)
curl -H "X-Mock-User: pharmacist" http://localhost:8080/patients

# Can read prescriptions
curl -H "X-Mock-User: pharmacist" http://localhost:8080/prescriptions

# Cannot create prescriptions (will fail)
curl -X POST -H "X-Mock-User: pharmacist" \
  http://localhost:8080/api/v1/prescriptions
# Returns: 403 Forbidden
```

**Can Access**:
- ✅ Read patient data
- ✅ Read prescriptions
- ✅ Dispense prescriptions
- ✅ View dashboard

**Cannot Access**:
- ❌ Create/update patients (missing `patient:write`)
- ❌ Create/update prescriptions (missing `prescription:write`)
- ❌ Approve prescriptions (missing `prescription:approve`)

---

### 4. Nurse User

**Key**: `nurse`

**Purpose**: Clinical staff with read-only access to patients and prescriptions

**Permissions**:
- `patient:read` - View patient data
- `prescription:read` - View prescriptions
- `nurse:role` - Nurse role identifier

**Usage**:
```bash
# Can read patients
curl -H "X-Mock-User: nurse" http://localhost:8080/patients

# Can read prescriptions
curl -H "X-Mock-User: nurse" http://localhost:8080/prescriptions

# Cannot view dashboard (will fail)
curl -H "X-Mock-User: nurse" http://localhost:8080/
# Returns: 403 Forbidden

# Cannot create patients (will fail)
curl -X POST -H "X-Mock-User: nurse" http://localhost:8080/api/v1/patients
# Returns: 403 Forbidden
```

**Can Access**:
- ✅ Read patient data
- ✅ Read prescriptions

**Cannot Access**:
- ❌ Create/update patients (missing `patient:write`)
- ❌ Create/update prescriptions (missing `prescription:write`)
- ❌ View dashboard (missing `dashboard:view`)
- ❌ Any write operations

---

### 5. Read-Only User

**Key**: `readonly`

**Purpose**: Users with minimal access for viewing data and dashboard

**Permissions**:
- `patient:read` - View patient data
- `prescription:read` - View prescriptions
- `dashboard:view` - View dashboard

**Usage**:
```bash
# Can read patients
curl -H "X-Mock-User: readonly" http://localhost:8080/patients

# Can read prescriptions  
curl -H "X-Mock-User: readonly" http://localhost:8080/prescriptions

# Can view dashboard
curl -H "X-Mock-User: readonly" http://localhost:8080/

# Cannot create anything (will fail)
curl -X POST -H "X-Mock-User: readonly" http://localhost:8080/api/v1/patients
# Returns: 403 Forbidden
```

**Can Access**:
- ✅ Read patient data
- ✅ Read prescriptions
- ✅ View dashboard

**Cannot Access**:
- ❌ Any write operations
- ❌ Create/update patients (missing `patient:write`)
- ❌ Create/update prescriptions (missing `prescription:write`)

---

## 📊 Quick Comparison Table

| Feature | Admin | Doctor | Pharmacist | Nurse | Read-Only |
|---------|-------|--------|------------|-------|-----------|
| View Patients | ✅ | ✅ | ✅ | ✅ | ✅ |
| Create/Edit Patients | ✅ | ✅ | ❌ | ❌ | ❌ |
| Export Patients | ✅ | ❌ | ❌ | ❌ | ❌ |
| Delete Patients | ✅ | ❌ | ❌ | ❌ | ❌ |
| View Prescriptions | ✅ | ✅ | ✅ | ✅ | ✅ |
| Create/Edit Prescriptions | ✅ | ✅ | ❌ | ❌ | ❌ |
| Approve Prescriptions | ✅ | ✅ | ❌ | ❌ | ❌ |
| Dispense Prescriptions | ✅ | ❌ | ✅ | ❌ | ❌ |
| View Dashboard | ✅ | ✅ | ✅ | ❌ | ✅ |

---

## 🧪 Testing Permission Scenarios

### Test Read Access
```bash
# All users should be able to read
for user in admin doctor pharmacist nurse readonly; do
  echo "Testing $user:"
  curl -H "X-Mock-User: $user" http://localhost:8080/patients
done
```

### Test Write Access
```bash
# Only admin and doctor should succeed
for user in admin doctor pharmacist nurse readonly; do
  echo "Testing $user:"
  curl -X POST -H "X-Mock-User: $user" \
    -H "Content-Type: application/json" \
    -d '{"name": "Test"}' \
    http://localhost:8080/api/v1/patients
done
```

### Test Permission Boundaries

**Doctor can approve but not dispense:**
```bash
# Should work (has prescription:approve)
curl -X POST -H "X-Mock-User: doctor" \
  http://localhost:8080/api/v1/prescriptions/123/approve

# Should fail 403 (missing prescription:dispense)
curl -X POST -H "X-Mock-User: doctor" \
  http://localhost:8080/api/v1/prescriptions/123/dispense
```

**Pharmacist can dispense but not approve:**
```bash
# Should work (has prescription:dispense)
curl -X POST -H "X-Mock-User: pharmacist" \
  http://localhost:8080/api/v1/prescriptions/123/dispense

# Should fail 403 (missing prescription:approve)
curl -X POST -H "X-Mock-User: pharmacist" \
  http://localhost:8080/api/v1/prescriptions/123/approve
```

---

## 🛠️ Adding Custom Mock Users

You can add custom mock users for specific test scenarios:

### In Code (for persistent custom users)

```go
// In main.go or test setup
package main

import "pharmacy-modernization-project-model/internal/platform/auth"

func main() {
    // ... app setup ...
    
    if auth.IsDevModeEnabled() {
        // Add a user with minimal permissions
        auth.AddMockUser("limited", &auth.User{
            ID:    "mock-limited-001",
            Email: "limited@dev.local",
            Name:  "Limited User",
            Permissions: []string{
                "patient:read",  // Only this permission
            },
        })
        
        // Add a user with specific permissions for edge case testing
        auth.AddMockUser("exporter", &auth.User{
            ID:    "mock-exporter-001",
            Email: "exporter@dev.local",
            Name:  "Data Exporter",
            Permissions: []string{
                "patient:read",
                "patient:export",
            },
        })
        
        // Add a user with no permissions (to test authorization failures)
        auth.AddMockUser("noperm", &auth.User{
            ID:    "mock-noperm-001",
            Email: "noperm@dev.local",
            Name:  "No Permissions",
            Permissions: []string{},
        })
    }
}
```

### Test Custom Users

```bash
# Check if custom users were added
curl http://localhost:8080/__dev/auth | jq '.mock_users[] | select(.key == "limited")'

# Test with custom user
curl -H "X-Mock-User: limited" http://localhost:8080/patients

# Test with no permissions (should get 403)
curl -H "X-Mock-User: noperm" http://localhost:8080/patients
```

---

## 📝 Permission Definitions

### Patient Permissions
- `patient:read` - View patient data
- `patient:write` - Create/update patients
- `patient:delete` - Delete patients
- `patient:export` - Export patient data

### Prescription Permissions
- `prescription:read` - View prescriptions
- `prescription:write` - Create/update prescriptions
- `prescription:approve` - Approve prescriptions
- `prescription:dispense` - Dispense prescriptions
- `prescription:cancel` - Cancel prescriptions

### Dashboard Permissions
- `dashboard:view` - View dashboard
- `dashboard:analytics` - View analytics
- `dashboard:reports` - View reports

### System Permissions
- `admin:all` - Full administrative access (grants all permissions)

### Role Permissions
- `doctor:role` - Doctor role identifier
- `pharmacist:role` - Pharmacist role identifier
- `nurse:role` - Nurse role identifier

---

## 🔧 Troubleshooting

### Mock Users Not Working

**Problem**: Getting 401 Unauthorized even with `X-Mock-User` header

**Solution**:
1. Check dev mode is enabled:
   ```yaml
   # internal/configs/app.yaml
   auth:
     dev_mode: true
   ```

2. Check the startup logs:
   ```
   ⚠️  AUTH DEV MODE ENABLED - Security bypassed with mock users
   ```

3. Verify the endpoint is available:
   ```bash
   curl http://localhost:8080/__dev/auth
   ```

### Getting 403 Forbidden

**Problem**: User is authenticated but getting 403

**Solution**: The user doesn't have the required permission
1. Check what permissions the route requires (see route code or docs)
2. Check what permissions the mock user has:
   ```bash
   curl http://localhost:8080/__dev/auth | jq '.mock_users[] | select(.key == "doctor")'
   ```
3. Use a different mock user with the required permissions

### Can't See Dev Endpoint

**Problem**: `/__dev/auth` returns 404

**Solution**: Dev mode is not enabled or endpoint not registered
1. Enable dev mode in config
2. Restart the application
3. Check logs for dev mode registration message

---

## 🎯 Best Practices

### 1. Use Appropriate User for Testing
- Use `admin` for testing happy paths (everything works)
- Use specific roles (`doctor`, `pharmacist`, `nurse`) for role-based testing
- Use `readonly` for testing read-only scenarios
- Create custom users for edge cases

### 2. Test Permission Boundaries
Always test:
- ✅ What the user CAN do (should succeed)
- ❌ What the user CANNOT do (should fail with 403)

### 3. Document Custom Users
If you add custom mock users, document them:
```go
// Custom mock user for testing export without read permission
auth.AddMockUser("exportonly", &auth.User{
    // ... 
    Permissions: []string{"patient:export"},  // Missing patient:read!
})
```

### 4. Use Dev Endpoint for Verification
Before testing, always verify:
```bash
curl http://localhost:8080/__dev/auth | jq '.mock_users[] | {key, permissions}'
```

---

## 📚 Related Documentation

- **[SECURITY_DEV_MODE.md](./SECURITY_DEV_MODE.md)** - Complete dev mode guide
- **[SECURITY_DEV_MODE_EXAMPLE.md](./SECURITY_DEV_MODE_EXAMPLE.md)** - Practical examples
- **[ROUTES_SECURITY_IMPLEMENTATION.md](../ROUTES_SECURITY_IMPLEMENTATION.md)** - Route security overview

---

## 🎉 Summary

**To check mock user permissions:**
1. **Easiest**: `curl http://localhost:8080/__dev/auth`
2. **Code**: Check `internal/platform/auth/dev_mode.go`
3. **Docs**: Refer to tables in this document

**Mock users available:**
- `admin` - Full access
- `doctor` - Patient + prescription management
- `pharmacist` - Read + dispense
- `nurse` - Read-only clinical
- `readonly` - Minimal read access

**Remember**: Dev mode is for development only - disable in production! 🔒

