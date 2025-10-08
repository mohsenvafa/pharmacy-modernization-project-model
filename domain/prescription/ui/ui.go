package prescription

import (
	presSvc "pharmacy-modernization-project-model/domain/prescription/service"
	"pharmacy-modernization-project-model/domain/prescription/ui/paths"
	prescriptionList "pharmacy-modernization-project-model/domain/prescription/ui/prescription_list"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	prescriptionsecurity "pharmacy-modernization-project-model/domain/prescription/security"
	"pharmacy-modernization-project-model/internal/platform/auth"
)

type PrescriptionDependencies struct {
	PrescriptionSvc presSvc.PrescriptionService
	Log             *zap.Logger
}

func MountUI(r chi.Router, deps *PrescriptionDependencies) {
	prescriptionListHandler := prescriptionList.NewPrescriptionListHandler(deps.PrescriptionSvc, deps.Log)

	r.Route(paths.BasePath, func(r chi.Router) {
		// All prescription UI routes require authentication
		// Uses dev mode if enabled, otherwise cookie-based auth
		r.Use(auth.RequireAuthWithDevMode())

		// All routes require prescription:read or healthcare role or admin
		r.Use(auth.RequirePermissionsMatchAny(prescriptionsecurity.ReadAccess))

		r.Get("/", prescriptionListHandler.Handler)
	})
}
