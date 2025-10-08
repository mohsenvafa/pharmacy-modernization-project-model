# Development Mode Implementation Summary

## ✅ Complete!

A simple, secure development mode has been added to bypass JWT authentication and use mock users for local testing.

---

## 🎯 What Was Added

### 1. Dev Mode Implementation

**File**: `internal/platform/auth/dev_mode.go`

**Features**:
- ✅ Dev mode toggle (enabled/disabled)
- ✅ 5 pre-configured mock users (admin, doctor, pharmacist, nurse, readonly)
- ✅ Custom mock user support (`AddMockUser()`)
- ✅ Mock user selection via `X-Mock-User` header
- ✅ Dev info endpoint (`/__dev/auth`)
- ✅ `RequireAuthWithDevMode()` middleware (auto-switches between mock and real JWT)
- ✅ Comprehensive logging

**Mock Users**:
| Key | Role | Permissions |
|-----|------|-------------|
| `admin` | Super Admin | `admin:all` |
| `doctor` | Doctor | `patient:*`, `prescription:*`, `doctor:role`, `dashboard:view` |
| `pharmacist` | Pharmacist | `patient:read`, `prescription:read/dispense`, `pharmacist:role`, `dashboard:view` |
| `nurse` | Nurse | `patient:read`, `prescription:read`, `nurse:role` |
| `readonly` | Read Only | `patient:read`, `prescription:read`, `dashboard:view` |

### 2. Configuration Updates

**Local Development** (`app.yaml`):
```yaml
auth:
  dev_mode: true  # Enable mock users
```

**Production** (`app.prod.yaml`):
```yaml
auth:
  dev_mode: false  # MUST be false
  jwt:
    secret: "${RX_AUTH_JWT_SECRET}"
```

### 3. Safety Checks

Added to `internal/app/wire.go`:
- ✅ Fatal error if dev mode enabled in production environment
- ✅ Warning log when dev mode is active
- ✅ Clear startup messages

### 4. Documentation

- ✅ **[SECURITY_DEV_MODE.md](./docs/SECURITY_DEV_MODE.md)** (456 lines)
  - Complete guide to dev mode
  - Configuration instructions
  - Mock user reference
  - Testing scenarios

- ✅ **[SECURITY_DEV_MODE_EXAMPLE.md](./docs/SECURITY_DEV_MODE_EXAMPLE.md)** (545 lines)
  - Practical examples
  - Route testing scenarios
  - GraphQL testing
  - Custom mock users
  - Browser integration
  - Automated testing

---

## 🚀 How to Use

### Enable Dev Mode

```yaml
# internal/configs/app.yaml
auth:
  dev_mode: true
```

### Update Routes to Support Dev Mode

**Before**:
```go
r.Use(auth.RequireAuth())
```

**After**:
```go
r.Use(auth.RequireAuthWithDevMode())  // Auto-switches between dev and production
```

### Test with Mock Users

```bash
# Default (uses admin)
curl http://localhost:8080/patients

# Specify user via header
curl -H "X-Mock-User: doctor" http://localhost:8080/patients
curl -H "X-Mock-User: pharmacist" http://localhost:8080/prescriptions
curl -H "X-Mock-User: nurse" http://localhost:8080/dashboard

# See available mock users
curl http://localhost:8080/__dev/auth
```

---

## 📊 Dev Mode vs Production

| Feature | Dev Mode | Production |
|---------|----------|------------|
| Auth Method | Mock users | Real JWT |
| Token Required | No | Yes |
| User Selection | `X-Mock-User` header | JWT claims |
| Permissions | Pre-defined sets | From JWT token |
| Startup Log | ⚠️ Warning | Normal |
| Dev Info Endpoint | Available | Disabled |

---

## 🧪 Testing Scenarios

### Test Permission Success

```bash
# ✅ Doctor can read patients (has patient:read)
curl -H "X-Mock-User: doctor" http://localhost:8080/patients

# ✅ Admin can do anything (has admin:all)
curl -H "X-Mock-User: admin" http://localhost:8080/patients/export
```

### Test Permission Failures

```bash
# ❌ Nurse cannot create patients (missing patient:write)
curl -X POST -H "X-Mock-User: nurse" http://localhost:8080/patients
# Returns: 403 Forbidden

# ❌ Doctor cannot export (has patient:read but NOT patient:export)
curl -H "X-Mock-User: doctor" http://localhost:8080/patients/export
# Returns: 403 Forbidden (needs BOTH permissions)
```

### Test MatchAll vs MatchAny

```bash
# MatchAny - needs ANY of [patient:read, admin:all]
curl -H "X-Mock-User: nurse" http://localhost:8080/patients
# ✅ Works (has patient:read)

# MatchAll - needs ALL of [patient:read, patient:export]
curl -H "X-Mock-User: doctor" http://localhost:8080/patients/export
# ❌ Fails (has patient:read but NOT patient:export)

curl -H "X-Mock-User: admin" http://localhost:8080/patients/export
# ✅ Works (admin:all grants everything)
```

---

## 🔒 Security Features

### 1. Production Safety

```go
// Fatal error if dev mode enabled in prod
if cfg.App.Env == "prod" && cfg.Auth.DevMode {
    log.Fatal("FATAL: Dev mode cannot be enabled in production")
}
```

