package app

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/pharmacy-modernization-project-model/internal/platform/httpx"
	"github.com/pharmacy-modernization-project-model/internal/platform/logging"

	prescapi "github.com/pharmacy-modernization-project-model/internal/domain/prescription/api"
	prescrepo "github.com/pharmacy-modernization-project-model/internal/domain/prescription/repository"
	prescservice "github.com/pharmacy-modernization-project-model/internal/domain/prescription/service"

	uiprescriptionModule "github.com/pharmacy-modernization-project-model/internal/domain/prescription/ui"

	patientModule "github.com/pharmacy-modernization-project-model/internal/domain/patient"
)

func (a *App) wire() error {
	// Logger
	logger := logging.NewLogger(a.Cfg)
	a.Logger = logger

	// Shared HTTP client (for future external calls)
	_ = httpx.NewClient(a.Cfg)

	// In-memory repos
	preRepo := prescrepo.NewPrescriptionMemoryRepository()

	// Services
	preSvc := prescservice.New(preRepo, logger.Base)

	// Router & middleware
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(logging.CorrelationID())
	r.Use(logging.ZapRequestLogger(logger.Base))
	r.Use(middleware.Timeout(60 * time.Second))

	// Static assets
	r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("web/public"))))

	patientModule.Module(r, &patientModule.ModuleDependencies{Logger: logger.Base})
	// API
	prescapi.Mount(r, prescapi.New(preSvc, logger.Base))

	// UI
	//uidashboardMdoule.MountUI(r, &uidashboardMdoule.DashboardDpendencies{PatientSvc: patSvc, PrescriptionSvc: preSvc})
	uiprescriptionModule.MountUI(r, &uiprescriptionModule.PrescriptionDependencies{PrescriptionSvc: preSvc, Log: logger.Base})
	a.Router = r
	return nil
}
