# GraphQL Resolver Patterns: Evolution Guide

This document explains the three-phase evolution of GraphQL resolver organization, with working examples in the codebase.

## 🎯 Overview

As your GraphQL domain grows, you'll need different resolver organization patterns. This guide shows three phases with real implementations in the `domain/patient/graphql/` folder.

---

## 📊 Phase Comparison

| Phase | When to Use | Files | Complexity | Patient Example |
|-------|-------------|-------|------------|-----------------|
| **Phase 1** | < 300 lines, 1-3 services | 1 file | Low | Not shown (too simple) |
| **Phase 2** | 300-500 lines, 3-5 services | 2-3 files | Medium | ✅ **ACTIVE** |
| **Phase 3** | 500+ lines, 5+ services | 5+ files | High | 📚 Reference in `resolvers/` |

---

## Phase 1: Single Resolver (Simple Domains)

### Structure
```
domain/patient/graphql/
├── schema.graphql
└── resolver.go         # Everything in one file
```

### When to Use
- ✅ New domains (MVP stage)
- ✅ Less than 300 lines of code
- ✅ 1-3 service dependencies
- ✅ Small team (1-2 developers)
- ✅ Simple CRUD operations

### Example Pattern
```go
package graphql

type PatientResolver struct {
    PatientService      patientservice.PatientService
    AddressService      patientservice.AddressService
    PrescriptionService prescriptionservice.PrescriptionService
    Logger              *zap.Logger
}

// Query: Get single patient
func (r *PatientResolver) Patient(ctx context.Context, id string) (*model.Patient, error) {
    return r.PatientService.GetByID(ctx, id)
}

// Query: List patients
func (r *PatientResolver) Patients(ctx context.Context, query *string, limit *int, offset *int) ([]model.Patient, error) {
    // Implementation
}

// Field: Resolve addresses
func (r *PatientResolver) Addresses(ctx context.Context, obj *model.Patient) ([]model.Address, error) {
    return r.AddressService.GetByPatientID(ctx, obj.ID)
}

// Field: Resolve prescriptions
func (r *PatientResolver) Prescriptions(ctx context.Context, obj *model.Patient) ([]model1.Prescription, error) {
    // Implementation
}
```

### Benefits
- ✅ Simple and straightforward
- ✅ Easy to navigate (everything in one place)
- ✅ Fast development
- ✅ Low cognitive overhead

### Limitations
- ❌ File grows large over time
- ❌ Mixed concerns (queries + fields + mutations)
- ❌ Hard to test individual concerns
- ❌ Merge conflicts as team grows

---

## Phase 2: Split by Sub-Resource ✅ (CURRENT)

### Structure
```
domain/patient/graphql/
├── schema.graphql
├── patient_resolver.go   # ✅ Patient operations
└── address_resolver.go   # ✅ Address operations (separated)
```

### When to Use
- ✅ Sub-resources have significant logic
- ✅ File approaching 300+ lines
- ✅ Want better organization
- ✅ 3-5 service dependencies
- ✅ Team of 2-4 developers

### Active Implementation

**📁 `domain/patient/graphql/patient_resolver.go`:**
```go
type PatientResolver struct {
    PatientService      patientservice.PatientService
    PrescriptionService prescriptionservice.PrescriptionService
    AddressResolver     *AddressResolver  // Delegates address operations
    Logger              *zap.Logger
}

func NewPatientResolver(
    patientSvc patientservice.PatientService,
    addressSvc patientservice.AddressService,
    prescriptionSvc prescriptionservice.PrescriptionService,
    logger *zap.Logger,
) *PatientResolver {
    return &PatientResolver{
        PatientService:      patientSvc,
        PrescriptionService: prescriptionSvc,
        AddressResolver:     NewAddressResolver(addressSvc, logger),  // Create sub-resolver
        Logger:              logger,
    }
}

// Delegates to AddressResolver
func (r *PatientResolver) Addresses(ctx context.Context, obj *model.Patient) ([]model.Address, error) {
    return r.AddressResolver.Addresses(ctx, obj)
}
```

