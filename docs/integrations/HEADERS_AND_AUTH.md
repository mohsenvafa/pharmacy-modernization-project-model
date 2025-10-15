# Headers and Authentication - Best Practices

## Overview

This guide covers best practices for handling headers and authentication tokens in API integrations.

---

## 1️⃣ **Passing Specific Headers**

### **Approach 1: Global Headers (All Endpoints)**

Use `HeaderProvider` when headers should be included in **all requests** for a service:

```go
// Create a header provider with static headers
headerProvider := httpclient.NewStaticHeaderProvider(map[string]string{
    "X-API-Version": "v1",
    "X-Client-ID":   "rxintake-app",
    "Accept":        "application/json",
})

// Create HTTP client with global headers
httpClient := httpclient.NewClient(
    httpclient.Config{
        Timeout:        30 * time.Second,
        ServiceName:    "external_apis",
        HeaderProvider: headerProvider, // ✅ Headers added to ALL requests
    },
    logger,
)
```

**When to use:**
- ✅ Headers needed for **all endpoints** of a service
- ✅ API version headers
- ✅ Client identification headers
- ✅ Common accept/content-type headers

---

### **Approach 2: Per-Request Headers (Specific Endpoints)**

Pass headers directly when calling an endpoint:

```go
// integration_wire.go - Create client without header provider
sharedHTTPClient := httpclient.NewClient(
    httpclient.Config{
        Timeout:     30 * time.Second,
        ServiceName: "external_apis",
        // No HeaderProvider - headers passed per-request
    },
    logger,
)
```

```go
// http_client.go - Pass headers in specific requests
func (c *HTTPClient) GetInvoice(ctx, prescriptionID) (*InvoiceResponse, error) {
    url := replacePathParams(c.endpoints.GetInvoiceEndpoint(), map[string]string{
        "prescriptionID": prescriptionID,
    })
    
    // ✅ Custom headers for this specific endpoint
    resp, err := c.client.Get(ctx, url, map[string]string{
        "Content-Type":  "application/json",
        "Accept":        "application/json",
        "X-Request-Type": "invoice-lookup",  // Endpoint-specific header
    })
    
    // ... rest
}
```

**When to use:**
- ✅ Headers needed for **specific endpoints only**
- ✅ Endpoint-specific metadata
- ✅ Different headers for different operations

---

### **Approach 3: Hybrid (Global + Per-Request)**

Combine both approaches - global headers for all requests, specific headers for some:

```go
// Global headers for all requests
headerProvider := httpclient.NewStaticHeaderProvider(map[string]string{
    "X-API-Version": "v1",
    "Accept":        "application/json",
})

httpClient := httpclient.NewClient(
    httpclient.Config{
        HeaderProvider: headerProvider, // Global headers
    },
    logger,
)
```

```go
// Per-request headers (will merge with global headers)
func (c *HTTPClient) CreateInvoice(ctx, req) (*CreateInvoiceResponse, error) {
    // ... build URL ...
    
    resp, err := c.client.Post(ctx, url, body, map[string]string{
        "X-Idempotency-Key": generateKey(), // ✅ Endpoint-specific
    })
    // Result: Headers will include both global AND endpoint-specific
}
```

**When to use:**
- ✅ Common headers needed everywhere + some endpoints need extras
- ✅ Most flexible approach

---

## 2️⃣ **Authorization Token - Best Practices**

### **Scenario: Separate Token Service**

You have a token service that fetches/refreshes access tokens:

```go
// Example: Your token service interface
type TokenService interface {
    GetAccessToken(ctx context.Context) (string, error)
}

// Your implementation (OAuth, JWT, etc.)
type OAuthTokenService struct {
    clientID     string
    clientSecret string
    tokenURL     string
    // ...
}

func (s *OAuthTokenService) GetAccessToken(ctx context.Context) (string, error) {
    // Your logic to get token (OAuth flow, API call, etc.)
    // ...
    return "eyJhbGc...", nil
}
```

---

### **Best Practice: Use Token Provider with Caching**

Create a token provider that wraps your token service:

