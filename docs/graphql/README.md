# GraphQL Documentation

Complete documentation for the RxIntake GraphQL API implementation.

## 📚 Documentation Files

### Getting Started
- **[GraphQL Implementation](./GRAPHQL_IMPLEMENTATION.md)** - Complete implementation guide
- **[Sample Queries](./GRAPHQL_SAMPLE_QUERIES.md)** - Example queries and mutations

### Architecture
- **[Architecture Diagrams](./GRAPHQL_ARCHITECTURE_DIAGRAMS.md)** - Visual architecture overview
- **[Domain Structure](./GRAPHQL_DOMAIN_STRUCTURE.md)** - How domains are organized
- **[Resolver Patterns](./GRAPHQL_RESOLVER_PATTERNS.md)** - Best practices for resolvers

### Development
- **[GraphQL Dev Mode](./GRAPHQL_DEV_MODE.md)** - Testing with GraphQL Playground and mock users
- **[Migration Phases](./GRAPHQL_MIGRATION_PHASES.md)** - Migration strategy and phases

## 🏗️ Architecture Overview

### Domain-Driven Design
Each domain (patient, prescription, etc.) has its own:
- **Schema** - GraphQL type definitions (`.graphql` files)
- **Resolvers** - Business logic implementation
- **Services** - Domain-specific operations
- **Models** - Data structures

### File Structure
```
domain/
  patient/
    graphql/
      schema.graphql       # Type definitions
      resolvers.go         # Resolver implementation
    service/
      patient_service.go   # Business logic
    models/
      patient.go           # Data models
```

## 🔌 Key Features

### Type Safety
- Generated types from schema using `gqlgen`
- Compile-time type checking
- Auto-complete in IDE

### Security Integration
- `@auth` directive for authentication
- `@permissionAny` for permission checks (ANY match)
- `@permissionAll` for permission checks (ALL match)

### N+1 Problem Solutions
- DataLoader pattern for batch loading
- Efficient database queries
- Relationship optimization

## 🚀 Quick Examples

### Schema Definition
```graphql
type Patient {
  id: ID!
  firstName: String!
  lastName: String!
  dateOfBirth: String!
}

type Query {
  patient(id: ID!): Patient @auth @permissionAny(requires: ["patient:read"])
  patients: [Patient!]! @auth @permissionAny(requires: ["patient:read"])
}
```

### Resolver Implementation
```go
func (r *queryResolver) Patient(ctx context.Context, id string) (*models.Patient, error) {
    return r.PatientService.GetByID(ctx, id)
}
```

### Sample Query
```graphql
query GetPatient {
  patient(id: "123") {
    id
    firstName
    lastName
    dateOfBirth
  }
}
```

## 🧪 Testing with Dev Mode

### GraphQL Playground
Access at `http://localhost:8080/playground`

### Set Mock User in HTTP Headers
```json
{
  "X-Mock-User": "doctor"
}
```

Available mock users:
- `admin` - Full access
- `doctor` - Patient and prescription access
- `pharmacist` - Prescription access
- `nurse` - Limited patient access
- `readonly` - Read-only access

## 🔧 Development Workflow

1. **Define Schema** - Add types to `schema.graphql`
2. **Generate Code** - Run `go run github.com/99designs/gqlgen generate`
3. **Implement Resolvers** - Add business logic
4. **Add Services** - Implement domain operations
5. **Test** - Use GraphQL Playground

## 📖 Related Documentation

- [Security Architecture](../security/SECURITY_ARCHITECTURE.md)
- [Dev Mode Guide](../security/SECURITY_DEV_MODE.md)
- [Architecture Overview](../architecture/ARCHITECTURE.md)

## 🛠️ Tools & Libraries

- **gqlgen** - Go GraphQL code generation
- **GraphQL Playground** - Interactive query interface
- **DataLoader** - Batch loading for N+1 prevention

