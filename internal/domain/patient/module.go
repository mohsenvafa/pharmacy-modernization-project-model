package patient

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	patientapi "github.com/pharmacy-modernization-project-model/internal/domain/patient/api"
	patientrepo "github.com/pharmacy-modernization-project-model/internal/domain/patient/repository"
	patientservice "github.com/pharmacy-modernization-project-model/internal/domain/patient/service"
	uipatient "github.com/pharmacy-modernization-project-model/internal/domain/patient/ui"
)

type ModuleDependencies struct {
	Logger *zap.Logger
}

type ModuleExport struct {
	PatientService patientservice.PatientService
	AddressService patientservice.AddressService
}

func Module(r chi.Router, deps *ModuleDependencies) ModuleExport {
	patRepo := patientrepo.NewPatientMemoryRepository()
	patSvc := patientservice.New(patRepo, deps.Logger)
	addrRepo := patientrepo.NewAddressMemoryRepository()
	addrSvc := patientservice.NewAddressService(addrRepo)

	patientapi.MountAPI(r, &patientapi.Dependencies{
		PatientService: patSvc,
		AddressService: addrSvc,
		Logger:         deps.Logger,
	})

	uipatient.MountUI(r, &uipatient.UiDpendencies{PatientSvc: patSvc, Log: deps.Logger})

	return ModuleExport{PatientService: patSvc, AddressService: addrSvc}
}
