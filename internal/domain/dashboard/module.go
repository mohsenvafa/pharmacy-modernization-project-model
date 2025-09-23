package dashboard

import (
	"github.com/go-chi/chi/v5"

	dashboardproviders "github.com/pharmacy-modernization-project-model/internal/domain/dashboard/providers"
	dashboardservice "github.com/pharmacy-modernization-project-model/internal/domain/dashboard/service"
	dashboardsvc "github.com/pharmacy-modernization-project-model/internal/domain/dashboard/ui"
)

type ModuleDependencies struct {
	PatientStats      dashboardproviders.PatientStatsProvider
	PrescriptionStats dashboardproviders.PrescriptionStatsProvider
}

type ModuleExport struct{}

func Module(r chi.Router, deps *ModuleDependencies) ModuleExport {
	service := dashboardservice.New(deps.PatientStats, deps.PrescriptionStats)
	dashboardsvc.MountUI(r, &dashboardsvc.DashboardUiDependencies{Service: service})
	return ModuleExport{}
}
