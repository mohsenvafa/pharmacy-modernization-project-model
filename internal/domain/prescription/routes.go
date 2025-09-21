package prescription

import (
	"github.com/go-chi/chi/v5"
	handlers "github.com/pharmacy-modernization-project-model/internal/domain/prescription/handlers"
)

func Mount(r chi.Router, api *handlers.API) {
	r.Route("/api/v1/prescriptions", func(r chi.Router) {
		r.Get("/", api.List)
		r.Get("/{id}", api.Get)
	})
}
