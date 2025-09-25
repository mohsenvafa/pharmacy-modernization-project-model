package prescription

import (
	presSvc "pharmacy-modernization-project-model/domain/prescription/service"
	prescriptionList "pharmacy-modernization-project-model/domain/prescription/ui/prescription_list"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type PrescriptionDependencies struct {
	PrescriptionSvc presSvc.PrescriptionService
	Log             *zap.Logger
}

func MountUI(r chi.Router, deps *PrescriptionDependencies) {
	prescriptionListHandler := prescriptionList.NewPrescriptionListHandler(deps.PrescriptionSvc, deps.Log)

	r.Route("/prescriptions", func(r chi.Router) {
		r.Get("/", prescriptionListHandler.Handler)
	})
}
