# Security Documentation

Comprehensive documentation for the RxIntake authentication and authorization system.

## 📚 Documentation Files

### Getting Started
- **[Security Quick Start](./SECURITY_QUICK_START.md)** - Quick setup guide for authentication
- **[Security README](./SECURITY_README.md)** - General security overview

### Architecture & Design
- **[Security Architecture](./SECURITY_ARCHITECTURE.md)** - Complete architecture overview with diagrams
- **[Security Diagrams](./SECURITY_DIAGRAMS.md)** - Visual representations of the security system
- **[Security Builder Pattern](./SECURITY_BUILDER.md)** - Authentication initialization pattern

### Implementation Guides
- **[Complete Security Summary](./COMPLETE_SECURITY_SUMMARY.md)** - Comprehensive implementation summary
- **[Security Implementation Summary](./SECURITY_IMPLEMENTATION_SUMMARY.md)** - Step-by-step implementation
- **[Routes Security Implementation](./ROUTES_SECURITY_IMPLEMENTATION.md)** - How routes are protected

### Development Mode
- **[Security Dev Mode](./SECURITY_DEV_MODE.md)** - Development mode guide
- **[Dev Mode Implementation](./DEV_MODE_IMPLEMENTATION_SUMMARY.md)** - Dev mode implementation details
- **[Dev Mode Example](./SECURITY_DEV_MODE_EXAMPLE.md)** - Usage examples
- **[Mock Users](./SECURITY_MOCK_USERS.md)** - Available mock users for testing

### User Interface
- **[Security User Display](./SECURITY_USER_DISPLAY.md)** - User information display components

## 🔑 Key Concepts

### Authentication (AuthN)
- JWT-based stateless authentication
- Token validation and extraction
- Cookie and header-based token sources
- Dev mode bypass for local testing

### Authorization (AuthZ)
- Permission-based access control
- Resource:action permission model (e.g., `patient:read`)
- Match strategies: `ANY` and `ALL`
- Route-level and UI-level permission checks

### Middleware
- `RequireAuth()` - Basic authentication
- `RequireAuthWithDevMode()` - Dev mode-aware auth
- `RequirePermissionsMatchAny()` - ANY permission matching
- `RequirePermissionsMatchAll()` - ALL permission matching

### GraphQL Directives
- `@auth` - Requires authentication
- `@permissionAny` - Requires ANY of specified permissions
- `@permissionAll` - Requires ALL specified permissions

## 🧪 Mock Users (Dev Mode)

Available test users:
- **admin** - Full access (`admin:all`)
- **doctor** - Patient and prescription management
- **pharmacist** - Prescription dispensing
- **nurse** - Patient read-only access
- **readonly** - Read-only access to all resources

## 🚀 Quick Examples

### Protecting a Route
```go
r.With(
    auth.RequireAuthWithDevMode(),
    auth.RequirePermissionsMatchAny([]string{"patient:read", "admin:all"}),
).Get("/patients", handler)
```

### GraphQL Field Protection
```graphql
type Query {
  patients: [Patient!]! 
    @auth 
    @permissionAny(requires: ["patient:read", "admin:all"])
}
```

### UI Permission Check
```go
@authComponents.IfHasAnyPermission(ctx, []string{"patient:read"}) {
    <div>Patient data here</div>
}
```

## 🔧 Configuration

Security is configured in `internal/configs/app.yaml`:
```yaml
auth:
  dev_mode: true # Only for development
  jwt:
    secret: "your-secret-key"
    issuer: "rxintake"
    audience: "rxintake"
```

## 📖 Related Documentation

- [Architecture Overview](../architecture/ARCHITECTURE.md)
- [GraphQL Dev Mode](../graphql/GRAPHQL_DEV_MODE.md)

