package components

import (
	"net/http"

	patSvc "github.com/pharmacy-modernization-project-model/internal/domain/patient/service"
	preSvc "github.com/pharmacy-modernization-project-model/internal/domain/prescription/service"
)

func DashboardHandler(ps patSvc.Service, rs preSvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pats, _ := ps.List(r.Context(), "", 1000, 0)
		pres, _ := rs.List(r.Context(), "Active", 1000, 0)
		_ = Dashboard(len(pats), len(pres)).Render(r.Context(), w)
	}
}
