# Caching Implementation Guide

## Overview

This project implements a flexible, multi-strategy caching layer with **Memcached** as the primary distributed cache. The caching system supports multiple strategies (Memcached, Redis, Memory, Hybrid) and can be easily injected into services throughout the application.

## Architecture

### Cache Strategies

1. **Memcached** (Primary) - Distributed cache for production use
2. **Redis** - Alternative distributed cache
3. **Memory** (Ristretto) - Ultra-fast in-memory cache
4. **Hybrid** - Combines memory (L1) + distributed cache (L2) for optimal performance

### Key Components

```
internal/platform/cache/
├── cache.go          # Cache interface and types
├── memcached.go      # Memcached implementation
├── redis.go          # Redis implementation
├── memory.go         # Ristretto in-memory implementation
├── hybrid.go         # Hybrid cache (memory + distributed)
├── factory.go        # Cache factory for creating instances
├── middleware.go     # Metrics and logging middleware
└── health.go         # Health check utilities

internal/app/
└── cache_builder.go  # Builder for creating cache instances
```

## Configuration

### YAML Configuration (`internal/configs/app.yaml`)

```yaml
cache:
  # Memcached as primary cache (distributed cache)
  memcached:
    addr: "localhost:11211"
    prefix: "rx:"
  
  # Redis as alternative distributed cache
  redis:
    addr: "localhost:6379"
    password: ""
    db: 0
    prefix: "rx:"
  
  # In-memory cache configuration (for hybrid or memory-only)
  memory:
    max_cost: 67108864  # 64MB (64<<20)
    buffer_items: 64
    metrics: true
    default_ttl: "30m"
```

### Docker/Podman Compose

The `podman/compose.yml` includes Memcached:

```yaml
services:
  memcached:
    image: memcached:1.6-alpine
    container_name: memcached
    restart: always
    ports:
      - "11211:11211"
    command: memcached -m 64 -c 1024 -I 4m
    healthcheck:
      test: ["CMD", "nc", "-z", "localhost", "11211"]
      interval: 10s
      timeout: 5s
      retries: 3
```

## Usage

### Starting the Cache Service

```bash
# Start Memcached and MongoDB
cd podman
podman-compose up -d

# Or with Docker
docker compose up -d
```

### Creating Cache Instances

The `CacheBuilder` provides flexible cache creation:

```go
// In wire.go or module initialization
cacheBuilder := NewCacheBuilder(config, logger)

// Create Memcached cache (primary)
primaryCache, err := cacheBuilder.BuildMemcachedCache("rx:")

// Create memory cache
memCache, err := cacheBuilder.BuildMemoryCache(67108864) // 64MB

// Create hybrid cache
hybridCache, err := cacheBuilder.BuildHybridCache(33554432) // 32MB memory

// Create Redis cache
redisCache, err := cacheBuilder.BuildRedisCache("localhost:6379", "api:")

// Create custom strategy cache
customCache, err := cacheBuilder.BuildCache("memcached", CacheInstanceConfig{
    Memcached: cache.MemcachedConfig{
        Addr:   "localhost:11211",
        Prefix: "custom:",
    },
})
```

### Service Integration

#### Patient Service Example

```go
type patientSvc struct {
    repo  repo.PatientRepository
    cache Cache
    log   *zap.Logger
}

func (s *patientSvc) GetByID(ctx context.Context, id string) (m.Patient, error) {
    cacheKey := fmt.Sprintf("patient:id:%s", id)

    // Try cache first
    if s.cache != nil {
        if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
            var patient m.Patient
            if err := json.Unmarshal(cached, &patient); err == nil {
                return patient, nil
            }
        }
    }

    // Cache miss - get from repository
    patient, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return m.Patient{}, err
    }

    // Cache the result
    if s.cache != nil {
        if data, err := json.Marshal(patient); err == nil {
            s.cache.Set(ctx, cacheKey, data, 30*time.Minute)
        }
    }

    return patient, nil
}
```

## Cache Key Strategy

### Key Patterns

```go
// Patient data
"patient:id:{id}"              // Individual patient by ID
"patient:list:{query}:{limit}:{offset}"  // Patient list query
"patient:count:{query}"        // Patient count query

// Prescription data
"prescription:id:{id}"         // Individual prescription by ID
"prescription:list:{status}:{limit}:{offset}"  // Prescription list
"prescription:count:status:{status}"  // Count by status

// Dashboard aggregates
"dashboard:summary"            // Dashboard summary data

// External API data
"iris:prescription:{id}"       // Cached external API response
"iris:invoice:{id}"            // Cached invoice data
```

### TTL Recommendations

```go
const (
    // Static/rarely changing data
    TTLPatientData    = 30 * time.Minute
    TTLAddressData    = 1 * time.Hour
    
    // Dynamic data
    TTLPrescriptionData = 15 * time.Minute
    TTLDashboardSummary = 2 * time.Minute
    
    // External API data (critical for performance)
    TTLIrisPrescription = 10 * time.Minute
    TTLIrisInvoice     = 15 * time.Minute
    
    // Count/aggregate data
    TTLCountData      = 5 * time.Minute
)
```

## Monitoring and Metrics

