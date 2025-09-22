package patient

import (
	"github.com/go-chi/chi/v5"
	patSvc "github.com/pharmacy-modernization-project-model/internal/domain/patient/service"
	pateitnList "github.com/pharmacy-modernization-project-model/internal/ui/patient/patient-list"
	"go.uber.org/zap"
)

type PatientDpendencies struct {
	PatientSvc patSvc.Service
	Log        *zap.Logger
}

func MountUI(r chi.Router, patientDpendencies *PatientDpendencies) {
	patientListPage := pateitnList.NewPatientListHandler(patientDpendencies.PatientSvc, patientDpendencies.Log)
	r.Route("/patients", func(r chi.Router) {
		r.Get("/", patientListPage.Handler)
	})
}
