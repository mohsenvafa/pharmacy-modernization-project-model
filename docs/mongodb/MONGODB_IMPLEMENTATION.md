# MongoDB Implementation Guide

This document describes the MongoDB implementation for the patient repository with best practices for efficiency, scalability, and connection management.

## Overview

The MongoDB implementation provides:
- **Connection Management**: Proper connection pooling and lifecycle management
- **Performance Optimization**: Optimized queries with proper indexing
- **Error Handling**: Comprehensive error handling and logging
- **Health Monitoring**: Database health checks and metrics
- **Scalability**: Support for horizontal scaling and read replicas

## Architecture

### Components

1. **Connection Manager** (`internal/platform/database/mongodb.go`)
   - Manages MongoDB connections with proper pooling
   - Handles connection lifecycle and health checks
   - Provides collection access

2. **Patient Repository** (`domain/patient/repository/patient_mongodb_repository.go`)
   - Implements the `PatientRepository` interface
   - Optimized queries with proper error handling
   - Comprehensive logging and metrics

3. **Error Handling** (`internal/platform/database/errors.go`)
   - Structured error types and handling
   - Retry logic for transient failures
   - Context-aware error reporting

4. **Health Monitoring** (`internal/platform/database/health.go`)
   - Database health checks
   - Connection pool monitoring
   - Performance metrics

5. **Metrics Collection** (`internal/platform/database/metrics.go`)
   - Operation metrics and performance tracking
   - Error rate monitoring
   - Connection pool statistics

## Configuration

### YAML Configuration

```yaml
database:
  mongodb:
    uri: "mongodb://localhost:27017"
    database: "pharmacy_modernization"
    collections:
      patients: "patients"
    connection:
      max_pool_size: 100
      min_pool_size: 5
      max_idle_time: "30m"
      connect_timeout: "10s"
      socket_timeout: "30s"
    options:
      retry_writes: true
      retry_reads: true
```

### Environment Variables

You can override configuration using environment variables:

```bash
export RX_DATABASE_MONGODB_URI="mongodb://localhost:27017"
export PM_DATABASE_MONGODB_DATABASE="pharmacy_modernization"
export RX_DATABASE_MONGODB_CONNECTION_MAX_POOL_SIZE=100
```

## Setup and Installation

### 1. Install MongoDB

```bash
# macOS
brew install mongodb-community

# Ubuntu/Debian
sudo apt-get install mongodb

# Podman (Recommended)
podman run -d -p 27017:27017 --name mongodb mongo:latest

# Docker (Alternative)
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

### 2. Setup Indexes

Run the index setup script:

```bash
# Using MongoDB shell
mongo scripts/setup_mongodb_indexes.js

# Using Go setup command
go run cmd/setup_mongodb/main.go
```

### 3. Start the Application

```bash
go run cmd/server/main.go
```

## Best Practices Implemented

### 1. Connection Management

- **Connection Pooling**: Configured with optimal pool size
- **Connection Lifecycle**: Proper establishment, health checks, and cleanup
- **Retry Logic**: Built-in retry mechanisms for transient failures
- **Timeout Management**: Configurable timeouts for different operations

### 2. Performance Optimization

#### Indexing Strategy

```javascript
// Text search index
db.patients.createIndex({ "name": "text" })

// State filtering index
db.patients.createIndex({ "state": 1 })

// Created date sorting index
db.patients.createIndex({ "created_at": -1 })

// Unique phone index
db.patients.createIndex({ "phone": 1 }, { unique: true })

// Unique ID index
db.patients.createIndex({ "id": 1 }, { unique: true })

// Compound indexes for complex queries
db.patients.createIndex({ "state": 1, "name": 1 })
db.patients.createIndex({ "created_at": -1, "state": 1 })
```

#### Query Optimization

- **Projection**: Limit returned fields when possible
- **Pagination**: Efficient skip/limit implementation
- **Sorting**: Proper index usage for sorting
- **Filtering**: Optimized filter conditions

### 3. Error Handling

#### Error Types

- `NotFound`: Document not found
- `DuplicateKey`: Duplicate key constraint violation
- `Timeout`: Operation timeout
- `NetworkError`: Network connectivity issues
- `ConnectionError`: Connection pool issues
- `ValidationError`: Data validation failures

#### Retry Logic

```go
// Automatic retry for transient errors
if err.IsRetryable() {
    // Implement exponential backoff
    time.Sleep(time.Duration(attempt) * time.Second)
    // Retry operation
}
```

### 4. Monitoring and Observability

#### Health Checks

```go
// Basic health check
if err := connMgr.HealthCheck(ctx); err != nil {
    log.Error("Database health check failed", err)
}

