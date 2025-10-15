package app

import (
	"go.uber.org/zap"

	"pharmacy-modernization-project-model/internal/app/builder"
	"pharmacy-modernization-project-model/internal/platform/cache"
)

func (a *App) wireCache() cache.Cache {
	// Create separate MongoDB connection for cache
	cacheMongoConnMgr, err := builder.CreateCacheMongoDBConnection(a.Cfg, a.Logger.Base)
	if err != nil {
		a.Logger.Base.Error("Failed to create cache MongoDB connection", zap.Error(err))
		// Continue without cache MongoDB - will fallback to memory cache
	}

	// Create cache builder
	cacheBuilder := builder.NewCacheBuilder(a.Cfg, a.Logger.Base)

	// Create primary cache (MongoDB or Memory)
	var primaryCache cache.Cache
	if cacheMongoConnMgr != nil {
		cacheCollection := builder.GetCacheMongoCollection(cacheMongoConnMgr, a.Cfg)
		primaryCache, err = cacheBuilder.BuildMongoDBCache(cacheCollection, "rx:")
		if err != nil {
			a.Logger.Base.Warn("Failed to create MongoDB cache, falling back to memory cache", zap.Error(err))
			primaryCache, _ = cacheBuilder.BuildMemoryCache(67108864) // 64MB fallback
		}
	} else {
		a.Logger.Base.Info("Cache MongoDB not configured, using memory cache")
		primaryCache, _ = cacheBuilder.BuildMemoryCache(67108864) // 64MB fallback
	}
	return primaryCache
}
