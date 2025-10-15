# Cache MongoDB Separation Guide

## Overview

The caching system uses a **completely separate MongoDB connection** from your main database. This provides flexibility, independence, and better control over your caching infrastructure.

## Architecture

```
Application
├── Main MongoDB Connection
│   ├── URI: mongodb://admin:admin123@localhost:27017
│   ├── Database: pharmacy_modernization
│   ├── Collections: patients, prescriptions, addresses
│   └── Purpose: Primary data storage
│
└── Cache MongoDB Connection (Independent)
    ├── URI: mongodb://admin:admin123@localhost:27017 (can be different!)
    ├── Database: pharmacy_modernization_cache (separate database)
    ├── Collection: cache
    └── Purpose: Distributed caching
```

## Configuration

### Same Server, Different Databases

```yaml
# Main database
database:
  mongodb:
    uri: "mongodb://admin:admin123@localhost:27017"
    database: "pharmacy_modernization"
    collections:
      patients: "patients"
      prescriptions: "prescriptions"

# Cache database (same server, different database)
cache:
  mongodb:
    uri: "mongodb://admin:admin123@localhost:27017"  # Same server
    database: "pharmacy_modernization_cache"         # Different database
    collection: "cache"
```

### Different MongoDB Clusters

```yaml
# Main database (production cluster)
database:
  mongodb:
    uri: "mongodb://prod-cluster:27017"
    database: "pharmacy_modernization"

# Cache database (cache cluster)
cache:
  mongodb:
    uri: "mongodb://cache-cluster:27017"  # Different cluster!
    database: "cache_db"
    collection: "cache"
```

### Shared Connection, Different Pool Sizes

```yaml
# Main database (larger pool for data operations)
database:
  mongodb:
    connection:
      max_pool_size: 100
      min_pool_size: 10

# Cache database (smaller pool for cache operations)
cache:
  mongodb:
    connection:
      max_pool_size: 50   # Smaller pool
      min_pool_size: 5
```

## Benefits of Separation

### 1. **Independent Scaling**
```yaml
# Scale cache MongoDB independently
cache:
  mongodb:
    uri: "mongodb://cache-replica-set:27017"
    connection:
      max_pool_size: 200  # High throughput for cache
```

### 2. **Different Retention Policies**
```bash
# Drop cache database without affecting data
mongosh
use pharmacy_modernization_cache
db.dropDatabase()  # ✅ Safe - only cache data lost
```

### 3. **Separate Backups**
```bash
# Backup main data frequently
mongodump --db pharmacy_modernization --out /backups/daily/

# Backup cache less frequently (or not at all)
# Cache can be rebuilt from main data
```

### 4. **Independent Monitoring**
```javascript
// Monitor main database
use pharmacy_modernization
db.stats()

// Monitor cache database separately
use pharmacy_modernization_cache
db.stats()
```

### 5. **Different Write Concerns**
```go
// Main database: Strong consistency
mainDB.SetWriteConcern(writeconcern.Majority())

// Cache database: Fast writes, eventual consistency
cacheDB.SetWriteConcern(writeconcern.W(1))
```

## Use Cases

### Development: Use Same Server

```yaml
cache:
  mongodb:
    uri: "mongodb://admin:admin123@localhost:27017"
    database: "pharmacy_modernization_cache"
```

### Production: Use Separate Cluster

```yaml
cache:
  mongodb:
    uri: "mongodb://cache-cluster.internal:27017"
    database: "cache"
    connection:
      max_pool_size: 200  # High throughput for cache
```

### High Availability: Use MongoDB Atlas

```yaml
cache:
  mongodb:
    uri: "mongodb+srv://user:pass@cache-cluster.mongodb.net"
    database: "cache"
```

## Migration Strategies

### Start: Share Main MongoDB

```yaml
# Use main MongoDB for cache initially
cache:
  mongodb:
    uri: "mongodb://admin:admin123@localhost:27017"  # Same
    database: "pharmacy_modernization_cache"          # Separate DB
```

### Later: Move to Dedicated Cache Cluster

```yaml
# Move cache to dedicated cluster when traffic grows
cache:
  mongodb:
    uri: "mongodb://cache-dedicated:27017"  # Dedicated cluster
    database: "cache"
```

## Database Operations

### View All Databases

```bash
mongosh mongodb://admin:admin123@localhost:27017

show dbs
# pharmacy_modernization
# pharmacy_modernization_cache  ← Cache database
```

### Check Cache Database

```javascript
use pharmacy_modernization_cache
db.cache.find().limit(10)
db.cache.stats()
db.cache.getIndexes()
```

### Clear Cache Database

```javascript
use pharmacy_modernization_cache
db.cache.drop()
// Cache will be recreated automatically with TTL index
```

### Monitor Connections

```javascript
db.currentOp()
db.serverStatus().connections
```

## Connection Management

Both connections are managed independently:

```go
// Main MongoDB connection
mongoConnMgr, err := CreateMongoDBConnection(config, logger)
// Connections: 100 max, 10 min

// Cache MongoDB connection (independent)
cacheMongoConnMgr, err := CreateCacheMongoDBConnection(config, logger)
// Connections: 50 max, 5 min
```

## Best Practices

### 1. Use Same URI, Different Database (Development)
```yaml
cache:
  mongodb:
    uri: "mongodb://admin:admin123@localhost:27017"  # Same server
    database: "pharmacy_modernization_cache"         # Different database
```

### 2. Use Separate Cluster (Production)
```yaml
cache:
  mongodb:
    uri: "mongodb://cache-cluster:27017"  # Different cluster
    database: "cache"
```

### 3. Configure Appropriate Pool Sizes
- Main DB: Larger pools (100+)
- Cache DB: Smaller pools (50+)

### 4. Monitor Both Separately
- Track cache hit rates
- Monitor cache database size
- Alert on main database issues separately from cache issues

## Troubleshooting

### Cache database not created?

The database is created automatically on first cache operation. To verify:

```bash
mongosh mongodb://admin:admin123@localhost:27017
show dbs
# Should see: pharmacy_modernization_cache
```

### Want to disable cache MongoDB?

Simply remove or comment out the cache MongoDB config:

```yaml
cache:
  # mongodb:  # Commented out
  #   uri: "..."
  memory:
    max_cost: 67108864
```

The system will automatically fallback to memory cache.

## Summary

✅ **Complete separation** between main DB and cache DB  
✅ **Flexible configuration** - same server or different clusters  
✅ **Independent scaling** - tune each connection pool separately  
✅ **Different retention** - drop cache without affecting data  
✅ **Production ready** - supports any MongoDB topology  

The cache MongoDB is a first-class citizen with its own configuration!

