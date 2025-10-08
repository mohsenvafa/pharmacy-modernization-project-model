package graphql

import (
	dashboardgraphql "pharmacy-modernization-project-model/domain/dashboard/graphql"
	patientgraphql "pharmacy-modernization-project-model/domain/patient/graphql"
	prescriptiongraphql "pharmacy-modernization-project-model/domain/prescription/graphql"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver aggregates all domain-specific GraphQL resolvers
type Resolver struct {
	PatientResolver      *patientgraphql.PatientResolver
	PrescriptionResolver *prescriptiongraphql.PrescriptionResolver
	DashboardResolver    *dashboardgraphql.DashboardResolver
}
