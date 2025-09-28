package ui

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	patientproviders "pharmacy-modernization-project-model/domain/patient/providers"
	patSvc "pharmacy-modernization-project-model/domain/patient/service"
	addresscomponents "pharmacy-modernization-project-model/domain/patient/ui/components/addresslist_server_side"
	patientprescriptioncomponents "pharmacy-modernization-project-model/domain/patient/ui/components/patient_prescriptions"
	patientdetail "pharmacy-modernization-project-model/domain/patient/ui/patient_detail"
	pateitnList "pharmacy-modernization-project-model/domain/patient/ui/patient_list"
)

type UiDpendencies struct {
	PatientSvc           patSvc.PatientService
	AddressSvc           patSvc.AddressService
	PrescriptionProvider patientproviders.PatientPrescriptionProvider
	Log                  *zap.Logger
}

func MountUI(r chi.Router, patientDpendencies *UiDpendencies) {
	patientListPage := pateitnList.NewPatientListHandler(patientDpendencies.PatientSvc, patientDpendencies.Log)
	addressListHandler := addresscomponents.NewAddressListComponentHandler(patientDpendencies.AddressSvc, patientDpendencies.Log)
	prescriptionListHandler := patientprescriptioncomponents.NewPrescriptionListComponentHandler(patientprescriptioncomponents.PrescriptionListDependencies{
		Provider: patientDpendencies.PrescriptionProvider,
		Log:      patientDpendencies.Log,
	})

	patientDetailPage := patientdetail.NewPatientDetailHandler(
		patientDpendencies.PatientSvc,
		addressListHandler,
		prescriptionListHandler,
		patientDpendencies.Log,
	)
	r.Route("/patients", func(r chi.Router) {
		r.Get("/", patientListPage.Handler)
		r.Get("/{patientID}", patientDetailPage.Handler)
	})
}
