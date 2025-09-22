package dashboard

import (
	"net/http"

	patSvc "github.com/pharmacy-modernization-project-model/internal/domain/patient/service"
	preSvc "github.com/pharmacy-modernization-project-model/internal/domain/prescription/service"
)

type DashboardPageHandler struct {
	patientsService      patSvc.PatientService
	prescriptionsService preSvc.Service
}

func NewDashboardPageHandler(patients patSvc.PatientService, prescriptions preSvc.Service) *DashboardPageHandler {
	return &DashboardPageHandler{patientsService: patients, prescriptionsService: prescriptions}
}

func (u *DashboardPageHandler) Handler(w http.ResponseWriter, r *http.Request) {
	pats, err := u.patientsService.List(r.Context(), "", 1000, 0)
	if err != nil {
		http.Error(w, "failed to load patients", http.StatusInternalServerError)
		return
	}

	pres, err := u.prescriptionsService.List(r.Context(), "Active", 1000, 0)
	if err != nil {
		http.Error(w, "failed to load prescriptions", http.StatusInternalServerError)
		return
	}

	page := DashboardPage(DashboardPageParam{
		NumberOfPatients:    len(pats),
		ActivePrescriptions: len(pres),
	})
	if err := page.Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render dashboard", http.StatusInternalServerError)
		return
	}
}
