# Architecture Documentation

System architecture and design patterns for the RxIntake application.

## 📚 Documentation Files

- **[Architecture Overview](./ARCHITECTURE.md)** - Complete system architecture and design principles

## 🏗️ System Architecture

### Domain-Driven Design (DDD)
The application is organized around business domains:
- **Patient** - Patient management and records
- **Prescription** - Prescription management and dispensing
- **Dashboard** - Analytics and overview

### Layered Architecture

```
┌─────────────────────────────────────┐
│         Presentation Layer          │
│   (UI, API Controllers, GraphQL)    │
├─────────────────────────────────────┤
│         Application Layer           │
│    (Use Cases, Orchestration)       │
├─────────────────────────────────────┤
│          Domain Layer               │
│   (Business Logic, Services)        │
├─────────────────────────────────────┤
│       Infrastructure Layer          │
│  (Database, Auth, External APIs)    │
└─────────────────────────────────────┘
```

## 🔑 Key Principles

### 1. Domain Independence
Each domain is self-contained with its own:
- Models
- Services
- Repositories
- UI Components
- API Controllers
- GraphQL Schema & Resolvers

### 2. Clean Separation of Concerns
- **UI** - Templ templates and HTMX for interactivity
- **API** - RESTful endpoints for programmatic access
- **GraphQL** - Flexible query interface
- **Services** - Business logic
- **Repositories** - Data access

### 3. Platform Services
Shared services in `internal/platform/`:
- **Auth** - Authentication & authorization
- **Database** - MongoDB connection & patterns
- **Logging** - Structured logging with Zap
- **Config** - Configuration management
- **HTTP** - HTTP utilities and middleware

## 🗂️ Directory Structure

```
rxintake_scaffold/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── app/             # Application wiring
│   ├── configs/         # Configuration files
│   ├── platform/        # Shared platform services
│   │   ├── auth/        # Authentication & authorization
│   │   ├── database/    # Database connection
│   │   ├── logging/     # Logging utilities
│   │   └── paths/       # Route registry
│   └── graphql/         # GraphQL server setup
├── domain/              # Business domains
│   ├── patient/
│   │   ├── api/         # REST API
│   │   ├── graphql/     # GraphQL schema & resolvers
│   │   ├── models/      # Data models
│   │   ├── repository/  # Data access
│   │   ├── service/     # Business logic
│   │   ├── security/    # Domain permissions
│   │   └── ui/          # UI components & routes
│   ├── prescription/
│   └── dashboard/
└── web/                 # Shared web components
    ├── components/
    │   ├── layouts/     # Page layouts
    │   ├── elements/    # Reusable UI elements
    │   ├── auth/        # Auth-related components
    │   └── user/        # User display components
    └── assets/          # Static assets
```

## 🛠️ Technology Stack

### Backend
- **Go** - Primary language
- **Chi** - HTTP router
- **MongoDB** - Database
- **gqlgen** - GraphQL code generation

### Frontend
- **Templ** - Type-safe Go templates
- **HTMX** - Dynamic UI without JavaScript
- **DaisyUI** - Tailwind CSS component library
- **TypeScript** - For complex client interactions

### Authentication
- **JWT** - Stateless authentication
- **Cookie & Header-based** - Flexible token sources
- **Permission-based AuthZ** - Fine-grained access control

## 🔄 Request Flow

### Web UI Request
```
Browser → Chi Router → Auth Middleware → Permission Check → Handler → Templ Template → HTML Response
```

### API Request
```
Client → Chi Router → Auth Middleware (Header) → Permission Check → Controller → Service → Repository → JSON Response
```

### GraphQL Request
```
Client → GraphQL Server → Auth Middleware → Directive Check → Resolver → Service → Repository → JSON Response
```

## 📖 Related Documentation

- [Security Architecture](../security/SECURITY_ARCHITECTURE.md)
- [GraphQL Implementation](../graphql/GRAPHQL_IMPLEMENTATION.md)
- [MongoDB Implementation](../mongodb/MONGODB_IMPLEMENTATION.md)

## 🎯 Design Goals

1. **Maintainability** - Clear structure, easy to understand
2. **Scalability** - Domains can grow independently
3. **Testability** - Easy to unit test and mock
4. **Security** - Built-in auth/authz at every level
5. **Developer Experience** - Type safety, hot reload, good docs

