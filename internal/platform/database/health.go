package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// HealthStatus represents the health status of the database
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusUnknown   HealthStatus = "unknown"
)

// HealthCheck represents a database health check result
type HealthCheck struct {
	Status    HealthStatus           `json:"status"`
	Message   string                 `json:"message"`
	Timestamp time.Time              `json:"timestamp"`
	Duration  time.Duration          `json:"duration"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Error     error                  `json:"error,omitempty"`
}

// HealthChecker provides database health checking capabilities
type HealthChecker struct {
	client *mongo.Client
	logger *zap.Logger
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(client *mongo.Client, logger *zap.Logger) *HealthChecker {
	return &HealthChecker{
		client: client,
		logger: logger,
	}
}

// Check performs a comprehensive health check
func (hc *HealthChecker) Check(ctx context.Context) *HealthCheck {
	start := time.Now()

	check := &HealthCheck{
		Timestamp: start,
		Details:   make(map[string]interface{}),
	}

	// Basic ping check
	if err := hc.pingCheck(ctx); err != nil {
		check.Status = HealthStatusUnhealthy
		check.Message = "Database ping failed"
		check.Error = err
		check.Duration = time.Since(start)
		return check
	}

	// Server status check
	serverStatus, err := hc.serverStatusCheck(ctx)
	if err != nil {
		check.Status = HealthStatusDegraded
		check.Message = "Server status check failed"
		check.Error = err
		check.Details["server_status_error"] = err.Error()
	} else {
		check.Details["server_status"] = serverStatus
	}

	// Connection pool check
	poolStats, err := hc.connectionPoolCheck(ctx)
	if err != nil {
		check.Status = HealthStatusDegraded
		check.Message = "Connection pool check failed"
		check.Error = err
		check.Details["pool_error"] = err.Error()
	} else {
		check.Details["connection_pool"] = poolStats
	}

	// Collection access check
	collectionAccess, err := hc.collectionAccessCheck(ctx)
	if err != nil {
		check.Status = HealthStatusDegraded
		check.Message = "Collection access check failed"
		check.Error = err
		check.Details["collection_error"] = err.Error()
	} else {
		check.Details["collection_access"] = collectionAccess
	}

	// Determine final status
	if check.Status == "" {
		check.Status = HealthStatusHealthy
		check.Message = "All health checks passed"
	}

	check.Duration = time.Since(start)

	// Log health check result
	hc.logger.Debug("Database health check completed",
		zap.String("status", string(check.Status)),
		zap.Duration("duration", check.Duration),
		zap.Error(check.Error))

	return check
}

// pingCheck performs a basic ping to the database
func (hc *HealthChecker) pingCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return hc.client.Ping(ctx, nil)
}

// serverStatusCheck retrieves server status information
func (hc *HealthChecker) serverStatusCheck(ctx context.Context) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var result bson.M
	err := hc.client.Database("admin").RunCommand(ctx, bson.M{"serverStatus": 1}).Decode(&result)
	if err != nil {
		return nil, err
	}

	// Extract relevant information
	status := map[string]interface{}{
		"uptime":      result["uptime"],
		"version":     result["version"],
		"connections": result["connections"],
	}

	return status, nil
}

// connectionPoolCheck checks connection pool statistics
func (hc *HealthChecker) connectionPoolCheck(ctx context.Context) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var result bson.M
	err := hc.client.Database("admin").RunCommand(ctx, bson.M{"connPoolStats": 1}).Decode(&result)
	if err != nil {
		return nil, err
	}

	// Extract connection pool information
	pool := map[string]interface{}{
		"total_created": result["totalCreated"],
		"total_closed":  result["totalClosed"],
		"current":       result["current"],
		"available":     result["available"],
	}

	return pool, nil
}

// collectionAccessCheck verifies access to collections
func (hc *HealthChecker) collectionAccessCheck(ctx context.Context) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// List collections to verify access
	collections, err := hc.client.Database("rxintake").ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	access := map[string]interface{}{
		"collections": collections,
		"count":       len(collections),
	}

	return access, nil
}

// QuickCheck performs a quick health check (ping only)
func (hc *HealthChecker) QuickCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return hc.client.Ping(ctx, nil)
}

// IsHealthy checks if the database is healthy
func (hc *HealthChecker) IsHealthy(ctx context.Context) bool {
	return hc.QuickCheck(ctx) == nil
}

// GetHealthStatus returns the current health status
func (hc *HealthChecker) GetHealthStatus(ctx context.Context) HealthStatus {
	if hc.IsHealthy(ctx) {
		return HealthStatusHealthy
	}
	return HealthStatusUnhealthy
}

// MonitorHealth starts a health monitoring routine
func (hc *HealthChecker) MonitorHealth(ctx context.Context, interval time.Duration, callback func(*HealthCheck)) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			hc.logger.Info("Health monitoring stopped")
			return
		case <-ticker.C:
			check := hc.Check(ctx)
			if callback != nil {
				callback(check)
			}

			// Log unhealthy status
			if check.Status != HealthStatusHealthy {
				hc.logger.Warn("Database health check failed",
					zap.String("status", string(check.Status)),
					zap.String("message", check.Message),
					zap.Error(check.Error))
			}
		}
	}
}
