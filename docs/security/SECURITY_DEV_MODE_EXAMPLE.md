# Development Mode - Practical Examples

This guide shows real-world examples of using dev mode in your RxIntake application.

---

## Quick Start

### 1. Enable Dev Mode

```yaml
# internal/configs/app.yaml
auth:
  dev_mode: true
```

### 2. Update Your Routes

```go
// domain/patient/ui/ui.go
func (ui *UI) RegisterRoutes(r chi.Router) {
    r.Route("/patients", func(r chi.Router) {
        // Use RequireAuthWithDevMode instead of RequireAuth
        r.Use(auth.RequireAuthWithDevMode())
        
        r.With(auth.RequirePermissionsMatchAny(
            patientsecurity.ReadAccess,
        )).Get("/", ui.listPatients)
    })
}
```

### 3. Start Your App

```bash
go run cmd/server/main.go
```

You'll see:
```
⚠️  AUTH DEV MODE ENABLED - Security bypassed with mock users
⚠️  AUTH DEV MODE ACTIVE - Do not use in production!
```

### 4. Test with Mock Users

```bash
# Use doctor (default if no header)
curl http://localhost:8080/patients

# Explicitly use pharmacist
curl -H "X-Mock-User: pharmacist" http://localhost:8080/patients

# Use nurse (limited permissions)
curl -H "X-Mock-User: nurse" http://localhost:8080/patients
```

---

## Example: Testing Patient Routes

### Route Setup

```go
// domain/patient/ui/ui.go
package ui

import (
    "github.com/go-chi/chi/v5"
    "internal/platform/auth"
    patientsecurity "domain/patient/security"
)

func (ui *UI) RegisterRoutes(r chi.Router) {
    r.Route("/patients", func(r chi.Router) {
        // Dev mode enabled automatically if configured
        r.Use(auth.RequireAuthWithDevMode())
        
        // List - needs patient:read OR admin:all
        r.With(auth.RequirePermissionsMatchAny(
            patientsecurity.ReadAccess,
        )).Get("/", ui.listPatients)
        
        // Create - needs patient:write OR admin:all
        r.With(auth.RequirePermissionsMatchAny(
            patientsecurity.WriteAccess,
        )).Post("/", ui.createPatient)
        
        // Export - needs BOTH patient:read AND patient:export
        r.With(auth.RequirePermissionsMatchAll(
            patientsecurity.ExportAccess,
        )).Get("/export", ui.exportPatients)
    })
}
```

### Testing Different Users

```bash
# ✅ Doctor can list (has patient:read)
curl -H "X-Mock-User: doctor" http://localhost:8080/patients
# Returns: 200 OK

# ✅ Nurse can list (has patient:read)
curl -H "X-Mock-User: nurse" http://localhost:8080/patients
# Returns: 200 OK

# ✅ Admin can list (has admin:all)
curl -H "X-Mock-User: admin" http://localhost:8080/patients
# Returns: 200 OK

# ✅ Doctor can create (has patient:write)
curl -X POST -H "X-Mock-User: doctor" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Patient"}' \
  http://localhost:8080/patients
# Returns: 200 OK

# ❌ Nurse CANNOT create (no patient:write)
curl -X POST -H "X-Mock-User: nurse" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Patient"}' \
  http://localhost:8080/patients
# Returns: 403 Forbidden

# ❌ Doctor CANNOT export (has patient:read but NOT patient:export)
curl -H "X-Mock-User: doctor" http://localhost:8080/patients/export
# Returns: 403 Forbidden

# ✅ Admin CAN export (has admin:all)
curl -H "X-Mock-User: admin" http://localhost:8080/patients/export
# Returns: 200 OK
```

---

## Example: Testing Prescription Routes

### Route Setup

```go
// domain/prescription/ui/ui.go
package ui

import (
    "github.com/go-chi/chi/v5"
    "internal/platform/auth"
    prescripsecurity "domain/prescription/security"
)

func (ui *UI) RegisterRoutes(r chi.Router) {
    r.Route("/prescriptions", func(r chi.Router) {
        r.Use(auth.RequireAuthWithDevMode())
        
        // Read - doctors, pharmacists, nurses, or admin
        r.With(auth.RequirePermissionsMatchAny(
            prescripsecurity.ReadAccess,
        )).Get("/", ui.listPrescriptions)
        
        // Create - only doctors or admin
        r.With(auth.RequirePermissionsMatchAny(
            prescripsecurity.WriteAccess,
        )).Post("/", ui.createPrescription)
        
        // Approve - needs BOTH prescription:write AND prescription:approve
        r.With(auth.RequirePermissionsMatchAll(
            prescripsecurity.ApproveAccess,
        )).Post("/{id}/approve", ui.approvePrescription)
        
        // Dispense - only pharmacists or admin
        r.With(auth.RequirePermissionsMatchAny(
            prescripsecurity.DispenseAccess,
        )).Post("/{id}/dispense", ui.dispensePrescription)
    })
}
```

