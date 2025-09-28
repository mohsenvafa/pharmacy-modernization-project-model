package providers

import (
	"context"

	commonmodel "pharmacy-modernization-project-model/domain/common/model"
)

type PatientPrescriptionProvider interface {
	PatientPrescriptionListByPatientID(ctx context.Context, patientID string) ([]commonmodel.PatientPrescription, error)
}
