# Security Quick Start Guide

This guide shows you how to implement authentication and authorization in your RxIntake application.

## Prerequisites

The security framework is already set up. You just need to:
1. Initialize the JWT config
2. Apply middleware to your routes
3. Add directives to GraphQL schemas
4. Wire up in your app

---

## Step 1: Initialize JWT Config in App Startup

Update your application initialization to set up JWT:

```go
// internal/app/app.go
package app

import (
    "internal/platform/auth"
    "internal/platform/config"
)

func NewApp(cfg *config.Config) *App {
    // Initialize JWT configuration
    auth.InitJWTConfig(auth.JWTConfig{
        Secret:     cfg.Auth.JWT.Secret,
        Issuer:     cfg.Auth.JWT.Issuer,
        Audience:   cfg.Auth.JWT.Audience,
        CookieName: cfg.Auth.JWT.Cookie.Name,
    })
    
    // ... rest of your app initialization
}
```

---

## Step 2: Protect HTTP Routes

### Web UI Routes (Cookie-based)

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
        // Require authentication via cookie
        r.Use(auth.RequireAuthFromCookie())
        
        // Public patient routes - needs read permission
        r.With(auth.RequirePermissionsMatchAny(
            patientsecurity.ReadAccess,
        )).Get("/", ui.listPatients)
        
        // Create/edit - needs write permission
        r.With(auth.RequirePermissionsMatchAny(
            patientsecurity.WriteAccess,
        )).Post("/", ui.createPatient)
        
        // Export - needs both read AND export permissions
        r.With(auth.RequirePermissionsMatchAll(
            patientsecurity.ExportAccess,
        )).Get("/export", ui.exportPatients)
    })
}
```

### REST API Routes (Header-based)

```go
// domain/patient/api/api.go
package api

import (
    "github.com/go-chi/chi/v5"
    "internal/platform/auth"
    patientsecurity "domain/patient/security"
)