**📁 `domain/patient/graphql/address_resolver.go`:**
```go
type AddressResolver struct {
    AddressService patientservice.AddressService
    Logger         *zap.Logger
}

func NewAddressResolver(addressSvc patientservice.AddressService, logger *zap.Logger) *AddressResolver {
    return &AddressResolver{
        AddressService: addressSvc,
        Logger:         logger,
    }
}

func (r *AddressResolver) Addresses(ctx context.Context, obj *model.Patient) ([]model.Address, error) {
    addresses, err := r.AddressService.GetByPatientID(ctx, obj.ID)
    if err != nil {
        r.Logger.Error("Failed to fetch addresses", zap.String("patient_id", obj.ID), zap.Error(err))
        return []model.Address{}, nil
    }
    return addresses, nil
}

// Address-specific computed fields
func (r *AddressResolver) FormattedAddress(ctx context.Context, obj *model.Address) (string, error) {
    return obj.Line1 + ", " + obj.City + ", " + obj.State + " " + obj.Zip, nil
}

func (r *AddressResolver) ValidateAddress(ctx context.Context, obj *model.Address) (bool, error) {
    // Complex validation logic
    return obj.Zip != "" && obj.City != "" && obj.State != "", nil
}
```

### Benefits
- ✅ Better organization (focused files)
- ✅ Clear boundaries (patient vs address)
- ✅ Sub-resource has dedicated space for complex logic
- ✅ Easier to test address logic separately
- ✅ Still simple to understand
- ✅ Scales well for medium domains

### Delegation Pattern
The key pattern is **composition + delegation**:
1. Main resolver creates sub-resolvers
2. Main resolver delegates to sub-resolvers
3. Sub-resolvers handle their specific logic

---

## Phase 3: Full Separation (Large Domains)

### Structure
```
domain/patient/graphql/
├── schema.graphql
├── resolvers/
│   ├── query/
│   │   └── patient_query_resolver.go       # Read operations
│   ├── mutation/
│   │   └── patient_mutation_resolver.go    # Write operations
│   ├── field/
│   │   ├── patient_field_resolver.go       # Complex patient fields
│   │   └── address_field_resolver.go       # Complex address fields
│   └── aggregator_example.go               # Combines all resolvers
```

### When to Use
- ✅ Domain exceeds 500 lines
- ✅ 5+ service dependencies
- ✅ Multiple teams working on same domain
- ✅ Need different dependencies for queries vs mutations
- ✅ Complex authorization or validation logic
- ✅ Planning microservices extraction

### Implementation Guide

**When implementing Phase 3, create:**
- `domain/your_domain/graphql/resolvers/query/` - Query operations
- `domain/your_domain/graphql/resolvers/mutation/` - Mutation operations
- `domain/your_domain/graphql/resolvers/field/` - Field resolvers
- `domain/your_domain/graphql/resolvers/aggregator.go` - Combines all

**Detailed patterns in:** `domain/patient/graphql/RESOLVER_PHASES.md`

### Pattern: Query Resolver
```go
// Focused on read operations
type PatientQueryResolver struct {
    PatientService patientservice.PatientService
    Logger         *zap.Logger
}

func (r *PatientQueryResolver) Patient(ctx context.Context, id string) (*model.Patient, error) {
    // Query logic
}

func (r *PatientQueryResolver) Patients(ctx context.Context, query *string, limit *int, offset *int) ([]model.Patient, error) {
    // Query logic
}
```

### Pattern: Mutation Resolver
```go
// Focused on write operations
type PatientMutationResolver struct {
    PatientService     patientservice.PatientService
    ValidationService  validationservice.ValidationService  // Different deps!
    EventPublisher     events.Publisher                     // Different deps!
    Logger             *zap.Logger
}

func (r *PatientMutationResolver) CreatePatient(ctx context.Context, input generated.CreatePatientInput) (*model.Patient, error) {
    // 1. Validate
    // 2. Create
    // 3. Publish event
    // 4. Audit log
}
```

### Pattern: Field Resolver
```go
// Focused on complex field resolution
type PatientFieldResolver struct {
    PatientService      patientservice.PatientService
    AddressService      patientservice.AddressService
    PrescriptionService prescriptionservice.PrescriptionService
    CacheService        cache.CacheService  // Different deps!
    Logger              *zap.Logger
}

func (r *PatientFieldResolver) Addresses(ctx context.Context, obj *model.Patient) ([]model.Address, error) {
    // Complex field resolution
}

func (r *PatientFieldResolver) ActivePrescriptionCount(ctx context.Context, obj *model.Patient) (int, error) {
    // Computed/aggregated field
}
```

### Pattern: Aggregator
```go
// Combines all specialized resolvers
type PatientResolverAggregator struct {
    QueryResolver        *query.PatientQueryResolver
    MutationResolver     *mutation.PatientMutationResolver
    FieldResolver        *field.PatientFieldResolver
    AddressFieldResolver *field.AddressFieldResolver
}

// Delegates to appropriate resolver
func (r *PatientResolverAggregator) Patient(ctx context.Context, id string) (*model.Patient, error) {
    return r.QueryResolver.Patient(ctx, id)
}

func (r *PatientResolverAggregator) CreatePatient(ctx context.Context, input generated.CreatePatientInput) (*model.Patient, error) {
    return r.MutationResolver.CreatePatient(ctx, input)
}
```

