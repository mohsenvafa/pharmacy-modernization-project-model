package app

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/pharmacy-modernization-project-model/internal/platform/httpx"
	"github.com/pharmacy-modernization-project-model/internal/platform/logging"

	patientapi "github.com/pharmacy-modernization-project-model/internal/domain/patient/handlers"
	patientroutes "github.com/pharmacy-modernization-project-model/internal/domain/patient"
	patientrepo "github.com/pharmacy-modernization-project-model/internal/domain/patient/repository"
	patientservice "github.com/pharmacy-modernization-project-model/internal/domain/patient/service"

	prescapi "github.com/pharmacy-modernization-project-model/internal/domain/prescription/handlers"
	prescroutes "github.com/pharmacy-modernization-project-model/internal/domain/prescription"
	prescrepo "github.com/pharmacy-modernization-project-model/internal/domain/prescription/repository"
	prescservice "github.com/pharmacy-modernization-project-model/internal/domain/prescription/service"

	uipatient "github.com/pharmacy-modernization-project-model/internal/ui/patient/handlers"
	uipatientroutes "github.com/pharmacy-modernization-project-model/internal/ui/patient"
	uipres "github.com/pharmacy-modernization-project-model/internal/ui/prescription/handlers"
	uipresroutes "github.com/pharmacy-modernization-project-model/internal/ui/prescription"
	uicommon "github.com/pharmacy-modernization-project-model/internal/ui/common/components"
)

func (a *App) wire() error {
	// Logger
	logger := logging.NewLogger(a.Cfg)
	a.Logger = logger

	// Shared HTTP client (for future external calls)
	_ = httpx.NewClient(a.Cfg)

	// In-memory repos
	patRepo := patientrepo.NewMemRepo()
	preRepo := prescrepo.NewMemRepo()

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
	uipatientroutes.MountUI(r, uipatient.New(patSvc, logger.Base))
	uipresroutes.MountUI(r, uipres.New(preSvc, logger.Base))

	// Root dashboard
	r.Get("/", uicommon.DashboardHandler(patSvc, preSvc))

	a.Router = r
	return nil
}