### Testing Scenarios

```bash
# ✅ All healthcare roles can read
curl -H "X-Mock-User: doctor" http://localhost:8080/prescriptions
curl -H "X-Mock-User: pharmacist" http://localhost:8080/prescriptions
curl -H "X-Mock-User: nurse" http://localhost:8080/prescriptions

# ✅ Doctor can create
curl -X POST -H "X-Mock-User: doctor" http://localhost:8080/prescriptions

# ❌ Pharmacist CANNOT create
curl -X POST -H "X-Mock-User: pharmacist" http://localhost:8080/prescriptions
# Returns: 403 Forbidden

# ❌ Nurse CANNOT create
curl -X POST -H "X-Mock-User: nurse" http://localhost:8080/prescriptions
# Returns: 403 Forbidden

# ✅ Doctor can approve (has both prescription:write AND prescription:approve)
curl -X POST -H "X-Mock-User: doctor" \
  http://localhost:8080/prescriptions/123/approve

# ✅ Pharmacist can dispense
curl -X POST -H "X-Mock-User: pharmacist" \
  http://localhost:8080/prescriptions/123/dispense

# ❌ Doctor CANNOT dispense (no prescription:dispense)
curl -X POST -H "X-Mock-User: doctor" \
  http://localhost:8080/prescriptions/123/dispense
# Returns: 403 Forbidden
```

---

## Example: GraphQL with Dev Mode

### Schema

```graphql
# domain/patient/graphql/schema.graphql

extend type Query {
    patient(id: ID!): Patient 
        @auth 
        @permissionAny(requires: ["patient:read", "admin:all"])
    
    patients: [Patient!]! 
        @auth 
        @permissionAny(requires: ["patient:read", "admin:all"])
}

extend type Mutation {
    createPatient(input: CreatePatientInput!): Patient! 
        @auth 
        @permissionAny(requires: ["patient:write", "admin:all"])
}
```

### Wrap GraphQL Endpoint

```go
// internal/graphql/server.go

// Wrap with dev mode support
router.Handle("/graphql", auth.RequireAuthWithDevMode()(graphqlHandler))
```

### Testing

```bash
# ✅ Doctor can query
curl -X POST http://localhost:8080/graphql \
  -H "X-Mock-User: doctor" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ patients { id name email } }"
  }'

# ✅ Nurse can query
curl -X POST http://localhost:8080/graphql \
  -H "X-Mock-User: nurse" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ patients { id name } }"
  }'

# ✅ Doctor can create
curl -X POST http://localhost:8080/graphql \
  -H "X-Mock-User: doctor" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { createPatient(input: {name: \"Test\"}) { id } }"
  }'

# ❌ Nurse CANNOT create (no patient:write)
curl -X POST http://localhost:8080/graphql \
  -H "X-Mock-User: nurse" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { createPatient(input: {name: \"Test\"}) { id } }"
  }'
# Returns GraphQL error with status 403
```

---

## Example: Custom Mock Users for Edge Cases

### Create Custom Users

```go
// In main.go or test setup
func main() {
    // ... app initialization
    
    if auth.IsDevModeEnabled() {
        // Add user with only export permission (edge case)
        auth.AddMockUser("exporter", &auth.User{
            ID:    "mock-exporter-001",
            Email: "exporter@dev.local",
            Name:  "Dev Exporter",
            Permissions: []string{
                "patient:export",  // Has export but NOT read
            },
        })
        
        // Add user with no permissions (testing authorization failures)
        auth.AddMockUser("noperm", &auth.User{
            ID:    "mock-noperm-001",
            Email: "noperm@dev.local",
            Name:  "No Permissions",
            Permissions: []string{},
        })
    }
    
    // ... rest of setup
}
```

### Test Edge Cases

```bash
# ❌ Exporter has patient:export but NOT patient:read
# Export requires BOTH patient:read AND patient:export
curl -H "X-Mock-User: exporter" http://localhost:8080/patients/export
# Returns: 403 Forbidden (needs ALL permissions)

# ❌ User with no permissions
curl -H "X-Mock-User: noperm" http://localhost:8080/patients
# Returns: 403 Forbidden
```

---

## Example: Browser Testing with Dev Mode

### Create User Switcher UI

```html
<!-- Add to your base layout in dev mode -->
{{if .DevMode}}
<div style="position: fixed; top: 10px; right: 10px; 
            background: #ffd700; padding: 10px; 
            border: 2px solid #333; z-index: 9999;">
    <strong>🔧 DEV MODE</strong><br>
    <select id="mock-user-select" onchange="switchMockUser(this.value)">
        <option value="admin">Admin</option>
        <option value="doctor" selected>Doctor</option>
        <option value="pharmacist">Pharmacist</option>
        <option value="nurse">Nurse</option>
        <option value="readonly">Read Only</option>
    </select>
</div>

<script>
// Store selected user
const currentUser = localStorage.getItem('mockUser') || 'doctor';
document.getElementById('mock-user-select').value = currentUser;

function switchMockUser(user) {
    localStorage.setItem('mockUser', user);
    location.reload();
}

// Add header to all HTMX requests
document.body.addEventListener('htmx:configRequest', function(evt) {
    const mockUser = localStorage.getItem('mockUser') || 'doctor';
    evt.detail.headers['X-Mock-User'] = mockUser;
});
</script>
{{end}}
```

