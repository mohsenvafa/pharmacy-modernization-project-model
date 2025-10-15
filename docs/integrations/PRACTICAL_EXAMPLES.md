# Practical Examples - Headers and Authentication

## Real-World Scenarios

### Scenario 1: Simple API Key (No Token Service)

**When:** External API uses a static API key for authentication

```go
// integration_wire.go
func New(deps Dependencies) Export {
    logger := deps.Logger.With(zap.String("layer", "integrations"))
    
    // ✅ Simple: Static API key from config
    headerProvider := httpclient.NewStaticHeaderProvider(map[string]string{
        "X-API-Key":     deps.Config.External.IRIS.APIKey,
        "X-Client-ID":   "rxintake-app",
        "X-API-Version": "v1",
    })
    
    // Create shared HTTP client with API key
    sharedHTTPClient := httpclient.NewClient(
        httpclient.Config{
            Timeout:        30 * time.Second,
            MaxIdleConns:   100,
            ServiceName:    "external_apis",
            HeaderProvider: headerProvider, // ✅ All requests get these headers
        },
        logger,
    )
    
    // Initialize integrations
    pharmacy := irispharmacy.Module(irispharmacy.ModuleDependencies{
        Config:     /* ... */,
        HTTPClient: sharedHTTPClient,
        // ...
    }).PharmacyClient
    
    billing := irisbilling.Module(irisbilling.ModuleDependencies{
        Config:     /* ... */,
        HTTPClient: sharedHTTPClient,
        // ...
    }).BillingClient
    
    return Export{
        PharmacyClient: pharmacy,
        BillingClient:  billing,
    }
}
```

**Result:** All API calls automatically include:
```
X-API-Key: your-key-here
X-Client-ID: rxintake-app
X-API-Version: v1
```

---

### Scenario 2: OAuth with Existing Token Service

**When:** You have a separate service that handles OAuth token management

#### Step 1: Define Token Service Adapter

```go
// internal/integrations/auth_adapter.go
package integrations

import (
    "context"
    "pharmacy-modernization-project-model/internal/platform/httpclient"
)

// TokenService is your existing interface (could be anywhere in your app)
type TokenService interface {
    GetAccessToken(ctx context.Context) (string, error)
    RefreshToken(ctx context.Context) (string, error)
}

// TokenServiceAdapter adapts your TokenService to httpclient.TokenProvider
type TokenServiceAdapter struct {
    tokenService TokenService
}

func NewTokenServiceAdapter(service TokenService) *TokenServiceAdapter {
    return &TokenServiceAdapter{
        tokenService: service,
    }
}

func (a *TokenServiceAdapter) GetToken(ctx context.Context) (string, error) {
    return a.tokenService.GetAccessToken(ctx)
}
```

#### Step 2: Use in Integration Wire

```go
// integration_wire.go
type Dependencies struct {
    Config       *config.Config
    Logger       *zap.Logger
    TokenService TokenService // ✅ Existing token service
}

func New(deps Dependencies) Export {
    logger := deps.Logger.With(zap.String("layer", "integrations"))
    
    var headerProvider httpclient.HeaderProvider
    
    // If token service is provided, use token-based auth
    if deps.TokenService != nil {
        // Adapt your token service
        tokenAdapter := NewTokenServiceAdapter(deps.TokenService)
        
        // Cache tokens for efficiency
        cachedProvider := httpclient.NewCachedTokenProvider(
            tokenAdapter,
            5*time.Minute, // Refresh 5 min before expiry
            logger,
        )
        
        // Create auth header provider
        headerProvider = httpclient.NewAuthHeaderProvider(
            cachedProvider,
            "Bearer",
            logger,
        )
        
        logger.Info("using token-based authentication for external APIs")
    } else {
        // Fallback to static API key if no token service
        if apiKey := deps.Config.External.APIKey; apiKey != "" {
            headerProvider = httpclient.NewStaticHeaderProvider(map[string]string{
                "X-API-Key": apiKey,
            })
            logger.Info("using API key authentication for external APIs")
        }
    }
    
    // Create shared HTTP client
    sharedHTTPClient := httpclient.NewClient(
        httpclient.Config{
            Timeout:        30 * time.Second,
            MaxIdleConns:   100,
            ServiceName:    "external_apis",
            HeaderProvider: headerProvider, // ✅ Auth handled automatically
        },
        logger,
    )
    
    // ... initialize integrations ...
}
```

