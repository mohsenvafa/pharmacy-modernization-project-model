# GraphQL Migration Phases

This document outlines the evolution path from a simple monolithic GraphQL implementation to a fully federated microservices architecture. Each phase builds on the previous one, allowing incremental adoption based on actual needs.

## Current State: Phase 1 (Implemented)

### Simple Monolithic GraphQL

**Architecture:**
```
Client
  ↓
Single Go Application (port 8080)
  ├── REST API (/api/*)
  ├── UI Handlers (/patients, /prescriptions, etc.)
  └── GraphQL (/graphql, /playground)
       ↓
    Resolvers (one per domain)
       ↓
    Services (existing business logic)
       ↓
    Repositories (MongoDB/Memory)
```

**Structure:**
```
domain/
  patient/
    api/              # REST API
    graphql/          # GraphQL resolvers
    service/          # Shared business logic
    repository/
    ui/
  prescription/
    api/
    graphql/
    service/
    ...
  dashboard/
    graphql/
    service/
    ...

internal/
  graphql/            # GraphQL server infrastructure
    ├── server.go          # GraphQL server setup
    ├── resolver.go        # Root resolver
    └── generated/         # gqlgen generated code
```

**Benefits:**
- ✅ Zero network overhead (all in-process)
- ✅ Simple deployment (single binary)
- ✅ Easy debugging
- ✅ Shared middleware and error handling
- ✅ Fast development cycles

**Limitations:**
- ❌ Cannot scale domains independently
- ❌ Single point of failure
- ❌ All domains must use same language/runtime
- ❌ Coupled deployment (all domains deploy together)

**When to Stay Here:**
- Team size < 20 developers
- Monolithic architecture is working well
- No need for polyglot services
- Simple deployment requirements
- Development velocity is more important than independent scaling

---

## Phase 2: Internal Subgraphs (Preparation)

### Monolith with Federation-Ready Schemas

Reorganize GraphQL schemas to follow federation patterns **without** splitting services. This prepares the codebase for future extraction.

**Architecture:**
```
Single Go Application (port 8080)
  └── Internal GraphQL Gateway
       ├── Patient Subgraph (in-process)
       ├── Prescription Subgraph (in-process)
       └── Dashboard Subgraph (in-process)
```

**Changes:**
1. Add federation directives to schemas (`@key`, `@shareable`, etc.)
2. Implement entity resolvers
3. Use gqlgen federation plugin
4. Test cross-domain references

**Example Schema Change:**

**Before (Phase 1):**
```graphql
type Patient {
  id: ID!
  firstName: String!
  prescriptions: [Prescription!]!  # Direct resolver
}
```

**After (Phase 2):**
```graphql
type Patient @key(fields: "id") {
  id: ID!
  firstName: String!
  # prescriptions moved to prescription subgraph
}

# In prescription subgraph
extend type Patient @key(fields: "id") {
  id: ID! @external
  prescriptions: [Prescription!]!
}
```

**Implementation Steps:**
1. Update schemas with federation directives
2. Split root resolver into subgraph resolvers
3. Add entity resolvers for each domain
4. Create internal "gateway" that routes within process
5. Update resolver code to handle entity resolution

**Benefits:**
- ✅ Schema is federation-ready
- ✅ Easy to extract services later
- ✅ Learn federation patterns safely
- ✅ Still runs as monolith

**When to Move to Phase 2:**
- Planning microservices migration
- Want to prepare codebase incrementally
- Testing federation concepts
- Multiple teams working on different domains

---

## Phase 3: Hybrid Architecture

### Extract First Service + Gateway

Extract one domain as a separate service while keeping others in the monolith. This validates the federation approach with minimal risk.

**Architecture:**
```
Client
  ↓
Apollo Gateway (port 8080)
  ├→ Patient Service (Go, port 8081) ← EXTRACTED
  └→ Core Service (Go, port 8082)     ← Monolith
       ├── Prescription subgraph
       └── Dashboard subgraph
```

