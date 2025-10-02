package paths

import "strings"

const (
	// Base path for patient domain
	BasePath = "/patients"

	// UI Routes
	ListPath      = BasePath + "/"
	DetailPath    = BasePath + "/{patientID}"
	ComponentPath = BasePath + "/components/patient-prescriptions-card"

	// Relative route patterns (for use within Route() blocks)
	ListRoute                             = "/"
	DetailRoute                           = "/{patientID}"
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

func PatientAddressAPIURL(patientID string) string {
	return APIPath + "/" + patientID + "/addresses"
}
