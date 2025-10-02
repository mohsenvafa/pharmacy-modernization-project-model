package paths

const (
	// Base path for prescription domain
	BasePath = "/prescriptions"

	// UI Routes
	ListPath = BasePath + "/"

	// API paths
	APIPath = "/api/v1/prescriptions"
)

// Helper functions for path generation
func PrescriptionListURL() string {
	return ListPath
}

func PrescriptionAPIURL() string {
	return APIPath
}