```go
// Step 1: Implement TokenProvider interface
type TokenServiceAdapter struct {
    tokenService TokenService
}

func (a *TokenServiceAdapter) GetToken(ctx context.Context) (string, error) {
    return a.tokenService.GetAccessToken(ctx)
}

// Step 2: Wrap with caching
cachedTokenProvider := httpclient.NewCachedTokenProvider(
    &TokenServiceAdapter{tokenService: yourTokenService},
    5*time.Minute, // Refresh 5 min before expiry
    logger,
)

// Step 3: Create auth header provider
authHeaderProvider := httpclient.NewAuthHeaderProvider(
    cachedTokenProvider,
    "Bearer", // Auth type
    logger,
)

// Step 4: Create HTTP client with auth headers
httpClient := httpclient.NewClient(
    httpclient.Config{
        Timeout:        30 * time.Second,
        ServiceName:    "external_apis",
        HeaderProvider: authHeaderProvider, // ✅ Automatic token refresh!
    },
    logger,
)
```

**How it works:**
1. First request: Fetches token from your service, caches it
2. Subsequent requests: Uses cached token (fast!)
3. Token expires soon: Automatically refreshes token
4. Token refresh fails: Error logged, request fails

---

### **Example Integration: Complete Flow**

```go
// 1. Your token service (existing service in your app)
type IrisTokenService struct {
    client   *http.Client
    tokenURL string
    apiKey   string
}

func (s *IrisTokenService) GetAccessToken(ctx context.Context) (string, error) {
    // Call your auth API to get token
    // POST to tokenURL with apiKey
    // Return access token
    return "eyJhbGc...", nil
}

// 2. In integration_wire.go
func New(deps Dependencies) Export {
    logger := deps.Logger.With(zap.String("layer", "integrations"))
    
    // Create token service (or get it from dependencies)
    tokenService := &IrisTokenService{
        tokenURL: deps.Config.External.Iris.AuthURL,
        apiKey:   deps.Config.External.Iris.APIKey,
    }
    
    // Wrap with caching for efficiency
    cachedTokenProvider := httpclient.NewCachedTokenProvider(
        &TokenServiceAdapter{tokenService: tokenService},
        5*time.Minute,
        logger,
    )
    
    // Create auth header provider
    authHeaderProvider := httpclient.NewAuthHeaderProvider(
        cachedTokenProvider,
        "Bearer",
        logger,
    )
    
    // Create shared HTTP client with auth headers
    sharedHTTPClient := httpclient.NewClient(
        httpclient.Config{
            Timeout:        30 * time.Second,
            MaxIdleConns:   100,
            ServiceName:    "external_apis",
            HeaderProvider: authHeaderProvider, // ✅ All requests get auth token!
        },
        logger,
    )
    
    // Use this client for all integrations
    pharmacy := irispharmacy.Module(irispharmacy.ModuleDependencies{
        HTTPClient: sharedHTTPClient, // ✅ Automatic auth!
        // ...
    })
    
    billing := irisbilling.Module(irisbilling.ModuleDependencies{
        HTTPClient: sharedHTTPClient, // ✅ Automatic auth!
        // ...
    })
    
    return Export{
        PharmacyClient: pharmacy.PharmacyClient,
        BillingClient:  billing.BillingClient,
    }
}
```

**Benefits:**
- ✅ Token automatically fetched on first request
- ✅ Token cached for performance
- ✅ Token auto-refreshed before expiry
- ✅ All API calls get auth header automatically
- ✅ No manual token management
- ✅ Thread-safe caching

---

## 🎯 **Pattern Comparison**

### **Pattern 1: Static Headers**

```go
headerProvider := httpclient.NewStaticHeaderProvider(map[string]string{
    "X-API-Key": "your-api-key",
    "X-Client":  "rxintake",
})
```

**Use when:**
- Headers don't change
- Simple API key auth
- Common headers for all requests

---

### **Pattern 2: Dynamic Headers (Custom Function)**

