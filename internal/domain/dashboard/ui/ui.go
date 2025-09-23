package dashboard

import (
	"github.com/go-chi/chi/v5"

	dashboardservice "github.com/pharmacy-modernization-project-model/internal/domain/dashboard/service"
	dashboardPage "github.com/pharmacy-modernization-project-model/internal/domain/dashboard/ui/dashboard_page"
)

type DashboardUiDependencies struct {
	Service dashboardservice.IDashboardService
}

func MountUI(r chi.Router, deps *DashboardUiDependencies) {
	handler := dashboardPage.NewDashboardPageHandler(deps.Service)
	r.Get("/", handler.Handler)
}
