package dashboard

import (
	"github.com/go-chi/chi/v5"
)

func MountUI(r chi.Router, ui *DashboardPageDI) {
	r.Get("/", ui.DashboardPageHandler)
}
