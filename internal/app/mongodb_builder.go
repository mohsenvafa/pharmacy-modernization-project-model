package app

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"pharmacy-modernization-project-model/internal/platform/config"
	"pharmacy-modernization-project-model/internal/platform/database"
)

// CreateMongoDBConnection creates a MongoDB connection manager based on configuration
func CreateMongoDBConnection(cfg *config.Config, logger *zap.Logger) *database.ConnectionManager {
	if cfg.Database.MongoDB.URI == "" {
		return nil
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
		logger.Error("Failed to initialize MongoDB connection", zap.Error(err))
		return nil
	}

	logger.Info("MongoDB connection established successfully")
	return connMgr
}

// GetPatientsCollection returns the patients collection from MongoDB connection manager
func GetPatientsCollection(mongoConnMgr *database.ConnectionManager) *mongo.Collection {
	if mongoConnMgr == nil {
		return nil
	}
	return mongoConnMgr.GetCollection("patients")
}
