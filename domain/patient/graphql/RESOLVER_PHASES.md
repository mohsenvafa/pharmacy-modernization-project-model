# Patient GraphQL Resolver Phases

This directory demonstrates three phases of resolver organization, showing how to evolve from simple to complex as your domain grows.

## 📁 Current Structure

```
domain/patient/graphql/
├── schema.graphql                    # GraphQL schema
├── resolver.go                       # ✅ ACTIVE: Phase 2 implementation
├── address_resolver.go               # ✅ ACTIVE: Phase 2 implementation
│
└── resolvers/                        # 📚 Phase 3 examples (reference)
    ├── query/
    │   └── patient_query_resolver.go
    ├── mutation/
    │   └── patient_mutation_resolver.go
    ├── field/
    │   ├── patient_field_resolver.go
    │   └── address_field_resolver.go
    └── aggregator_example.go
```

---

## Phase 1: Simple Single Resolver (Small Domains)

**Use when:**
- Domain is < 300 lines
- 1-3 services
- Small team (1-2 developers)

**Structure:**
```
domain/patient/graphql/
├── schema.graphql
└── resolver.go         # Everything in one file
```

**Example:**
```go
type PatientResolver struct {
    PatientService      patientservice.PatientService
    AddressService      patientservice.AddressService
    PrescriptionService prescriptionservice.PrescriptionService
    Logger              *zap.Logger
}

// All queries, mutations, and fields in one file
func (r *PatientResolver) Patient(ctx context.Context, id string) (*model.Patient, error) { }
func (r *PatientResolver) Addresses(ctx context.Context, obj *model.Patient) ([]model.Address, error) { }
func (r *PatientResolver) CreatePatient(ctx context.Context, input) (*model.Patient, error) { }
```

**Benefits:**
- ✅ Simple
- ✅ Easy to navigate
- ✅ Fast development

**Limitations:**
- ❌ Large file as domain grows
- ❌ Mixed concerns

---

## Phase 2: Split by Sub-Resource ✅ (CURRENT)

**Use when:**
- Sub-resources have significant logic
- Want better organization
- File approaching 300+ lines

**Structure:**
```
domain/patient/graphql/
├── schema.graphql
├── patient_resolver.go   # ✅ Patient operations
└── address_resolver.go   # ✅ Address operations
```

**Active Implementation:**

**`patient_resolver.go`:**
```go
type PatientResolver struct {
    PatientService      patientservice.PatientService
    PrescriptionService prescriptionservice.PrescriptionService
    AddressResolver     *AddressResolver  // Delegates address ops
    Logger              *zap.Logger
}

func (r *PatientResolver) Addresses(ctx context.Context, obj *model.Patient) ([]model.Address, error) {
    return r.AddressResolver.Addresses(ctx, obj)  // Delegate!
}
```

**`address_resolver.go`:**
```go
type AddressResolver struct {
    AddressService patientservice.AddressService
    Logger         *zap.Logger
}

func (r *AddressResolver) Addresses(ctx context.Context, obj *model.Patient) ([]model.Address, error) {
    // Address-specific logic here
}

func (r *AddressResolver) ValidateAddress(ctx context.Context, obj *model.Address) (bool, error) {
    // Complex validation logic
}
```

**Benefits:**
- ✅ Better organization
- ✅ Clear boundaries
- ✅ Focused responsibilities
- ✅ Still simple to understand

**When to use:**
- ✅ Current patient domain size
- ✅ Address has specific logic
- ✅ Want to separate concerns without over-engineering

---

## Phase 3: Full Separation (Large Domains)

**Use when:**
- Domain is 500+ lines
- 5+ different services
- Multiple teams working on domain
- Need different dependencies for queries vs mutations

**Structure:**
```
domain/patient/graphql/
├── schema.graphql
├── resolvers/
│   ├── query/
│   │   └── patient_query_resolver.go      # Read operations
│   ├── mutation/
│   │   └── patient_mutation_resolver.go   # Write operations
│   ├── field/
│   │   ├── patient_field_resolver.go      # Complex patient fields
│   │   └── address_field_resolver.go      # Complex address fields
│   └── aggregator_example.go              # Combines all resolvers
```

**Example Files (Reference Implementation):**

See:
- `resolvers/query/patient_query_resolver.go`
- `resolvers/mutation/patient_mutation_resolver.go`
- `resolvers/field/patient_field_resolver.go`
- `resolvers/field/address_field_resolver.go`
- `resolvers/aggregator_example.go`

