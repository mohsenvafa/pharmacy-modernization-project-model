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
	authplatform "pharmacy-modernization-project-model/internal/platform/auth"
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

	// Create GraphQL server with auth directives
	config := generated.Config{
		Resolvers: resolver,
		Directives: generated.DirectiveRoot{
			Auth:          authplatform.AuthDirective(),
			PermissionAny: authplatform.PermissionAnyDirective(),
			PermissionAll: authplatform.PermissionAllDirective(),
		},
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(config))

	// Mount GraphQL endpoint with auth middleware (to set user in context)
	// Uses dev mode if enabled, otherwise requires real JWT
	r.Handle(paths.GraphQLPath, authplatform.RequireAuthWithDevMode()(srv))
	r.Handle(paths.GraphQLPlayground, playground.Handler("GraphQL Playground", paths.GraphQLPath))

	deps.Logger.Info("GraphQL server mounted",
		zap.String("endpoint", paths.GraphQLPath),
		zap.String("playground", paths.GraphQLPlayground))
}
