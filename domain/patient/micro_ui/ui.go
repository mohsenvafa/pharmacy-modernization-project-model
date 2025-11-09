package microui

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	patientinfo "pharmacy-modernization-project-model/domain/patient/micro_ui/patient_info"
	patientservice "pharmacy-modernization-project-model/domain/patient/service"
)

type Dependencies struct {
	PatientSvc patientservice.PatientService
	Log        *zap.Logger
}

func Mount(r chi.Router, deps *Dependencies) {
	handler := patientinfo.NewHandler(&patientinfo.Dependencies{
		Service: deps.PatientSvc,
		Log:     deps.Log,
	})

	r.Route("/micro-ui/patients", func(r chi.Router) {
		r.Get("/{patientID}", handler.Handle)
		r.Options("/{patientID}", handler.Handle)
	})
}
