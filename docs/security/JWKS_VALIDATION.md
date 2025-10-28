# JWKS Signature Validation

This document explains how to use JWKS (JSON Web Key Set) for JWT signature validation using `github.com/MicahParks/keyfunc`.

## Overview

The authentication system uses JWKS for JWT signature validation, which provides:

- **Public Key Cryptography** - Only the auth provider needs the private key
- **Automatic Key Rotation** - Keys can be rotated without service restarts
- **Better Security** - No shared secrets to compromise
- **Industry Standard** - Used by OAuth 2.0, OIDC, and major auth providers

## Configuration

### Environment Variables

```bash
# JWKS Configuration (required)
RX_AUTH_JWT_JWKS_URL="https://your-auth-provider.com/.well-known/jwks.json"
RX_AUTH_JWT_JWKS_CACHE=15  # Cache duration in minutes (default: 15)
RX_AUTH_JWT_SIGNING_METHODS="RS256,ES256"  # Comma-separated list

# JWT Claims Validation
RX_AUTH_JWT_ISSUER="your-issuer"
RX_AUTH_JWT_AUDIENCE="your-audience"
RX_AUTH_JWT_CLIENT_IDS="client1,client2"
```

### Configuration Files

```yaml
# internal/configs/app.yaml
auth:
  jwt:
    issuer: ["PharmacyModernization"]
    audience: ["PharmacyModernization"]
    client_ids: ["default-client", "mobile-app", "web-app"]
    cookie:
      name: "auth_token"
      secure: false  # Set to true in production (HTTPS only)
      httponly: true
      max_age: 3600  # 1 hour in seconds
    # JWKS Configuration
    jwks_url: ""  # REQUIRED: Set via RX_AUTH_JWT_JWKS_URL
    jwks_cache: 15  # Cache duration in minutes (default: 15)
    signing_methods: ["RS256", "ES256"]  # Allowed signing methods
```

## Usage

### Programmatic Configuration

```go
import "pharmacy-modernization-project-model/internal/platform/auth"

// Initialize with JWKS
err := auth.NewBuilder().
    WithJWTConfig(
        []string{"your-issuer"},
        []string{"your-audience"},
        []string{"your-client"},
        "auth_token",
    ).
    WithJWKSConfig(
        "https://your-auth-provider.com/.well-known/jwks.json",
        15,  // Cache for 15 minutes
        []string{"RS256", "ES256"},  // Allowed signing methods
    ).
    WithDevMode(false).
    WithEnvironment("prod").
    WithLogger(logger).
    Build()

if err != nil {
    log.Fatal("Failed to initialize auth:", err)
}
```

## Supported Signing Methods

- **RS256** - RSA with SHA-256
- **RS384** - RSA with SHA-384  
- **RS512** - RSA with SHA-512
- **ES256** - ECDSA with SHA-256
- **ES384** - ECDSA with SHA-384
- **ES512** - ECDSA with SHA-512

## JWT Token Requirements

Tokens must include:

1. **Key ID (`kid`)** in the header
2. **Valid signature** using one of the allowed signing methods
3. **Standard claims** (iss, aud, exp, etc.)

Example token header:
```json
{
  "alg": "RS256",
  "typ": "JWT",
  "kid": "key-id-123"
}
```

## Error Handling

The system handles JWKS errors gracefully:

- **JWKS fetch failures** → Returns validation error with details
- **Key not found** → Returns validation error
- **Invalid signing method** → Returns validation error
- **Network issues** → Uses cached keys, logs errors

## Security Considerations

1. **HTTPS Only** - Always use HTTPS for JWKS URLs in production
2. **Key Rotation** - The system automatically refreshes keys from JWKS
3. **Cache Duration** - Set appropriate cache duration (15 minutes default)
4. **Signing Methods** - Restrict to required methods only
5. **Token Validation** - Always validate issuer, audience, and expiration

## Testing

```go
func TestJWKSValidation(t *testing.T) {
    // Test with mock JWKS endpoint
    config := JWTConfig{
        JWKSURL: "https://mock-jwks.example.com/.well-known/jwks.json",
        SigningMethods: []string{"RS256"},
    }
    
    InitJWTConfig(config)
    
    // Test token validation...
}
```

## Troubleshooting

### Common Issues

**Issue: "JWKS URL not configured"**
- Set `RX_AUTH_JWT_JWKS_URL` environment variable
- Verify JWKS URL is accessible

**Issue: "key not found in JWKS"**
- Check if token `kid` matches keys in JWKS
- Verify JWKS URL is accessible
- Check if key rotation occurred

**Issue: "signing method not allowed"**
- Add the signing method to `signing_methods` config
- Verify token algorithm matches configuration

**Issue: JWKS fetch errors**
- Check network connectivity to JWKS URL
- Verify JWKS URL returns valid JSON
- Check if JWKS endpoint requires authentication

### Debug Logging

Enable debug logging to see JWKS operations:

```yaml
logging:
  level: debug
```

This will show:
- JWKS fetch attempts
- Key refresh operations  
- Validation errors
- Error details
