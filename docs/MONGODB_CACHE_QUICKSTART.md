# MongoDB Cache Quick Start

If you already have MongoDB running and want to use it for caching without adding Memcached or Redis, here's how:

## Why Use MongoDB for Caching?

**Pros:**
- ✅ No additional service to run
- ✅ Persistent cache (survives restarts)
- ✅ Automatic TTL expiration
- ✅ Simple setup
- ✅ Good for development/testing

**Cons:**
- ⚠️ Slower than Memcached/Redis (but still fast!)
- ⚠️ Uses MongoDB connection pool
- ⚠️ Not ideal for very high traffic

## Quick Setup

### 1. Switch to MongoDB Cache

In `internal/app/wire.go`, comment out Memcached and use MongoDB:

```go
// Create cache builder
cacheBuilder := NewCacheBuilder(a.Cfg, logger.Base)

// Use MongoDB cache instead of Memcached
cacheCollection := GetCacheCollection(mongoConnMgr)
primaryCache, err := cacheBuilder.BuildMongoDBCache(cacheCollection, "rx:")
if err != nil {
    logger.Base.Warn("Failed to create MongoDB cache, falling back to memory cache", zap.Error(err))
    primaryCache, _ = cacheBuilder.BuildMemoryCache(67108864) // 64MB fallback
}
```

### 2. That's it!

No configuration changes needed. The cache will use your existing MongoDB connection.

## How It Works

The MongoDB cache:

1. **Creates a `cache` collection** in your database
2. **Automatically creates TTL index** for expiration
3. **Stores cache entries** with automatic cleanup:

```javascript
// MongoDB document structure
{
  "_id": "rx:patient:id:123",      // Cache key
  "value": BinData(...),            // Cached data
  "expireAt": ISODate("2024-...")   // Auto-delete time
}
```

## Verify It's Working

### Check the cache collection:

```bash
mongosh mongodb://admin:admin123@localhost:27017

use pharmacy_modernization
db.cache.find()
db.cache.stats()
```

### Check TTL index:

```javascript
db.cache.getIndexes()
// Should see: { "expireAt": 1 } with expireAfterSeconds: 0
```

### Monitor cache hits in logs:

```bash
# Look for cache debug logs
tail -f logs/app.log | grep "cache"
```

## Performance Comparison

With MongoDB cache vs Memcached:

| Operation | MongoDB Cache | Memcached |
|-----------|---------------|-----------|
| Get (cold) | ~5-10ms | ~1-2ms |
| Get (hot) | ~2-5ms | ~0.5-1ms |
| Set | ~5-10ms | ~1-2ms |
| Persistence | ✅ Yes | ❌ No |

**MongoDB is plenty fast for most applications!**

## When to Switch Back to Memcached

Switch to Memcached when:
- You have >1000 requests/second
- You need sub-millisecond latency
- You're running multiple app instances
- Cache hit rate is >90%

## Clean Up Cache Collection

If you want to clear the cache:

```javascript
db.cache.drop()
```

The collection and index will be recreated automatically on next cache operation.

## Switching to Hybrid (Memory + MongoDB)

For even better performance, use hybrid cache:

```go
// In wire.go - use hybrid with MongoDB as shared tier
primaryCache, err := cacheBuilder.BuildHybridCache(33554432) // 32MB memory
```

This gives you:
- ⚡ Ultra-fast in-memory L1 cache
- 💾 Persistent MongoDB L2 cache
- 🎯 Best of both worlds!

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

1. **Consider switching to Memcached** for high traffic
2. **Use hybrid cache** for better performance
3. **Optimize MongoDB**: Ensure proper indexes and configuration

## Summary

MongoDB cache is perfect for:
- 🚀 Quick development/testing
- 📦 Small to medium applications  
- 💾 When persistence is valuable
- 🎯 Minimizing infrastructure complexity

Just uncomment 2 lines in `wire.go` and you're done!

