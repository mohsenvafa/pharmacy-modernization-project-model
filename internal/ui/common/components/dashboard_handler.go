package components

import (
	"net/http"

	patSvc "github.com/pharmacy-modernization-project-model/internal/domain/patient/service"
	preSvc "github.com/pharmacy-modernization-project-model/internal/domain/prescription/service"
)

func DashboardHandler(ps patSvc.Service, rs preSvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pats, err := ps.List(r.Context(), "", 1000, 0)
		if err != nil {
			http.Error(w, "failed to load patients", http.StatusInternalServerError)
			return
		}

		pres, err := rs.List(r.Context(), "Active", 1000, 0)
		if err != nil {
			http.Error(w, "failed to load prescriptions", http.StatusInternalServerError)
			return
		}

		page := BaseLayout("Dashboard", Dashboard(len(pats), len(pres)))
		if err := page.Render(r.Context(), w); err != nil {
			http.Error(w, "failed to render dashboard", http.StatusInternalServerError)
			return
		}
	}
}
