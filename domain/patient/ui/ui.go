package ui

import (
	"github.com/go-chi/chi/v5"

	addresscomponents "pharmacy-modernization-project-model/domain/patient/ui/components/address_list"
	patientprescriptioncomponents "pharmacy-modernization-project-model/domain/patient/ui/components/patient_prescriptions"
	contracts "pharmacy-modernization-project-model/domain/patient/ui/contracts"
	"pharmacy-modernization-project-model/domain/patient/ui/paths"
	patientdetail "pharmacy-modernization-project-model/domain/patient/ui/patient_detail"
	patientlist "pharmacy-modernization-project-model/domain/patient/ui/patient_list"
	patientsearch "pharmacy-modernization-project-model/domain/patient/ui/patient_search"

	patientsecurity "pharmacy-modernization-project-model/domain/patient/security"
	"pharmacy-modernization-project-model/internal/platform/auth"
)

func MountUI(r chi.Router, dep *contracts.UiDependencies) {
	patientListComponent := patientlist.NewPatientListComponent(dep)
	patientSearchComponent := patientsearch.NewPatientSearchPageComponent(dep)
	addressListComponent := addresscomponents.NewAddressListComponent(dep)
	prescriptionListComponent := patientprescriptioncomponents.NewPrescriptionListComponent(dep)
	patientDetailComponent := patientdetail.NewPatientDetailComponent(dep, addressListComponent, prescriptionListComponent)

	r.Route(paths.BasePath, func(r chi.Router) {
		// All patient UI routes require authentication (cookie-based for web)
		r.Use(auth.RequireAuthFromCookie())

		// All routes require patient:read permission or admin access
		r.Use(auth.RequirePermissionsMatchAny(patientsecurity.ReadAccess))

		r.Get(paths.ListRoute, patientListComponent.Handler)
		r.Get(paths.SearchRoute, patientSearchComponent.Handler)
		r.Get(paths.PatientPrescriptionCardComponentRoute, prescriptionListComponent.Handler)
		r.Get(paths.DetailRoute, patientDetailComponent.Handler)
	})
}
