package ui

import (
	patientproviders "pharmacy-modernization-project-model/domain/patient/providers"
	patSvc "pharmacy-modernization-project-model/domain/patient/service"

	"go.uber.org/zap"
)

type UiDependencies struct {
	PatientSvc           patSvc.PatientService
	AddressSvc           patSvc.AddressService
	PrescriptionProvider patientproviders.PatientPrescriptionProvider
	InvoiceProvider      patientproviders.PatientInvoiceProvider
	Log                  *zap.Logger
}