**Project Structure:**
```
projects/
  ├── gateway/                    # NEW: Apollo Gateway
  │   ├── main.go
  │   ├── go.mod
  │   └── config/
  │       └── supergraph.yaml
  │
  ├── patient-service/            # NEW: Extracted service
  │   ├── cmd/server/main.go
  │   ├── domain/patient/
  │   ├── internal/
  │   └── go.mod
  │
  └── core-service/               # EXISTING: Monolith
      ├── cmd/server/main.go
      ├── domain/
      │   ├── prescription/
      │   └── dashboard/
      └── go.mod
```

**Migration Steps:**

1. **Setup Gateway:**
   ```bash
   # Create gateway project
   mkdir gateway
   cd gateway
   go mod init pharmacy/gateway
   go get github.com/apollographql/router
   ```

2. **Extract Patient Service:**
   - Copy `domain/patient/` to new service
   - Copy shared platform code
   - Update imports
   - Configure separate database connection
   - Deploy on different port

3. **Configure Gateway:**
   ```yaml
   # gateway/config/supergraph.yaml
   subgraphs:
     patient:
       routing_url: http://patient-service:8081/graphql
     core:
       routing_url: http://core-service:8082/graphql
   ```

4. **Update Deployment:**
   - Deploy gateway (port 8080)
   - Deploy patient-service (port 8081)
   - Deploy core-service (port 8082)
   - Update load balancer to point to gateway

**Benefits:**
- ✅ Validate federation with one service
- ✅ Reduced blast radius if issues occur
- ✅ Learn operational complexities
- ✅ Independent deployment for extracted service
- ✅ Can scale patient service independently

**Challenges:**
- ⚠️ Need to manage 3 deployments
- ⚠️ Network latency between services
- ⚠️ Distributed tracing required
- ⚠️ More complex local development

**When to Move to Phase 3:**
- One domain has different scaling needs
- Ready to split teams
- Proven federation patterns work
- Have DevOps capacity for multiple services

---

## Phase 4: Full Federation

### All Domains as Separate Services

Extract all domains into independent microservices with Apollo Gateway orchestrating them.

**Architecture:**
```
Client
  ↓
Apollo Gateway (port 8080)
  ├→ Patient Service (Go, port 8081)
  ├→ Prescription Service (Go, port 8082)
  ├→ Dashboard Service (Go, port 8083)
  └→ Future Services...
```

**Full Microservices Structure:**
```
projects/
  ├── gateway/
  │   ├── main.go
  │   └── config/supergraph.yaml
  │
  ├── patient-service/
  │   ├── cmd/server/main.go
  │   ├── domain/patient/
  │   ├── internal/platform/
  │   └── database/patients_db
  │
  ├── prescription-service/
  │   ├── cmd/server/main.go
  │   ├── domain/prescription/
  │   ├── internal/platform/
  │   └── database/prescriptions_db
  │
  └── dashboard-service/
      ├── cmd/server/main.go
      ├── domain/dashboard/
      └── internal/platform/
```

**Cross-Service Communication:**

**Option A: GraphQL Only (Recommended)**
- Services only talk through gateway
- Client queries drive all data fetching
- No service-to-service calls

**Option B: Direct Service Calls**
```go
// In prescription-service
type HTTPPatientClient struct {
    baseURL string
}

func (c *HTTPPatientClient) GetPatient(id string) (*Patient, error) {
    resp, err := http.Get(c.baseURL + "/api/patients/" + id)
    // ... handle response
}
```

**Option C: Event-Driven**
```go
// Patient service publishes
eventBus.Publish("patient.created", PatientCreatedEvent{ID: "123"})

// Prescription service subscribes
eventBus.Subscribe("patient.created", func(event PatientCreatedEvent) {
    // Update local cache or trigger action
})
```

**Benefits:**
- ✅ Full service independence
- ✅ Polyglot possible (Go, Node.js, Python)
- ✅ Independent scaling per domain
- ✅ Team autonomy
- ✅ Isolated failures

