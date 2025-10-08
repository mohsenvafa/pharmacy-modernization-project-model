# RxIntake Quick Reference

Quick links to commonly used documentation.

## 🚀 Getting Started

| Task | Document |
|------|----------|
| Understand system architecture | [Architecture Overview](./architecture/ARCHITECTURE.md) |
| Set up authentication | [Security Quick Start](./security/SECURITY_QUICK_START.md) |
| Use dev mode for testing | [Dev Mode Guide](./security/SECURITY_DEV_MODE.md) |
| Work with GraphQL | [GraphQL Implementation](./graphql/GRAPHQL_IMPLEMENTATION.md) |

## 🔐 Security & Authentication

| Task | Document |
|------|----------|
| Understand auth system | [Security Architecture](./security/SECURITY_ARCHITECTURE.md) |
| Protect routes | [Routes Security Implementation](./security/ROUTES_SECURITY_IMPLEMENTATION.md) |
| Use mock users | [Security Mock Users](./security/SECURITY_MOCK_USERS.md) |
| Initialize auth system | [Security Builder](./security/SECURITY_BUILDER.md) |
| Add UI permission checks | [Security README](./security/README.md) |

## 🔌 GraphQL

| Task | Document |
|------|----------|
| Create GraphQL schemas | [GraphQL Implementation](./graphql/GRAPHQL_IMPLEMENTATION.md) |
| Write resolvers | [Resolver Patterns](./graphql/GRAPHQL_RESOLVER_PATTERNS.md) |
| Test with Playground | [GraphQL Dev Mode](./graphql/GRAPHQL_DEV_MODE.md) |
| Example queries | [Sample Queries](./graphql/GRAPHQL_SAMPLE_QUERIES.md) |
| Understand architecture | [Architecture Diagrams](./graphql/GRAPHQL_ARCHITECTURE_DIAGRAMS.md) |

## 🗄️ Database

| Task | Document |
|------|----------|
| Set up MongoDB | [MongoDB Implementation](./mongodb/MONGODB_IMPLEMENTATION.md) |
| Repository patterns | [MongoDB README](./mongodb/README.md) |

## 📦 Frontend

| Task | Document |
|------|----------|
| Add TypeScript components | [Adding TypeScript Components](./typescript/ADDING_TYPESCRIPT_COMPONENTS.md) |
| TypeScript patterns | [TypeScript README](./typescript/README.md) |

## 🧪 Development Workflow

### Start Dev Server
```bash
go run cmd/server/main.go
```

### Access Dev Mode Endpoints
- **Auth Info**: http://localhost:8080/__dev/auth
- **Switch User**: http://localhost:8080/__dev/switch?user=doctor
- **GraphQL Playground**: http://localhost:8080/playground

### Available Mock Users
| User | Permissions | Use Case |
|------|-------------|----------|
| `admin` | `admin:all` | Full system access |
| `doctor` | Patient + Prescription management | Medical staff |
| `pharmacist` | Prescription dispensing | Pharmacy staff |
| `nurse` | Patient read-only | Limited access |
| `readonly` | Read all | Auditing/reporting |

### Switch Users in Browser
```
http://localhost:8080/__dev/switch?user=nurse
```

### Switch Users in GraphQL Playground
Add HTTP header:
```json
{
  "X-Mock-User": "doctor"
}
```

## 🔧 Common Tasks

### Protect a Route
```go
r.With(
    auth.RequireAuthWithDevMode(),
    auth.RequirePermissionsMatchAny([]string{"patient:read", "admin:all"}),
).Get("/patients", handler)
```

### Add GraphQL Permission
```graphql
type Query {
  patients: [Patient!]! 
    @auth 
    @permissionAny(requires: ["patient:read", "admin:all"])
}
```

### Check Permission in UI
```go
@authComponents.IfHasAnyPermission(ctx, []string{"patient:read"}) {
    <div>Patient data here</div>
}
```

### Query Database
```go
collection := db.Collection("patients")
err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&patient)
```

## 📖 Documentation Index

- [Main Documentation Index](./README.md)
- [Security Documentation](./security/README.md)
- [GraphQL Documentation](./graphql/README.md)
- [Architecture Documentation](./architecture/README.md)
- [MongoDB Documentation](./mongodb/README.md)
- [TypeScript Documentation](./typescript/README.md)

## 🆘 Need Help?

1. Check the relevant topic folder README
2. Look for specific implementation guides
3. Review example code in the documentation
4. Check the codebase for similar patterns

## 🔄 Keep Documentation Updated

When you make changes to the codebase:
1. Update relevant documentation
2. Add examples for new features
3. Keep quick reference current
4. Update diagrams if architecture changes

