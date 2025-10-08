# GraphQL Domain-Based Structure

This document explains how GraphQL is organized by domain, mirroring the REST API structure.

## 📂 Structure Overview

### Before (Centralized):
```
internal/
  graphql/
    ├── schema.graphql           # ALL types
    ├── schema.resolvers.go      # ALL resolvers
    ├── resolver.go
    └── server.go
```

### After (Domain-Based):
```
domain/
  patient/
    ├── api/                     # REST API
    ├── graphql/                 # GraphQL API ✨
    │   ├── schema.graphql       # Patient types
    │   └── resolver.go          # Patient resolvers
    ├── service/
    └── ui/

  prescription/
    ├── api/                     # REST API
    ├── graphql/                 # GraphQL API ✨
    │   ├── schema.graphql       # Prescription types
    │   └── resolver.go          # Prescription resolvers
    └── service/

  dashboard/
    ├── graphql/                 # GraphQL API ✨
    │   ├── schema.graphql       # Dashboard types
    │   └── resolver.go          # Dashboard resolvers
    └── service/

internal/
  graphql/
    ├── schema.graphql           # Root types only (Query, Mutation)
    ├── schema.resolvers.go      # Delegates to domain resolvers
    ├── resolver.go              # Aggregates domain resolvers
    ├── server.go                # Server setup
    └── generated/               # Auto-generated code
```

---

## 🎯 Benefits of Domain-Based Structure

### 1. **Consistency with REST API**
```
domain/patient/
  ├── api/           # REST endpoints
  ├── graphql/       # GraphQL resolvers  ← Same level!
  ├── service/       # Shared business logic
  └── ui/            # Server-rendered pages
```

All delivery layers (REST, GraphQL, UI) are siblings within each domain.

### 2. **Domain Ownership**
Each domain team owns:
- ✅ Business logic (services)
- ✅ REST API endpoints
- ✅ GraphQL schema & resolvers
- ✅ UI handlers

No need to coordinate with a central GraphQL team!

### 3. **Independent Evolution**
```bash
# Change patient schema
vim domain/patient/graphql/schema.graphql

# Regenerate
gqlgen generate

# Only patient resolvers are affected!
```

Changes to one domain don't affect other domains.

### 4. **Clearer Boundaries**
```
domain/patient/graphql/
  ├── schema.graphql       # 40 lines (patient types only)
  └── resolver.go          # 120 lines (patient resolvers only)

domain/prescription/graphql/
  ├── schema.graphql       # 35 lines (prescription types only)
  └── resolver.go          # 90 lines (prescription resolvers only)
```

Instead of one massive 300+ line resolver file!

### 5. **Easy to Extract** (Future)
When you move to microservices, each domain is already isolated:
```bash
# Extract patient domain to its own service
mv domain/patient/ patient-service/
# GraphQL, REST, and services move together!
```

---

## 📝 Schema Organization

### Root Schema (`internal/graphql/schema.graphql`)

Defines only root types and common scalars:

```graphql
# Common scalars
scalar Time

# Root Query type (domains extend this)
type Query {
  _empty: String
}

# Root Mutation type (domains extend this)
type Mutation {
  _empty: String
}
```

### Domain Schemas

Each domain extends the root types:

**`domain/patient/graphql/schema.graphql`:**
```graphql
type Patient {
  id: ID!
  name: String!
  # ...
}

extend type Query {
  patient(id: ID!): Patient
  patients(query: String, limit: Int, offset: Int): [Patient!]!
}
```

**`domain/prescription/graphql/schema.graphql`:**
```graphql
type Prescription {
  id: ID!
  drug: String!
  # ...
}

extend type Query {
  prescription(id: ID!): Prescription
  prescriptions(status: String, limit: Int, offset: Int): [Prescription!]!
}
```

**Key Point:** Uses `extend type Query` to add to the root Query type!

---

## 🔧 Resolver Organization

### Domain Resolvers

Each domain has its own resolver struct:

```go
// domain/patient/graphql/resolver.go
package graphql

type PatientResolver struct {
    PatientService      patientservice.PatientService
    AddressService      patientservice.AddressService
    PrescriptionService prescriptionservice.PrescriptionService
    Logger              *zap.Logger
}

func (r *PatientResolver) Patient(ctx context.Context, id string) (*model.Patient, error) {
    return r.PatientService.GetByID(ctx, id)
}

func (r *PatientResolver) Patients(ctx context.Context, query *string, limit *int, offset *int) ([]model.Patient, error) {
    // ... implementation
}
```