#### Step 3: Update App Wire

```go
// app/wire.go
func (a *App) wire() error {
    logger := logging.NewLogger(a.Cfg)
    
    // Initialize your token service (if you have one)
    var tokenService integrations.TokenService
    if a.Cfg.External.OAuth.Enabled {
        tokenService = auth.NewOAuthTokenService(
            a.Cfg.External.OAuth.ClientID,
            a.Cfg.External.OAuth.ClientSecret,
            a.Cfg.External.OAuth.TokenURL,
            logger.Base,
        )
    }
    
    // Initialize integrations with token service
    integration := integrations.New(integrations.Dependencies{
        Config:       a.Cfg,
        Logger:       logger.Base,
        TokenService: tokenService, // ✅ Pass your token service
    })
    
    // ... rest of wire ...
}
```

**Logs you'll see:**
```
INFO  using token-based authentication for external APIs
INFO  shared http client created for external API integrations
INFO  fetching new access token
INFO  access token refreshed expires_at=2025-10-14T16:30:00Z
DEBUG using cached access token expires_at=2025-10-14T16:30:00Z
DEBUG using cached access token expires_at=2025-10-14T16:30:00Z
...
INFO  fetching new access token  (after 55 minutes)
INFO  access token refreshed expires_at=2025-10-14T17:25:00Z
```

---

### Scenario 3: Different Headers Per Service

**When:** Billing needs auth, Pharmacy just needs API key

```go
// integration_wire.go
func New(deps Dependencies) Export {
    logger := deps.Logger.With(zap.String("layer", "integrations"))
    
    // Create base HTTP client (no auth)
    baseHTTPClient := httpclient.NewClient(
        httpclient.Config{
            Timeout:      30 * time.Second,
            MaxIdleConns: 100,
            ServiceName:  "external_apis",
            // No HeaderProvider - will add per-service
        },
        logger,
    )
    
    // Pharmacy: Just needs API key
    pharmacyHeaders := httpclient.NewStaticHeaderProvider(map[string]string{
        "X-API-Key": deps.Config.External.Pharmacy.APIKey,
    })
    
    pharmacyClient := httpclient.NewClient(
        httpclient.Config{
            Timeout:        30 * time.Second,
            ServiceName:    "iris_pharmacy",
            HeaderProvider: pharmacyHeaders, // ✅ Pharmacy-specific headers
        },
        logger,
    )
    
    // Billing: Needs OAuth token
    if deps.TokenService != nil {
        tokenAdapter := NewTokenServiceAdapter(deps.TokenService)
        cachedProvider := httpclient.NewCachedTokenProvider(tokenAdapter, 5*time.Minute, logger)
        billingHeaders := httpclient.NewAuthHeaderProvider(cachedProvider, "Bearer", logger)
        
        billingClient := httpclient.NewClient(
            httpclient.Config{
                Timeout:        30 * time.Second,
                ServiceName:    "iris_billing",
                HeaderProvider: billingHeaders, // ✅ Billing-specific auth
            },
            logger,
        )
        
        pharmacy := irispharmacy.Module(irispharmacy.ModuleDependencies{
            HTTPClient: pharmacyClient, // ✅ Uses API key
            // ...
        }).PharmacyClient
        
        billing := irisbilling.Module(irisbilling.ModuleDependencies{
            HTTPClient: billingClient, // ✅ Uses OAuth token
            // ...
        }).BillingClient
        
        return Export{
            PharmacyClient: pharmacy,
            BillingClient:  billing,
        }
    }
    
    // ... fallback ...
}
```

---

### Scenario 4: Context-Based Headers (Request ID, User ID)

**When:** Headers depend on request context

