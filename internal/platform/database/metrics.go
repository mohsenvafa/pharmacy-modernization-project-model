package database

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// Metrics represents database operation metrics
type Metrics struct {
	Operations  map[string]*OperationMetrics `json:"operations"`
	Connections ConnectionMetrics            `json:"connections"`
	LastUpdated time.Time                    `json:"last_updated"`
	mu          sync.RWMutex
}

// OperationMetrics represents metrics for a specific operation
type OperationMetrics struct {
	Count         int64         `json:"count"`
	TotalDuration time.Duration `json:"total_duration"`
	AvgDuration   time.Duration `json:"avg_duration"`
	MinDuration   time.Duration `json:"min_duration"`
	MaxDuration   time.Duration `json:"max_duration"`
	Errors        int64         `json:"errors"`
	LastOperation time.Time     `json:"last_operation"`
}

// ConnectionMetrics represents connection pool metrics
type ConnectionMetrics struct {
	Active    int `json:"active"`
	Idle      int `json:"idle"`
	Total     int `json:"total"`
	Available int `json:"available"`
}

// MetricsCollector collects and manages database metrics
type MetricsCollector struct {
	metrics *Metrics
	logger  *zap.Logger
	client  *mongo.Client
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(client *mongo.Client, logger *zap.Logger) *MetricsCollector {
	return &MetricsCollector{
		metrics: &Metrics{
			Operations:  make(map[string]*OperationMetrics),
			Connections: ConnectionMetrics{},
			LastUpdated: time.Now(),
		},
		logger: logger,
		client: client,
	}
}

// RecordOperation records metrics for a database operation
func (mc *MetricsCollector) RecordOperation(operation string, duration time.Duration, err error) {
	mc.metrics.mu.Lock()
	defer mc.metrics.mu.Unlock()

	op, exists := mc.metrics.Operations[operation]
	if !exists {
		op = &OperationMetrics{
			MinDuration: duration,
			MaxDuration: duration,
		}
		mc.metrics.Operations[operation] = op
	}

	// Update metrics
	op.Count++
	op.TotalDuration += duration
	op.AvgDuration = op.TotalDuration / time.Duration(op.Count)
	op.LastOperation = time.Now()

	if duration < op.MinDuration {
		op.MinDuration = duration
	}
	if duration > op.MaxDuration {
		op.MaxDuration = duration
	}

	if err != nil {
		op.Errors++
	}

	mc.metrics.LastUpdated = time.Now()
}

// GetMetrics returns the current metrics
func (mc *MetricsCollector) GetMetrics() *Metrics {
	mc.metrics.mu.RLock()
	defer mc.metrics.mu.RUnlock()

	// Create a copy to avoid race conditions
	metricsCopy := &Metrics{
		Operations:  make(map[string]*OperationMetrics),
		Connections: mc.metrics.Connections,
		LastUpdated: mc.metrics.LastUpdated,
	}

	for op, metrics := range mc.metrics.Operations {
		metricsCopy.Operations[op] = &OperationMetrics{
			Count:         metrics.Count,
			TotalDuration: metrics.TotalDuration,
			AvgDuration:   metrics.AvgDuration,
			MinDuration:   metrics.MinDuration,
			MaxDuration:   metrics.MaxDuration,
			Errors:        metrics.Errors,
			LastOperation: metrics.LastOperation,
		}
	}

	return metricsCopy
}

// GetOperationMetrics returns metrics for a specific operation
func (mc *MetricsCollector) GetOperationMetrics(operation string) *OperationMetrics {
	mc.metrics.mu.RLock()
	defer mc.metrics.mu.RUnlock()

	op, exists := mc.metrics.Operations[operation]
	if !exists {
		return &OperationMetrics{}
	}

	// Return a copy
	return &OperationMetrics{
		Count:         op.Count,
		TotalDuration: op.TotalDuration,
		AvgDuration:   op.AvgDuration,
		MinDuration:   op.MinDuration,
		MaxDuration:   op.MaxDuration,
		Errors:        op.Errors,
		LastOperation: op.LastOperation,
	}
}

// GetConnectionMetrics returns current connection metrics
func (mc *MetricsCollector) GetConnectionMetrics(ctx context.Context) ConnectionMetrics {
	// This would typically query the MongoDB server for connection pool stats
	// For now, return basic metrics
	return ConnectionMetrics{
		Active:    0, // Would be populated from server stats
		Idle:      0,
		Total:     0,
		Available: 0,
	}
}

// ResetMetrics resets all metrics
func (mc *MetricsCollector) ResetMetrics() {
	mc.metrics.mu.Lock()
	defer mc.metrics.mu.Unlock()

	mc.metrics.Operations = make(map[string]*OperationMetrics)
	mc.metrics.Connections = ConnectionMetrics{}
	mc.metrics.LastUpdated = time.Now()

	mc.logger.Info("Database metrics reset")
}

// GetTotalOperations returns the total number of operations
func (mc *MetricsCollector) GetTotalOperations() int64 {
	mc.metrics.mu.RLock()
	defer mc.metrics.mu.RUnlock()

	var total int64
	for _, op := range mc.metrics.Operations {
		total += op.Count
	}
	return total
}

// GetTotalErrors returns the total number of errors
func (mc *MetricsCollector) GetTotalErrors() int64 {
	mc.metrics.mu.RLock()
	defer mc.metrics.mu.RUnlock()

	var total int64
	for _, op := range mc.metrics.Operations {
		total += op.Errors
	}
	return total
}

// GetErrorRate returns the error rate as a percentage
func (mc *MetricsCollector) GetErrorRate() float64 {
	totalOps := mc.GetTotalOperations()
	if totalOps == 0 {
		return 0
	}

	totalErrors := mc.GetTotalErrors()
	return float64(totalErrors) / float64(totalOps) * 100
}

// GetAverageOperationDuration returns the average duration across all operations
func (mc *MetricsCollector) GetAverageOperationDuration() time.Duration {
	mc.metrics.mu.RLock()
	defer mc.metrics.mu.RUnlock()

	var totalDuration time.Duration
	var totalCount int64

	for _, op := range mc.metrics.Operations {
		totalDuration += op.TotalDuration
		totalCount += op.Count
	}

	if totalCount == 0 {
		return 0
	}

	return totalDuration / time.Duration(totalCount)
}

// LogMetrics logs the current metrics
func (mc *MetricsCollector) LogMetrics() {
	metrics := mc.GetMetrics()

	mc.logger.Info("Database metrics summary",
		zap.Int64("total_operations", mc.GetTotalOperations()),
		zap.Int64("total_errors", mc.GetTotalErrors()),
		zap.Float64("error_rate_percent", mc.GetErrorRate()),
		zap.Duration("avg_operation_duration", mc.GetAverageOperationDuration()),
		zap.Time("last_updated", metrics.LastUpdated))

	// Log per-operation metrics
	for operation, opMetrics := range metrics.Operations {
		mc.logger.Debug("Operation metrics",
			zap.String("operation", operation),
			zap.Int64("count", opMetrics.Count),
			zap.Duration("avg_duration", opMetrics.AvgDuration),
			zap.Int64("errors", opMetrics.Errors),
			zap.Time("last_operation", opMetrics.LastOperation))
	}
}
