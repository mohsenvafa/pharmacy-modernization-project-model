package patient

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	patientapi "pharmacy-modernization-project-model/domain/patient/api"
	patientbuilder "pharmacy-modernization-project-model/domain/patient/builder"
	patientproviders "pharmacy-modernization-project-model/domain/patient/providers"
	patientservice "pharmacy-modernization-project-model/domain/patient/service"
	uipatient "pharmacy-modernization-project-model/domain/patient/ui"
	uipatientContracts "pharmacy-modernization-project-model/domain/patient/ui/contracts"
)

type ModuleDependencies struct {
	Logger                   *zap.Logger
	PrescriptionProvider     patientproviders.PatientPrescriptionProvider
	PatientsMongoCollection  *mongo.Collection
	AddressesMongoCollection *mongo.Collection
	CacheService             interface { // Cache interface
		Get(ctx context.Context, key string) ([]byte, error)
		Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
		Delete(ctx context.Context, key string) error
		Close() error
	}
}

type ModuleExport struct {
	PatientService patientservice.PatientService
	AddressService patientservice.AddressService
}

func Module(r chi.Router, deps *ModuleDependencies) ModuleExport {
	patRepo := patientbuilder.CreatePatientRepository(deps.Logger, deps.PatientsMongoCollection)
	addrRepo := patientbuilder.CreateAddressRepository(deps.Logger, deps.AddressesMongoCollection)

	patSvc := patientservice.New(patRepo, deps.CacheService, deps.Logger)
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
