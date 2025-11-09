package microui

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	prescriptioninfo "pharmacy-modernization-project-model/domain/prescription/micro_ui/prescription_info"
	prescriptionsvc "pharmacy-modernization-project-model/domain/prescription/service"
)

// Dependencies contains the services required by the micro UI layer for prescriptions.
type Dependencies struct {
	PrescriptionSvc prescriptionsvc.PrescriptionService
	Log             *zap.Logger
}

// Mount registers micro UI routes for the prescription domain.
func Mount(r chi.Router, deps *Dependencies) {
	handler := prescriptioninfo.NewHandler(&prescriptioninfo.Dependencies{
		Service: deps.PrescriptionSvc,
		Log:     deps.Log,
	})

	r.Route("/micro-ui/prescriptions", func(r chi.Router) {
		r.Get("/{prescriptionID}", handler.Handle)
	})
}
