package dashboard

import (
	"github.com/go-chi/chi/v5"
	handlers "github.com/pharmacy-modernization-project-model/internal/ui/dashboard/handlers"
)

func MountUI(r chi.Router, ui *handlers.UI) {
	r.Get("/", ui.Dashboard)
}
