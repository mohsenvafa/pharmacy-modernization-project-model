# Authentication Builder Pattern

## Overview

The `auth.Builder` provides a clean, fluent API for initializing the authentication system. It keeps your application wiring code simple and focused.

---

## Basic Usage

### In Application Initialization

```go
// internal/app/wire.go
package app

import (
    "pharmacy-modernization-project-model/internal/platform/auth"
)

func (a *App) wire() error {
    logger := logging.NewLogger(a.Cfg)
    
    // Initialize authentication with builder
    if err := auth.NewBuilder().
        WithJWTConfig(
            a.Cfg.Auth.JWT.Secret,
            a.Cfg.Auth.JWT.Issuer,
            a.Cfg.Auth.JWT.Audience,
            a.Cfg.Auth.JWT.Cookie.Name,
        ).
        WithDevMode(a.Cfg.Auth.DevMode).
        WithEnvironment(a.Cfg.App.Env).
        WithLogger(logger.Base).
        Build(); err != nil {
        return err
    }
    
    // Continue with rest of initialization...
    return nil
}
```

---

## Builder Methods

### `NewBuilder()`
Creates a new authentication builder instance.

```go
builder := auth.NewBuilder()
```

### `WithJWTConfig(secret, issuer, audience, cookieName string)`
Sets JWT configuration parameters.

```go
builder.WithJWTConfig(
    "your-secret-key",
    "rxintake",
    "rxintake",
    "auth_token",
)
```

**Parameters:**
- `secret` - JWT signing secret (use environment variable in production)
- `issuer` - Expected token issuer
- `audience` - Expected token audience
- `cookieName` - Name of the authentication cookie

### `WithDevMode(enabled bool)`
Enables or disables development mode with mock users.

```go
builder.WithDevMode(true)  // Enable for local development
builder.WithDevMode(false) // Disable for production
```

### `WithEnvironment(env string)`
Sets the application environment for safety checks.

```go
builder.WithEnvironment("dev")  // Development
builder.WithEnvironment("prod") // Production
```

**Safety:** If `env` is "prod" and dev mode is enabled, `Build()` will return an error.

### `WithLogger(logger *zap.Logger)`
Sets the logger for authentication messages.

```go
builder.WithLogger(logger.Base)
```

Logs include:
- Auth initialization info
- Dev mode warnings
- Configuration details

### `Build()`
Builds and initializes the authentication system. Returns an error if configuration is invalid.

```go
if err := builder.Build(); err != nil {
    return fmt.Errorf("failed to initialize auth: %w", err)
}
```

**Returns:**
- `nil` - Success
- `error` - Configuration error or safety violation

### `MustBuild()`
Builds and initializes, panicking on error. Use only when you're certain configuration is valid.

```go
builder.MustBuild() // Panics if error
```

---

## Usage Examples

### Example 1: Minimal Setup

```go
// Minimal configuration without logger
err := auth.NewBuilder().
    WithJWTConfig("secret", "issuer", "audience", "auth_token").
    WithDevMode(false).
    WithEnvironment("dev").
    Build()

if err != nil {
    log.Fatal(err)
}
```

### Example 2: Development Setup

```go
// Development with dev mode enabled
err := auth.NewBuilder().
    WithJWTConfig(
        "dev-secret-key",
        "rxintake-dev",
        "rxintake",
        "auth_token",
    ).
    WithDevMode(true).  // Enable mock users
    WithEnvironment("dev").
    WithLogger(logger).
    Build()

if err != nil {
    return err
}
```

**Output logs:**
```
⚠️  AUTH DEV MODE ACTIVE - Do not use in production!
Authentication initialized | mode=development dev_mode=true issuer=rxintake-dev
```

### Example 3: Production Setup

```go
// Production with strict checks
err := auth.NewBuilder().
    WithJWTConfig(
        os.Getenv("JWT_SECRET"), // From environment
        "rxintake",
        "rxintake",
        "auth_token",
    ).
    WithDevMode(false). // Dev mode disabled
    WithEnvironment("prod").
    WithLogger(logger).
    Build()

if err != nil {
    logger.Fatal("Failed to initialize auth", zap.Error(err))
}
```

**Output logs:**
```
Authentication initialized | mode=production dev_mode=false issuer=rxintake
```

### Example 4: From Configuration

```go
// Build from config file
err := auth.NewBuilder().
    WithJWTConfig(
        cfg.Auth.JWT.Secret,
        cfg.Auth.JWT.Issuer,
        cfg.Auth.JWT.Audience,
        cfg.Auth.JWT.Cookie.Name,
    ).
    WithDevMode(cfg.Auth.DevMode).
    WithEnvironment(cfg.App.Env).
    WithLogger(logger.Base).
    Build()
```

### Example 5: Testing Setup

```go
// Test setup with dev mode
func setupTestAuth(t *testing.T) {
    err := auth.NewBuilder().
        WithJWTConfig("test-secret", "test", "test", "test_token").
        WithDevMode(true).
        WithEnvironment("test").
        Build()
    
    if err != nil {
        t.Fatal(err)
    }
}
```

---

## Error Handling

### Production Safety Error

```go
err := auth.NewBuilder().
    WithJWTConfig("secret", "issuer", "audience", "token").
    WithDevMode(true).        // Dev mode enabled
    WithEnvironment("prod").  // Production environment
    Build()

// Returns error:
// "FATAL: Dev mode cannot be enabled in production environment (env=prod, dev_mode=true)"
```

### Handling Errors

```go
if err := builder.Build(); err != nil {
    // Log and handle error
    logger.Error("Auth initialization failed", zap.Error(err))
    return err
}
```

---

