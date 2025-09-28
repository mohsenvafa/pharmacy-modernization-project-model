package patient

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	patientapi "pharmacy-modernization-project-model/domain/patient/api"
	patientproviders "pharmacy-modernization-project-model/domain/patient/providers"
	patientrepo "pharmacy-modernization-project-model/domain/patient/repository"
	patientservice "pharmacy-modernization-project-model/domain/patient/service"
	uipatient "pharmacy-modernization-project-model/domain/patient/ui"
	uipatientContracts "pharmacy-modernization-project-model/domain/patient/ui/contracts"
)

type ModuleDependencies struct {
	Logger               *zap.Logger
	PrescriptionProvider patientproviders.PatientPrescriptionProvider
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

	uipatient.MountUI(r, &uipatientContracts.UiDependencies{
		PatientSvc:           patSvc,
		AddressSvc:           addrSvc,
		PrescriptionProvider: deps.PrescriptionProvider,
		Log:                  deps.Logger,
	})

	return ModuleExport{PatientService: patSvc, AddressService: addrSvc}
}
