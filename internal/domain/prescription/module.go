package prescription

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	prescriptionapi "pharmacy-modernization-project-model/internal/domain/prescription/api"
	prescriptionrepo "pharmacy-modernization-project-model/internal/domain/prescription/repository"
	prescriptionservice "pharmacy-modernization-project-model/internal/domain/prescription/service"
	uiprescription "pharmacy-modernization-project-model/internal/domain/prescription/ui"
	irisbilling "pharmacy-modernization-project-model/internal/integrations/iris_billing"
	irispharmacy "pharmacy-modernization-project-model/internal/integrations/iris_pharmacy"
)

type ModuleDependencies struct {
	Logger         *zap.Logger
	PharmacyClient irispharmacy.Client
	BillingClient  irisbilling.Client
}

type ModuleExport struct {
	PrescriptionService prescriptionservice.PrescriptionService
}

func Module(r chi.Router, deps *ModuleDependencies) ModuleExport {
	repo := prescriptionrepo.NewPrescriptionMemoryRepository()
	pharmacyClient := deps.PharmacyClient
	if pharmacyClient == nil {
		pharmacyClient = irispharmacy.NewMockClient(map[string]irispharmacy.GetPrescriptionResponse{}, deps.Logger)
	}
	billingClient := deps.BillingClient
	if billingClient == nil {
		billingClient = irisbilling.NewMockClient(map[string]irisbilling.GetInvoiceResponse{}, deps.Logger)
	}

	svc := prescriptionservice.New(repo, deps.Logger, pharmacyClient, billingClient)

	prescriptionapi.MountApi(r, prescriptionapi.New(svc, deps.Logger))
	uiprescription.MountUI(r, &uiprescription.PrescriptionDependencies{PrescriptionSvc: svc, Log: deps.Logger})

	return ModuleExport{PrescriptionService: svc}
}