// Comprehensive health check
healthCheck := healthChecker.Check(ctx)
if healthCheck.Status != HealthStatusHealthy {
    // Handle unhealthy database
}
```

#### Metrics Collection

- **Operation Metrics**: Count, duration, error rates
- **Connection Metrics**: Pool statistics, active connections
- **Performance Metrics**: Average response times, throughput

## Usage Examples

### Basic Operations

```go
// Create repository
repo := NewPatientMongoRepository(collection, logger)

// List patients with pagination
patients, err := repo.List(ctx, "search query", 10, 0)

// Get patient by ID
patient, err := repo.GetByID(ctx, "P001")

// Create new patient
newPatient, err := repo.Create(ctx, patient)

// Update patient
updatedPatient, err := repo.Update(ctx, "P001", patient)
```

### Advanced Operations

```go
// Bulk insert for better performance
patients := []m.Patient{...}
err := repo.BulkInsert(ctx, patients)

// Find by state with pagination
patients, err := repo.FindByState(ctx, "California", 10, 0)

// Health check
if err := repo.HealthCheck(ctx); err != nil {
    // Handle health check failure
}
```

## Migration from Memory Repository

The implementation supports gradual migration:

1. **Fallback Support**: Falls back to memory repository if MongoDB is unavailable
2. **Feature Flags**: Can be toggled via configuration
3. **Data Migration**: Scripts available for data migration

### Migration Steps

1. **Setup MongoDB**: Install and configure MongoDB
2. **Update Configuration**: Add MongoDB configuration
3. **Deploy Application**: Deploy with MongoDB support
4. **Migrate Data**: Run data migration scripts
5. **Verify**: Test all operations

## Performance Considerations

### Connection Pool Sizing

- **Max Pool Size**: 100 connections (configurable)
- **Min Pool Size**: 5 connections (configurable)
- **Idle Timeout**: 30 minutes (configurable)

### Query Performance

- **Index Usage**: All queries use appropriate indexes
- **Pagination**: Efficient skip/limit implementation
- **Projection**: Minimal data transfer
- **Caching**: Optional Redis integration

### Monitoring

- **Query Performance**: Track operation durations
- **Error Rates**: Monitor error rates and types
- **Connection Pool**: Monitor pool utilization
- **Health Status**: Regular health checks

## Security Considerations

### Connection Security

- **TLS/SSL**: Encrypted connections
- **Authentication**: MongoDB authentication
- **Authorization**: Role-based access control

### Data Security

- **Input Validation**: Proper validation of user inputs
- **Query Injection Prevention**: Parameterized queries
- **Access Control**: Database-level permissions

## Troubleshooting

### Common Issues

1. **Connection Timeouts**
   - Check network connectivity
   - Verify MongoDB server status
   - Adjust timeout configurations

2. **Performance Issues**
   - Check index usage with `explain()`
   - Monitor connection pool utilization
   - Review query patterns

3. **Error Handling**
   - Check error logs for specific error types
   - Verify data validation
   - Review connection pool settings

### Debugging

```go
// Enable debug logging
logger := zap.NewDevelopment()

// Check connection status
if err := connMgr.Ping(ctx); err != nil {
    log.Error("MongoDB ping failed", err)
}

// Monitor metrics
metrics := metricsCollector.GetMetrics()
log.Info("Database metrics", metrics)
```

## Future Enhancements

1. **Read Replicas**: Support for reading from secondary nodes
2. **Sharding**: Horizontal scaling support
3. **Caching**: Redis integration for frequently accessed data
4. **Analytics**: Advanced query analytics and optimization
5. **Backup**: Automated backup and recovery procedures

## Conclusion

This MongoDB implementation provides a robust, scalable, and efficient solution for patient data management with comprehensive error handling, monitoring, and performance optimization. The architecture supports both current needs and future growth while maintaining high availability and performance standards.
