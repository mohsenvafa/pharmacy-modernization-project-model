package handlers

import (
	"net/http"

	patSvc "github.com/pharmacy-modernization-project-model/internal/domain/patient/service"
	preSvc "github.com/pharmacy-modernization-project-model/internal/domain/prescription/service"
	dashboard "github.com/pharmacy-modernization-project-model/internal/ui/dashboard/components"
)

type DashboardPageDI struct {
	patients      patSvc.Service
	prescriptions preSvc.Service
}

func NewDashboardPage(patients patSvc.Service, prescriptions preSvc.Service) *DashboardPageDI {
	return &DashboardPageDI{patients: patients, prescriptions: prescriptions}
}

func (u *DashboardPageDI) DashboardPage(w http.ResponseWriter, r *http.Request) {
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

	page := dashboard.DashboardPage(len(pats), len(pres))
	if err := page.Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render dashboard", http.StatusInternalServerError)
		return
	}
}
