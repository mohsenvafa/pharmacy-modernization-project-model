# Caching Implementation Guide

## Overview

This project implements a simple, flexible caching layer with **MongoDB** as the primary distributed cache and **Ristretto** for ultra-fast in-memory caching. The system uses your existing MongoDB connection, so no additional services are required.

## Architecture

### Cache Strategies

1. **MongoDB** (Primary) - Persistent distributed cache using existing MongoDB
2. **Memory** (Ristretto) - Ultra-fast in-memory cache
3. **Hybrid** - Combines memory (L1) + MongoDB (L2) for optimal performance

### Key Components

```
internal/platform/cache/
├── cache.go          # Cache interface and types
├── mongodb.go        # MongoDB implementation with TTL
├── memory.go         # Ristretto in-memory implementation
├── hybrid.go         # Hybrid cache (memory + MongoDB)
├── factory.go        # Cache factory for creating instances
├── middleware.go     # Metrics and logging middleware
└── health.go         # Health check utilities

internal/app/
├── builder/
│   ├── cache.go      # Cache builder for creating cache instances
│   └── mongodb.go    # MongoDB builder for main database
├── cache.wire.go     # Cache wiring logic
└── mongodb.wire.go   # MongoDB wiring logic
```

## Configuration

### YAML Configuration (`internal/configs/app.yaml`)

```yaml
cache:
  # MongoDB cache configuration (independent from main database)
  mongodb:
    uri: "mongodb://admin:admin123@localhost:27017"
    database: "pharmacy_modernization_cache"  # Separate database for cache
    collection: "cache"
    connection:
      max_pool_size: 50
      min_pool_size: 5
      max_idle_time: "30m"
      connect_timeout: "10s"
      socket_timeout: "30s"
  
  # In-memory cache configuration
  memory:
    max_cost: 67108864  # 64MB (64<<20)
    buffer_items: 64
    metrics: true
    default_ttl: "30m"
```

### Docker/Podman Compose

Only MongoDB is required (already in your `compose.yml`):

```yaml
services:
  mongodb:
    image: mongo:6.0
    container_name: mongodb
    restart: always
    ports:
      - "27017:27017"
```

## Usage

### Starting the Cache Service

```bash
# Start MongoDB
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

// Create MongoDB cache (uses separate cache MongoDB connection)
cacheMongoConnMgr, _ := CreateCacheMongoDBConnection(config, logger)
cacheCollection := GetCacheMongoCollection(cacheMongoConnMgr, config)
mongoCache, err := cacheBuilder.BuildMongoDBCache(cacheCollection, "rx:")

// Create memory cache
memCache, err := cacheBuilder.BuildMemoryCache(67108864) // 64MB

// Create hybrid cache (memory + MongoDB)
hybridCache, err := cacheBuilder.BuildHybridCache(33554432) // 32MB memory
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

## MongoDB Cache Details

### How It Works

MongoDB cache uses a dedicated `cache` collection with automatic TTL expiration:

```go
type cacheDocument struct {
    Key      string    `bson:"_id"`       // Cache key (unique)
    Value    []byte    `bson:"value"`     // Cached data
    ExpireAt time.Time `bson:"expireAt"`  // TTL expiration time
}
```

### TTL Index

The MongoDB cache automatically creates a TTL index on the `expireAt` field:

```javascript
db.cache.createIndex({ "expireAt": 1 }, { expireAfterSeconds: 0 })
```

This tells MongoDB to automatically delete documents when `expireAt` time is reached.

### Database and Collection Structure

**Separate databases for clear separation:**

```
MongoDB Server (localhost:27017)
├── pharmacy_modernization/        # Main database
│   ├── patients
│   ├── prescriptions
│   └── addresses
└── pharmacy_modernization_cache/  # Cache database (separate)
    └── cache                       # Cache collection (auto-expires)
