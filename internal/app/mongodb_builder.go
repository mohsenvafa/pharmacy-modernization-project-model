package app

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"pharmacy-modernization-project-model/internal/platform/config"
	"pharmacy-modernization-project-model/internal/platform/database"
	platformErrors "pharmacy-modernization-project-model/internal/platform/errors"
)

// CreateMongoDBConnection creates a MongoDB connection manager based on configuration
func CreateMongoDBConnection(cfg *config.Config, logger *zap.Logger) (*database.ConnectionManager, error) {

	return nil, nil
	// Validate configuration
	if cfg.Database.MongoDB.URI == "" {
		logger.Warn("MongoDB URI not configured, skipping MongoDB connection")
		return nil, platformErrors.NewConfigurationError("database", "mongodb.uri", "MongoDB URI is not configured")
	}

	if cfg.Database.MongoDB.Database == "" {
		logger.Error("MongoDB database name is required")
		return nil, platformErrors.NewConfigurationError("database", "mongodb.database", "MongoDB database name is required")
	}

	mongoConfig := database.MongoDBConfig{
		URI:      cfg.Database.MongoDB.URI,
		Database: cfg.Database.MongoDB.Database,
		Collections: map[string]string{
			"patients": cfg.Database.MongoDB.Collections.Patients,
		},
		Connection: database.ConnectionConfig{
			MaxPoolSize:    cfg.Database.MongoDB.Connection.MaxPoolSize,
			MinPoolSize:    cfg.Database.MongoDB.Connection.MinPoolSize,
			MaxIdleTime:    cfg.Database.MongoDB.Connection.MaxIdleTime,
			ConnectTimeout: cfg.Database.MongoDB.Connection.ConnectTimeout,
			SocketTimeout:  cfg.Database.MongoDB.Connection.SocketTimeout,
		},
		Options: database.OptionsConfig{
			RetryWrites: cfg.Database.MongoDB.Options.RetryWrites,
			RetryReads:  cfg.Database.MongoDB.Options.RetryReads,
		},
	}

	connMgr, err := database.NewConnectionManager(mongoConfig, logger)
	if err != nil {
		logger.Error("Failed to initialize MongoDB connection",
			zap.Error(err),
			zap.String("uri", cfg.Database.MongoDB.URI),
			zap.String("database", cfg.Database.MongoDB.Database))
		return nil, platformErrors.NewConfigurationError("database", "mongodb.connection", "Failed to initialize MongoDB connection")
	}

	logger.Info("MongoDB connection established successfully",
		zap.String("database", cfg.Database.MongoDB.Database))
	return connMgr, nil

}

// GetPatientsCollection returns the patients collection from MongoDB connection manager
func GetPatientsCollection(mongoConnMgr *database.ConnectionManager) *mongo.Collection {
	if mongoConnMgr == nil {
		return nil
	}
	return mongoConnMgr.GetCollection("patients")
}