## Benefits of Builder Pattern

### ✅ Clean Wire Code

**Before:**
```go
auth.InitJWTConfig(auth.JWTConfig{
    Secret:     cfg.Auth.JWT.Secret,
    Issuer:     cfg.Auth.JWT.Issuer,
    Audience:   cfg.Auth.JWT.Audience,
    CookieName: cfg.Auth.JWT.Cookie.Name,
})
auth.InitDevMode(cfg.Auth.DevMode)
if cfg.App.Env == "prod" && cfg.Auth.DevMode {
    logger.Fatal("FATAL: Dev mode in production")
}
if cfg.Auth.DevMode {
    logger.Warn("Dev mode active")
}
```

**After:**
```go
if err := auth.NewBuilder().
    WithJWTConfig(
        cfg.Auth.JWT.Secret,
        cfg.Auth.JWT.Issuer,
        cfg.Auth.JWT.Audience,
        cfg.Auth.JWT.Cookie.Name,
    ).
    WithDevMode(cfg.Auth.DevMode).
    WithEnvironment(cfg.App.Env).
    WithLogger(logger.Base).
    Build(); err != nil {
    return err
}
```

### ✅ Fluent API

Chain methods for readable configuration:

```go
auth.NewBuilder().
    WithJWTConfig(...).
    WithDevMode(true).
    WithEnvironment("dev").
    WithLogger(logger).
    Build()
```

### ✅ Centralized Validation

All safety checks and validation in one place.

### ✅ Error Handling

Returns errors instead of panicking (except `MustBuild()`).

### ✅ Testable

Easy to create different configurations for testing.

### ✅ Extensible

Easy to add new configuration options without changing existing code.

---

## Advanced Patterns

### Pattern 1: Configuration Validation

```go
func buildAuth(cfg *config.Config, logger *zap.Logger) error {
    builder := auth.NewBuilder().
        WithEnvironment(cfg.App.Env).
        WithLogger(logger)
    
    // Validate config
    if cfg.Auth.JWT.Secret == "" {
        return errors.New("JWT secret is required")
    }
    
    builder.WithJWTConfig(
        cfg.Auth.JWT.Secret,
        cfg.Auth.JWT.Issuer,
        cfg.Auth.JWT.Audience,
        cfg.Auth.JWT.Cookie.Name,
    ).WithDevMode(cfg.Auth.DevMode)
    
    return builder.Build()
}
```

### Pattern 2: Environment-Specific Setup

```go
func buildAuth(cfg *config.Config, logger *zap.Logger) error {
    builder := auth.NewBuilder().
        WithLogger(logger).
        WithEnvironment(cfg.App.Env)
    
    switch cfg.App.Env {
    case "dev", "development":
        builder.
            WithJWTConfig("dev-secret", "dev-issuer", "rxintake", "auth_token").
            WithDevMode(true)
    
    case "staging":
        builder.
            WithJWTConfig(os.Getenv("JWT_SECRET"), "staging-issuer", "rxintake", "auth_token").
            WithDevMode(false)
    
    case "prod", "production":
        builder.
            WithJWTConfig(os.Getenv("JWT_SECRET"), "rxintake", "rxintake", "auth_token").
            WithDevMode(false)
    
    default:
        return fmt.Errorf("unknown environment: %s", cfg.App.Env)
    }
    
    return builder.Build()
}
```

### Pattern 3: Conditional Dev Mode

```go
// Only enable dev mode if explicitly set AND not in production
devModeEnabled := cfg.Auth.DevMode && cfg.App.Env != "prod"

err := auth.NewBuilder().
    WithJWTConfig(...).
    WithDevMode(devModeEnabled).
    WithEnvironment(cfg.App.Env).
    WithLogger(logger).
    Build()
```

---

## Testing with Builder

### Test Helper

```go
// test_helpers.go
package testutil

import (
    "testing"
    "pharmacy-modernization-project-model/internal/platform/auth"
    "go.uber.org/zap"
)

func SetupTestAuth(t *testing.T) {
    logger := zap.NewNop()
    
    err := auth.NewBuilder().
        WithJWTConfig("test-secret", "test", "test", "test_token").
        WithDevMode(true).
        WithEnvironment("test").
        WithLogger(logger).
        Build()
    
    if err != nil {
        t.Fatal(err)
    }
}
```

### In Tests

```go
func TestPatientRoutes(t *testing.T) {
    // Setup
    testutil.SetupTestAuth(t)
    
    // Test with mock user
    req := httptest.NewRequest("GET", "/patients", nil)
    req.Header.Set("X-Mock-User", "doctor")
    
    // ... rest of test
}
```

---

## Migration Guide

If you have existing auth initialization code:

### Old Code

```go
auth.InitJWTConfig(auth.JWTConfig{
    Secret:     secret,
    Issuer:     issuer,
    Audience:   audience,
    CookieName: cookieName,
})
auth.InitDevMode(devMode)
```

### New Code

```go
err := auth.NewBuilder().
    WithJWTConfig(secret, issuer, audience, cookieName).
    WithDevMode(devMode).
    WithEnvironment(env).
    WithLogger(logger).
    Build()

if err != nil {
    return err
}
```

---

## Summary

The auth builder provides:

- ✅ **Fluent API** - Readable, chainable methods
- ✅ **Safety checks** - Prevents dev mode in production
- ✅ **Clean code** - Keeps wire.go focused
- ✅ **Error handling** - Returns errors instead of panicking
- ✅ **Logging** - Integrated logging support
- ✅ **Testable** - Easy to configure for tests
- ✅ **Extensible** - Easy to add new options

Use it for clean, safe, and maintainable authentication initialization!