### Cache Statistics

```go
// Get cache stats
stats := cache.Stats()

// Stats structure
type CacheStats struct {
    Hits      int64   // Number of cache hits
    Misses    int64   // Number of cache misses
    Evictions int64   // Number of evictions
    Errors    int64   // Number of errors
    HitRate   float64 // Hit rate (hits / total requests)
    Size      int64   // Current cache size
    MaxSize   int64   // Maximum cache size
}
```

### Health Checks

```go
healthChecker := cache.NewCacheHealthChecker(cacheService)
err := healthChecker.Check(ctx)
```

## Best Practices

### 1. Cache-Aside Pattern

Always check cache first, then fallback to data source:

```go
// 1. Check cache
if cached, err := cache.Get(ctx, key); err == nil {
    return cached, nil
}

// 2. Get from source
data, err := source.Get(ctx, id)
if err != nil {
    return nil, err
}

// 3. Update cache
cache.Set(ctx, key, data, ttl)
return data, nil
```

### 2. Cache Invalidation

Invalidate cache on writes:

```go
func (s *service) Update(ctx context.Context, id string, data Data) error {
    // Update data source
    if err := s.repo.Update(ctx, id, data); err != nil {
        return err
    }

    // Invalidate related cache entries
    s.cache.Delete(ctx, fmt.Sprintf("data:id:%s", id))
    s.cache.Delete(ctx, "data:list:*")  // Invalidate list queries
    
    return nil
}
```

### 3. Graceful Degradation

Handle cache failures gracefully:

```go
if s.cache != nil {
    if cached, err := s.cache.Get(ctx, key); err == nil {
        return cached, nil
    }
    // Log cache miss but continue
    s.log.Debug("Cache miss", zap.String("key", key))
}

// Continue with data source query
return s.repo.Get(ctx, id)
```

### 4. Appropriate TTLs

- **Frequently changing data**: 2-5 minutes
- **Moderately changing data**: 15-30 minutes
- **Rarely changing data**: 1+ hours
- **Count queries**: 5 minutes (cheaper to recompute)
- **External API calls**: 10-15 minutes (balance freshness vs. API costs)

## Advanced Scenarios

### Different Caches for Different Layers

```go
// In wire.go
serviceCache, _ := cacheBuilder.BuildMemoryCache(33554432)     // 32MB fast memory
apiCache, _ := cacheBuilder.BuildMemcachedCache("api:")        // Distributed for API responses
externalCache, _ := cacheBuilder.BuildHybridCache(16777216)    // 16MB hybrid for external APIs

// Pass to modules
patientModule.Module(router, &patient.ModuleDependencies{
    ServiceCache:  serviceCache,   // For business logic
    APICache:      apiCache,        // For API responses
    ExternalCache: externalCache,   // For external API calls
})
```

### Custom Cache Strategy per Use Case

```go
// Fast in-memory for frequently accessed data
frequentCache, _ := cacheBuilder.BuildMemoryCache(67108864)

// Distributed for shared data across instances
sharedCache, _ := cacheBuilder.BuildMemcachedCache("shared:")

// Hybrid for best of both worlds
hybridCache, _ := cacheBuilder.BuildHybridCache(33554432)
```

## Troubleshooting

### Cache Not Working

1. **Check Memcached is running**:
   ```bash
   podman ps | grep memcached
   # or
   telnet localhost 11211
   ```

2. **Check configuration**:
   ```bash
   # Verify cache config in app.yaml
   cat internal/configs/app.yaml | grep -A 10 "cache:"
   ```

3. **Enable debug logging**:
   ```yaml
   logging:
     level: debug
   ```

### Performance Issues

1. **Monitor cache hit rates**:
   ```go
   stats := cache.Stats()
   log.Info("Cache stats", 
       zap.Float64("hit_rate", stats.HitRate),
       zap.Int64("hits", stats.Hits),
       zap.Int64("misses", stats.Misses))
   ```

2. **Adjust TTLs** based on hit rates and data freshness requirements

3. **Consider hybrid strategy** for frequently accessed data

## Migration Guide

### From No Cache to Memcached

1. **Start Memcached**:
   ```bash
   cd podman && podman-compose up -d memcached
   ```

2. **Add cache configuration** to `app.yaml`

3. **Inject cache into services** via module dependencies

4. **Update service methods** to use cache-aside pattern

5. **Monitor and tune** TTLs based on usage patterns

## Dependencies

Add to `go.mod`:

```go
require (
    github.com/bradfitz/gomemcache v0.0.0-20230905024940-24af94b03874
    github.com/dgraph-io/ristretto v0.1.1
    github.com/go-redis/redis/v8 v8.11.5
)
```

## Summary

The caching implementation provides:

- ✅ **Flexible strategy selection** (Memcached, Redis, Memory, Hybrid)
- ✅ **Memcached as primary** distributed cache
- ✅ **Easy service injection** via builder pattern
- ✅ **Built-in metrics** and health checks
- ✅ **Graceful fallbacks** for resilience
- ✅ **Simple configuration** via YAML
- ✅ **Production-ready** with proper error handling

The system is designed for maximum flexibility while keeping implementation simple and maintainable.

