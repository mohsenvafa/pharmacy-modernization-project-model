package graphql

import (
	"context"

	"go.uber.org/zap"

	dashboardservice "pharmacy-modernization-project-model/domain/dashboard/service"
	"pharmacy-modernization-project-model/internal/graphql/generated"
)

// DashboardResolver handles all Dashboard domain GraphQL operations
type DashboardResolver struct {
	DashboardService dashboardservice.IDashboardService
	Logger           *zap.Logger
}

// NewDashboardResolver creates a new dashboard resolver
func NewDashboardResolver(
	dashboardSvc dashboardservice.IDashboardService,
	logger *zap.Logger,
) *DashboardResolver {
	return &DashboardResolver{
		DashboardService: dashboardSvc,
		Logger:           logger,
	}
}

// ============================================================================
// Query Resolvers
// ============================================================================

// DashboardStats resolves the dashboardStats query
func (r *DashboardResolver) DashboardStats(ctx context.Context) (*generated.DashboardStats, error) {
	summary, err := r.DashboardService.Summary(ctx)
	if err != nil {
		r.Logger.Error("Failed to fetch dashboard stats", zap.Error(err))
		return nil, err
	}

	return &generated.DashboardStats{
		TotalPatients:       summary.TotalPatients,
		ActivePrescriptions: summary.ActivePrescriptions,
	}, nil
}