**Challenges:**
- ⚠️ Complex deployment orchestration
- ⚠️ Distributed transactions (avoid if possible)
- ⚠️ Network reliability issues
- ⚠️ Distributed tracing essential
- ⚠️ Testing requires all services
- ⚠️ Higher operational overhead

**When to Move to Phase 4:**
- Large engineering organization (50+ devs)
- Clear domain boundaries
- Different scaling requirements per domain
- Team wants full autonomy
- Have mature DevOps/SRE practices

---

## Advanced Patterns (Phase 5+)

### DataLoader (Performance Optimization)

**Problem:**
```graphql
query {
  prescriptions {           # Returns 100 items
    id
    patient {               # 100 separate queries!
      firstName
    }
  }
}
```

**Solution: Batching + Caching**
```go
// internal/graphql/dataloader/patient_loader.go
package dataloader

import (
    "context"
    "github.com/graph-gophers/dataloader/v7"
)

type PatientLoader struct {
    loader *dataloader.Loader[string, *model.Patient]
}

func NewPatientLoader(svc service.PatientService) *PatientLoader {
    batchFn := func(ctx context.Context, keys []string) []*dataloader.Result[*model.Patient] {
        // Batch fetch all patients in one call
        patients, err := svc.GetByIDs(ctx, keys)
        
        // Map results back to keys
        results := make([]*dataloader.Result[*model.Patient], len(keys))
        for i, key := range keys {
            results[i] = &dataloader.Result[*model.Patient]{
                Data: patients[key],
                Error: err,
            }
        }
        return results
    }
    
    return &PatientLoader{
        loader: dataloader.NewBatchedLoader(batchFn),
    }
}

func (l *PatientLoader) Load(ctx context.Context, id string) (*model.Patient, error) {
    return l.loader.Load(ctx, id)()
}
```

**Usage in Resolver:**
```go
// domain/prescription/graphql/resolver.go
func (r *prescriptionResolver) Patient(ctx context.Context, obj *model.Prescription) (*model.Patient, error) {
    // Get loader from context
    loader := dataloader.GetPatientLoader(ctx)
    return loader.Load(ctx, obj.PatientID)
}
```

**When to Add:**
- Queries show N+1 patterns in logs/APM
- Response times > 500ms due to repeated queries
- High database load from duplicate queries

---

### Managed Federation (Apollo GraphOS)

**What It Is:**
Cloud service that manages schema composition, routing, and observability.

**Features:**
- 🔄 Schema registry and versioning
- 📊 Query analytics and performance
- 🚀 Zero-downtime schema deployments
- 🔍 Schema validation and checks
- 📈 Field-level usage metrics

**Migration:**
1. Sign up for Apollo GraphOS
2. Publish schemas to registry
3. Replace self-hosted gateway with Apollo Router
4. Configure managed federation

**When to Add:**
- Multiple teams deploying independently
- Need schema governance
- Want advanced analytics
- Require contract-based schema evolution

---

### Service Mesh (Advanced Networking)

**What It Is:**
Infrastructure layer that handles service-to-service communication.

**Tools:** Istio, Linkerd, Consul Connect

**Features:**
- 🔒 mTLS between services
- 🔄 Automatic retries and circuit breaking
- 📊 Distributed tracing
- 🚦 Traffic splitting for canary deployments

**When to Add:**
- 10+ microservices
- Complex networking requirements
- Need fine-grained traffic control
- Security compliance requires mTLS

---

## Decision Framework

### Should You Move to Next Phase?

Use this checklist to decide if you're ready:

**Move to Phase 2 (Internal Subgraphs) if:**
- [ ] Planning to extract services in 6-12 months
- [ ] Want to prepare incrementally
- [ ] Have time to learn federation
- [ ] No urgent production needs