### Root Resolver (Aggregator)

The root resolver aggregates all domain resolvers:

```go
// internal/graphql/resolver.go
package graphql

type Resolver struct {
    PatientResolver      *patientgraphql.PatientResolver
    PrescriptionResolver *prescriptiongraphql.PrescriptionResolver
    DashboardResolver    *dashboardgraphql.DashboardResolver
}
```

### Delegation Pattern

The generated resolvers delegate to domain resolvers:

```go
// internal/graphql/schema.resolvers.go (auto-generated, you implement delegation)

func (r *queryResolver) Patient(ctx context.Context, id string) (*model.Patient, error) {
    // Delegate to patient domain resolver
    return r.PatientResolver.Patient(ctx, id)
}

func (r *queryResolver) Prescription(ctx context.Context, id string) (*model1.Prescription, error) {
    // Delegate to prescription domain resolver
    return r.PrescriptionResolver.Prescription(ctx, id)
}
```

---

## 🚀 Workflow: Adding New GraphQL Field

### Example: Add `activeOrderCount` field to Patient

**Step 1: Update Domain Schema**
```bash
vim domain/patient/graphql/schema.graphql
```

```graphql
type Patient {
  id: ID!
  name: String!
  activeOrderCount: Int!  # NEW
  # ...
}
```

**Step 2: Regenerate**
```bash
gqlgen generate
```

