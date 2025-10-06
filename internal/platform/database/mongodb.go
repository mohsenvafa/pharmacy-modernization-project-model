package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// MongoDBConfig represents MongoDB configuration
type MongoDBConfig struct {
	URI         string
	Database    string
	Collections map[string]string
	Connection  ConnectionConfig
	Options     OptionsConfig
}

// ConnectionConfig represents connection pool configuration
type ConnectionConfig struct {
	MaxPoolSize    uint64
	MinPoolSize    uint64
	MaxIdleTime    string
	ConnectTimeout string
	SocketTimeout  string
}

// OptionsConfig represents MongoDB client options
type OptionsConfig struct {
	RetryWrites bool
	RetryReads  bool
}

// ConnectionManager manages MongoDB connections with proper lifecycle
type ConnectionManager struct {
	client   *mongo.Client
	database *mongo.Database
	config   MongoDBConfig
	logger   *zap.Logger
}

// NewConnectionManager creates a new MongoDB connection manager
func NewConnectionManager(config MongoDBConfig, logger *zap.Logger) (*ConnectionManager, error) {
	cm := &ConnectionManager{
		config: config,
		logger: logger,
	}

	if err := cm.connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return cm, nil
}

// connect establishes connection to MongoDB with proper configuration
func (cm *ConnectionManager) connect() error {
	// Parse timeouts
	maxIdleTime, err := time.ParseDuration(cm.config.Connection.MaxIdleTime)
	if err != nil {
		return fmt.Errorf("invalid max_idle_time: %w", err)
	}

	connectTimeout, err := time.ParseDuration(cm.config.Connection.ConnectTimeout)
	if err != nil {
		return fmt.Errorf("invalid connect_timeout: %w", err)
	}

	socketTimeout, err := time.ParseDuration(cm.config.Connection.SocketTimeout)
	if err != nil {
		return fmt.Errorf("invalid socket_timeout: %w", err)
	}

	// Configure client options
	clientOptions := options.Client().
		ApplyURI(cm.config.URI).
		SetMaxPoolSize(cm.config.Connection.MaxPoolSize).
		SetMinPoolSize(cm.config.Connection.MinPoolSize).
		SetMaxConnIdleTime(maxIdleTime).
		SetConnectTimeout(connectTimeout).
		SetSocketTimeout(socketTimeout).
		SetRetryWrites(cm.config.Options.RetryWrites).
		SetRetryReads(cm.config.Options.RetryReads)

	// Create client
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	cm.client = client
	cm.database = client.Database(cm.config.Database)

	// Test connection
	if err := cm.Ping(ctx); err != nil {
		cm.Close()
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	cm.logger.Info("Successfully connected to MongoDB",
		zap.String("database", cm.config.Database),
		zap.Uint64("max_pool_size", cm.config.Connection.MaxPoolSize),
		zap.Uint64("min_pool_size", cm.config.Connection.MinPoolSize))

	return nil
}

// GetCollection returns a MongoDB collection by name
func (cm *ConnectionManager) GetCollection(name string) *mongo.Collection {
	collectionName := cm.config.Collections[name]
	if collectionName == "" {
		collectionName = name // fallback to name if not configured
	}
	return cm.database.Collection(collectionName)
}

// Ping tests the connection to MongoDB
func (cm *ConnectionManager) Ping(ctx context.Context) error {
	if cm.client == nil {
		return fmt.Errorf("MongoDB client is not initialized")
	}
	return cm.client.Ping(ctx, nil)
}

// Close closes the MongoDB connection
func (cm *ConnectionManager) Close() error {
	if cm.client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := cm.client.Disconnect(ctx); err != nil {
		cm.logger.Error("Failed to disconnect from MongoDB", zap.Error(err))
		return err
	}

	cm.logger.Info("MongoDB connection closed")
	return nil
}

// GetClient returns the MongoDB client (for advanced operations)
func (cm *ConnectionManager) GetClient() *mongo.Client {
	return cm.client
}

// GetDatabase returns the MongoDB database
func (cm *ConnectionManager) GetDatabase() *mongo.Database {
	return cm.database
}

// HealthCheck performs a comprehensive health check
func (cm *ConnectionManager) HealthCheck(ctx context.Context) error {
	if cm.client == nil {
		return fmt.Errorf("MongoDB client is not initialized")
	}

	// Ping the database
	if err := cm.client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("MongoDB ping failed: %w", err)
	}

	// Check if we can access the database
	if err := cm.database.RunCommand(ctx, map[string]interface{}{"ping": 1}).Err(); err != nil {
		return fmt.Errorf("MongoDB database access failed: %w", err)
	}

	return nil
}
