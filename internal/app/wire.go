package app

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/pharmacy-modernization-project-model/internal/platform/httpx"
	"github.com/pharmacy-modernization-project-model/internal/platform/logging"

	patientroutes "github.com/pharmacy-modernization-project-model/internal/domain/patient"
	patientapi "github.com/pharmacy-modernization-project-model/internal/domain/patient/handlers"
	patientrepo "github.com/pharmacy-modernization-project-model/internal/domain/patient/repository"
	patientservice "github.com/pharmacy-modernization-project-model/internal/domain/patient/service"

	prescroutes "github.com/pharmacy-modernization-project-model/internal/domain/prescription"
	prescapi "github.com/pharmacy-modernization-project-model/internal/domain/prescription/handlers"
	prescrepo "github.com/pharmacy-modernization-project-model/internal/domain/prescription/repository"
	prescservice "github.com/pharmacy-modernization-project-model/internal/domain/prescription/service"

	uidashboardMdoule "github.com/pharmacy-modernization-project-model/internal/ui/dashboard"
	uipatientMdoule "github.com/pharmacy-modernization-project-model/internal/ui/patient"
	uiprescriptionModule "github.com/pharmacy-modernization-project-model/internal/ui/prescription"
)

func (a *App) wire() error {
	// Logger
	logger := logging.NewLogger(a.Cfg)
	a.Logger = logger

	// Shared HTTP client (for future external calls)
	_ = httpx.NewClient(a.Cfg)

	// In-memory repos
	patRepo := patientrepo.NewPatientMemoryRepository()
	preRepo := prescrepo.NewPrescriptionMemoryRepository()

	// Services
	patSvc := patientservice.New(patRepo, logger.Base)
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

	// API
	patientroutes.Mount(r, patientapi.New(patSvc, logger.Base))
	prescroutes.Mount(r, prescapi.New(preSvc, logger.Base))

	// UI
	uidashboardMdoule.MountUI(r, &uidashboardMdoule.DashboardDpendencies{PatientSvc: patSvc, PrescriptionSvc: preSvc})
	uipatientMdoule.MountUI(r, &uipatientMdoule.PatientDpendencies{PatientSvc: patSvc, Log: logger.Base})
	uiprescriptionModule.MountUI(r, &uiprescriptionModule.PrescriptionDependencies{PrescriptionSvc: preSvc, Log: logger.Base})
	a.Router = r
	return nil
}
