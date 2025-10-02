package paths

const (
	// Base path for dashboard domain
	BasePath = "/"

	// UI Routes
	DashboardPath = BasePath
)

// Helper functions for path generation
func DashboardURL() string {
	return DashboardPath
}
