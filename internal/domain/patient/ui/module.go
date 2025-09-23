package patient

import (
	"github.com/go-chi/chi/v5"
	patSvc "github.com/pharmacy-modernization-project-model/internal/domain/patient/service"
	patientdetail "github.com/pharmacy-modernization-project-model/internal/domain/patient/ui/patient_detail"
	pateitnList "github.com/pharmacy-modernization-project-model/internal/domain/patient/ui/patient_list"
	"go.uber.org/zap"
)

type PatientDpendencies struct {
	PatientSvc patSvc.PatientService
	Log        *zap.Logger
}

func MountUI(r chi.Router, patientDpendencies *PatientDpendencies) {
	patientListPage := pateitnList.NewPatientListHandler(patientDpendencies.PatientSvc, patientDpendencies.Log)
	patientDetailPage := patientdetail.NewPatientDetailHandler(patientDpendencies.PatientSvc, patientDpendencies.Log)
	r.Route("/patients", func(r chi.Router) {
		r.Get("/", patientListPage.Handler)
		r.Get("/{patientID}", patientDetailPage.Handler)
	})
}
