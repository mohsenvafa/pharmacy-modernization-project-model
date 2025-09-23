package prescription

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	prescriptionapi "github.com/pharmacy-modernization-project-model/internal/domain/prescription/api"
	prescriptionrepo "github.com/pharmacy-modernization-project-model/internal/domain/prescription/repository"
	prescriptionservice "github.com/pharmacy-modernization-project-model/internal/domain/prescription/service"
	uiprescription "github.com/pharmacy-modernization-project-model/internal/domain/prescription/ui"
)

type ModuleDependencies struct {
	Logger *zap.Logger
}

func Module(r chi.Router, deps *ModuleDependencies) prescriptionservice.PrescriptionService {
	repo := prescriptionrepo.NewPrescriptionMemoryRepository()
	svc := prescriptionservice.New(repo, deps.Logger)

	prescriptionapi.MountApi(r, prescriptionapi.New(svc, deps.Logger))
	uiprescription.MountUI(r, &uiprescription.PrescriptionDependencies{PrescriptionSvc: svc, Log: deps.Logger})

	return svc
}