### Usage

1. Open your app in browser
2. See dev mode switcher in top-right
3. Select different users to test permissions
4. HTMX requests automatically use selected user

---

## Example: Automated Tests

### Test Helper

```go
// internal/platform/auth/testing.go (optional helper file)
package auth

func WithMockUser(t *testing.T, userKey string, permissions []string) {
    if !IsDevModeEnabled() {
        InitDevMode(true)
    }
    
    AddMockUser(userKey, &User{
        ID:          "test-" + userKey,
        Email:       userKey + "@test.local",
        Name:        "Test " + userKey,
        Permissions: permissions,
    })
}
```

### Integration Tests

```go
// domain/patient/ui/ui_test.go
package ui

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "internal/platform/auth"
)

func TestListPatients_WithDoctor(t *testing.T) {
    // Enable dev mode
    auth.InitDevMode(true)
    
    // Create request with doctor
    req := httptest.NewRequest("GET", "/patients", nil)
    req.Header.Set("X-Mock-User", "doctor")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Doctor should be able to list patients
    assert.Equal(t, http.StatusOK, w.Code)
}

func TestListPatients_WithReadonly(t *testing.T) {
    auth.InitDevMode(true)
    
    req := httptest.NewRequest("GET", "/patients", nil)
    req.Header.Set("X-Mock-User", "readonly")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Readonly should also be able to list patients
    assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreatePatient_WithNurse(t *testing.T) {
    auth.InitDevMode(true)
    
    req := httptest.NewRequest("POST", "/patients", 
        strings.NewReader(`{"name": "Test"}`))
    req.Header.Set("X-Mock-User", "nurse")
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Nurse should NOT be able to create patients
    assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestExportPatients_EdgeCase(t *testing.T) {
    auth.InitDevMode(true)
    
    // Create user with only export permission (not read)
    auth.AddMockUser("exportonly", &auth.User{
        ID:          "test-export",
        Email:       "export@test.local",
        Name:        "Export Only",
        Permissions: []string{"patient:export"}, // Missing patient:read!
    })
    
    req := httptest.NewRequest("GET", "/patients/export", nil)
    req.Header.Set("X-Mock-User", "exportonly")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Should fail - needs BOTH patient:read AND patient:export
    assert.Equal(t, http.StatusForbidden, w.Code)
}
```

---

## Tips & Best Practices

### 1. Always Check Dev Mode Status

```go
if auth.IsDevModeEnabled() {
    // Add dev-specific features
    router.Get("/__dev/auth", auth.DevAuthInfo)
}
```

### 2. Visual Indicators

Show dev mode status prominently in UI:
- Banner: "Development Mode - Mock Users Active"
- Different color scheme
- User switcher widget

### 3. Log Mock User

Check logs to see which user is active:
```
DEV AUTH: Using mock user 'doctor' (doctor@dev.local) with permissions: [patient:read patient:write ...]
```

### 4. Test Permission Boundaries

Create users that test edge cases:
- User with minimal permissions
- User with no permissions
- User with only one of required permissions (for MatchAll)

### 5. Document Your Test Scenarios

Keep a list of test scenarios for your team:
```markdown
## Patient Management Test Scenarios

1. List patients:
   - ✅ doctor, nurse, pharmacist, admin
   - ❌ readonly (depends on config)

2. Create patient:
   - ✅ doctor, admin
   - ❌ nurse, pharmacist, readonly

3. Export patients:
   - ✅ admin (has admin:all)
   - ❌ doctor (missing patient:export)
```

---

## Switching Between Dev and Production

### Local Development
```yaml
# app.yaml
auth:
  dev_mode: true
```

### Production
```yaml
# app.prod.yaml
auth:
  dev_mode: false
  jwt:
    secret: "${RX_AUTH_JWT_SECRET}"
```

### Run with Production Config
```bash
RX_APP_ENV=prod go run cmd/server/main.go
```

The app will automatically:
- Use real JWT validation
- Reject requests without valid tokens
- Log: "AUTH 401: No token"

---

## Summary

Dev mode makes local development easy:
- ✅ No need to generate JWT tokens
- ✅ Quick permission testing
- ✅ 5 ready-to-use mock users
- ✅ Custom users for edge cases
- ✅ Works with HTTP, GraphQL, and HTMX
- ✅ Perfect for automated tests
- ✅ Automatic safety checks for production

Just remember: **Dev mode = Development ONLY!** 🚨

