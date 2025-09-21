package dashboard

import (
	"github.com/go-chi/chi/v5"
	patSvc "github.com/pharmacy-modernization-project-model/internal/domain/patient/service"
	preSvc "github.com/pharmacy-modernization-project-model/internal/domain/prescription/service"
)

type DashboardDpendencies struct {
	PatientSvc      patSvc.Service
	PrescriptionSvc preSvc.Service
}

func MountUI(r chi.Router, dashboardDpendencies *DashboardDpendencies) {
	dashboardPage := NewDashboardPage(dashboardDpendencies.PatientSvc, dashboardDpendencies.PrescriptionSvc)

	r.Get("/", dashboardPage.DashboardPageHandler)
}
