package prescription

import (
	"github.com/go-chi/chi/v5"
	presSvc "github.com/pharmacy-modernization-project-model/internal/domain/prescription/service"
	prescriptionList "github.com/pharmacy-modernization-project-model/internal/ui/prescription/prescription-list"
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
