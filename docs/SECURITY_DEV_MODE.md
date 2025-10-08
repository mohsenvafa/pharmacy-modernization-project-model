# Security Development Mode

## Overview

Development mode allows you to **bypass JWT authentication** and use **mock users** for local development and testing. This makes it easy to test different permission scenarios without needing real JWT tokens.

⚠️ **WARNING**: Dev mode should **NEVER** be enabled in production!

---

## Configuration

### Enable Dev Mode

In `internal/configs/app.yaml`:

```yaml
auth:
  dev_mode: true  # Enable dev mode
  jwt:
    # ... JWT config (still needed for when dev_mode is false)
```

### Disable Dev Mode (Production)

In `internal/configs/app.prod.yaml`:

```yaml
auth:
  dev_mode: false  # MUST be false in production
  jwt:
    secret: "${RX_AUTH_JWT_SECRET}"
    # ...
```

---

## Mock Users

When dev mode is enabled, you have access to 5 pre-configured mock users:

| Key | Email | Role | Permissions |
|-----|-------|------|-------------|
| `admin` | admin@dev.local | Super Admin | `admin:all` |
| `doctor` | doctor@dev.local | Doctor | `patient:read`, `patient:write`, `prescription:read`, `prescription:write`, `prescription:approve`, `doctor:role`, `dashboard:view` |
| `pharmacist` | pharmacist@dev.local | Pharmacist | `patient:read`, `prescription:read`, `prescription:dispense`, `pharmacist:role`, `dashboard:view` |
| `nurse` | nurse@dev.local | Nurse | `patient:read`, `prescription:read`, `nurse:role` |
| `readonly` | readonly@dev.local | Read Only | `patient:read`, `prescription:read`, `dashboard:view` |

---

## Usage

### In Code - Use Dev Mode Middleware

Update your routes to use `RequireAuthWithDevMode()` instead of `RequireAuth()`:

```go
// Before (production only)
r.Use(auth.RequireAuth())

// After (supports dev mode)
r.Use(auth.RequireAuthWithDevMode())
```

**The middleware automatically:**
- Uses mock users when `dev_mode: true`
- Uses real JWT validation when `dev_mode: false`

### Default User

By default, if no mock user is specified, the `admin` user is used.

### Select Mock User via Header

Use the `X-Mock-User` header to select a different mock user:

```bash
# Use doctor
curl -H "X-Mock-User: doctor" http://localhost:8080/patients

# Use pharmacist
curl -H "X-Mock-User: pharmacist" http://localhost:8080/prescriptions

# Use readonly
curl -H "X-Mock-User: readonly" http://localhost:8080/dashboard
```

### In Browser

Set the header in browser DevTools or use a browser extension:

```javascript
// In browser console
fetch('/patients', {
    headers: {
        'X-Mock-User': 'doctor'
    }
});
```

Or use **ModHeader** or similar browser extension to set `X-Mock-User` header.

---

## Dev Info Endpoint

When dev mode is enabled, a special endpoint shows available mock users:

```bash
curl http://localhost:8080/__dev/auth
```

**Response:**
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
      "permissions": ["patient:read", "patient:write", ...]
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

### Register Dev Info Endpoint

In your router setup:

```go
// Add dev info endpoint
if auth.IsDevModeEnabled() {
    r.Get("/__dev/auth", auth.DevAuthInfo)
}
```

---

## Custom Mock Users

You can add custom mock users for testing specific scenarios:

```go
// In your test setup or main.go
if auth.IsDevModeEnabled() {
    auth.AddMockUser("limited", &auth.User{
        ID:    "mock-limited-001",
        Email: "limited@dev.local",
        Name:  "Limited User",
        Permissions: []string{
            "patient:read",  // Can only read patients
        },
    })
    
    auth.AddMockUser("superpharmacist", &auth.User{
        ID:    "mock-superpharmacist-001",
        Email: "superpharmacist@dev.local",
        Name:  "Super Pharmacist",
        Permissions: []string{
            "patient:read",
            "patient:write",
            "prescription:read",
            "prescription:write",
            "prescription:approve",
            "prescription:dispense",
            "pharmacist:role",
        },
    })
}
```

Then use them:

```bash
curl -H "X-Mock-User: limited" http://localhost:8080/patients
curl -H "X-Mock-User: superpharmacist" http://localhost:8080/prescriptions
```

---

## Testing Permission Scenarios

### Test READ access

```bash
# Should work (has patient:read)
curl -H "X-Mock-User: nurse" http://localhost:8080/patients

# Should fail 403 (no patient:write)
curl -X POST -H "X-Mock-User: nurse" http://localhost:8080/patients
```

### Test WRITE access

```bash
# Should work (has patient:write)
curl -X POST -H "X-Mock-User: doctor" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Patient"}' \
  http://localhost:8080/patients
```

### Test MatchAll permissions

```bash
# Requires BOTH patient:read AND patient:export
# Doctor has patient:read but NOT patient:export → Should fail 403
curl -H "X-Mock-User: doctor" http://localhost:8080/patients/export

# Admin has admin:all → Should work
curl -H "X-Mock-User: admin" http://localhost:8080/patients/export
```

### Test MatchAny permissions

```bash
# Requires patient:read OR admin:all
# Nurse has patient:read → Should work
curl -H "X-Mock-User: nurse" http://localhost:8080/patients

# Admin has admin:all → Should work
curl -H "X-Mock-User: admin" http://localhost:8080/patients
```

---

## GraphQL Testing

Dev mode works with GraphQL too:

