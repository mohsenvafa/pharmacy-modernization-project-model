package dashboard

import (
	"github.com/go-chi/chi/v5"

	dashboardservice "pharmacy-modernization-project-model/domain/dashboard/service"
	dashboardPage "pharmacy-modernization-project-model/domain/dashboard/ui/dashboard_page"
	"pharmacy-modernization-project-model/domain/dashboard/ui/paths"

	dashboardsecurity "pharmacy-modernization-project-model/domain/dashboard/security"
	"pharmacy-modernization-project-model/internal/platform/auth"
)

type DashboardUiDependencies struct {
	Service dashboardservice.IDashboardService
}

func MountUI(r chi.Router, deps *DashboardUiDependencies) {
	handler := dashboardPage.NewDashboardPageHandler(deps.Service)

	// Dashboard requires authentication and dashboard:view permission
	r.With(
		auth.RequireAuthFromCookie(),
		auth.RequirePermissionsMatchAny(dashboardsecurity.ViewAccess),
	).Get(paths.DashboardPath, handler.Handler)
}
