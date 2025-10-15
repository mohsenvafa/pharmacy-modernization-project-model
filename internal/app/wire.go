package app

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"pharmacy-modernization-project-model/internal/app/builder"
	"pharmacy-modernization-project-model/internal/integrations"
	"pharmacy-modernization-project-model/internal/platform/auth"
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
	if err := a.wireAuth(); err != nil {
		return err
	}

	// Create main MongoDB connection
	mongoConnMgr := a.wireMongodb()

	// Create primary cache (MongoDB or Memory)
	primaryCache := a.wireCache()

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

	// Initialize integrations layer (handles its own HTTP client internally)
	integration := integrations.New(integrations.Dependencies{
		Config: a.Cfg,
		Logger: logger.Base,
	})

	// Prescription Module
	prescriptionMod := prescriptionModule.Module(r, &prescriptionModule.ModuleDependencies{
		Logger:                       logger.Base,
		PharmacyClient:               integration.PharmacyClient,
		BillingClient:                integration.BillingClient,
		PrescriptionsMongoCollection: builder.GetPrescriptionsCollection(mongoConnMgr),
		CacheService:                 primaryCache,
	})

	// Patient Module
	var patientModDeps = &patientModule.ModuleDependencies{
		Logger:                   logger.Base,
		PrescriptionProvider:     prescriptionMod.PrescriptionService,
		PatientsMongoCollection:  builder.GetPatientsCollection(mongoConnMgr),
		AddressesMongoCollection: builder.GetAddressesCollection(mongoConnMgr),
		CacheService:             primaryCache,
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
