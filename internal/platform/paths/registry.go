package paths

const (
	// Main navigation paths
	DashboardPath     = "/"
	PatientsPath      = "/patients"
	PrescriptionsPath = "/prescriptions"

	// Static assets
	AssetsPath = "/assets/"

	// API versions
	APIV1Prefix = "/api/v1"
)

// NavigationPaths provides a structured way to access navigation paths
var NavigationPaths = struct {
	Dashboard     string
	Patients      string
	Prescriptions string
}{
	Dashboard:     DashboardPath,
	Patients:      PatientsPath,
	Prescriptions: PrescriptionsPath,
}
