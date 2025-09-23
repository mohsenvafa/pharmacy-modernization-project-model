package dashboard

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	dashboardsvc "github.com/pharmacy-modernization-project-model/internal/domain/dashboard/ui"
	patSvc "github.com/pharmacy-modernization-project-model/internal/domain/patient/service"
	presSvc "github.com/pharmacy-modernization-project-model/internal/domain/prescription/service"
)

type ModuleDependencies struct {
	Logger          *zap.Logger
	PatientSvc      patSvc.PatientService
	PrescriptionSvc presSvc.PrescriptionService
}

func Module(r chi.Router, deps *ModuleDependencies) {
	dashboardsvc.MountUI(r, &dashboardsvc.DashboardDpendencies{
		PatientSvc:      deps.PatientSvc,
		PrescriptionSvc: deps.PrescriptionSvc,
	})
}
