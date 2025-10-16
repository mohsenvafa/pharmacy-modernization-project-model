# Environment Configuration

This project uses environment variables for configuration, stored in a `.env` file at the project root.

## Setup

1. **Copy the example file:**
   ```bash
   cp .env.example .env
   ```

2. **Update credentials** (optional for local development):
   The default credentials are already set for local development. Update them for production deployments.

## Environment Variables

### MongoDB Credentials (for Podman & Seed Script)
These variables are used by Podman containers and the seed script:

- `MONGO_ROOT_USERNAME` - MongoDB admin username
- `MONGO_ROOT_PASSWORD` - MongoDB admin password
- `MONGO_DATABASE` - MongoDB database name (default: `pharmacy_modernization`)

### Application Configuration (RX_ prefix - MUST BE UPPERCASE)
These are **REQUIRED** and read by the Go application using Viper:

- `RX_DATABASE_MONGODB_URI` - Full MongoDB connection URI for the main database
- `RX_CACHE_MONGODB_URI` - Full MongoDB connection URI for caching
- `RX_AUTH_JWT_SECRET` - JWT signing secret (minimum 32 characters, required even in dev mode)

**Critical:** 
- Environment variable names **MUST be UPPERCASE** (e.g., `RX_DATABASE_MONGODB_URI`, not `rx_database_mongodb_uri`)
- The `app.yaml` file now has empty URI fields
- The application will **fail to start** if these environment variables are not set

**Why Uppercase?** Viper automatically converts YAML keys to uppercase when looking for environment variables. See "How Viper Maps Variables" section below.

## How Viper Maps Environment Variables

Viper automatically maps environment variables to configuration keys using this process:

```
YAML key:             database.mongodb.uri
                           ↓
Replace dots with _:  database_mongodb_uri
                           ↓
Convert to UPPERCASE: DATABASE_MONGODB_URI
                           ↓
Add RX_ prefix:       RX_DATABASE_MONGODB_URI
```

**Example mappings:**
- `app.port` → `RX_APP_PORT`
- `auth.jwt.secret` → `RX_AUTH_JWT_SECRET`
- `database.mongodb.uri` → `RX_DATABASE_MONGODB_URI`
- `cache.mongodb.uri` → `RX_CACHE_MONGODB_URI`

**Important:** Environment variables are case-sensitive on Unix/Linux/macOS. They **MUST be UPPERCASE**.

```bash
# ✅ CORRECT - Will work
export RX_DATABASE_MONGODB_URI="mongodb://..."

# ❌ WRONG - Will NOT work
export rx_database_mongodb_uri="mongodb://..."
export Rx_Database_MongoDB_Uri="mongodb://..."
```

## How It Works

### 1. Podman Compose
The `podman/compose.yml` file reads the `.env` file and uses environment variables to set MongoDB credentials:

```yaml
env_file:
  - ../.env
environment:
  MONGO_INITDB_ROOT_USERNAME: ${MONGO_ROOT_USERNAME}
  MONGO_INITDB_ROOT_PASSWORD: ${MONGO_ROOT_PASSWORD}
```

### 2. Makefile
The `podman/Makefile` loads environment variables to display connection information and run MongoDB shell commands.

### 3. Go Application
The application uses Viper to read configuration from:
1. YAML config files in `internal/configs/`
2. Environment variables with the `RX_` prefix (UPPERCASE required)

Viper configuration in `config.go`:
```go
v.SetEnvPrefix("RX")                              // Add RX_ prefix
v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Replace . with _
v.AutomaticEnv()                                   // Enable automatic env var lookup
```

This means any YAML key like `database.mongodb.uri` automatically maps to `RX_DATABASE_MONGODB_URI`.

### 4. Seed Script
The `cmd/seed/main.go` reads environment variables directly:
- `MONGO_ROOT_USERNAME`
- `MONGO_ROOT_PASSWORD`
- `MONGO_DATABASE`

## Security

⚠️ **Important:**
- The `.env` file is in `.gitignore` and should **never** be committed
- Use `.env.example` as a template (this IS committed)
- For production, use secure secrets management (e.g., AWS Secrets Manager, HashiCorp Vault)
- Change default passwords before deploying to any environment

