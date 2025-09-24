package ui

import (
	patSvc "pharmacy-modernization-project-model/internal/domain/patient/service"
	patientdetail "pharmacy-modernization-project-model/internal/domain/patient/ui/patient_detail"
	pateitnList "pharmacy-modernization-project-model/internal/domain/patient/ui/patient_list"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UiDpendencies struct {
	PatientSvc patSvc.PatientService
	Log        *zap.Logger
}

func MountUI(r chi.Router, patientDpendencies *UiDpendencies) {
	patientListPage := pateitnList.NewPatientListHandler(patientDpendencies.PatientSvc, patientDpendencies.Log)
	patientDetailPage := patientdetail.NewPatientDetailHandler(patientDpendencies.PatientSvc, patientDpendencies.Log)
	r.Route("/patients", func(r chi.Router) {
		r.Get("/", patientListPage.Handler)
		r.Get("/{patientID}", patientDetailPage.Handler)
	})
}
