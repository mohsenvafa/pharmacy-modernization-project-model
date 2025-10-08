package app

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"pharmacy-modernization-project-model/internal/integrations"
	"pharmacy-modernization-project-model/internal/platform/auth"
	"pharmacy-modernization-project-model/internal/platform/httpx"
	"pharmacy-modernization-project-model/internal/platform/logging"
	"pharmacy-modernization-project-model/internal/platform/paths"

	dashboardModule "pharmacy-modernization-project-model/domain/dashboard"
	patientModule "pharmacy-modernization-project-model/domain/patient"
	prescriptionModule "pharmacy-modernization-project-model/domain/prescription"
	"pharmacy-modernization-project-model/internal/graphql"
)

func (a *App) wire() error {
	// Logger
	logger := logging.NewLogger(a.Cfg)
	a.Logger = logger

	// Initialize authentication system
	if err := auth.NewBuilder().
		WithJWTConfig(
			a.Cfg.Auth.JWT.Secret,
			a.Cfg.Auth.JWT.Issuer,
			a.Cfg.Auth.JWT.Audience,
			a.Cfg.Auth.JWT.Cookie.Name,
		).
		WithDevMode(a.Cfg.Auth.DevMode).
		WithEnvironment(a.Cfg.App.Env).
		WithLogger(logger.Base).
		Build(); err != nil {
		return err
	}

	// Shared HTTP client (for future external calls)
	_ = httpx.NewClient(a.Cfg)

	mongoConnMgr, err := CreateMongoDBConnection(a.Cfg, logger.Base)
	if err != nil {
		logger.Base.Error("Failed to create MongoDB connection", zap.Error(err))
		// Continue without MongoDB - will use memory repository as fallback
	}

	// Router & middleware
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(logging.CorrelationID())
	r.Use(logging.ZapRequestLogger(logger.Base))
	r.Use(middleware.Timeout(60 * time.Second))

	// Static assets
	r.Handle(paths.AssetsPath+"*", http.StripPrefix(paths.AssetsPath, http.FileServer(http.Dir("web/public"))))

	// Register dev mode endpoints (only when dev mode is enabled)
	auth.RegisterDevEndpoints(r, logger.Base)

	// Integrations
	integration := integrations.New(integrations.Dependencies{
		Config: a.Cfg,
		Logger: logger.Base,
	})

	// Prescription Module
	prescriptionMod := prescriptionModule.Module(r, &prescriptionModule.ModuleDependencies{
		Logger:                       logger.Base,
		PharmacyClient:               integration.Pharmacy,
		BillingClient:                integration.Billing,
		PrescriptionsMongoCollection: GetPrescriptionsCollection(mongoConnMgr),
	})

	// Patient Module
	var patientModDeps = &patientModule.ModuleDependencies{
		Logger:                   logger.Base,
		PrescriptionProvider:     prescriptionMod.PrescriptionService,
		PatientsMongoCollection:  GetPatientsCollection(mongoConnMgr),
		AddressesMongoCollection: GetAddressesCollection(mongoConnMgr),
	}

	patientMod := patientModule.Module(r, patientModDeps)

	// Dashboard Module
	dashboardMod := dashboardModule.Module(r, &dashboardModule.ModuleDependencies{
		PatientStats:      patientMod.PatientService,
		PrescriptionStats: prescriptionMod.PrescriptionService,
	})

	// GraphQL API
	graphql.MountGraphQL(r, &graphql.Dependencies{
		PatientService:      patientMod.PatientService,
		AddressService:      patientMod.AddressService,
		PrescriptionService: prescriptionMod.PrescriptionService,
		DashboardService:    dashboardMod.DashboardService,
		Logger:              logger.Base,
	})

	a.Router = r
	return nil
}
