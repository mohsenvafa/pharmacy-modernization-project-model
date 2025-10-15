package paths

import "strings"

const (
	// Base path for patient domain
	BasePath = "/patients"

	// UI Routes
	ListPath      = BasePath + "/"
	SearchPath    = BasePath + "/search"
	DetailPath    = BasePath + "/{patientID}"
	EditPath      = BasePath + "/{patientID}/edit"
	ComponentPath = BasePath + "/components/patient-prescriptions-card"

	// Relative route patterns (for use within Route() blocks)
	ListRoute                             = "/"
	SearchRoute                           = "/search"
	DetailRoute                           = "/{patientID}"
	EditRoute                             = "/{patientID}/edit"
	PatientPrescriptionCardComponentRoute = "/components/patient-prescriptions-card"

	// API paths
	APIPath = "/api/v1/patients"

	// Address sub-routes
	AddressSubRoute = "/{patientID}/addresses"
)

// Helper functions for path generation with parameters
func PatientDetailURL(patientID string) string {
	return strings.Replace(DetailPath, "{patientID}", patientID, 1)
}

func PatientEditURL(patientID string) string {
	return strings.Replace(EditPath, "{patientID}", patientID, 1)
}

func PatientAddressAPIURL(patientID string) string {
	return APIPath + "/" + patientID + "/addresses"
}

func PatientListURL() string {
	return ListPath
}
