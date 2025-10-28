# Security Architecture

## Overview

RxIntake uses JWT-based authentication and permission-based authorization. The system is designed to be simple, secure, and declarative - similar to .NET's `[Authorize]` attribute pattern.

## Table of Contents

- [Key Principles](#key-principles)
- [Architecture Diagrams](#architecture-diagrams)
- [Components](#components)
- [Authentication Flow](#authentication-flow)
- [Authorization Flow](#authorization-flow)
- [Permission Model](#permission-model)
- [Usage Examples](#usage-examples)
- [Security Best Practices](#security-best-practices)

---

## Key Principles

1. **JWT-Based Authentication** - Stateless, scalable, no session storage required
2. **Permission-Based Authorization** - Simple string-based permissions (e.g., `patient:read`)
3. **Declarative Security** - Auth requirements defined at route/field level, not in handlers
4. **Dual Token Support** - Works with both HTTP cookies (web) and Authorization headers (API)
5. **Clear Error Codes** - 401 for auth failures, 403 for permission failures
6. **No User Domain** - Auth handled externally; JWT contains all needed info

---

## Architecture Diagrams

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     External Auth Service                       │
│                  (Issues JWT with permissions)                  │
└─────────────────────┬───────────────────────────────────────────┘
                      │
                      │ JWT Token
                      ↓
         ┌────────────────────────────┐
         │   RxIntake Application     │
         │                            │
         │  ┌──────────────────────┐  │
         │  │  Auth Middleware     │  │
         │  │  - Validates JWT     │  │
         │  │  - Extracts User     │  │
         │  │  - Sets in Context   │  │
         │  └──────────┬───────────┘  │
         │             │              │
         │             ↓              │
         │  ┌──────────────────────┐  │
         │  │ Permission Middleware│  │
         │  │  - Checks Permissions│  │
         │  │  - 403 if denied     │  │
         │  └──────────┬───────────┘  │
         │             │              │
         │             ↓              │
         │  ┌──────────────────────┐  │
         │  │   Handler/Resolver   │  │
         │  │  (Business Logic)    │  │
         │  └──────────────────────┘  │
         └────────────────────────────┘
```

### Request Flow

```
┌──────────────┐
│   Client     │
│ (Browser/API)│
└──────┬───────┘
       │
       │ 1. Request with JWT
       │    (Cookie or Auth Header)
       ↓
┌──────────────────────┐
│ RequireAuth()        │
│ Middleware           │
└──────┬───────────────┘
       │
       │ 2. Extract & Validate JWT
       │
       ├─── Invalid/Missing → 401 Unauthorized
       │
       │ 3. Parse Claims → User Object
       │
       │ 4. Set User in Context
       ↓
┌──────────────────────┐
│ RequirePermission()  │
│ Middleware           │
└──────┬───────────────┘
       │
       │ 5. Get User from Context
       │
       │ 6. Check Permissions
       │
       ├─── Missing Permissions → 403 Forbidden
       │
       │ 7. Permissions OK
       ↓
┌──────────────────────┐
│ Handler/Resolver     │
│ (Business Logic)     │
└──────────────────────┘
```

### Component Structure

```
rxintake_scaffold/
│
├── internal/platform/auth/          ← Core auth platform
│   ├── models.go                    ← User, JWTClaims, TokenSource
│   ├── jwt.go                       ← JWT validation & parsing
│   ├── context.go                   ← User context helpers
│   ├── permissions.go               ← Permission checking logic
│   ├── middleware.go                ← HTTP middleware
│   └── directives.go                ← GraphQL directives
│
├── domain/*/security/               ← Domain-specific permissions
│   └── permissions.go               ← Permission constants & sets
│
└── internal/configs/
    └── app.yaml                     ← JWT configuration
```

---

## Components

### 1. JWT Handler (`jwt.go`)

**Responsibilities:**
- Validate JWT signature
- Parse claims
- Extract token from request (header or cookie)
- Create tokens (for testing/login)

**Key Functions:**
```go
ValidateToken(tokenString string) (*User, error)
ExtractToken(r *http.Request, source TokenSource) (string, error)
```

### 2. Context Helper (`context.go`)

**Responsibilities:**
- Store/retrieve user in request context
- Provide type-safe access to current user

**Key Functions:**
```go
SetUser(ctx context.Context, user *User) context.Context
GetCurrentUser(ctx context.Context) (*User, error)
```

### 3. Permission Checker (`permissions.go`)

**Responsibilities:**
- Check if user has specific permissions
- Support AND/OR permission logic

**Key Functions:**
```go
HasAllPermissions(userPerms, required []string) bool
HasAnyPermission(userPerms, required []string) bool
HasAllPermissionsCtx(ctx context.Context, required []string) bool
HasAnyPermissionCtx(ctx context.Context, required []string) bool
```

### 4. HTTP Middleware (`middleware.go`)

**Responsibilities:**
- Validate authentication (401 on failure)
- Check permissions (403 on failure)
- Support different token sources

**Key Functions:**
```go
RequireAuth() func(http.Handler) http.Handler
RequireAuthFromHeader() func(http.Handler) http.Handler
RequireAuthFromCookie() func(http.Handler) http.Handler
RequirePermission(permission string) func(http.Handler) http.Handler
RequirePermissionsMatchAll(permissions []string) func(http.Handler) http.Handler
RequirePermissionsMatchAny(permissions []string) func(http.Handler) http.Handler
```

### 5. GraphQL Directives (`directives.go`)

**Responsibilities:**
- Enforce auth on GraphQL fields
- Declarative permission checking

**Directives:**
```graphql
@auth                                  # Must be authenticated
@permissionAny(requires: [...])       # Needs any of the permissions
@permissionAll(requires: [...])       # Needs all of the permissions
```

### 6. Domain Security (`domain/*/security/permissions.go`)

**Responsibilities:**
- Define domain-specific permissions as constants
- Provide reusable permission sets

**Example:**
```go
const (
    PermissionRead  = "patient:read"
    PermissionWrite = "patient:write"
)

var ReadAccess = []string{PermissionRead, "admin:all"}
```

---

## Authentication Flow

### Scenario 1: Web Browser (Cookie-based)

```
1. User logs in via external auth service
   ↓
2. Auth service returns JWT
   ↓
3. Application sets JWT in HTTP-only cookie
   Set-Cookie: auth_token=eyJhbGc...; HttpOnly; Secure; SameSite=Strict
   ↓
4. Browser automatically sends cookie with every request
   Cookie: auth_token=eyJhbGc...
   ↓
5. RequireAuth() middleware extracts & validates token
   ↓
6. User loaded into context
```

### Scenario 2: API Client (Header-based)

```
1. Client authenticates and receives JWT
   ↓
2. Client stores JWT securely
   ↓
3. Client includes JWT in Authorization header
   Authorization: Bearer eyJhbGc...
   ↓
4. RequireAuthFromHeader() middleware extracts & validates token
   ↓
5. User loaded into context
```

### JWT Token Structure

```json
{
  "user_id": "12345",
  "email": "doctor@example.com",
  "name": "Dr. John Smith",
  "permissions": [
    "patient:read",
    "patient:write",
    "prescription:read",
    "prescription:write",
    "doctor:role"
  ],
  "iss": "rxintake",
  "aud": "rxintake",
  "exp": 1728691200,
  "iat": 1728604800
}
```

---

## Authorization Flow

### Permission Check Flow

```
Request with valid JWT
    ↓
User in Context
    ↓
Permission Middleware
    ↓
    ├─ RequirePermission("patient:read")
    │  └─ Check: user has "patient:read"? → YES/NO
    │
    ├─ RequirePermissionsMatchAll(["patient:read", "patient:export"])
    │  └─ Check: user has BOTH? → YES/NO
    │
    └─ RequirePermissionsMatchAny(["patient:read", "admin:all"])
       └─ Check: user has EITHER? → YES/NO
```

### Error Responses

**401 Unauthorized** (Authentication Failures)
```json
{
  "error": "unauthorized",
  "message": "Invalid or expired token",
  "status": 401
}
```

**403 Forbidden** (Permission Failures)
```json
{
  "error": "forbidden",
  "message": "User requires at least one of the following permissions",
  "required_permissions": ["patient:read", "admin:all"],
  "match": "any",
  "status": 403
}
```

---

## Permission Model

### Permission Format

```
<resource>:<action>
```

**Examples:**
- `patient:read` - Read patient data
- `patient:write` - Create/update patients
- `prescription:approve` - Approve prescriptions
- `admin:all` - Full administrative access

### Permission Hierarchy

```
Domain Permissions:
├── patient:read
├── patient:write
├── patient:delete
├── patient:export
│
├── prescription:read
├── prescription:write
├── prescription:approve
├── prescription:dispense
│
└── dashboard:view
    dashboard:analytics
    dashboard:reports

System-Wide Permissions:
├── admin:all (super admin)
└── <role>:role (e.g., doctor:role, pharmacist:role)
```

### Permission Sets

Each domain defines reusable permission sets:

```go
// domain/patient/security/permissions.go

// ReadAccess - user needs ANY of these
var ReadAccess = []string{"patient:read", "admin:all"}

// ExportAccess - user needs ALL of these
var ExportAccess = []string{"patient:read", "patient:export"}
```

---

## Usage Examples

### HTTP Routes (Web UI)

```go
// domain/patient/ui/ui.go
import (
    "internal/platform/auth"
    patientsecurity "domain/patient/security"
)

func (ui *UI) RegisterRoutes(r chi.Router) {
    r.Route("/patients", func(r chi.Router) {
        // All routes require cookie-based auth
        r.Use(auth.RequireAuthFromCookie())
        
        // List - needs ANY of ReadAccess permissions
        r.With(auth.RequirePermissionsMatchAny(
            patientsecurity.ReadAccess,
        )).Get("/", ui.listPatients)
        
        // Create - needs ANY of WriteAccess permissions
        r.With(auth.RequirePermissionsMatchAny(
            patientsecurity.WriteAccess,
        )).Post("/", ui.createPatient)
        
        // Export - needs ALL of ExportAccess permissions
        r.With(auth.RequirePermissionsMatchAll(
            patientsecurity.ExportAccess,
        )).Get("/export", ui.exportPatients)
    })
}
```

### REST API Routes

```go
// domain/patient/api/api.go
func (api *API) RegisterRoutes(r chi.Router) {
    r.Route("/api/v1/patients", func(r chi.Router) {
        // All routes require header-based auth
        r.Use(auth.RequireAuthFromHeader())
        
        r.With(auth.RequirePermissionsMatchAny(
            patientsecurity.ReadAccess,
        )).Get("/", api.listPatients)
        
        r.With(auth.RequirePermissionsMatchAny(
            patientsecurity.WriteAccess,
        )).Post("/", api.createPatient)
    })
}
```

### GraphQL Schema

```graphql
# domain/patient/graphql/schema.graphql

type Query {
    # User needs to be authenticated AND have any of these permissions
    patient(id: ID!): Patient 
        @auth 
        @permissionAny(requires: ["patient:read", "admin:all"])
    
    # User needs ALL of these permissions
    exportPatients: [Patient!]! 
        @auth 
        @permissionAll(requires: ["patient:read", "patient:export"])
}

type Mutation {
    createPatient(input: CreatePatientInput!): Patient! 
        @auth 
        @permissionAny(requires: ["patient:write", "admin:all"])
}
```

### GraphQL Resolver

```go
// domain/patient/graphql/patient_resolver.go

// No manual permission checks needed - directives handle it!
func (r *queryResolver) Patient(ctx context.Context, id string) (*model.Patient, error) {
    // Directive already verified auth & permissions
    // Just implement business logic
    return r.PatientService.GetByID(ctx, id)
}

// Can still access user info for audit logging
func (r *mutationResolver) CreatePatient(ctx context.Context, input model.CreatePatientInput) (*model.Patient, error) {
    user, _ := auth.GetCurrentUser(ctx)
    log.Printf("User %s creating patient", user.Email)
    
    return r.PatientService.Create(ctx, input)
}
```

### Inside Handlers (Optional Conditional Logic)

```go
func (ui *UI) patientDetailPage(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    user, _ := auth.GetCurrentUser(ctx)
    
    // Conditional UI based on permissions
    canEdit := auth.HasAnyPermissionCtx(ctx, []string{"patient:write", "admin:all"})
    canExport := auth.HasAllPermissionsCtx(ctx, []string{"patient:read", "patient:export"})
    
    // Pass to template
    component := PatientDetail(patient, canEdit, canExport)
    component.Render(ctx, w)
}
```

---

## Security Best Practices

### 1. JWT Storage

✅ **DO:**
- Store JWT in HTTP-only cookies for web applications
- Use `Secure` flag in production (HTTPS only)
- Set `SameSite=Strict` to prevent CSRF
- Use short expiration times (1 hour)

❌ **DON'T:**
- Store JWT in localStorage (XSS vulnerability)
- Store JWT in sessionStorage
- Use long expiration times

### 2. JWT Configuration

**Development:**
```yaml
auth:
  jwt:
    secret: "dev-secret-key"
    cookie:
      secure: false  # HTTP OK for localhost
```

**Production:**
```yaml
auth:
  jwt:
    jwks_url: "${RX_AUTH_JWT_JWKS_URL}"  # From environment variable
    cookie:
      secure: true  # HTTPS only
```

### 3. Permission Design

✅ **DO:**
- Use clear, descriptive permission names
- Follow `resource:action` format
- Define permissions as constants
- Group related permissions in reusable sets
- Use `MatchAny` for hierarchical access (e.g., admin can do anything)
- Use `MatchAll` for multi-stage operations (e.g., approve + sign)

❌ **DON'T:**
- Hardcode permission strings in routes
- Create overly granular permissions initially
- Use unclear permission names

### 4. Error Handling

- Return **401** for missing/invalid tokens (authentication failure)
- Return **403** for valid tokens with insufficient permissions (authorization failure)
- Provide clear error messages in API responses
- Log auth failures for security monitoring

### 5. Token Validation

Always validate:
- Token signature
- Expiration time (`exp`)
- Issuer (`iss`)
- Audience (`aud`)
- Not before time (`nbf`)

### 6. Testing

```bash
# Test 401 - No token
curl -v https://api.example.com/api/v1/patients
# Expect: 401 Unauthorized

# Test 403 - Valid token, missing permission
curl -v -H "Authorization: Bearer $TOKEN" \
  https://api.example.com/api/v1/patients
# Expect: 403 Forbidden

# Test 200 - Valid token with permission
curl -v -H "Authorization: Bearer $VALID_TOKEN" \
  https://api.example.com/api/v1/patients
# Expect: 200 OK
```

---

## Configuration

### JWT Settings

```yaml
# internal/configs/app.yaml
auth:
  jwt:
    secret: "your-secret-key"           # Signing key (use env var in prod)
    issuer: "rxintake"                  # Expected issuer
    audience: "rxintake"                # Expected audience
    cookie:
      name: "auth_token"                # Cookie name
      secure: true                      # HTTPS only
      httponly: true                    # No JavaScript access
      max_age: 3600                     # 1 hour
```

### Environment Variables

Production should use environment variables:

```bash
export RX_AUTH_JWT_JWKS_URL="https://your-auth-provider.com/.well-known/jwks.json"
export RX_AUTH_JWT_ISSUER="your-auth-service"
```

---

## Troubleshooting

### Common Issues

**Issue: Always getting 401**
- Check token is being sent (cookie or header)
- Verify JWKS URL is accessible
- Check token hasn't expired
- Validate issuer/audience claims

**Issue: Getting 403 instead of 401**
- Token is valid but user lacks permissions
- Check user's permissions array in JWT
- Verify permission constants match exactly

**Issue: GraphQL directives not working**
- Ensure directives are registered in `gqlgen.yml`
- Run `go run github.com/99designs/gqlgen generate`
- Check directive implementation in server setup

---

## Summary

The RxIntake security architecture provides:

1. **Simple JWT-based authentication** - No session management needed
2. **Declarative authorization** - Permissions defined at route/field level
3. **Flexible token support** - Works with cookies (web) and headers (API)
4. **Clear error semantics** - 401 vs 403 distinction
5. **Domain-owned permissions** - Each domain defines its security model
6. **Type-safe permission constants** - No hardcoded strings
7. **Reusable permission sets** - DRY principle for common combinations

This design is production-ready, secure, and easy to maintain!