### Benefits
- ✅ Maximum separation of concerns
- ✅ Team scalability (query team, mutation team)
- ✅ Different dependencies per concern
- ✅ Easier isolated testing
- ✅ Clear evolution to microservices
- ✅ CQRS-ready architecture

### Limitations
- ❌ More complexity
- ❌ More files to navigate
- ❌ Overkill for small/medium domains
- ❌ More boilerplate code

---

## 🔄 Evolution Decision Tree

```
Start: Is domain new?
  ├─ YES → Phase 1 (single resolver.go)
  └─ NO → Continue

Check: Is file > 300 lines OR sub-resource has complex logic?
  ├─ YES → Phase 2 (split sub-resource)
  └─ NO → Stay Phase 1

Check: Is total > 500 lines OR multiple teams OR different dependencies?
  ├─ YES → Phase 3 (full separation)
  └─ NO → Stay Phase 2
```

---

## 📖 Patient Domain Example (All Phases)

### Current State: Phase 2 ✅

```
domain/patient/graphql/
├── schema.graphql                    # Types
├── patient_resolver.go               # ✅ ACTIVE: Patient operations
├── address_resolver.go               # ✅ ACTIVE: Address operations
├── README.md                         # Quick start guide
└── RESOLVER_PHASES.md                # Complete phase patterns
```

**Why Phase 2?**
- ✅ Patient domain is 200-300 lines
- ✅ Address has specific logic (validation, formatting)
- ✅ Clear separation without over-engineering
- ✅ Perfect balance for current needs

**When to move to Phase 3?**
- ⏳ File grows beyond 500 lines
- ⏳ Adding complex mutations (create, update, delete)
- ⏳ Multiple teams work on patient domain
- ⏳ Queries need cache service, mutations need event publisher

---

## 🎓 Best Practices

### Across All Phases

**✅ DO:**
- Keep business logic in services, not resolvers
- Use delegation pattern (resolver → sub-resolver)
- Log errors at resolver level
- Return null for not-found, error for failures
- Test resolvers with mock services

**❌ DON'T:**
- Put business logic in resolvers
- Have resolvers call each other directly (use aggregator)
- Skip error handling or logging
- Return errors for not-found (return null instead)
- Create deep resolver hierarchies

### Phase-Specific

**Phase 1:**
- ✅ Keep file under 300 lines
- ✅ Move to Phase 2 when sub-resources grow

**Phase 2:**
- ✅ Create sub-resolver when > 50 lines of related logic
- ✅ Use composition (main resolver creates sub-resolvers)
- ✅ Move to Phase 3 when total > 500 lines

**Phase 3:**
- ✅ Group by concern (query/mutation/field)
- ✅ Use aggregator pattern
- ✅ Different dependencies per resolver type
- ✅ Consider DataLoader for N+1 queries

---

## 🚀 Quick Reference

### Using Current Implementation (Phase 2)

```go
// In internal/graphql/server.go
patientResolver := patientgraphql.NewPatientResolver(
    deps.PatientService,
    deps.AddressService,
    deps.PrescriptionService,
    deps.Logger,
)
```

### Switching to Phase 3 (When Needed)

```go
// In internal/graphql/server.go
patientResolver := patientresolvers.NewPatientResolverAggregator(
    deps.PatientService,
    deps.AddressService,
    deps.PrescriptionService,
    deps.Logger,
)
```

**Same interface, different internal organization!** 🎯

---

## 📚 Additional Resources

- **Patient Examples**: `domain/patient/graphql/RESOLVER_PHASES.md`
- **Architecture Diagrams**: `docs/GRAPHQL_ARCHITECTURE_DIAGRAMS.md`
- **Domain Structure**: `docs/GRAPHQL_DOMAIN_STRUCTURE.md`
- **Implementation Guide**: `docs/GRAPHQL_IMPLEMENTATION.md`

---

## 💡 Summary

| Phase | Lines | Files | Status | Use Case |
|-------|-------|-------|--------|----------|
| **1** | < 300 | 1 | 📖 Documented | New/small domains |
| **2** | 300-500 | 2-3 | ✅ **Active** | **Current - Perfect!** |
| **3** | 500+ | 5+ | 📖 Documented | Large enterprise domains |

**The patient domain shows the evolution:**
- Phase 1: Documented pattern (simple, single file)
- Phase 2: Active implementation ✅ (patient_resolver.go + address_resolver.go)
- Phase 3: Documented pattern (create when needed)

**Start simple, evolve as needed!** 🚀