```bash
# Use doctor mock user
curl -X POST http://localhost:8080/graphql \
  -H "X-Mock-User: doctor" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ patients { id name } }"
  }'

# Use nurse (limited permissions)
curl -X POST http://localhost:8080/graphql \
  -H "X-Mock-User: nurse" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { createPatient(input: {name: \"Test\"}) { id } }"
  }'
# Should return 403 error
```

---

## Integration with HTMX/Templ

For browser-based testing, you can set a default mock user in your dev environment:

```go
// In development, add middleware to set default mock user
if auth.IsDevModeEnabled() {
    r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // If no mock user header, default to doctor for UI testing
            if r.Header.Get("X-Mock-User") == "" {
                r.Header.Set("X-Mock-User", "doctor")
            }
            next.ServeHTTP(w, r)
        })
    })
}
```

Or create a dev UI to switch users:

```html
<!-- Dev mode user switcher -->
<div id="dev-user-switcher" style="position: fixed; top: 10px; right: 10px; background: yellow; padding: 10px;">
    <strong>DEV MODE</strong>
    <select onchange="switchUser(this.value)">
        <option value="admin">Admin</option>
        <option value="doctor">Doctor</option>
        <option value="pharmacist">Pharmacist</option>
        <option value="nurse">Nurse</option>
        <option value="readonly">Read Only</option>
    </select>
</div>

<script>
function switchUser(user) {
    // Store in localStorage and reload
    localStorage.setItem('mockUser', user);
    location.reload();
}

// Apply stored mock user to all HTMX requests
document.body.addEventListener('htmx:configRequest', function(evt) {
    const mockUser = localStorage.getItem('mockUser') || 'doctor';
    evt.detail.headers['X-Mock-User'] = mockUser;
});
</script>
```

---

## Automated Testing

Use dev mode in your automated tests:

```go
func TestPatientAPI(t *testing.T) {
    // Enable dev mode
    auth.InitDevMode(true)
    
    // Create test user with specific permissions
    auth.AddMockUser("test", &auth.User{
        ID:    "test-001",
        Email: "test@test.com",
        Name:  "Test User",
        Permissions: []string{"patient:read"},
    })
    
    // Make request
    req := httptest.NewRequest("GET", "/patients", nil)
    req.Header.Set("X-Mock-User", "test")
    
    w := httptest.NewRecorder()
    handler.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}

func TestUnauthorized(t *testing.T) {
    auth.InitDevMode(true)
    
    // User with NO permissions
    auth.AddMockUser("noperm", &auth.User{
        ID:    "test-002",
        Email: "noperm@test.com",
        Name:  "No Permissions",
        Permissions: []string{},  // Empty permissions
    })
    
    req := httptest.NewRequest("GET", "/patients", nil)
    req.Header.Set("X-Mock-User", "noperm")
    
    w := httptest.NewRecorder()
    handler.ServeHTTP(w, req)
    
    // Should get 403 Forbidden
    assert.Equal(t, http.StatusForbidden, w.Code)
}
```

---

## Logging

When dev mode is enabled, auth middleware logs which mock user is being used:

```
⚠️  AUTH DEV MODE ENABLED - Security bypassed with mock users
DEV AUTH: Using mock user 'doctor' (doctor@dev.local) with permissions: [patient:read patient:write ...]
```

This helps you verify which user is being used in each request.

---

## Security Checks

### Prevent Production Use

Add a startup check to ensure dev mode is not accidentally enabled in production:

```go
if cfg.App.Env == "prod" && cfg.Auth.DevMode {
    log.Fatal("FATAL: Dev mode cannot be enabled in production environment")
}
```

### Environment Detection

```go
// In wire.go or main.go
if a.Cfg.Auth.DevMode && a.Cfg.App.Env == "prod" {
    logger.Base.Fatal("Dev mode is enabled in production - aborting")
}

if a.Cfg.Auth.DevMode {
    logger.Base.Warn("⚠️  AUTH DEV MODE ACTIVE - Do not use in production!")
}
```

---

## Best Practices

1. **Always disable in production**
   - Set `dev_mode: false` in `app.prod.yaml`
   - Add environment checks

2. **Use descriptive mock user names**
   - Makes logs easier to understand
   - Clearly indicates test vs real data

3. **Test permission boundaries**
   - Create users with minimal permissions
   - Verify 403 responses for unauthorized actions

4. **Document test scenarios**
   - Keep a list of test users and their purposes
   - Share with team for consistent testing

5. **Use in automated tests**
   - Faster than generating real JWTs
   - More predictable permissions

6. **Visual indicators in UI**
   - Show "DEV MODE" banner
   - Display current mock user

---

## Migration to Production

When moving to production:

1. ✅ Set `auth.dev_mode: false` in production config
2. ✅ Set up real JWT issuer (Auth0, Okta, etc.)
3. ✅ Configure `RX_AUTH_JWT_SECRET` environment variable
4. ✅ Remove any dev mode UI elements
5. ✅ Test with real JWT tokens
6. ✅ Verify dev mode is disabled on startup

---

## Summary

Dev mode provides:
- ✅ **No JWT tokens needed** for local development
- ✅ **5 pre-configured mock users** with different permission sets
- ✅ **Easy permission testing** via `X-Mock-User` header
- ✅ **Custom mock users** for specific test scenarios
- ✅ **Automatic switching** between dev and production
- ✅ **Clean logging** to see which user is active
- ✅ **Perfect for automated tests** - fast and predictable

**Remember**: Dev mode = Development ONLY! Never in production! 🚨

