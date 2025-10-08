package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	dashboardgraphql "pharmacy-modernization-project-model/domain/dashboard/graphql"
	dashboardservice "pharmacy-modernization-project-model/domain/dashboard/service"
	patientgraphql "pharmacy-modernization-project-model/domain/patient/graphql"
	patientservice "pharmacy-modernization-project-model/domain/patient/service"
	prescriptiongraphql "pharmacy-modernization-project-model/domain/prescription/graphql"
	prescriptionservice "pharmacy-modernization-project-model/domain/prescription/service"
	"pharmacy-modernization-project-model/internal/graphql/generated"
	"pharmacy-modernization-project-model/internal/platform/paths"
)

// Dependencies holds all the service dependencies needed for GraphQL resolvers
type Dependencies struct {
	PatientService      patientservice.PatientService
	AddressService      patientservice.AddressService
	PrescriptionService prescriptionservice.PrescriptionService
	DashboardService    dashboardservice.IDashboardService
	Logger              *zap.Logger
}

// MountGraphQL mounts GraphQL endpoints on the provided router
// Paths are defined in internal/platform/paths/registry.go
func MountGraphQL(r chi.Router, deps *Dependencies) {
	// Create domain-specific resolvers
	patientResolver := patientgraphql.NewPatientResolver(
		deps.PatientService,
		deps.AddressService,
		deps.PrescriptionService,
		deps.Logger,
	)

	prescriptionResolver := prescriptiongraphql.NewPrescriptionResolver(
		deps.PrescriptionService,
		deps.PatientService,
		deps.Logger,
	)

	dashboardResolver := dashboardgraphql.NewDashboardResolver(
		deps.DashboardService,
		deps.Logger,
	)

	// Aggregate domain resolvers into root resolver
	resolver := &Resolver{
		PatientResolver:      patientResolver,
		PrescriptionResolver: prescriptionResolver,
		DashboardResolver:    dashboardResolver,
	}

	// Create GraphQL server
	config := generated.Config{Resolvers: resolver}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(config))

	// Mount endpoints using paths from registry
	r.Handle(paths.GraphQLPath, srv)
	r.Handle(paths.GraphQLPlayground, playground.Handler("GraphQL Playground", paths.GraphQLPath))

	deps.Logger.Info("GraphQL server mounted",
		zap.String("endpoint", paths.GraphQLPath),
		zap.String("playground", paths.GraphQLPlayground))
}
