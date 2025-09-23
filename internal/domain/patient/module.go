package patient

import (
	"github.com/go-chi/chi/v5"
	patientapi "github.com/pharmacy-modernization-project-model/internal/domain/patient/api"
	patientrepo "github.com/pharmacy-modernization-project-model/internal/domain/patient/repository"
	patientservice "github.com/pharmacy-modernization-project-model/internal/domain/patient/service"
	uipatient "github.com/pharmacy-modernization-project-model/internal/domain/patient/ui"
	"go.uber.org/zap"
)

type ModuleDependencies struct {
	Logger *zap.Logger
}

func Module(r chi.Router, deps *ModuleDependencies) patientservice.PatientService {
	patRepo := patientrepo.NewPatientMemoryRepository()
	patSvc := patientservice.New(patRepo, deps.Logger)

	patientapi.MountApi(r, patientapi.New(patSvc, deps.Logger))
	uipatient.MountUI(r, &uipatient.UiDpendencies{PatientSvc: patSvc, Log: deps.Logger})

	return patSvc
}
