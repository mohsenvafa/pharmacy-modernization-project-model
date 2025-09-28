package ui

import (
	"github.com/go-chi/chi/v5"

	addresscomponents "pharmacy-modernization-project-model/domain/patient/ui/components/addresslist_server_side"
	patientprescriptioncomponents "pharmacy-modernization-project-model/domain/patient/ui/components/patient_prescriptions"
	contracts "pharmacy-modernization-project-model/domain/patient/ui/contracts"
	patientdetail "pharmacy-modernization-project-model/domain/patient/ui/patient_detail"
	pateitnList "pharmacy-modernization-project-model/domain/patient/ui/patient_list"
)

func MountUI(r chi.Router, dep *contracts.UiDependencies) {
	patientListPage := pateitnList.NewPatientListHandler(
		dep.PatientSvc,
		dep.Log)

	addressListHandler := addresscomponents.NewAddressListComponentHandler(
		dep.AddressSvc,
		dep.Log)

	prescriptionListHandler := patientprescriptioncomponents.NewPrescriptionListComponent(dep)

	patientDetailPage := patientdetail.NewPatientDetailHandler(
		dep,
		addressListHandler,
		prescriptionListHandler,
	)

	r.Route("/patients", func(r chi.Router) {
		r.Get("/", patientListPage.Handler)
		r.Get("/components/patient-prescriptions-card", prescriptionListHandler.Handler)
		r.Get("/{patientID}", patientDetailPage.Handler)
	})
}