**Benefits:**
- ✅ Maximum separation of concerns
- ✅ Team scalability (query team, mutation team, field team)
- ✅ Different dependencies per type
- ✅ Easier testing in isolation
- ✅ Clear evolution path for microservices

**Limitations:**
- ❌ More complex
- ❌ More files to navigate
- ❌ Overkill for small domains

---

## 📊 Decision Matrix

| Metric | Phase 1 | Phase 2 ✅ | Phase 3 |
|--------|---------|-----------|---------|
| **Total Lines** | < 300 | 300-500 | 500+ |
| **Files** | 1 | 2-3 | 5+ |
| **Services** | 1-3 | 3-5 | 5+ |
| **Team Size** | 1-2 | 2-4 | 4+ |
| **Complexity** | Low | Medium | High |
| **Recommended** | MVP/Small | **Current** | Enterprise |

---

## 🔄 Evolution Path

### Start Here (Phase 1)
```
resolver.go (150 lines)
```

### Grow to Phase 2 (Current) ✅
```
When file hits 250+ lines:
→ Split address into address_resolver.go

patient_resolver.go (180 lines)
address_resolver.go (70 lines)
```

### Evolve to Phase 3 (If Needed)
```
When combined > 500 lines OR multiple teams:
→ Split into query/mutation/field folders

resolvers/
  query/patient_query_resolver.go (150 lines)
  mutation/patient_mutation_resolver.go (120 lines)
  field/patient_field_resolver.go (100 lines)
  field/address_field_resolver.go (80 lines)
  aggregator_example.go (100 lines)
```

---

## 🎯 Current Implementation (Phase 2)

**Active Files:**
- ✅ `patient_resolver.go` - Patient queries and field resolution
- ✅ `address_resolver.go` - Address-specific operations

**Why Phase 2 is Perfect Right Now:**
1. ✅ Patient operations are focused (patient queries + prescriptions)
2. ✅ Address has its own logic (validation, formatting)
3. ✅ Clean separation without over-engineering
4. ✅ Easy to understand and maintain
5. ✅ Ready to evolve to Phase 3 if needed

**When to Move to Phase 3:**
- ⏳ When patient_resolver.go grows beyond 300 lines
- ⏳ When adding complex mutations (create, update, delete)
- ⏳ When multiple teams work on patient domain
- ⏳ When queries need different dependencies than mutations

---

## 📚 Implementation Guide for Phase 3

When your domain reaches Phase 3 size, create this structure:

```bash
mkdir -p domain/patient/graphql/resolvers/{query,mutation,field}
```

Then implement following the patterns documented above:

1. **Query Resolver** - Create `resolvers/query/patient_query_resolver.go`
2. **Mutation Resolver** - Create `resolvers/mutation/patient_mutation_resolver.go`
3. **Field Resolvers** - Create `resolvers/field/patient_field_resolver.go`
4. **Aggregator** - Create `resolvers/aggregator.go` to combine them

Each file should follow the delegation pattern shown in the examples above.

---

## 🎓 Best Practices

### ✅ DO:
- Start with Phase 1 for new domains
- Move to Phase 2 when sub-resources have significant logic
- Move to Phase 3 when file size or team size demands it
- Keep business logic in services, not resolvers
- Use delegation pattern at all phases

### ❌ DON'T:
- Jump straight to Phase 3 for small domains
- Mix business logic in resolvers
- Create resolvers that call other resolvers (except via aggregator)
- Forget to update this README when you change patterns

---

## 🚀 Quick Start

### Using Phase 2 (Current):
```go
// In internal/graphql/server.go
patientResolver := patientgraphql.NewPatientResolver(
    deps.PatientService,
    deps.AddressService,
    deps.PrescriptionService,
    deps.Logger,
)
```

### Switching to Phase 3:
```go
// In internal/graphql/server.go
patientResolver := patientgraphql.NewPatientResolverAggregator(
    deps.PatientService,
    deps.AddressService,
    deps.PrescriptionService,
    deps.Logger,
)
```

Same interface, different implementation! 🎯

---

## 📖 Summary

- **Phase 1**: Everything in one file (simple domains)
- **Phase 2**: Split by sub-resource ✅ **(Current - Perfect for patient domain)**
- **Phase 3**: Full separation by query/mutation/field (large domains)

**Implement Phase 3 structure when your domain needs it, following the patterns described above!**

