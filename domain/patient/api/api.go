package api

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	controllers "pharmacy-modernization-project-model/domain/patient/api/controllers"
	"pharmacy-modernization-project-model/domain/patient/service"
)

type Dependencies struct {
	PatientService service.PatientService
	AddressService service.AddressService
	Logger         *zap.Logger
}

func MountAPI(r chi.Router, deps *Dependencies) {
	patientController := controllers.NewPatientController(deps.PatientService, deps.Logger)
	addressController := controllers.NewAddressController(deps.AddressService, deps.Logger)

	r.Route("/api/v1/patients", func(router chi.Router) {
		patientController.RegisterRoutes(router)
		router.Route("/{patientID}/addresses", func(addressRouter chi.Router) {
			addressController.RegisterRoutes(addressRouter)
		})
	})
}