## Example `.env` File

### Local Development
```env
# MongoDB Configuration (for Podman & seed script)
MONGO_ROOT_USERNAME=
MONGO_ROOT_PASSWORD=
MONGO_DATABASE=pharmacy_modernization

# Application MongoDB URIs (REQUIRED)
RX_DATABASE_MONGODB_URI=mongodb://:@localhost:27017
RX_CACHE_MONGODB_URI=mongodb://:@localhost:27017

# JWT Secret (REQUIRED - even in dev mode)
# Generate with: openssl rand -base64 32
RX_AUTH_JWT_SECRET=dev-secret-key-change-in-production-minimum-32-chars
```

### Production
```env
# MongoDB Configuration
MONGO_ROOT_USERNAME=
MONGO_ROOT_PASSWORD=
MONGO_DATABASE=pharmacy_modernization

# Application MongoDB URIs (REQUIRED)
RX_DATABASE_MONGODB_URI=mongodb://:@prod-mongo.example.com:27017
RX_CACHE_MONGODB_URI=mongodb://r:@cache-mongo.example.com:27017

# Security (REQUIRED)
# Generate with: openssl rand -base64 32
RX_AUTH_JWT_SECRET=<your-secret-key-minimum-32-chars>

# Production Environment
RX_APP_ENV=prod
RX_LOGGING_LEVEL=info
RX_LOGGING_FORMAT=json
```

## Overriding Configuration

You can override any configuration value using environment variables with the `RX_` prefix (UPPERCASE):

```bash
# Override app port
export RX_APP_PORT=9090

# Override JWT secret
export RX_AUTH_JWT_SECRET=my-super-secret-key

# Override MongoDB URI
export RX_DATABASE_MONGODB_URI=mongodb://user:pass@prod-mongo:27017
```

**Remember:** All `RX_*` environment variables must be UPPERCASE. The naming convention follows the config structure with dots replaced by underscores:
- `app.port` → `RX_APP_PORT`
- `auth.jwt.secret` → `RX_AUTH_JWT_SECRET`
- `database.mongodb.uri` → `RX_DATABASE_MONGODB_URI`
- `logging.level` → `RX_LOGGING_LEVEL`

## Configuration Priority

Viper loads configuration in this order (later sources override earlier ones):

1. Default values in `app.yaml`
2. Environment-specific file (`app.dev.yaml` or `app.prod.yaml`)
3. Environment variables with `RX_` prefix (highest priority)

Example: If `app.yaml` has `uri: ""` and you set `RX_DATABASE_MONGODB_URI=mongodb://...`, the environment variable wins.

## Troubleshooting

**Issue:** Application fails to start with "database URI is empty"
- **Solution:** Ensure `.env` file exists and `RX_DATABASE_MONGODB_URI` is set
- **Solution:** Verify variable name is UPPERCASE: `RX_DATABASE_MONGODB_URI` not `rx_database_mongodb_uri`
- **Solution:** Run `source .env` or restart your application to load environment variables

**Issue:** MongoDB container fails to start
- **Solution:** Check that `.env` file exists and contains valid `MONGO_ROOT_USERNAME` and `MONGO_ROOT_PASSWORD`

**Issue:** Application can't connect to MongoDB
- **Solution:** Ensure `RX_DATABASE_MONGODB_URI` matches the credentials in `.env`
- **Solution:** Verify MongoDB is running: `make podman-up`

**Issue:** Seed script fails
- **Solution:** Ensure MongoDB is running (`make podman-up`) 
- **Solution:** Check that `MONGO_ROOT_USERNAME`, `MONGO_ROOT_PASSWORD`, and `MONGO_DATABASE` are set in `.env`

**Issue:** Environment variables not loading
- **Solution:** Verify variable names are UPPERCASE (e.g., `RX_DATABASE_MONGODB_URI`, not `rx_database_mongodb_uri`)
- **Solution:** Make sure you're running commands from the project root
- **Solution:** The Makefile loads `.env` automatically, but shell commands need `source .env`
- **Solution:** Check for typos in variable names - they must match the Viper mapping exactly

