package app

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/pharmacy-modernization-project-model/internal/integrations"
	"github.com/pharmacy-modernization-project-model/internal/platform/httpx"
	"github.com/pharmacy-modernization-project-model/internal/platform/logging"

	dashboardModule "github.com/pharmacy-modernization-project-model/internal/domain/dashboard"
	patientModule "github.com/pharmacy-modernization-project-model/internal/domain/patient"
	prescriptionModule "github.com/pharmacy-modernization-project-model/internal/domain/prescription"
)

func (a *App) wire() error {
	// Logger
	logger := logging.NewLogger(a.Cfg)
	a.Logger = logger

	// Shared HTTP client (for future external calls)
	_ = httpx.NewClient(a.Cfg)

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

	// Integrations
	integration := integrations.New(integrations.Dependencies{
		Config: a.Cfg,
		Logger: logger.Base,
	})

	// Patient Module
	patientMod := patientModule.Module(r, &patientModule.ModuleDependencies{Logger: logger.Base})

	// Prescription Module
	prescriptionMod := prescriptionModule.Module(r, &prescriptionModule.ModuleDependencies{
		Logger:         logger.Base,
		PharmacyClient: integration.Pharmacy,
		BillingClient:  integration.Billing,
	})

	// Dashboard Module
	dashboardModule.Module(r, &dashboardModule.ModuleDependencies{
		PatientStats:      patientMod.PatientService,
		PrescriptionStats: prescriptionMod.PrescriptionService,
	})

	a.Router = r
	return nil
}