```go
headerProvider := httpclient.HeaderProviderFunc(func(ctx context.Context) (map[string]string, error) {
    // Custom logic to build headers
    correlationID := getCorrelationIDFromContext(ctx)
    
    return map[string]string{
        "X-Correlation-ID": correlationID,
        "X-Timestamp":      time.Now().Format(time.RFC3339),
    }, nil
})
```

**Use when:**
- Headers depend on context
- Need to generate headers dynamically
- Headers include request-specific data

---

### **Pattern 3: Token-Based Auth (Recommended for OAuth/JWT)**

```go
// Your token service
type TokenService interface {
    GetAccessToken(ctx context.Context) (string, error)
}

// Wrap with caching
cachedProvider := httpclient.NewCachedTokenProvider(
    tokenServiceAdapter,
    5*time.Minute,
    logger,
)

// Create auth header provider
authProvider := httpclient.NewAuthHeaderProvider(
    cachedProvider,
    "Bearer",
    logger,
)

// Use in HTTP client
httpClient := httpclient.NewClient(
    httpclient.Config{
        HeaderProvider: authProvider,
    },
    logger,
)
```

**Use when:**
- OAuth 2.0 authentication
- JWT tokens with expiry
- Tokens need refresh
- Multiple services share same auth

---

## 💡 **Real-World Examples**

### **Example 1: API Key Authentication**

```go
// Simple API key that doesn't change
func New(deps Dependencies) Export {
    headerProvider := httpclient.NewStaticHeaderProvider(map[string]string{
        "X-API-Key": deps.Config.External.Billing.APIKey,
    })
    
    httpClient := httpclient.NewClient(
        httpclient.Config{
            HeaderProvider: headerProvider,
            ServiceName:    "external_apis",
        },
        logger,
    )
    
    // All requests will include X-API-Key header
}
```

---

### **Example 2: OAuth 2.0 Client Credentials**

```go
// Token service that calls OAuth endpoint
type OAuthTokenService struct {
    clientID     string
    clientSecret string
    tokenURL     string
    httpClient   *http.Client
}

func (s *OAuthTokenService) GetAccessToken(ctx context.Context) (string, error) {
    // Call OAuth token endpoint
    resp, err := s.httpClient.PostForm(s.tokenURL, url.Values{
        "grant_type":    {"client_credentials"},
        "client_id":     {s.clientID},
        "client_secret": {s.clientSecret},
    })
    // ... parse response, return token
}

// In integration_wire.go
func New(deps Dependencies) Export {
    // Create OAuth service
    oauthService := &OAuthTokenService{
        clientID:     deps.Config.External.OAuth.ClientID,
        clientSecret: deps.Config.External.OAuth.ClientSecret,
        tokenURL:     deps.Config.External.OAuth.TokenURL,
        httpClient:   http.DefaultClient,
    }
    
    // Adapter to TokenProvider
    tokenProvider := &TokenServiceAdapter{
        tokenService: oauthService,
    }
    
    // Wrap with caching
    cachedProvider := httpclient.NewCachedTokenProvider(
        tokenProvider,
        5*time.Minute,
        logger,
    )
    
    // Create auth header provider
    authProvider := httpclient.NewAuthHeaderProvider(
        cachedProvider,
        "Bearer",
        logger,
    )
    
    // Create HTTP client
    httpClient := httpclient.NewClient(
        httpclient.Config{
            HeaderProvider: authProvider,
            ServiceName:    "external_apis",
        },
        logger,
    )
    
    // Use for integrations - all will have OAuth token!
}
```

---

### **Example 3: Multiple Header Types**

Combine multiple header providers:

```go
// Composite header provider
type CompositeHeaderProvider struct {
    providers []httpclient.HeaderProvider
}

func (p *CompositeHeaderProvider) GetHeaders(ctx context.Context) (map[string]string, error) {
    headers := make(map[string]string)
    
    for _, provider := range p.providers {
        h, err := provider.GetHeaders(ctx)
        if err != nil {
            return nil, err
        }
        for key, value := range h {
            headers[key] = value
        }
    }
    
    return headers, nil
}

// Usage
composite := &CompositeHeaderProvider{
    providers: []httpclient.HeaderProvider{
        httpclient.NewStaticHeaderProvider(map[string]string{
            "X-API-Version": "v1",
            "X-Client-ID":   "rxintake",
        }),
        authHeaderProvider, // From token service
        httpclient.HeaderProviderFunc(func(ctx context.Context) (map[string]string, error) {
            return map[string]string{
                "X-Request-ID": getRequestIDFromContext(ctx),
            }, nil
        }),
    },
}

httpClient := httpclient.NewClient(
    httpclient.Config{
        HeaderProvider: composite,
    },
    logger,
)
```

