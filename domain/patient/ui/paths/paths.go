package paths

import "strings"

const (
	// Base path for patient domain
	BasePath = "/patients"

	// UI Routes
	ListPath      = BasePath + "/"
	DetailPath    = BasePath + "/{patientID}"
	ComponentPath = BasePath + "/components/patient-prescriptions-card"

	// API paths
	APIPath = "/api/v1/patients"
)

// Helper functions for path generation
func PatientDetailURL(patientID string) string {
	return strings.Replace(DetailPath, "{patientID}", patientID, 1)
}

func PatientListURL() string {
	return ListPath
}

func PatientComponentURL() string {
	return ComponentPath
}

func PatientAPIURL() string {
	return APIPath
}

func PatientAddressAPIURL(patientID string) string {
	return APIPath + "/" + patientID + "/addresses"
}
