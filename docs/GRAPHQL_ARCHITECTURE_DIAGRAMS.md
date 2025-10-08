# GraphQL Architecture Diagrams

Visual representations of how GraphQL works in the rxintake_scaffold application.

## Table of Contents
1. [High-Level Request Flow](#high-level-request-flow)
2. [Detailed Component Architecture](#detailed-component-architecture)
3. [Query Resolution Flow](#query-resolution-flow)
4. [Nested Query Resolution](#nested-query-resolution)
5. [Comparison: REST vs GraphQL](#comparison-rest-vs-graphql)
6. [Code Generation Workflow](#code-generation-workflow)

---

## High-Level Request Flow

```mermaid
sequenceDiagram
    participant Client
    participant Router
    participant GraphQL Server
    participant Resolver
    participant Service
    participant Repository
    participant Database

    Client->>Router: POST /graphql<br/>{query: "{ patient(id: \"123\") { name } }"}
    Router->>GraphQL Server: Handle GraphQL Request
    GraphQL Server->>GraphQL Server: Parse & Validate Query
    GraphQL Server->>Resolver: Execute Query
    Resolver->>Service: GetByID(ctx, "123")
    Service->>Repository: GetByID(ctx, "123")
    Repository->>Database: SELECT * FROM patients WHERE id=123
    Database-->>Repository: Patient Row
    Repository-->>Service: Patient Model
    Service-->>Resolver: Patient Model
    Resolver-->>GraphQL Server: GraphQL Response
    GraphQL Server-->>Router: JSON Response
    Router-->>Client: {"data": {"patient": {"name": "John Doe"}}}
```

---

## Detailed Component Architecture

```mermaid
graph TB
    subgraph "Client Layer"
        WebApp[Web App]
        Mobile[Mobile App]
        CLI[CLI Tool]
    end

    subgraph "HTTP Layer (Chi Router)"
        Router[Chi Router<br/>:8080]
        MW[Middleware<br/>- RequestID<br/>- Logging<br/>- Recovery<br/>- Timeout]
    end

    subgraph "GraphQL Layer (internal/graphql/)"
        Server[GraphQL Server<br/>server.go<br/>/graphql endpoint]
        Playground[Playground Handler<br/>/playground]
        Generated[Generated Code<br/>gqlgen generated/]
        Schema[Schema Definition<br/>schema.graphql]
        RootResolver[Root Resolver<br/>resolver.go]
        Resolvers[Query/Mutation Resolvers<br/>schema.resolvers.go]
    end

    subgraph "Domain Services"
        PatientSvc[Patient Service<br/>domain/patient/service/]
        PrescriptionSvc[Prescription Service<br/>domain/prescription/service/]
        DashboardSvc[Dashboard Service<br/>domain/dashboard/service/]
        AddressSvc[Address Service<br/>domain/patient/service/]
    end

    subgraph "Data Layer"
        PatientRepo[Patient Repository]
        PrescriptionRepo[Prescription Repository]
        AddressRepo[Address Repository]
    end

    subgraph "Storage"
        MongoDB[(MongoDB)]
        Memory[(In-Memory)]
    end

    WebApp -->|GraphQL Query| Router
    Mobile -->|GraphQL Query| Router
    CLI -->|GraphQL Query| Router
    
    Router --> MW
    MW --> Server
    MW --> Playground
    
    Server --> Generated
    Server --> RootResolver
    Generated --> Schema
    RootResolver --> Resolvers
    
    Resolvers -->|Calls| PatientSvc
    Resolvers -->|Calls| PrescriptionSvc
    Resolvers -->|Calls| DashboardSvc
    Resolvers -->|Calls| AddressSvc
    
    PatientSvc --> PatientRepo
    PrescriptionSvc --> PrescriptionRepo
    AddressSvc --> AddressRepo
    
    PatientRepo --> MongoDB
    PrescriptionRepo --> Memory
    AddressRepo --> Memory
    
    style Server fill:#e1f5ff
    style Resolvers fill:#e1f5ff
    style Schema fill:#fff4e1
```

---

## Query Resolution Flow

### Simple Query: Get Single Patient

```mermaid
flowchart LR
    subgraph "1. Query Arrives"
        Query["query {<br/>  patient(id: \"123\") {<br/>    name<br/>    phone<br/>  }<br/>}"]
    end
    
    subgraph "2. GraphQL Server"
        Parse[Parse Query]
        Validate[Validate Against Schema]
        Plan[Create Execution Plan]
    end
    
    subgraph "3. Root Resolver"
        QueryResolver[queryResolver.Patient<br/>id: \"123\"]
    end
    
    subgraph "4. Service Layer"
        PatientSvc[PatientService.GetByID<br/>ctx, \"123\"]
    end
    
    subgraph "5. Repository"
        PatientRepo[PatientRepository.GetByID<br/>ctx, \"123\"]
    end
    
    subgraph "6. Database"
        DB[(MongoDB/<br/>Memory)]
    end
    
    subgraph "7. Response"
        Response["{<br/>  data: {<br/>    patient: {<br/>      name: \"John Doe\"<br/>      phone: \"555-1234\"<br/>    }<br/>  }<br/>}"]
    end
    
    Query --> Parse --> Validate --> Plan
    Plan --> QueryResolver
    QueryResolver --> PatientSvc
    PatientSvc --> PatientRepo
    PatientRepo --> DB
    DB -.->|Patient Data| PatientRepo
    PatientRepo -.->|Patient Model| PatientSvc
    PatientSvc -.->|Patient Model| QueryResolver
    QueryResolver -.->|GraphQL Object| Response
```

---

## Nested Query Resolution

### Complex Query: Patient with Addresses and Prescriptions

```mermaid
sequenceDiagram
    autonumber
    
    participant Client
    participant GraphQL
    participant QueryResolver
    participant PatientResolver
    participant PatientSvc
    participant AddressSvc
    participant PrescriptionSvc
    
    Note over Client: Query with nested fields
    Client->>GraphQL: query {<br/>  patient(id: "123") {<br/>    name<br/>    addresses { city }<br/>    prescriptions { drug }<br/>  }<br/>}
    
    GraphQL->>QueryResolver: Resolve patient(id: "123")
    QueryResolver->>PatientSvc: GetByID(ctx, "123")
    PatientSvc-->>QueryResolver: Patient{ID: "123", Name: "John"}
    QueryResolver-->>GraphQL: Patient Object
    
    Note over GraphQL: Now resolve patient.addresses field
    GraphQL->>PatientResolver: Addresses(ctx, Patient{ID:"123"})
    PatientResolver->>AddressSvc: GetByPatientID(ctx, "123")
    AddressSvc-->>PatientResolver: []Address{{City: "NYC"}, {City: "LA"}}
    PatientResolver-->>GraphQL: []Address
    
    Note over GraphQL: Now resolve patient.prescriptions field
    GraphQL->>PatientResolver: Prescriptions(ctx, Patient{ID:"123"})
    PatientResolver->>PrescriptionSvc: List(ctx, "", 100, 0)
    PrescriptionSvc-->>PatientResolver: []Prescription (filtered by patientID)
    PatientResolver-->>GraphQL: []Prescription
    
    GraphQL-->>Client: {<br/>  data: {<br/>    patient: {<br/>      name: "John",<br/>      addresses: [{city: "NYC"}, {city: "LA"}],<br/>      prescriptions: [{drug: "Aspirin"}]<br/>    }<br/>  }<br/>}
```

### N+1 Query Problem Illustrated

```mermaid
graph TD
    subgraph "Query"
        Q["query {<br/>  patients {<br/>    name<br/>    prescriptions { drug }<br/>  }<br/>}"]
    end
    
    subgraph "Query 1: Get All Patients"
        P1[Patient 1]
        P2[Patient 2]
        P3[Patient 3]
        PN[Patient N]
    end
    
    subgraph "N Queries: Get Prescriptions Per Patient"
        RX1[Get Prescriptions<br/>for Patient 1]
        RX2[Get Prescriptions<br/>for Patient 2]
        RX3[Get Prescriptions<br/>for Patient 3]
        RXN[Get Prescriptions<br/>for Patient N]
    end
    
    Q --> P1 & P2 & P3 & PN
    P1 --> RX1
    P2 --> RX2
    P3 --> RX3
    PN --> RXN
    
    style Q fill:#ffcccc
    style RX1 fill:#ffcccc
    style RX2 fill:#ffcccc
    style RX3 fill:#ffcccc
    style RXN fill:#ffcccc
    
    Note1[⚠️ This creates 1 + N database calls<br/>Solution: Add DataLoader in Phase 2+]
    RXN -.-> Note1
```

---

## Comparison: REST vs GraphQL

### REST API Flow

```mermaid
sequenceDiagram
    participant Client
    participant REST as REST API<br/>/api/v1
    participant PatientCtrl as Patient Controller
    participant PrescriptionCtrl as Prescription Controller
    participant PatientSvc as Patient Service
    participant PrescriptionSvc as Prescription Service
    
    Note over Client: Need patient + prescriptions
    
    Client->>REST: GET /api/v1/patients/123
    REST->>PatientCtrl: GetPatient(123)
    PatientCtrl->>PatientSvc: GetByID(123)
    PatientSvc-->>PatientCtrl: Patient
    PatientCtrl-->>REST: JSON (ALL patient fields)
    REST-->>Client: Patient JSON
    
    Note over Client: Need separate request for prescriptions
    
    Client->>REST: GET /api/v1/prescriptions?patientId=123
    REST->>PrescriptionCtrl: ListPrescriptions
    PrescriptionCtrl->>PrescriptionSvc: List(...)
    PrescriptionSvc-->>PrescriptionCtrl: Prescriptions
    PrescriptionCtrl-->>REST: JSON (ALL prescription fields)
    REST-->>Client: Prescriptions JSON
    
    Note over Client: 2 HTTP requests<br/>Over-fetching: Got all fields<br/>Client must combine data
```

### GraphQL Flow

```mermaid
sequenceDiagram
    participant Client
    participant GraphQL as GraphQL API<br/>/graphql
    participant Resolvers
    participant PatientSvc as Patient Service
    participant PrescriptionSvc as Prescription Service
    
    Note over Client: Need patient + prescriptions
    
    Client->>GraphQL: POST /graphql<br/>{<br/>  patient(id: "123") {<br/>    name phone<br/>    prescriptions { drug dose }<br/>  }<br/>}
    GraphQL->>Resolvers: Execute Query
    Resolvers->>PatientSvc: GetByID(123)
    PatientSvc-->>Resolvers: Patient
    Resolvers->>PrescriptionSvc: List(patientId=123)
    PrescriptionSvc-->>Resolvers: Prescriptions
    Resolvers-->>GraphQL: Composed Response
    GraphQL-->>Client: {<br/>  patient: {<br/>    name, phone,<br/>    prescriptions: [drug, dose]<br/>  }<br/>}
    
    Note over Client: 1 HTTP request<br/>No over-fetching: Got only requested fields<br/>Data already combined
```

---

## Code Generation Workflow

```mermaid
flowchart TB
    subgraph "1. Define Schema"
        Schema[internal/graphql/schema.graphql<br/>---<br/>type Patient {<br/>  id: ID!<br/>  name: String!<br/>}]
    end
    
    subgraph "2. Configure gqlgen"
        Config[gqlgen.yml<br/>---<br/>schema: internal/graphql/*.graphql<br/>resolver: internal/graphql<br/>autobind: domain models]
    end
    
    subgraph "3. Run Generator"
        Cmd[$ gqlgen generate]
    end
    
    subgraph "4. Generated Code"
        GenCode[internal/graphql/generated/<br/>---<br/>- generated.go<br/>- models_gen.go]
        ResolverStubs[internal/graphql/<br/>schema.resolvers.go<br/>---<br/>func queryResolver Patient<br/>  panic not implemented]
    end
    
    subgraph "5. Implement Resolvers"
        ImplResolvers[internal/graphql/<br/>schema.resolvers.go<br/>---<br/>func queryResolver Patient {<br/>  return r.PatientService.GetByID<br/>}]
    end
    
    subgraph "6. Wire Dependencies"
        Wire[internal/app/wire.go<br/>---<br/>graphql.MountGraphQL r,<br/>  PatientService,<br/>  PrescriptionService,<br/>  ...]
    end
    
    subgraph "7. Server Ready"
        Server[GraphQL Server Running<br/>---<br/>✓ /graphql endpoint<br/>✓ /playground endpoint]
    end
    
    Schema --> Config
    Config --> Cmd
    Cmd --> GenCode
    Cmd --> ResolverStubs
    ResolverStubs --> ImplResolvers
    ImplResolvers --> Wire
    Wire --> Server
    
    style Cmd fill:#d4f1d4
    style ImplResolvers fill:#ffe6cc
    style Server fill:#d4e6f1
```

---

## Dependency Injection Flow

```mermaid
graph TD
    subgraph "Application Startup (internal/app/wire.go)"
        Main[main.go starts app]
        Wire[wire function]
    end
    
    subgraph "Module Initialization"
        PatientMod[Patient Module]
        PrescriptionMod[Prescription Module]
        DashboardMod[Dashboard Module]
    end
    
    subgraph "Services Created"
        PatientSvc[Patient Service]
        AddressSvc[Address Service]
        PrescriptionSvc[Prescription Service]
        DashboardSvc[Dashboard Service]
    end
    
    subgraph "GraphQL Setup (internal/graphql/)"
        Deps[GraphQL Dependencies{<br/>  PatientService<br/>  AddressService<br/>  PrescriptionService<br/>  DashboardService<br/>  Logger<br/>}]
        Resolver[Root Resolver<br/>resolver.go]
        Server[GraphQL Server<br/>server.go]
        Mount[Mount to Router<br/>/graphql<br/>/playground]
    end
    
    Main --> Wire
    Wire --> PatientMod
    Wire --> PrescriptionMod
    Wire --> DashboardMod
    
    PatientMod --> PatientSvc
    PatientMod --> AddressSvc
    PrescriptionMod --> PrescriptionSvc
    DashboardMod --> DashboardSvc
    
    PatientSvc --> Deps
    AddressSvc --> Deps
    PrescriptionSvc --> Deps
    DashboardSvc --> Deps
    
    Deps --> Resolver
    Resolver --> Server
    Server --> Mount
    
    style Deps fill:#e1f5ff
    style Resolver fill:#e1f5ff
    style Mount fill:#d4f1d4
```

---

## File Structure & Responsibilities

```
internal/graphql/
├── schema.graphql              📄 Schema Definition
│   ├─ Defines all types (Patient, Prescription, etc.)
│   ├─ Defines all queries (patient, patients, etc.)
│   └─ Defines all mutations (future)
│
├── server.go                   🖥️  Server Setup
│   ├─ MountGraphQL() function
│   ├─ Creates resolver with dependencies
│   ├─ Configures gqlgen handler
│   └─ Mounts /graphql and /playground endpoints
│
├── resolver.go                 🔧 Dependency Container
│   ├─ Resolver struct with all services
│   ├─ Injected by wire.go
│   └─ Never regenerated (you maintain this)
│
├── schema.resolvers.go         ⚙️  Resolver Implementations
│   ├─ Query resolvers (patient, patients, etc.)
│   ├─ Field resolvers (patient.addresses, etc.)
│   ├─ Mutation resolvers (future)
│   ├─ Calls domain services
│   └─ Regenerated with stubs, you implement logic
│
└── generated/                  🤖 Auto-Generated Code
    ├── generated.go            (GraphQL executor)
    └── models_gen.go           (Type mappings)
```

---

## Execution Model

### How a Query Gets Executed

```
1. Client sends query
   ↓
2. Chi Router receives POST /graphql
   ↓
3. Middleware chain executes
   - Request ID generation
   - Correlation ID
   - Logging
   - Panic recovery
   - Timeout
   ↓
4. GraphQL handler (from gqlgen) receives request
   ↓
5. Parse query into AST (Abstract Syntax Tree)
   ↓
6. Validate query against schema
   - Type checking
   - Field existence
   - Argument validation
   ↓
7. Create execution plan
   - Determine which resolvers to call
   - Determine execution order
   ↓
8. Execute resolvers in order
   - Root query resolver (e.g., patient)
   - Field resolvers (e.g., addresses, prescriptions)
   - Each resolver calls domain service
   ↓
9. Collect all resolved data
   ↓
10. Format as JSON response
   ↓
11. Send response to client
```

---

## Data Flow Diagram

```mermaid
flowchart LR
    subgraph Client
        Query["{<br/>  patient(id: \"123\") {<br/>    name<br/>    prescriptions {<br/>      drug<br/>    }<br/>  }<br/>}"]
    end
    
    subgraph GraphQL Layer
        Parse[Parse &<br/>Validate]
        QueryRes[Patient<br/>Query Resolver]
        FieldRes[Prescriptions<br/>Field Resolver]
    end
    
    subgraph Service Layer
        PatientSvc[Patient<br/>Service]
        PrescriptionSvc[Prescription<br/>Service]
    end
    
    subgraph Data Layer
        PatientRepo[(Patient<br/>Repository)]
        PrescriptionRepo[(Prescription<br/>Repository)]
    end
    
    Query -->|HTTP POST| Parse
    Parse --> QueryRes
    QueryRes -->|GetByID| PatientSvc
    PatientSvc -->|Query| PatientRepo
    PatientRepo -.->|Patient| PatientSvc
    PatientSvc -.->|Patient Model| QueryRes
    
    QueryRes -.->|Patient Object| FieldRes
    FieldRes -->|List| PrescriptionSvc
    PrescriptionSvc -->|Query| PrescriptionRepo
    PrescriptionRepo -.->|Prescriptions| PrescriptionSvc
    PrescriptionSvc -.->|Prescriptions| FieldRes
    
    FieldRes -.->|Complete Response| Query
    
    style Parse fill:#e1f5ff
    style QueryRes fill:#e1f5ff
    style FieldRes fill:#e1f5ff
```

---

## Summary

### Key Takeaways

1. **GraphQL sits as a delivery layer** - Same level as REST API and UI
2. **Resolvers are thin** - They just call services, no business logic
3. **Services are shared** - REST, GraphQL, and UI all use the same services
4. **Type-safe** - gqlgen generates code from schema
5. **Flexible queries** - Clients request exactly what they need
6. **Nested resolution** - GraphQL automatically resolves nested fields
7. **Single endpoint** - All queries go through `/graphql`

### Benefits of This Architecture

- ✅ No code duplication (services shared with REST)
- ✅ Type safety (compile-time checks)
- ✅ Easy to add new fields (update schema → regenerate)
- ✅ Client flexibility (request only needed fields)
- ✅ Single HTTP request for complex data
- ✅ Same middleware as REST (logging, errors, etc.)
- ✅ Simple to understand and maintain

