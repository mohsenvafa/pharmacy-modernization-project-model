package dashboard

import (
	"github.com/go-chi/chi/v5"

	dashboardservice "pharmacy-modernization-project-model/domain/dashboard/service"
	dashboardPage "pharmacy-modernization-project-model/domain/dashboard/ui/dashboard_page"
	"pharmacy-modernization-project-model/domain/dashboard/ui/paths"
)

type DashboardUiDependencies struct {
	Service dashboardservice.IDashboardService
}

func MountUI(r chi.Router, deps *DashboardUiDependencies) {
	handler := dashboardPage.NewDashboardPageHandler(deps.Service)
	r.Get(paths.DashboardPath, handler.Handler)
}
