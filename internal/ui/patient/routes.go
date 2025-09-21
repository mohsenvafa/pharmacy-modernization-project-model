package patient

import (
	"github.com/go-chi/chi/v5"
	handlers "github.com/pharmacy-modernization-project-model/internal/ui/patient/handlers"
)

func MountUI(r chi.Router, ui *handlers.UI) {
	r.Route("/patients", func(r chi.Router) {
		r.Get("/", ui.List)
	})
}
