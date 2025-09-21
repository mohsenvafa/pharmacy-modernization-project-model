package patient

import (
	"github.com/go-chi/chi/v5"
	handlers "github.com/pharmacy-modernization-project-model/internal/domain/patient/handlers"
)

func Mount(r chi.Router, api *handlers.API) {
	r.Route("/api/v1/patients", func(r chi.Router) {
		r.Get("/", api.List)
		r.Get("/{id}", api.Get)
	})
}
