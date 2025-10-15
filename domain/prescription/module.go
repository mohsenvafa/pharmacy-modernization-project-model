package prescription

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	prescriptionapi "pharmacy-modernization-project-model/domain/prescription/api"
	prescriptionbuilder "pharmacy-modernization-project-model/domain/prescription/builder"
	prescriptionservice "pharmacy-modernization-project-model/domain/prescription/service"
	uiprescription "pharmacy-modernization-project-model/domain/prescription/ui"
	irisbilling "pharmacy-modernization-project-model/internal/integrations/iris_billing"
	irispharmacy "pharmacy-modernization-project-model/internal/integrations/iris_pharmacy"
	"pharmacy-modernization-project-model/internal/platform/cache"
)

type ModuleDependencies struct {
	Logger                       *zap.Logger
	PharmacyClient               irispharmacy.PharmacyClient
	BillingClient                irisbilling.BillingClient
	PrescriptionsMongoCollection *mongo.Collection
	CacheService                 cache.Cache
}

type ModuleExport struct {
	PrescriptionService prescriptionservice.PrescriptionService
}

func Module(r chi.Router, deps *ModuleDependencies) ModuleExport {
	repo := prescriptionbuilder.CreatePrescriptionRepository(deps.Logger, deps.PrescriptionsMongoCollection)
	pharmacyClient := deps.PharmacyClient
	if pharmacyClient == nil {
		pharmacyClient = irispharmacy.NewMockClient(deps.Logger)
	}
	billingClient := deps.BillingClient
	if billingClient == nil {
		billingClient = irisbilling.NewMockClient(deps.Logger)
	}

	svc := prescriptionservice.New(repo, deps.CacheService, deps.Logger, pharmacyClient, billingClient)

	prescriptionapi.MountAPI(r, &prescriptionapi.Dependencies{Service: svc, Logger: deps.Logger})
	uiprescription.MountUI(r, &uiprescription.PrescriptionDependencies{PrescriptionSvc: svc, Log: deps.Logger})

	return ModuleExport{PrescriptionService: svc}
}