### 2. Clear Logging

```
⚠️  AUTH DEV MODE ENABLED - Security bypassed with mock users
DEV AUTH: Using mock user 'doctor' (doctor@dev.local) with permissions: [patient:read patient:write ...]
```

### 3. Automatic Switching

`RequireAuthWithDevMode()` automatically:
- Uses mock users when `dev_mode: true`
- Uses real JWT when `dev_mode: false`
- No code changes needed

---

## 🎨 Custom Mock Users

### Add Custom Users for Testing

```go
if auth.IsDevModeEnabled() {
    // User with minimal permissions
    auth.AddMockUser("limited", &auth.User{
        ID:    "mock-limited-001",
        Email: "limited@dev.local",
        Name:  "Limited User",
        Permissions: []string{"patient:read"},
    })
    
    // User with no permissions
    auth.AddMockUser("noperm", &auth.User{
        ID:    "mock-noperm-001",
        Email: "noperm@dev.local",
        Name:  "No Permissions",
        Permissions: []string{},
    })
}
```

### Use in Tests

```bash
curl -H "X-Mock-User: limited" http://localhost:8080/patients
curl -H "X-Mock-User: noperm" http://localhost:8080/patients  # 403
```

---

## 📱 Integration Examples

### HTTP Routes

```go
func (ui *UI) RegisterRoutes(r chi.Router) {
    r.Route("/patients", func(r chi.Router) {
        r.Use(auth.RequireAuthWithDevMode())  // ← Dev mode support
        
        r.With(auth.RequirePermissionsMatchAny(
            patientsecurity.ReadAccess,
        )).Get("/", ui.listPatients)
    })
}
```

### GraphQL

```go
// Wrap GraphQL endpoint
router.Handle("/graphql", auth.RequireAuthWithDevMode()(graphqlHandler))
```

```bash
# Test with mock user
curl -X POST http://localhost:8080/graphql \
  -H "X-Mock-User: doctor" \
  -H "Content-Type: application/json" \
  -d '{"query": "{ patients { id name } }"}'
```

### Automated Tests

```go
func TestPatientAccess(t *testing.T) {
    auth.InitDevMode(true)
    
    req := httptest.NewRequest("GET", "/patients", nil)
    req.Header.Set("X-Mock-User", "nurse")
    
    w := httptest.NewRecorder()
    handler.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
```

---

## 🎯 Benefits

### For Developers
- ✅ **No JWT tokens needed** for local development
- ✅ **Quick permission testing** - just change header
- ✅ **5 ready-to-use personas** covering common scenarios
- ✅ **Add custom users** for edge cases
- ✅ **Clear logging** shows which user is active
- ✅ **Fast iteration** - no token generation delays

### For Testing
- ✅ **Predictable** - same mock users every time
- ✅ **Fast** - no external auth service calls
- ✅ **Flexible** - create users for specific test scenarios
- ✅ **Easy to debug** - clear error messages

### For Production
- ✅ **Safe** - fatal error prevents accidental production use
- ✅ **Clean switch** - single config flag
- ✅ **No code changes** - works with same middleware
- ✅ **Automatic** - transparently uses real JWT when disabled

---

## 📚 Documentation Reference

1. **[SECURITY_DEV_MODE.md](./docs/SECURITY_DEV_MODE.md)**
   - Configuration guide
   - Mock user reference
   - Usage patterns
   - Best practices

2. **[SECURITY_DEV_MODE_EXAMPLE.md](./docs/SECURITY_DEV_MODE_EXAMPLE.md)**
   - Practical examples
   - Route testing
   - GraphQL testing
   - Custom users
   - Browser integration

3. **[SECURITY_IMPLEMENTATION_SUMMARY.md](./SECURITY_IMPLEMENTATION_SUMMARY.md)**
   - Complete security overview
   - Updated with dev mode info

---

## ✅ Implementation Checklist

When using dev mode:

- [x] Set `auth.dev_mode: true` in `app.yaml`
- [x] Update routes to use `RequireAuthWithDevMode()`
- [x] (Optional) Add `/__dev/auth` endpoint
- [x] (Optional) Create custom mock users for your scenarios
- [x] Test with different mock users
- [x] Verify 403 responses for unauthorized actions

When going to production:

- [ ] Set `auth.dev_mode: false` in `app.prod.yaml`
- [ ] Set `RX_AUTH_JWT_SECRET` environment variable
- [ ] Verify startup logs show no dev mode warning
- [ ] Test with real JWT tokens
- [ ] Confirm `/__dev/auth` returns 404

---

## 🎉 Summary

Dev mode provides a **simple, secure way** to develop and test your application locally without dealing with JWT tokens:

- **5 mock users** cover common scenarios
- **Custom users** for edge cases
- **One config flag** to enable/disable
- **Automatic switching** via `RequireAuthWithDevMode()`
- **Safety checks** prevent production accidents
- **Perfect for testing** - fast and predictable
- **Works everywhere** - HTTP, GraphQL, HTMX
- **Well documented** - clear examples and guides

**Start using it now:**
```yaml
# app.yaml
auth:
  dev_mode: true
```

```bash
curl -H "X-Mock-User: doctor" http://localhost:8080/patients
```

That's it! No JWT tokens, no external services, just simple mock users for rapid development! 🚀

