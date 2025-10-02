package ui

import (
	"github.com/go-chi/chi/v5"

	addresscomponents "pharmacy-modernization-project-model/domain/patient/ui/components/address_list"
	patientprescriptioncomponents "pharmacy-modernization-project-model/domain/patient/ui/components/patient_prescriptions"
	contracts "pharmacy-modernization-project-model/domain/patient/ui/contracts"
	"pharmacy-modernization-project-model/domain/patient/ui/paths"
	patientdetail "pharmacy-modernization-project-model/domain/patient/ui/patient_detail"
	patientlist "pharmacy-modernization-project-model/domain/patient/ui/patient_list"
)

func MountUI(r chi.Router, dep *contracts.UiDependencies) {
	patientListComponent := patientlist.NewPatientListComponent(dep)
	addressListComponent := addresscomponents.NewAddressListComponent(dep)
	prescriptionListComponent := patientprescriptioncomponents.NewPrescriptionListComponent(dep)
	patientDetailComponent := patientdetail.NewPatientDetailComponent(dep, addressListComponent, prescriptionListComponent)

	r.Route(paths.BasePath, func(r chi.Router) {
		r.Get("/", patientListComponent.Handler)
		r.Get("/components/patient-prescriptions-card", prescriptionListComponent.Handler)
		r.Get("/{patientID}", patientDetailComponent.Handler)
	})
}