func (api *API) RegisterRoutes(r chi.Router) {
    r.Route("/api/v1/patients", func(r chi.Router) {
        // Require authentication via Authorization header
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

---

## Step 3: Protect GraphQL Resolvers

### Add Directives to Schema

```graphql
# domain/patient/graphql/schema.graphql

extend type Query {
    # Authenticated + needs any of these permissions
    patient(id: ID!): Patient 
        @auth 
        @permissionAny(requires: ["patient:read", "admin:all"])
    
    patients: [Patient!]! 
        @auth 
        @permissionAny(requires: ["patient:read", "admin:all"])
    
    # Needs ALL permissions
    exportPatients: [Patient!]! 
        @auth 
        @permissionAll(requires: ["patient:read", "patient:export"])
}

extend type Mutation {
    createPatient(input: CreatePatientInput!): Patient! 
        @auth 
        @permissionAny(requires: ["patient:write", "admin:all"])
}
```

### Register Directives in GraphQL Server

```go
// internal/graphql/server.go
package graphql

import (
    "internal/graphql/generated"
    "internal/platform/auth"
    "github.com/99designs/gqlgen/graphql/handler"
)

func NewGraphQLHandler(resolver *Resolver) http.Handler {
    srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
        Resolvers: resolver,
        Directives: generated.DirectiveRoot{
            Auth:          auth.AuthDirective(),
            PermissionAny: auth.PermissionAnyDirective(),
            PermissionAll: auth.PermissionAllDirective(),
        },
    }))
    
    return srv
}
```

### Wrap GraphQL Endpoint with Auth Middleware

```go
// Make sure the GraphQL endpoint extracts JWT and sets user in context
router.Handle("/graphql", auth.RequireAuth()(graphqlHandler))
```

### Keep Resolvers Clean

```go
// domain/patient/graphql/patient_resolver.go

// NO manual permission checks needed!
func (r *queryResolver) Patient(ctx context.Context, id string) (*model.Patient, error) {
    // Directives already verified auth & permissions
    return r.PatientService.GetByID(ctx, id)
}

// Can still access user for logging
func (r *mutationResolver) CreatePatient(ctx context.Context, input model.CreatePatientInput) (*model.Patient, error) {
    user, _ := auth.GetCurrentUser(ctx)
    log.Printf("User %s creating patient", user.Email)
    return r.PatientService.Create(ctx, input)
}
```

---

## Step 4: Regenerate GraphQL Code

After adding directives to your schemas:

```bash
go run github.com/99designs/gqlgen generate
```

---

## Step 5: Testing

### Create a Test JWT Token

For development/testing, create tokens with the helper:

```go
// Create test user
testUser := &auth.User{
    ID:    "test-123",
    Email: "doctor@test.com",
    Name:  "Dr. Test",
    Permissions: []string{
        "patient:read",
        "patient:write",
        "prescription:read",
        "doctor:role",
    },
}

# Test with your auth provider's token
# Use tokens created by Auth0, Okta, or your auth provider

fmt.Println("Test Token:", token)
```

### Test with curl

```bash
# Test API endpoint
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/v1/patients

# Test GraphQL
curl -X POST http://localhost:8080/graphql \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"query": "{ patients { id name } }"}'
```

### Test with Browser

Set a cookie manually for testing:

```javascript
// In browser console
document.cookie = "auth_token=YOUR_TOKEN; path=/";
```

Then navigate to protected pages.

---

## Common Patterns

### Pattern 1: Public + Protected Routes

```go
func SetupRoutes(r chi.Router) {
    // Public routes (no auth)
    r.Get("/", homePage)
    r.Get("/login", loginPage)
    
    // Protected routes (auth required)
    r.Group(func(r chi.Router) {
        r.Use(auth.RequireAuthFromCookie())
        r.Get("/dashboard", dashboardPage)
        r.Get("/patients", patientsPage)
    })
}
```

### Pattern 2: Different Permissions for Same Resource

```go
r.Route("/patients", func(r chi.Router) {
    r.Use(auth.RequireAuth())
    
    // Read operations - any of these permissions
    r.With(auth.RequirePermissionsMatchAny(
        []string{"patient:read", "admin:all"},
    )).Get("/", listPatients)
    
    // Write operations - any of these permissions
    r.With(auth.RequirePermissionsMatchAny(
        []string{"patient:write", "admin:all"},
    )).Post("/", createPatient)
    
    // Delete - needs both write AND delete
    r.With(auth.RequirePermissionsMatchAll(
        []string{"patient:write", "patient:delete"},
    )).Delete("/{id}", deletePatient)
})
```

### Pattern 3: Conditional UI Based on Permissions

```go
func (ui *UI) patientDetailPage(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // Check what user can do
    canEdit := auth.HasAnyPermissionCtx(ctx, 
        []string{"patient:write", "admin:all"})
    canDelete := auth.HasAllPermissionsCtx(ctx, 
        []string{"patient:write", "patient:delete"})
    
    // Render template with conditional buttons
    component := PatientDetail(patient, canEdit, canDelete)
    component.Render(ctx, w)
}
```

---

## Debugging

### Enable Auth Logging

Auth middleware automatically logs:
- `AUTH 401` - Authentication failures
- `AUTH OK` - Successful authentication
- `AUTHZ 403` - Authorization failures

Check your logs to see what's happening:

```
AUTH 401: GET /patients - No token: no Authorization header found
AUTH OK: GET /patients - User: Dr. Smith (doctor@example.com)
AUTHZ 403: POST /patients/export - User doctor@example.com lacks all permissions: [patient:read patient:export]
```

### Common Issues

**Always getting 401:**
- Token not being sent
- Token expired
- Wrong secret key
- Invalid issuer/audience

**Getting 403:**
- Token valid but user lacks permissions
- Check permission spelling
- Verify user's permissions array in JWT

---

## Security Checklist

Before deploying to production:

- [ ] Set `auth.jwt.secret` via environment variable
- [ ] Set `auth.jwt.cookie.secure: true` (HTTPS only)
- [ ] Use short token expiration (1 hour max)
- [ ] Review all route permissions
- [ ] Test unauthorized access (should get 401/403)
- [ ] Add rate limiting on auth endpoints
- [ ] Set up monitoring for auth failures

---

## Next Steps

1. ✅ Read [SECURITY_ARCHITECTURE.md](./SECURITY_ARCHITECTURE.md) for full details
2. ✅ Define permissions for your domains in `domain/*/security/permissions.go`
3. ✅ Apply middleware to your routes
4. ✅ Add directives to GraphQL schemas
5. ✅ Test with different permission combinations
6. ✅ Set up proper JWT issuing in production

Your authentication and authorization system is now ready! 🎉