**Move to Phase 3 (First Extraction) if:**
- [ ] One domain has 3x+ different load than others
- [ ] Team is growing (10+ developers)
- [ ] Clear domain boundaries exist
- [ ] Have CI/CD pipeline for multiple services
- [ ] DevOps/SRE capacity exists

**Move to Phase 4 (Full Federation) if:**
- [ ] Multiple domains need independent scaling
- [ ] Different programming languages needed
- [ ] Large team (50+ developers)
- [ ] Mature microservices practices
- [ ] Budget for operational complexity

**Add DataLoader if:**
- [ ] GraphQL queries > 500ms
- [ ] Database shows N+1 query patterns
- [ ] Slow nested queries in production

**Add Managed Federation if:**
- [ ] 5+ subgraph services
- [ ] Multiple teams deploying independently
- [ ] Need schema analytics
- [ ] Want zero-downtime deployments

---

## Anti-Patterns to Avoid

### ❌ Premature Federation

**Don't:**
- Split to microservices "because everyone does it"
- Extract before domain boundaries are clear
- Federate with < 5 developers

**Instead:**
- Stay in Phase 1 until you have a specific reason
- Validate boundaries in monolith first

### ❌ Chatty Resolvers

**Don't:**
```graphql
query {
  patient(id: "123") {
    prescriptions {
      # Each calls prescription service
      medication { id, name }  # Calls medication service
      pharmacy { id, name }    # Calls pharmacy service
    }
  }
}
```
This creates a cascade of 4+ service calls!

**Instead:**
- Batch with DataLoader
- Denormalize data for read models
- Use BFF (Backend for Frontend) pattern

### ❌ Distributed Transactions

**Don't:**
```go
// In prescription service
tx.Begin()
CreatePrescription(...)      // Local
UpdatePatientRecord(...)     // Remote call to patient service ⚠️
tx.Commit()                  // Can't rollback remote!
```

**Instead:**
- Use eventual consistency
- Event-driven architecture
- Saga pattern for complex workflows

### ❌ Shared Database

**Don't:**
```
Patient Service → Shared PostgreSQL ← Prescription Service
```

**Instead:**
- Each service owns its data
- Cross-service queries through GraphQL
- Use events for data synchronization

---

## Testing Strategy by Phase

### Phase 1 (Monolith)
```go
// Simple integration test
func TestGraphQL_GetPatient(t *testing.T) {
    app := setupTestApp()
    query := `{ patient(id: "123") { firstName } }`
    
    resp := executeQuery(app, query)
    assert.Equal(t, "John", resp.Data.Patient.FirstName)
}
```

### Phase 3-4 (Federation)
```go
// Gateway integration test
func TestFederation_PatientWithPrescriptions(t *testing.T) {
    // Start mock services
    patientService := startMockPatientService()
    prescriptionService := startMockPrescriptionService()
    gateway := startGateway()
    
    query := `{ patient(id: "123") { firstName, prescriptions { medication } } }`
    
    resp := executeQuery(gateway, query)
    assert.NotEmpty(t, resp.Data.Patient.Prescriptions)
}
```

---

## Cost Analysis

### Phase 1: Monolith
- **Development:** 1 week initial setup
- **Deployment:** 1 server/container
- **Operations:** Minimal (same as current)
- **Complexity:** Low

### Phase 3: Hybrid
- **Development:** 2-3 weeks extraction + gateway
- **Deployment:** 3 servers/containers (gateway + 2 services)
- **Operations:** Moderate (monitoring, distributed tracing)
- **Complexity:** Medium

### Phase 4: Full Federation
- **Development:** 1-2 months full migration
- **Deployment:** 4+ servers/containers
- **Operations:** High (service mesh, observability, alerting)
- **Complexity:** High

---

## Conclusion

**Start with Phase 1** (current implementation). Only move to later phases when you have clear, measurable pain points that justify the added complexity.

**Key Principle:** Add complexity only when the benefits outweigh the costs.

Most teams will be successful staying in Phase 1 for years. The architecture is designed to make migration possible if needed, but not required.

