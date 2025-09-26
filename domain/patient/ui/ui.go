package ui

import (
	patSvc "pharmacy-modernization-project-model/domain/patient/service"
	patientcomponents "pharmacy-modernization-project-model/domain/patient/ui/components/addresslist_server_side"
	patientdetail "pharmacy-modernization-project-model/domain/patient/ui/patient_detail"
	pateitnList "pharmacy-modernization-project-model/domain/patient/ui/patient_list"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UiDpendencies struct {
	PatientSvc patSvc.PatientService
	AddressSvc patSvc.AddressService
	Log        *zap.Logger
}

func MountUI(r chi.Router, patientDpendencies *UiDpendencies) {
	patientListPage := pateitnList.NewPatientListHandler(patientDpendencies.PatientSvc, patientDpendencies.Log)
	addressListHandler := patientcomponents.NewAddressListComponentHandler(patientDpendencies.AddressSvc, patientDpendencies.Log)

	patientDetailPage := patientdetail.NewPatientDetailHandler(
		patientDpendencies.PatientSvc,
		addressListHandler,
		patientDpendencies.Log,
	)
	r.Route("/patients", func(r chi.Router) {
		r.Get("/", patientListPage.Handler)
		r.Get("/{patientID}", patientDetailPage.Handler)
	})
}
