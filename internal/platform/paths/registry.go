package paths

const (
	// Main navigation paths
	DashboardPath     = "/"
	PatientsPath      = "/patients"
	PrescriptionsPath = "/prescriptions"
	PatientSearchPath = "/patients/search"

	// Static assets
	AssetsPath  = "/assets/"
	AppCSSPath  = "/assets/app.css"
	HtmxJSPath  = "/assets/vendor/htmx.min.js"
	ThemeJSPath = "/assets/vendor/theme-change.js"
	MainJSPath  = "/assets/js/dist/main.js"

	// API versions
	APIV1Prefix = "/api/v1"

	// GraphQL API
	GraphQLPath       = "/graphql"
	GraphQLPlayground = "/playground"
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
