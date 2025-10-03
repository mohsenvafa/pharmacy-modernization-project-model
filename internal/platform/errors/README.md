# Shared Repository Error Handling

This package provides shared error handling utilities that can be used across all repository implementations in different domains.

## Usage

### Basic Error Handling

```go
import platformErrors "pharmacy-modernization-project-model/internal/platform/errors"

// Create a new repository error
err := platformErrors.NewRepositoryError(
    platformErrors.ErrorTypeNotFound,
    "User not found",
    originalError,
).WithContext("user_id", userID)

// Check error types
if repoErr, ok := err.(*platformErrors.RepositoryError); ok {
    if repoErr.IsNotFound() {
        // Handle not found error
    }
    if repoErr.IsRetryable() {
        // Handle retryable error
    }
}
```

### MongoDB-Specific Error Handling

```go
// Handle MongoDB errors
mongoErr := platformErrors.HandleMongoError("FindOne", err)
```

### Generic Database Error Handling

```go
// Handle any database errors (PostgreSQL, MySQL, etc.)
dbErr := platformErrors.HandleDatabaseError("Insert", err)
```

### Error Types

- `ErrorTypeNotFound` - Record not found
- `ErrorTypeDuplicateKey` - Duplicate key constraint violation
- `ErrorTypeTimeout` - Operation timeout
- `ErrorTypeNetworkError` - Network connectivity issues
- `ErrorTypeConnection` - Connection pool issues
- `ErrorTypeValidation` - Data validation failures
- `ErrorTypeDatabaseError` - General database errors
- `ErrorTypeUnknown` - Unknown errors

### Retryable Errors

```go
// Create a retryable error
retryErr := platformErrors.NewRetryableError(
    platformErrors.ErrorTypeTimeout,
    "Operation timed out",
    originalError,
    3, // max retries
    1000, // retry delay in milliseconds
)

// Check if should retry
if retryErr.ShouldRetry(attempt) {
    // Retry the operation
}
```

## Examples for Different Domains

### Patient Repository
```go
func (r *PatientRepository) GetByID(ctx context.Context, id string) (Patient, error) {
    patient, err := r.db.FindByID(id)
    if err != nil {
        return Patient{}, platformErrors.HandleDatabaseError("GetByID", err)
    }
    return patient, nil
}
```

### Prescription Repository
```go
func (r *PrescriptionRepository) Create(ctx context.Context, p Prescription) (Prescription, error) {
    if err := r.db.Insert(p); err != nil {
        return Prescription{}, platformErrors.HandleDatabaseError("Create", err)
    }
    return p, nil
}
```

### Address Repository
```go
func (r *AddressRepository) Update(ctx context.Context, id string, addr Address) (Address, error) {
    if err := r.db.Update(id, addr); err != nil {
        return Address{}, platformErrors.HandleDatabaseError("Update", err)
    }
    return addr, nil
}
```

## Benefits

1. **Consistency**: All repositories use the same error handling patterns
2. **Reusability**: Error handling logic is shared across domains
3. **Maintainability**: Changes to error handling affect all repositories
4. **Type Safety**: Strongly typed error types with helper methods
5. **Context**: Rich error context for better debugging
6. **Retry Logic**: Built-in support for retryable operations
