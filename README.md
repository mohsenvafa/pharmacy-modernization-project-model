# PharmacyModernization (Patient & Prescription Management - Scaffold)

## Dev

### Prerequisites
- **Install**: Go 1.24, `templ`, Node.js 18+, `npm`
- **Windows users**: PowerShell 5.1+ (included with Windows 10/11) or PowerShell Core
  - If you encounter execution policy issues, run: `Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser`
- **macOS/Linux users**: `make` utility

### Setup (One-time)

**macOS/Linux:**
1. `cp .env.example .env` (copy environment variables and fill in values)
2. `make setup` (downloads Tailwind binary and installs npm dependencies)

**Windows (PowerShell):**
1. Copy `.env.example` to `.env` (copy environment variables and fill in values)
2. `.\make.ps1 setup` (downloads Tailwind binary and installs npm dependencies)

### Running the Application

**macOS/Linux:**
- **Install gqlgen**: `make graphql-install` (one-time, for GraphQL code generation)
- **Run IRIS Mock Server**: `make mock-iris` (runs mock APIs on port 8881 - **start this first!**)
- **Run server**: `make dev` (runs all watchers, uses `app.yaml` defaults)

**Windows:**
- **Install gqlgen**: `.\make.ps1 graphql-install` (one-time, for GraphQL code generation)
- **Run IRIS Mock Server**: `.\make.ps1 mock-iris` (runs mock APIs on port 8881 - **start this first!**)
- **Run server**: `.\make.ps1 dev` (runs all watchers, uses `app.yaml` defaults)

- **Open**: http://localhost:8080

**Note:** The app uses `app.yaml` for development by default. For production, set `RX_APP_ENV=prod` to load `app.prod.yaml` with secure defaults.

## IRIS Mock Server (Development APIs)
The IRIS Mock Server simulates external pharmacy and billing APIs for local development:
- **Start mock server**: 
  - macOS/Linux: `make mock-iris`
  - Windows: `.\make.ps1 mock-iris`
  - (runs on port 8881)
- **Verify server**: Visit http://localhost:8881/ for welcome page
- **APIs provided**:
  - 📍 Pharmacy API: http://localhost:8881/pharmacy/v1
  - 📍 Billing API: http://localhost:8881/billing/v1
  - 📍 Stargate Auth: http://localhost:8881/oauth

**Note**: Start the mock server before running the main application for full functionality.

## MongoDB Setup
- **Start MongoDB**: 
  - macOS/Linux: `make podman-up`
  - Windows: `.\make.ps1 podman-up` or `.\podman\make.ps1 podman-up`
  - (starts MongoDB + Memcached containers)
- **Stop MongoDB**: 
  - macOS/Linux: `make podman-down`
  - Windows: `.\make.ps1 podman-down` or `.\podman\make.ps1 podman-down`
- **View logs**: 
  - macOS/Linux: `make podman-logs`
  - Windows: `.\make.ps1 podman-logs` or `.\podman\make.ps1 podman-logs`
- **First-time seeding**: 
  - macOS/Linux: `make podman-up && go run ./cmd/seed`
  - Windows: `.\make.ps1 podman-up; go run ./cmd/seed` or `.\podman\make.ps1 podman-up; go run ./cmd/seed`

For more MongoDB commands (restart, clean, shell, seed), see `podman/README.md` or run `.\podman\make.ps1 help` on Windows.

**Connection Details:** (configured in `.env` file)
- **Host**: localhost:27017
- **Username**: `admin` (from `MONGO_ROOT_USERNAME`)
- **Password**: (from `MONGO_ROOT_PASSWORD`)
- **Database**: (from `MONGO_DATABASE`)
- **Connection String**: Values from `.env` file are used to construct the MongoDB URI

### GraphQL Development
- **Generate code**: 
  - macOS/Linux: `make graphql-generate`
  - Windows: `.\make.ps1 graphql-generate`
  - (run after schema changes)
- GraphQL code is already generated and committed, no need to regenerate unless you modify schemas

## Features
- **REST API**: `/api/v1/*` - RESTful endpoints for all domains
- **GraphQL API**: `/graphql` - Flexible query interface for all data
  - GraphQL Playground: `/playground` (development UI)
  - See `docs/GRAPHQL_IMPLEMENTATION.md` for usage guide
- **UI**: Server-rendered pages with Templ, HTMX, and Tailwind CSS
- **TypeScript**: Component-based client-side interactions

## Command Reference

| Command | macOS/Linux | Windows |
|---------|-------------|---------|
| Setup | `make setup` | `.\make.ps1 setup` |
| Run dev server | `make dev` | `.\make.ps1 dev` |
| Run mock server | `make mock-iris` | `.\make.ps1 mock-iris` |
| Start MongoDB | `make podman-up` | `.\make.ps1 podman-up` or `.\podman\make.ps1 podman-up` |
| Stop MongoDB | `make podman-down` | `.\make.ps1 podman-down` or `.\podman\make.ps1 podman-down` |
| GraphQL generate | `make graphql-generate` | `.\make.ps1 graphql-generate` |
| Build TypeScript | `make build-ts` | `.\make.ps1 build-ts` |
| Watch TypeScript | `make watch-ts` | `.\make.ps1 watch-ts` |
| Tailwind watch | `make tailwind-watch` | `.\make.ps1 tailwind-watch` |
| Show help | `make help` | `.\make.ps1 help` |

## Notes
- Feature-based modules under `domain/*` with API, GraphQL, service, repository, and UI layers.
- Viper YAML config in `internal/configs/` with env overrides (RX_*).
- Zap logging with request/correlation IDs.
- Tailwind source lives in `web/styles/input.css`; `make tailwind-watch` (or `.\make.ps1 tailwind-watch` on Windows) rebuilds `web/public/app.css` via the standalone Tailwind CLI with DaisyUI.
- GraphQL code generation with `make graphql-generate` (or `.\make.ps1 graphql-generate` on Windows) (run after schema changes).

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
