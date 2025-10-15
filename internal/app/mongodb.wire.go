package app

import (
	"go.uber.org/zap"

	"pharmacy-modernization-project-model/internal/app/builder"
	"pharmacy-modernization-project-model/internal/platform/database"
)

func (a *App) wireMongodb() *database.ConnectionManager {
	// Create main MongoDB connection
	mongoConnMgr, err := builder.CreateMongoDBConnection(a.Cfg, a.Logger.Base)
	if err != nil {
		a.Logger.Base.Error("Failed to create MongoDB connection", zap.Error(err))
		// Continue without MongoDB - will use memory repository as fallback
	}

	return mongoConnMgr
}
