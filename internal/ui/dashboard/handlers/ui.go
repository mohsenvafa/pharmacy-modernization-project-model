package handlers

import (
	"net/http"

	patSvc "github.com/pharmacy-modernization-project-model/internal/domain/patient/service"
	preSvc "github.com/pharmacy-modernization-project-model/internal/domain/prescription/service"
	layouts "github.com/pharmacy-modernization-project-model/internal/ui/common/layouts"
	dashboard "github.com/pharmacy-modernization-project-model/internal/ui/dashboard/components"
)

type UI struct {
	patients      patSvc.Service
	prescriptions preSvc.Service
}

func New(patients patSvc.Service, prescriptions preSvc.Service) *UI {
	return &UI{patients: patients, prescriptions: prescriptions}
}

func (u *UI) Dashboard(w http.ResponseWriter, r *http.Request) {
	pats, err := u.patients.List(r.Context(), "", 1000, 0)
	if err != nil {
		http.Error(w, "failed to load patients", http.StatusInternalServerError)
		return
	}

	pres, err := u.prescriptions.List(r.Context(), "Active", 1000, 0)
	if err != nil {
		http.Error(w, "failed to load prescriptions", http.StatusInternalServerError)
		return
	}

	page := layouts.BaseLayout("Dashboard", dashboard.Dashboard(len(pats), len(pres)))
	if err := page.Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render dashboard", http.StatusInternalServerError)
		return
	}
}