```go
// Create dynamic header provider
contextHeaderProvider := httpclient.HeaderProviderFunc(func(ctx context.Context) (map[string]string, error) {
    headers := make(map[string]string)
    
    // Extract from context
    if requestID := ctx.Value("request_id"); requestID != nil {
        headers["X-Request-ID"] = requestID.(string)
    }
    
    if userID := ctx.Value("user_id"); userID != nil {
        headers["X-User-ID"] = userID.(string)
    }
    
    return headers, nil
})

// Combine with auth
compositeProvider := &CompositeHeaderProvider{
    providers: []httpclient.HeaderProvider{
        authHeaderProvider,    // OAuth token
        contextHeaderProvider, // Request-specific headers
    },
}

httpClient := httpclient.NewClient(
    httpclient.Config{
        HeaderProvider: compositeProvider,
    },
    logger,
)
```

---

## 🔐 **Token Invalidation**

If a token becomes invalid (401 Unauthorized), you can force refresh:

```go
// In your HTTP client implementation
func (c *HTTPClient) GetInvoice(ctx, prescriptionID) (*InvoiceResponse, error) {
    url := replacePathParams(c.endpoints.GetInvoiceEndpoint(), map[string]string{
        "prescriptionID": prescriptionID,
    })
    
    var response InvoiceResponse
    err := c.client.GetJSON(ctx, url, &response)
    
    if err != nil && isUnauthorizedError(err) {
        c.logger.Warn("received 401, invalidating cached token")
        
        // If you have access to the token provider, invalidate it
        // This would require passing it as a dependency
        // tokenProvider.InvalidateToken()
        
        // Retry the request (will fetch fresh token)
        err = c.client.GetJSON(ctx, url, &response)
    }
    
    return &response, err
}
```

---

## ✅ **Recommended Approach**

For most applications:

```go
// 1. Define your token service interface
type TokenService interface {
    GetAccessToken(ctx context.Context) (string, error)
}

// 2. Create adapter
tokenAdapter := NewTokenServiceAdapter(yourTokenService)

// 3. Add caching
cachedProvider := httpclient.NewCachedTokenProvider(
    tokenAdapter,
    5*time.Minute,
    logger,
)

// 4. Create auth header provider
authProvider := httpclient.NewAuthHeaderProvider(
    cachedProvider,
    "Bearer",
    logger,
)

// 5. Use in HTTP client
httpClient := httpclient.NewClient(
    httpclient.Config{
        HeaderProvider: authProvider,
    },
    logger,
)

// 6. All API calls automatically get fresh tokens! ✅
```

**Benefits:**
- ✅ Automatic token management
- ✅ Efficient caching
- ✅ Thread-safe
- ✅ Auto-refresh before expiry
- ✅ Works with existing token service
- ✅ No changes needed in endpoint implementations

---

## 📊 **Performance**

### Without Caching:
```
Request 1: Fetch token (200ms) + API call (100ms) = 300ms
Request 2: Fetch token (200ms) + API call (100ms) = 300ms
Request 3: Fetch token (200ms) + API call (100ms) = 300ms
Total: 900ms for 3 requests
```

### With Caching:
```
Request 1: Fetch token (200ms) + API call (100ms) = 300ms
Request 2: Use cached token + API call (100ms) = 100ms
Request 3: Use cached token + API call (100ms) = 100ms
Total: 500ms for 3 requests (45% faster!)
```

---

## 🎉 **Summary**

**Two Best Practices Implemented:**

### 1. **Headers for All/Some Endpoints:**
- ✅ **Global headers**: Use `HeaderProvider` in HTTP client config
- ✅ **Endpoint-specific**: Pass headers in individual calls
- ✅ **Both**: Combine for maximum flexibility

### 2. **Authorization with Token Service:**
- ✅ **Create adapter**: Wrap your token service
- ✅ **Add caching**: Use `CachedTokenProvider`
- ✅ **Auth headers**: Use `AuthHeaderProvider`
- ✅ **Pass to client**: Set as `HeaderProvider`
- ✅ **Automatic refresh**: Tokens refreshed before expiry

**Result: Clean, efficient, production-ready authentication!** 🚀