---

## 🔒 **Authentication Patterns**

### **Pattern 1: API Key (Simplest)**

```go
// Static API key
headerProvider := httpclient.NewStaticHeaderProvider(map[string]string{
    "X-API-Key": "your-api-key",
})
```

**Pros:**
- Simple
- No token refresh needed
- Fast

**Cons:**
- Less secure
- Key doesn't expire
- Can't revoke easily

---

### **Pattern 2: Bearer Token with Auto-Refresh (Recommended)**

```go
// Token service (your existing service)
type YourTokenService interface {
    GetAccessToken(ctx context.Context) (string, error)
}

// Adapter
type TokenServiceAdapter struct {
    service YourTokenService
}

func (a *TokenServiceAdapter) GetToken(ctx context.Context) (string, error) {
    return a.service.GetAccessToken(ctx)
}

// Cached provider
cachedProvider := httpclient.NewCachedTokenProvider(
    &TokenServiceAdapter{service: yourTokenService},
    5*time.Minute, // Refresh 5 min before expiry
    logger,
)

// Auth header provider
authProvider := httpclient.NewAuthHeaderProvider(
    cachedProvider,
    "Bearer",
    logger,
)

// Use in HTTP client
httpClient := httpclient.NewClient(
    httpclient.Config{
        HeaderProvider: authProvider,
    },
    logger,
)
```

**Pros:**
- ✅ Automatic token refresh
- ✅ Thread-safe caching
- ✅ Prevents unnecessary token fetches
- ✅ Works with existing token service
- ✅ More secure (tokens expire)

**Cons:**
- Slightly more complex setup
- Need token service

---

### **Pattern 3: Per-Request Token (If Needed)**

```go
// http_client.go - Get fresh token per request
func (c *HTTPClient) GetInvoice(ctx, prescriptionID) (*InvoiceResponse, error) {
    // Get token for this specific request
    token, err := c.tokenService.GetAccessToken(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get token: %w", err)
    }
    
    url := replacePathParams(c.endpoints.GetInvoiceEndpoint(), map[string]string{
        "prescriptionID": prescriptionID,
    })
    
    resp, err := c.client.Get(ctx, url, map[string]string{
        "Authorization": "Bearer " + token,
    })
    // ...
}
```

**When to use:**
- Tokens are request-specific
- Short-lived tokens
- Different tokens per endpoint

---

## 🏗️ **Recommended Architecture**

### **Setup in integration_wire.go**

```go
package integrations

import (
    "pharmacy-modernization-project-model/internal/platform/httpclient"
    // ...
)

type Dependencies struct {
    Config       *config.Config
    Logger       *zap.Logger
    TokenService TokenService // ✅ Pass token service as dependency
}

func New(deps Dependencies) Export {
    logger := deps.Logger.With(zap.String("layer", "integrations"))
    
    // Option 1: Static headers (API Key)
    var headerProvider httpclient.HeaderProvider
    if deps.Config.External.APIKey != "" {
        headerProvider = httpclient.NewStaticHeaderProvider(map[string]string{
            "X-API-Key": deps.Config.External.APIKey,
        })
    }
    
    // Option 2: Token-based auth (if token service provided)
    if deps.TokenService != nil {
        // Adapter
        tokenProvider := &TokenServiceAdapter{
            tokenService: deps.TokenService,
        }
        
        // Cache tokens
        cachedProvider := httpclient.NewCachedTokenProvider(
            tokenProvider,
            5*time.Minute,
            logger,
        )
        
        // Auth headers
        headerProvider = httpclient.NewAuthHeaderProvider(
            cachedProvider,
            "Bearer",
            logger,
        )
    }
    
    // Create shared HTTP client
    sharedHTTPClient := httpclient.NewClient(
        httpclient.Config{
            Timeout:        30 * time.Second,
            MaxIdleConns:   100,
            ServiceName:    "external_apis",
            HeaderProvider: headerProvider, // ✅ Auth for all requests
        },
        logger,
    )
    
    // Initialize integrations with authenticated client
    pharmacy := irispharmacy.Module(irispharmacy.ModuleDependencies{
        Config:     /* ... */,
        HTTPClient: sharedHTTPClient, // ✅ Already has auth!
        // ...
    }).PharmacyClient
    
    billing := irisbilling.Module(irisbilling.ModuleDependencies{
        Config:     /* ... */,
        HTTPClient: sharedHTTPClient, // ✅ Already has auth!
        // ...
    }).BillingClient
    
    return Export{
        PharmacyClient: pharmacy,
        BillingClient:  billing,
    }
}
```

