# PharmacyModernization (Patient & Prescription Management - Scaffold)

## Dev
- **Install**: Go 1.24, `templ`, Node.js 18+, `npm`
- **Setup once**: `make setup` (downloads Tailwind binary and installs npm dependencies)
- **Install gqlgen**: `make graphql-install` (one-time, for GraphQL code generation)
- **Run server**: `make dev` (runs all watchers)
- **Open**: http://localhost:8080

### GraphQL Development
- **Generate code**: `make graphql-generate` (run after schema changes)
- GraphQL code is already generated and committed, no need to regenerate unless you modify schemas

## Features
- **REST API**: `/api/v1/*` - RESTful endpoints for all domains
- **GraphQL API**: `/graphql` - Flexible query interface for all data
  - GraphQL Playground: `/playground` (development UI)
  - See `docs/GRAPHQL_IMPLEMENTATION.md` for usage guide
- **UI**: Server-rendered pages with Templ, HTMX, and Tailwind CSS
- **TypeScript**: Component-based client-side interactions

## Notes
- Feature-based modules under `domain/*` with API, GraphQL, service, repository, and UI layers.
- Viper YAML config in `internal/configs/` with env overrides (RX_*).
- Zap logging with request/correlation IDs.
- Tailwind source lives in `web/styles/input.css`; `make tailwind-watch` rebuilds `web/public/app.css` via the standalone Tailwind CLI with DaisyUI.
- GraphQL code generation with `make graphql-generate` (run after schema changes).

## Documentation
- **Architecture**: `ARCHITECTURE.md` - Overall system design and patterns
- **GraphQL**:
  - `docs/GRAPHQL_DOMAIN_STRUCTURE.md` - **Domain-based organization** (start here!)
  - `docs/GRAPHQL_RESOLVER_PATTERNS.md` - **Multi-resolver patterns** (3 phases with examples!)
  - `docs/GRAPHQL_ARCHITECTURE_DIAGRAMS.md` - Visual diagrams of how GraphQL works
  - `docs/GRAPHQL_IMPLEMENTATION.md` - How to use GraphQL API
  - `docs/GRAPHQL_MIGRATION_PHASES.md` - Evolution path to federation
- **MongoDB**: `docs/MONGODB_IMPLEMENTATION.md` - Database integration guide
- **TypeScript Components**: `docs/ADDING_TYPESCRIPT_COMPONENTS.md` - Client-side development
