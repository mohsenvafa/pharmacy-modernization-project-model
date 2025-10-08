# Security Implementation - Overview

This document provides a quick overview of the authentication and authorization system implemented in RxIntake.

## 📚 Documentation

- **[SECURITY_ARCHITECTURE.md](./SECURITY_ARCHITECTURE.md)** - Complete architecture details with diagrams
- **[SECURITY_QUICK_START.md](./SECURITY_QUICK_START.md)** - Step-by-step implementation guide

## 🎯 What's Implemented

### ✅ Core Authentication & Authorization

- **JWT-based authentication** - Stateless, scalable
- **Dual token support** - HTTP cookies (web) + Authorization headers (API)
- **Permission-based authorization** - Simple string permissions (e.g., `patient:read`)
- **Declarative security** - Define auth at route/field level (like .NET `[Authorize]`)
- **Proper HTTP status codes** - 401 for auth failures, 403 for permission failures

### ✅ Components Created

```
internal/platform/auth/
├── models.go           ← User model, JWT claims, token source types
├── jwt.go              ← JWT validation, parsing, token extraction
├── context.go          ← User context helpers
├── permissions.go      ← Permission checking logic (AND/OR)
├── middleware.go       ← HTTP middleware functions
└── directives.go       ← GraphQL directives (@auth, @permissionAny, @permissionAll)

domain/*/security/
├── patient/security/permissions.go        ← Patient permissions
├── prescription/security/permissions.go   ← Prescription permissions
└── dashboard/security/permissions.go      ← Dashboard permissions

internal/configs/
└── app.yaml            ← JWT configuration added
```

## 🚀 Quick Usage

### HTTP Routes

```go
import (
    "internal/platform/auth"
    patientsecurity "domain/patient/security"
)

// Web UI (cookie-based)
r.Use(auth.RequireAuthFromCookie())
r.With(auth.RequirePermissionsMatchAny(
    patientsecurity.ReadAccess,
)).Get("/patients", handler)

// REST API (header-based)
r.Use(auth.RequireAuthFromHeader())
r.With(auth.RequirePermissionsMatchAll(
    patientsecurity.ExportAccess,
)).Get("/patients/export", handler)
```

### GraphQL Schema

```graphql
type Query {
    patient(id: ID!): Patient 
        @auth 
        @permissionAny(requires: ["patient:read", "admin:all"])
}
```

### GraphQL Resolvers

```go
// No manual permission checks needed!
func (r *queryResolver) Patient(ctx context.Context, id string) (*model.Patient, error) {
    return r.PatientService.GetByID(ctx, id)
}
```

## 🔧 Middleware Functions

### Authentication (Returns 401 on failure)

- `RequireAuth()` - Auto-detect token source (header or cookie)
- `RequireAuthFromHeader()` - Only accept Authorization header
- `RequireAuthFromCookie()` - Only accept cookie

### Authorization (Returns 403 on failure)

- `RequirePermission(permission)` - Single permission
- `RequirePermissionsMatchAll(permissions)` - User needs ALL permissions
- `RequirePermissionsMatchAny(permissions)` - User needs ANY permission

## 📋 Permission Format

Permissions follow the `<resource>:<action>` pattern:

```
patient:read
patient:write
patient:delete
prescription:approve
prescription:dispense
admin:all
```

Each domain defines its permissions in `domain/*/security/permissions.go`

## 🎨 Permission Patterns

### ANY - Hierarchical Access
```go
// Admin OR specific permission
var ReadAccess = []string{"patient:read", "admin:all"}

r.With(auth.RequirePermissionsMatchAny(ReadAccess))
```

### ALL - Multi-Stage Requirements
```go
// Needs BOTH permissions
var ExportAccess = []string{"patient:read", "patient:export"}

r.With(auth.RequirePermissionsMatchAll(ExportAccess))
```

## 🔑 JWT Token Structure

```json
{
  "user_id": "12345",
  "email": "doctor@example.com",
  "name": "Dr. Smith",
  "permissions": [
    "patient:read",
    "patient:write",
    "prescription:read",
    "doctor:role"
  ],
  "iss": "rxintake",
  "aud": "rxintake",
  "exp": 1728691200
}
```

## 🛡️ Security Features

- ✅ HTTP-only cookies (XSS protection)
- ✅ Secure flag for HTTPS
- ✅ SameSite protection (CSRF protection)
- ✅ Short token expiration (1 hour default)
- ✅ Signature validation
- ✅ Issuer/audience validation
- ✅ Comprehensive logging

## 📊 Error Responses

**401 Unauthorized** (Authentication Failure)
```json
{
  "error": "unauthorized",
  "message": "Invalid or expired token",
  "status": 401
}
```

**403 Forbidden** (Permission Failure)
```json
{
  "error": "forbidden",
  "message": "User requires at least one of the following permissions",
  "required_permissions": ["patient:read", "admin:all"],
  "match": "any",
  "status": 403
}
```

## ⚙️ Configuration

```yaml
# internal/configs/app.yaml
auth:
  jwt:
    secret: "dev-secret-key-change-in-production"
    issuer: "rxintake"
    audience: "rxintake"
    cookie:
      name: "auth_token"
      secure: false  # true in production
      httponly: true
      max_age: 3600
```

## 🧪 Testing

```bash
# Test with curl
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/v1/patients

# Create test token in Go
token, _ := auth.CreateToken(&auth.User{
    ID: "test",
    Email: "test@example.com",
    Permissions: []string{"patient:read"},
}, 1)
```

## 📖 Next Steps

1. **Read** [SECURITY_ARCHITECTURE.md](./SECURITY_ARCHITECTURE.md) for complete details
2. **Follow** [SECURITY_QUICK_START.md](./SECURITY_QUICK_START.md) for implementation
3. **Initialize** JWT config in your app
4. **Apply** middleware to your routes
5. **Add** directives to GraphQL schemas
6. **Test** with different permissions

## 🎯 Key Design Principles

1. **Simple** - No over-engineering, easy to understand
2. **Declarative** - Auth defined at route/field level, not in handlers
3. **Flexible** - Works for web browsers and API clients
4. **Secure** - Industry-standard JWT practices
5. **Maintainable** - Domain-owned permissions, no hardcoded strings
6. **Type-safe** - Permission constants prevent typos

## 🤝 Benefits

- ✅ No session management needed
- ✅ Scales horizontally
- ✅ Works with external auth services
- ✅ Clear separation of concerns
- ✅ Easy to test
- ✅ Production-ready

---

For detailed information, see:
- [SECURITY_ARCHITECTURE.md](./SECURITY_ARCHITECTURE.md) - Full architecture
- [SECURITY_QUICK_START.md](./SECURITY_QUICK_START.md) - Implementation guide