---

## 📝 **Configuration**

### **YAML Config for API Key:**

```yaml
external:
  api_key: "${IRIS_API_KEY}"  # From environment variable
  
  pharmacy:
    endpoints:
      get_prescription: "https://api.iris.com/pharmacy/v1/prescriptions/{prescriptionID}"
```

### **YAML Config for OAuth:**

```yaml
external:
  oauth:
    client_id: "${OAUTH_CLIENT_ID}"
    client_secret: "${OAUTH_CLIENT_SECRET}"
    token_url: "https://auth.iris.com/oauth/token"
  
  pharmacy:
    endpoints:
      get_prescription: "https://api.iris.com/pharmacy/v1/prescriptions/{prescriptionID}"
```

---

## 🎯 **Best Practices Summary**

### **For Headers:**

1. **Global headers** (all requests): Use `HeaderProvider` in HTTP client config
2. **Endpoint-specific headers**: Pass in individual method calls
3. **Hybrid**: Combine both for flexibility

### **For Authentication:**

1. **API Key** (simple): Use `StaticHeaderProvider`
2. **OAuth/JWT** (recommended): Use `CachedTokenProvider` + `AuthHeaderProvider`
3. **Per-request tokens**: Pass directly in method calls

### **For Token Services:**

1. ✅ **DO**: Create an adapter to wrap your existing token service
2. ✅ **DO**: Use `CachedTokenProvider` to cache tokens
3. ✅ **DO**: Pass token service as dependency to integration layer
4. ✅ **DO**: Let the HTTP client handle token refresh automatically
5. ❌ **DON'T**: Fetch tokens in every endpoint implementation
6. ❌ **DON'T**: Store tokens in global variables

---

## 🚀 **Quick Start Examples**

### **Simple API Key:**

```go
httpClient := httpclient.NewClient(
    httpclient.Config{
        HeaderProvider: httpclient.NewStaticHeaderProvider(map[string]string{
            "X-API-Key": apiKey,
        }),
    },
    logger,
)
```

### **OAuth with Token Service:**

```go
authProvider := httpclient.NewAuthHeaderProvider(
    httpclient.NewCachedTokenProvider(
        tokenServiceAdapter,
        5*time.Minute,
        logger,
    ),
    "Bearer",
    logger,
)

httpClient := httpclient.NewClient(
    httpclient.Config{
        HeaderProvider: authProvider,
    },
    logger,
)
```

### **Per-Endpoint Headers:**

```go
// Just pass headers when calling
resp, err := client.Get(ctx, url, map[string]string{
    "X-Custom-Header": "value",
})
```

---

## 🎉 **Summary**

**Headers:**
- Global headers → Use `HeaderProvider` in client config
- Per-endpoint → Pass in method call
- Both → Combine approaches

**Authentication:**
- Use `CachedTokenProvider` for token caching
- Use `AuthHeaderProvider` for auth headers
- Pass your token service as dependency
- Let the framework handle token refresh

**Result:** Clean, maintainable, efficient authentication! 🚀

