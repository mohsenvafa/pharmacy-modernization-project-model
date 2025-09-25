package api

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"pharmacy-modernization-project-model/domain/prescription/api/controllers"
	"pharmacy-modernization-project-model/domain/prescription/service"
)

type Dependencies struct {
	Service service.PrescriptionService
	Logger  *zap.Logger
}

func MountAPI(r chi.Router, deps *Dependencies) {
	controller := controllers.NewPrescriptionController(deps.Service, deps.Logger)

	r.Route("/api/v1/prescriptions", func(router chi.Router) {
		controller.RegisterRoutes(router)
	})
}
