# Configuration Guide

This guide explains how the application configuration works and how to configure it for different environments.

## Configuration Strategy

The application uses a **hybrid approach** combining:
1. **YAML files** for non-sensitive, environment-specific defaults
2. **Environment variables** for sensitive credentials and runtime overrides

This provides the best of both worlds:
- ✅ Secure credential management
- ✅ Clear environment differences
- ✅ Minimal environment variables
- ✅ Easy to override any setting

## Configuration Files

### `app.yaml` (Base Configuration)
- Default configuration optimized for **local development**
- Contains all available settings with dev-friendly defaults
- Settings:
  - `auth.dev_mode: true` - Mock authentication
  - `cookie.secure: false` - Works without HTTPS
  - `logging: debug, console` - Verbose logging
  - `use_mock: true` - Mock external services

### `app.prod.yaml` (Production Overrides)
- Loaded when `RX_APP_ENV=prod` is set
- Overrides base config with **production-safe defaults**
- Critical differences:
  - `auth.dev_mode: false` ⚠️ **Real JWT authentication required**
  - `cookie.secure: true` ⚠️ **HTTPS only**
  - `logging: info, json` - Production logging
  - `use_mock: false` - Real external services
  - Higher connection pool sizes

### `.env` (Sensitive Credentials)
- **Never committed to git** (in `.gitignore`)
- Contains all sensitive credentials
- Required variables:
  - `RX_DATABASE_MONGODB_URI`
  - `RX_CACHE_MONGODB_URI`
  - `RX_AUTH_JWT_SECRET` (production)

## How Configuration Loading Works

```
┌─────────────────────────────────────────────────────┐
│ 1. Load app.yaml (base defaults)                    │
└─────────────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────────────┐
│ 2. If RX_APP_ENV is set, merge app.{env}.yaml      │
│    Example: RX_APP_ENV=prod → loads app.prod.yaml  │
└─────────────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────────────┐
│ 3. Apply environment variables (RX_* prefix)        │
│    These override everything                         │
└─────────────────────────────────────────────────────┘
```

**Priority (highest to lowest):**
1. Environment variables (`RX_*`)
2. Environment-specific YAML (`app.prod.yaml`)
3. Base YAML (`app.yaml`)

## Environment Setup

### Local Development

```bash
# 1. Copy .env template
cp .env.example .env

# 2. Fill in credentials
vim .env  # Add MongoDB URIs

# 3. Run application (uses app.yaml defaults)
make dev
```

Your `.env` file:
```bash
MONGO_ROOT_USERNAME=
MONGO_ROOT_PASSWORD=
MONGO_DATABASE=pharmacy_modernization
RX_DATABASE_MONGODB_URI=mongodb://:@localhost:27017
RX_CACHE_MONGODB_URI=mongodb://:@localhost:27017
```

### Production Deployment

```bash
# Set environment to production
export RX_APP_ENV=prod

# Set required credentials
export RX_DATABASE_MONGODB_URI="mongodb://prod_user:secure_pass@prod-mongo:27017"
export RX_CACHE_MONGODB_URI="mongodb://cache_user:secure_pass@cache-mongo:27017"
export RX_AUTH_JWT_SECRET="your-strong-jwt-secret-minimum-32-chars"

# Run application
./server
```

Or using `.env` file:
```bash
# .env for production
RX_APP_ENV=prod
RX_DATABASE_MONGODB_URI=mongodb://prod_user:secure_pass@prod-mongo:27017
RX_CACHE_MONGODB_URI=mongodb://cache_user:secure_pass@cache-mongo:27017
RX_AUTH_JWT_SECRET=your-strong-jwt-secret-minimum-32-chars
```

### Staging/QA Environment

You can create additional environment files:

```bash
# Create app.staging.yaml
cp internal/configs/app.prod.yaml internal/configs/app.staging.yaml

# Edit for staging-specific settings
vim internal/configs/app.staging.yaml

# Use it
export RX_APP_ENV=staging
```

## Configuration Variables

### Required for All Environments

| Variable | Description | Example |
|----------|-------------|---------|
| `RX_DATABASE_MONGODB_URI` | Main database connection | `mongodb://user:pass@host:27017` |
| `RX_CACHE_MONGODB_URI` | Cache database connection | `mongodb://user:pass@host:27017` |

### Required for Production

| Variable | Description | Example |
|----------|-------------|---------|
| `RX_AUTH_JWT_SECRET` | JWT signing secret (min 32 chars) | `your-secret-key-here` |
| `RX_APP_ENV` | Environment name | `prod` |

