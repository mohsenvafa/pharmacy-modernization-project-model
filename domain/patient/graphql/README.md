# Patient Domain GraphQL

This directory demonstrates GraphQL resolver organization with **Phase 2 implementation** and **documented patterns** for evolution.

## 📁 Directory Structure

```
domain/patient/graphql/
├── schema.graphql          # GraphQL types
├── patient_resolver.go     # ✅ Patient operations (Phase 2)
├── address_resolver.go     # ✅ Address operations (Phase 2)
├── README.md               # This guide
└── RESOLVER_PHASES.md      # Complete evolution patterns & documentation
```

---

## 🎯 Implementation Status

### Phase 1: Documented Pattern (Simple Single Resolver)
**Pattern documented in `RESOLVER_PHASES.md`**
- Everything in one `domain_resolver.go` file
- Use when: < 300 lines, 1-3 services, small team
- See documentation for implementation pattern

---

### ✅ Phase 2: Active Implementation (Current)

**Files:**
- `patient_resolver.go` - Patient queries and prescription field resolution
- `address_resolver.go` - Address-specific operations

**Pattern: Composition + Delegation**

```go
// patient_resolver.go
type PatientResolver struct {
    PatientService      patientservice.PatientService
    PrescriptionService prescriptionservice.PrescriptionService
    AddressResolver     *AddressResolver  // Composed sub-resolver
    Logger              *zap.Logger
}

func (r *PatientResolver) Addresses(ctx context.Context, obj *model.Patient) ([]model.Address, error) {
    return r.AddressResolver.Addresses(ctx, obj)  // Delegation
}

// address_resolver.go
type AddressResolver struct {
    AddressService patientservice.AddressService
    Logger         *zap.Logger
}

func (r *AddressResolver) Addresses(ctx context.Context, obj *model.Patient) ([]model.Address, error) {
    // Actual implementation
}
```

**Why Phase 2?**
- ✅ Patient domain is 200-300 lines
- ✅ Address has specific logic
- ✅ Clean separation
- ✅ Not over-engineered

**Use this as your template for other domains!**

---

### Phase 3: Documented Pattern (For Large Domains)

**Pattern documented in `RESOLVER_PHASES.md`**

**Structure (create when needed):**
```
domain/patient/graphql/
└── resolvers/
    ├── query/
    │   └── patient_query_resolver.go      # Read operations
    ├── mutation/
    │   └── patient_mutation_resolver.go   # Write operations
    ├── field/
    │   ├── patient_field_resolver.go      # Complex patient fields
    │   └── address_field_resolver.go      # Complex address fields
    └── aggregator.go                      # Combines all resolvers
```

**Pattern: Full Separation by Concern**
- Separate files for queries, mutations, and fields
- Each resolver can have different dependencies
- Aggregator pattern combines them

**When to implement Phase 3:**
- ⏳ Domain grows beyond 500 lines
- ⏳ Multiple teams work on domain
- ⏳ Need different dependencies (cache for queries, events for mutations)
- ⏳ Planning microservices extraction

**Implementation details in `RESOLVER_PHASES.md`**

---

## 📖 Documentation

### In This Directory
- **`RESOLVER_PHASES.md`** - Detailed guide for all three phases with implementation patterns

### Project-Wide
- **`docs/GRAPHQL_RESOLVER_PATTERNS.md`** - Complete evolution guide
- **`docs/GRAPHQL_DOMAIN_STRUCTURE.md`** - Domain-based organization
- **`docs/GRAPHQL_ARCHITECTURE_DIAGRAMS.md`** - Visual diagrams
- **`docs/GRAPHQL_IMPLEMENTATION.md`** - API usage guide

---

## 🚀 Quick Start

### Using Current Implementation (Phase 2)

```go
// Already wired in internal/graphql/server.go
patientResolver := patientgraphql.NewPatientResolver(
    deps.PatientService,
    deps.AddressService,
    deps.PrescriptionService,
    deps.Logger,
)
```

### Implementing Phase 3 (When Needed)

1. **Create structure:**
```bash
mkdir -p domain/patient/graphql/resolvers/{query,mutation,field}
```

2. **Follow patterns in `RESOLVER_PHASES.md`** to implement:
   - Query resolver
   - Mutation resolver
   - Field resolvers
   - Aggregator

3. **Wire in server.go** - Same interface, different internal organization!

---

## 🎓 Best Practices

### For Your Next Domain

**Start Simple:**
```
1. Create domain/your_domain/graphql/
2. Add schema.graphql
3. Add your_domain_resolver.go (Phase 1)
4. Build and test
```

**Evolve When Needed:**
```
When file > 250 lines:
→ Extract sub-resource into separate resolver (Phase 2)

When total > 500 lines OR multiple teams:
→ Use Phase 3 pattern from RESOLVER_PHASES.md
```

### Keep Resolvers Thin
✅ **DO:**
- Call services
- Log errors
- Handle not-found gracefully
- Delegate to sub-resolvers

❌ **DON'T:**
- Put business logic in resolvers
- Make resolvers call each other directly
- Skip error handling
- Return errors for not-found (return null)

---

## 📊 Evolution Path

```
Phase 1 (Simple)
    ↓ When file > 300 lines
Phase 2 (Current) ✅
    ↓ When total > 500 lines OR multiple teams
Phase 3 (Documented) 📖
```

---

## 🎯 Summary

This directory demonstrates:

1. **Phase 2 Active Implementation** ✅
   - `patient_resolver.go` + `address_resolver.go`
   - Clean separation, delegation pattern
   - Use this as your template!

2. **Complete Evolution Documentation** 📖
   - Phase 1, 2, and 3 patterns in `RESOLVER_PHASES.md`
   - Implementation guides for each phase
   - Decision criteria for when to evolve

3. **Clean Codebase** 🧹
   - No unused example code
   - Only active implementation
   - Documentation shows the way forward

**For your new domains:**
- Copy Phase 2 pattern (patient_resolver.go + sub_resolver.go)
- Reference `RESOLVER_PHASES.md` when you need Phase 3
- Follow the evolution path: Phase 1 → Phase 2 → Phase 3

**Clean, documented, and ready to scale! 🚀**