```

**Why separate databases?**
- ✅ Clear separation of concerns
- ✅ Independent backup/restore strategies
- ✅ Different retention policies
- ✅ Easy to drop cache without affecting data
- ✅ Independent monitoring and metrics

### Benefits

- ✅ **No additional service** required (uses existing MongoDB)
- ✅ **Persistent cache** (survives restarts)
- ✅ **Automatic TTL expiration** via MongoDB indexes
- ✅ **Multi-instance support** (shared cache across app instances)
- ✅ **Simple setup** (just use existing MongoDB connection)

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

### Verify Cache in MongoDB

```bash
mongosh mongodb://admin:admin123@localhost:27017

use pharmacy_modernization
db.cache.find().pretty()
db.cache.stats()
db.cache.getIndexes()
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

## Strategy Comparison

| Strategy | Speed | Persistence | Multi-Instance | Setup | Use Case |
|----------|-------|-------------|----------------|-------|----------|
| **Memory** | ⚡⚡⚡⚡⚡ | ❌ | ❌ | Easy | Single-instance, temporary |
| **MongoDB** | ⚡⚡⚡ | ✅ | ✅ | Easy | Multi-instance, persistent |
| **Hybrid** | ⚡⚡⚡⚡⚡ | ✅ | ✅ | Easy | Best performance + persistence |

## Advanced Scenarios

### Different Caches for Different Layers

```go
// In wire.go
serviceCache, _ := cacheBuilder.BuildMemoryCache(33554432)     // 32MB fast memory

cacheCollection := GetCacheCollection(mongoConnMgr)
apiCache, _ := cacheBuilder.BuildMongoDBCache(cacheCollection, "api:")  // Persistent for API

// Pass to modules
patientModule.Module(router, &patient.ModuleDependencies{
    ServiceCache:  serviceCache,   // For business logic
    APICache:      apiCache,        // For API responses
})
```

### Hybrid Cache for Best Performance

```go
// Combines ultra-fast memory L1 with persistent MongoDB L2
hybridCache, _ := cacheBuilder.BuildHybridCache(33554432) // 32MB memory

// Automatically:
// - Checks memory first (microseconds)
// - Falls back to MongoDB (milliseconds)
// - Backfills memory on MongoDB hits
```

## Troubleshooting

### Cache not working?

1. **Check MongoDB connection**: Verify MongoDB is running
2. **Check logs**: Look for cache creation errors
3. **Check TTL index**: Run `db.cache.getIndexes()`

### Cache growing too large?

1. **Check TTLs**: Make sure they're appropriate
2. **Monitor collection size**: `db.cache.stats()`
3. **Reduce cache size**: Lower TTLs or use more aggressive eviction

### Performance issues?

1. **Use hybrid cache** for better performance
2. **Optimize MongoDB**: Ensure proper indexes and configuration
3. **Monitor hit rates**: Adjust TTLs based on hit rates

## Clean Up Cache

If you want to clear the cache:

```javascript
db.cache.drop()
```

The collection and index will be recreated automatically on next cache operation.

## Migration from Memcached/Redis

If you previously used Memcached or Redis, MongoDB cache is a drop-in replacement:

1. **Remove Memcached/Redis** from compose.yml (already done)
2. **Code stays the same** - just the implementation changes
3. **Benefits**: One less service to manage, persistent cache

## Dependencies

Add to `go.mod`:

```go
require (
    github.com/dgraph-io/ristretto v0.1.1
    go.mongodb.org/mongo-driver v1.12.1
)
```

## Summary

The caching implementation provides:

- ✅ **Simple architecture** (MongoDB + Memory only)
- ✅ **MongoDB as primary** persistent cache
- ✅ **No additional services** required
- ✅ **Easy service injection** via builder pattern
- ✅ **Built-in metrics** and health checks
- ✅ **Graceful fallbacks** for resilience
- ✅ **Simple configuration** via YAML
- ✅ **Production-ready** with proper error handling

The system is designed for maximum simplicity while providing excellent performance and persistence.
