package api

import (
	"github.com/go-chi/chi/v5"
)

func Mount(r chi.Router, api *API) {
	r.Route("/api/v1/patients", func(r chi.Router) {
		r.Get("/", api.List)
		r.Get("/{id}", api.Get)
	})
}
