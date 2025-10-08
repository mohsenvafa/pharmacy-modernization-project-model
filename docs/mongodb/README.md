# MongoDB Documentation

Documentation for MongoDB implementation and data access patterns.

## 📚 Documentation Files

- **[MongoDB Implementation](./MONGODB_IMPLEMENTATION.md)** - Complete MongoDB setup and patterns

## 🗄️ Database Overview

### Connection Management
- Singleton connection pool
- Automatic reconnection
- Context-aware operations
- Graceful shutdown

### Database Structure
```
rxintake (database)
├── patients (collection)
├── prescriptions (collection)
└── [other domain collections]
```

## 🔑 Key Patterns

### Repository Pattern
Each domain has its own repository for data access:

```go
type PatientRepository interface {
    GetByID(ctx context.Context, id string) (*models.Patient, error)
    GetAll(ctx context.Context) ([]*models.Patient, error)
    Create(ctx context.Context, patient *models.Patient) error
    Update(ctx context.Context, patient *models.Patient) error
    Delete(ctx context.Context, id string) error
}
```

### Context-Aware Operations
All database operations accept `context.Context`:
- Request tracing
- Timeout management
- Cancellation propagation
- User context passing

### Error Handling
- Proper MongoDB error checking
- User-friendly error messages
- Logging of database errors
- Transaction support

## 🚀 Quick Examples

### Connection Setup
```go
db := database.GetDB()
collection := db.Collection("patients")
```

### Basic CRUD Operations

#### Create
```go
result, err := collection.InsertOne(ctx, patient)
```

#### Read
```go
err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&patient)
```

#### Update
```go
filter := bson.M{"_id": id}
update := bson.M{"$set": bson.M{"firstName": "John"}}
_, err := collection.UpdateOne(ctx, filter, update)
```

#### Delete
```go
_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
```

### Query Patterns

#### Find with Filter
```go
filter := bson.M{"lastName": "Smith"}
cursor, err := collection.Find(ctx, filter)
```

#### Projection
```go
opts := options.Find().SetProjection(bson.M{
    "firstName": 1,
    "lastName": 1,
})
cursor, err := collection.Find(ctx, filter, opts)
```

#### Sorting
```go
opts := options.Find().SetSort(bson.M{"lastName": 1})
cursor, err := collection.Find(ctx, filter, opts)
```

#### Pagination
```go
opts := options.Find().
    SetLimit(10).
    SetSkip(20)
cursor, err := collection.Find(ctx, filter, opts)
```

## 🔧 Configuration

MongoDB is configured in `internal/configs/app.yaml`:
```yaml
database:
  mongodb:
    uri: "${RX_MONGODB_URI}"
    database: "rxintake"
    timeout: 30
```

Environment variables:
```bash
RX_MONGODB_URI=mongodb://localhost:27017
```

## 📊 Data Models

### Best Practices
1. **Use BSON tags** - Map Go fields to MongoDB fields
2. **ObjectID for IDs** - Use MongoDB's native ID type
3. **Embedded vs Referenced** - Choose based on access patterns
4. **Indexes** - Create indexes for frequently queried fields

### Example Model
```go
type Patient struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    FirstName   string            `bson:"firstName"`
    LastName    string            `bson:"lastName"`
    DateOfBirth time.Time         `bson:"dateOfBirth"`
    CreatedAt   time.Time         `bson:"createdAt"`
    UpdatedAt   time.Time         `bson:"updatedAt"`
}
```

## 🧪 Testing

### Test Database
Use a separate database for tests:
```go
testDB := "rxintake_test"
```

### Cleanup
Always clean up test data:
```go
defer collection.Drop(ctx)
```

## 📖 Related Documentation

- [Architecture Overview](../architecture/ARCHITECTURE.md)
- [GraphQL Implementation](../graphql/GRAPHQL_IMPLEMENTATION.md)

## 🛠️ Tools

- **MongoDB Compass** - GUI for MongoDB
- **mongosh** - MongoDB shell
- **mongo-go-driver** - Official Go driver

## 🔍 Troubleshooting

### Common Issues

**Connection Timeout**
```
Solution: Check MongoDB URI and network connectivity
```

**Authentication Failed**
```
Solution: Verify username/password in connection string
```

**Document Not Found**
```
Solution: Check if ID exists and query filter is correct
```