### Optional Overrides

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `RX_APP_PORT` | Server port | `8080` | `9090` |
| `RX_LOGGING_LEVEL` | Log level | `debug` (dev), `info` (prod) | `warn` |
| `RX_LOGGING_FORMAT` | Log format | `console` (dev), `json` (prod) | `json` |
| `RX_LOGGING_OUTPUT` | Log output | `file` (dev), `both` (prod) | `console` |
| `RX_AUTH_DEV_MODE` | Enable dev auth | `true` (dev), `false` (prod) | `false` |
| `RX_DATABASE_MONGODB_DATABASE` | Database name | `pharmacy_modernization` | `custom_db` |

## Environment Variable Naming

Viper automatically maps YAML keys to environment variables:

```
YAML key:             database.mongodb.uri
                           ↓
Replace dots with _:  database_mongodb_uri
                           ↓
Convert to UPPERCASE: DATABASE_MONGODB_URI
                           ↓
Add RX_ prefix:       RX_DATABASE_MONGODB_URI
```

**Important:** Environment variables must be UPPERCASE!

✅ `RX_DATABASE_MONGODB_URI`  
❌ `rx_database_mongodb_uri`  
❌ `Rx_Database_MongoDB_Uri`

## Key Differences: Dev vs Prod

| Setting | Development | Production |
|---------|-------------|------------|
| Auth Mode | `dev_mode: true` (mock users) | `dev_mode: false` (real JWT) |
| Cookies | `secure: false` (HTTP OK) | `secure: true` (HTTPS only) |
| Logging Level | `debug` | `info` |
| Log Format | `console` (readable) | `json` (parseable) |
| Log Output | `file` | `both` (console + file) |
| Mock Services | `true` | `false` |
| MongoDB Pool | min: 5, max: 100 | min: 10, max: 100 |

## Security Best Practices

### ✅ DO:
- Store credentials in `.env` or secrets manager
- Use `RX_APP_ENV=prod` for production
- Set strong `RX_AUTH_JWT_SECRET` (min 32 chars)
- Enable HTTPS in production (required for secure cookies)
- Use separate MongoDB users for dev/prod
- Rotate secrets regularly

### ❌ DON'T:
- Commit `.env` to git
- Use dev mode in production
- Hardcode credentials in YAML files
- Use the same secrets across environments
- Use default passwords in production

## Kubernetes/Container Deployment

### Using Environment Variables
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: app
        env:
        - name: RX_APP_ENV
          value: "prod"
        - name: RX_DATABASE_MONGODB_URI
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: mongodb-uri
        - name: RX_AUTH_JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: jwt-secret
```

### Using ConfigMap (non-sensitive)
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  RX_APP_ENV: "prod"
  RX_LOGGING_LEVEL: "info"
  RX_APP_PORT: "8080"
```

## Troubleshooting

### App won't start - "database URI is empty"
**Solution:** Set `RX_DATABASE_MONGODB_URI` in `.env` or environment

### Authentication not working in production
**Solution:** Ensure `RX_APP_ENV=prod` is set (this disables dev_mode)

### Logs are too verbose in production
**Solution:** Check `RX_APP_ENV=prod` is set (changes to `info` level)

### Can't connect over HTTP in production
**Solution:** This is expected - production requires HTTPS for secure cookies

### Environment variables not loading
**Solution:** 
- Verify UPPERCASE naming
- Check `.env` file exists and is in the project root
- Use `source .env` before running the app manually
- Makefile loads `.env` automatically

## Examples

### Example 1: Local Development
```bash
# .env
RX_DATABASE_MONGODB_URI=mongodb://admin:admin123@localhost:27017
RX_CACHE_MONGODB_URI=mongodb://admin:admin123@localhost:27017

# Run
make dev  # Uses app.yaml defaults
```

### Example 2: Production with Custom Port
```bash
# .env
RX_APP_ENV=prod
RX_APP_PORT=9090
RX_DATABASE_MONGODB_URI=mongodb://prod:secure@prod-db:27017
RX_CACHE_MONGODB_URI=mongodb://prod:secure@cache-db:27017
RX_AUTH_JWT_SECRET=super-secret-key-minimum-32-characters

# Run
./server
```

### Example 3: Production with Debug Logging (Temporary)
```bash
# Override production logging for debugging
export RX_APP_ENV=prod
export RX_LOGGING_LEVEL=debug  # Overrides prod default
export RX_DATABASE_MONGODB_URI="mongodb://..."
# ... other vars

./server
```

## Migration from Old Configuration

If you're migrating from hardcoded configs:

1. **Identify credentials** in YAML files
2. **Move to `.env`** with `RX_` prefix
3. **Set URIs to empty** in YAML files
4. **Test locally** before deploying
5. **Update deployment scripts** to set environment variables

## Further Reading

- [ENVIRONMENT_SETUP.md](../ENVIRONMENT_SETUP.md) - Detailed environment variable guide
- [.env.example](../.env.example) - Template with all variables
- [Security Documentation](security/) - Security best practices

