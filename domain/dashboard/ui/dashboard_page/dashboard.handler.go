package dashboard_page

import (
	"net/http"

	dashboardservice "pharmacy-modernization-project-model/domain/dashboard/service"
)

type DashboardPageHandler struct {
	service dashboardservice.IDashboardService
}

func NewDashboardPageHandler(service dashboardservice.IDashboardService) *DashboardPageHandler {
	return &DashboardPageHandler{service: service}
}

func (u *DashboardPageHandler) Handler(w http.ResponseWriter, r *http.Request) {
	summary, err := u.service.Summary(r.Context())
	if err != nil {
		http.Error(w, "failed to load dashboard data", http.StatusInternalServerError)
		return
	}

	page := DashboardPage(DashboardPageParam{
		NumberOfPatients:    summary.TotalPatients,
		ActivePrescriptions: summary.ActivePrescriptions,
	})
	if err := page.Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render dashboard", http.StatusInternalServerError)
		return
	}
}