This creates a stub in `domain/patient/graphql/resolver.go` (if it doesn't exist, add it manually).

**Step 3: Implement in Domain Resolver**
```bash
vim domain/patient/graphql/resolver.go
```

```go
func (r *PatientResolver) ActiveOrderCount(ctx context.Context, obj *model.Patient) (int, error) {
    return r.OrderService.CountActive(ctx, obj.ID)
}
```

**Step 4: Update Root Resolver (if needed)**

Usually auto-handled, but if you get compile errors:

```bash
vim internal/graphql/schema.resolvers.go
```

```go
func (r *patientResolver) ActiveOrderCount(ctx context.Context, obj *model.Patient) (int, error) {
    return r.PatientResolver.ActiveOrderCount(ctx, obj)
}
```

**Done!** ✅

---

## 🔄 Workflow: Adding New Domain

### Example: Add Order Domain

**Step 1: Create Domain Structure**
```bash
mkdir -p domain/order/graphql
```

**Step 2: Create Schema**
```bash
vim domain/order/graphql/schema.graphql
```

```graphql
type Order {
  id: ID!
  patientID: ID!
  status: OrderStatus!
  total: Float!
}

enum OrderStatus {
  PENDING
  PROCESSING
  SHIPPED
  DELIVERED
}

extend type Query {
  order(id: ID!): Order
  orders(patientID: ID): [Order!]!
}
```

**Step 3: Create Resolver**
```bash
vim domain/order/graphql/resolver.go
```

```go
package graphql

import (
    "context"
    "go.uber.org/zap"
    orderservice "pharmacy-modernization-project-model/domain/order/service"
    "pharmacy-modernization-project-model/domain/order/contracts/model"
)

type OrderResolver struct {
    OrderService orderservice.OrderService
    Logger       *zap.Logger
}

func NewOrderResolver(orderSvc orderservice.OrderService, logger *zap.Logger) *OrderResolver {
    return &OrderResolver{
        OrderService: orderSvc,
        Logger:       logger,
    }
}

func (r *OrderResolver) Order(ctx context.Context, id string) (*model.Order, error) {
    return r.OrderService.GetByID(ctx, id)
}

func (r *OrderResolver) Orders(ctx context.Context, patientID *string) ([]model.Order, error) {
    // ... implementation
}
```

**Step 4: Add to Root Resolver**
```bash
vim internal/graphql/resolver.go
```

```go
type Resolver struct {
    PatientResolver      *patientgraphql.PatientResolver
    PrescriptionResolver *prescriptiongraphql.PrescriptionResolver
    DashboardResolver    *dashboardgraphql.DashboardResolver
    OrderResolver        *ordergraphql.OrderResolver  // NEW
}
```

**Step 5: Wire in Server**
```bash
vim internal/graphql/server.go
```

```go
func MountGraphQL(r chi.Router, deps *Dependencies) {
    // ... existing resolvers

    orderResolver := ordergraphql.NewOrderResolver(
        deps.OrderService,
        deps.Logger,
    )

    resolver := &Resolver{
        PatientResolver:      patientResolver,
        PrescriptionResolver: prescriptionResolver,
        DashboardResolver:    dashboardResolver,
        OrderResolver:        orderResolver,  // NEW
    }

    // ... rest of setup
}
```

**Step 6: Regenerate & Implement**
```bash
gqlgen generate
```

Then implement delegation in `internal/graphql/schema.resolvers.go`.

**Done!** ✅

---

## 📊 Comparison: Before vs After

| Aspect | Before (Centralized) | After (Domain-Based) |
|--------|---------------------|---------------------|
| **Schema Files** | 1 file (300+ lines) | 4 files (40-50 lines each) |
| **Resolver Files** | 1 file (500+ lines) | 4 files (100-150 lines each) |
| **Domain Isolation** | ❌ All mixed | ✅ Clear boundaries |
| **Consistency** | ❌ GraphQL separate from REST | ✅ Same structure as REST |
| **Team Ownership** | ❌ Central GraphQL team | ✅ Domain teams |
| **Change Impact** | ❌ Affects all domains | ✅ Isolated to one domain |
| **Microservices Ready** | ❌ Hard to extract | ✅ Easy to extract |

---

## 🎨 Visual Architecture

```
                 Client Request
                       ↓
              /graphql endpoint
                       ↓
         GraphQL Server (internal/graphql/server.go)
                       ↓
         Root Resolver (internal/graphql/resolver.go)
                       ↓
         ┌──────────────┼──────────────┐
         ↓              ↓              ↓
    Patient        Prescription    Dashboard
    Resolver       Resolver        Resolver
         ↓              ↓              ↓
    Patient        Prescription    Dashboard
    Service        Service         Service
         ↓              ↓              ↓
    Patient        Prescription    Dashboard
    Repository     Repository      Repository
```

---

## 🔑 Key Design Principles

### 1. **Delegation Pattern**
Root resolvers are thin delegators:
```go
// Root resolver (thin)
func (r *queryResolver) Patient(ctx context.Context, id string) (*model.Patient, error) {
    return r.PatientResolver.Patient(ctx, id)  // Delegate!
}

// Domain resolver (thick)
func (r *PatientResolver) Patient(ctx context.Context, id string) (*model.Patient, error) {
    // Real implementation with logging, error handling, etc.
    patient, err := r.PatientService.GetByID(ctx, id)
    if err != nil {
        r.Logger.Error("Failed to fetch patient", zap.String("id", id), zap.Error(err))
        return nil, err
    }
    return &patient, nil
}
```

### 2. **Schema Extension**
Domains extend root types:
```graphql
# Root schema defines base
type Query {
  _empty: String
}

# Domains extend it
extend type Query {
  patient(id: ID!): Patient
}

extend type Query {
  prescription(id: ID!): Prescription
}
```

### 3. **Service Reuse**
Domain resolvers call the same services as REST API:
```go
// Both use the same service!
type PatientResolver struct {
    PatientService patientservice.PatientService  // ← Same interface
}

type PatientController struct {
    PatientService patientservice.PatientService  // ← Same interface
}
```

---

## 🚦 Best Practices

### ✅ DO:
- Keep domain schemas focused on that domain's types
- Put all business logic in services, not resolvers
- Use the same services for REST and GraphQL
- Add logging and error handling in domain resolvers
- Test domain resolvers independently

### ❌ DON'T:
- Put business logic in resolvers
- Access other domain's services directly (use providers pattern)
- Create circular dependencies between domain resolvers
- Skip error handling or logging
- Duplicate code between REST and GraphQL

---

## 🎓 Summary

**Key Takeaways:**

1. **GraphQL mirrors REST structure** - Both are delivery layers within domains
2. **Domain ownership** - Each team owns their GraphQL schema and resolvers
3. **Delegation pattern** - Root resolvers delegate to domain resolvers
4. **Schema extension** - Domains extend root Query/Mutation types
5. **Service reuse** - Same services power REST, GraphQL, and UI

**This architecture:**
- ✅ Scales with team growth
- ✅ Makes changes isolated and safe
- ✅ Prepares for microservices migration
- ✅ Maintains consistency across delivery layers
- ✅ Keeps code organized and maintainable

**Next Steps:**
- See `GRAPHQL_IMPLEMENTATION.md` for query examples
- See `GRAPHQL_MIGRATION_PHASES.md` for future evolution paths
- See `GRAPHQL_ARCHITECTURE_DIAGRAMS.md` for visual diagrams

