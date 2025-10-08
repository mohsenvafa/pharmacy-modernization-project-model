package graphql

import (
	"context"

	"go.uber.org/zap"

	patientmodel "pharmacy-modernization-project-model/domain/patient/contracts/model"
	patientservice "pharmacy-modernization-project-model/domain/patient/service"
	"pharmacy-modernization-project-model/domain/prescription/contracts/model"
	prescriptionservice "pharmacy-modernization-project-model/domain/prescription/service"
	"pharmacy-modernization-project-model/internal/graphql/generated"
	"pharmacy-modernization-project-model/internal/platform/errors"
)

// PrescriptionResolver handles all Prescription domain GraphQL operations
type PrescriptionResolver struct {
	PrescriptionService prescriptionservice.PrescriptionService
	PatientService      patientservice.PatientService
	Logger              *zap.Logger
}

// NewPrescriptionResolver creates a new prescription resolver
func NewPrescriptionResolver(
	prescriptionSvc prescriptionservice.PrescriptionService,
	patientSvc patientservice.PatientService,
	logger *zap.Logger,
) *PrescriptionResolver {
	return &PrescriptionResolver{
		PrescriptionService: prescriptionSvc,
		PatientService:      patientSvc,
		Logger:              logger,
	}
}

// ============================================================================
// Query Resolvers
// ============================================================================

// Prescription resolves the prescription query
func (r *PrescriptionResolver) Prescription(ctx context.Context, id string) (*model.Prescription, error) {
	prescription, err := r.PrescriptionService.GetByID(ctx, id)
	if err != nil {
		r.Logger.Error("Failed to fetch prescription",
			zap.String("prescription_id", id),
			zap.Error(err))

		if errors.IsNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &prescription, nil
}

// Prescriptions resolves the prescriptions query
func (r *PrescriptionResolver) Prescriptions(ctx context.Context, status *string, limit *int, offset *int) ([]model.Prescription, error) {
	statusFilter := ""
	if status != nil {
		statusFilter = *status
	}

	lim := 50
	if limit != nil {
		lim = *limit
	}

	off := 0
	if offset != nil {
		off = *offset
	}

	prescriptions, err := r.PrescriptionService.List(ctx, statusFilter, lim, off)
	if err != nil {
		r.Logger.Error("Failed to list prescriptions",
			zap.String("status", statusFilter),
			zap.Int("limit", lim),
			zap.Int("offset", off),
			zap.Error(err))
		return nil, err
	}

	return prescriptions, nil
}

// ============================================================================
// Field Resolvers
// ============================================================================

// Patient resolves the patient field on Prescription
func (r *PrescriptionResolver) Patient(ctx context.Context, obj *model.Prescription) (*patientmodel.Patient, error) {
	patient, err := r.PatientService.GetByID(ctx, obj.PatientID)
	if err != nil {
		r.Logger.Error("Failed to fetch patient for prescription",
			zap.String("prescription_id", obj.ID),
			zap.String("patient_id", obj.PatientID),
			zap.Error(err))
		return nil, err
	}
	return &patient, nil
}

// Status resolves the status field on Prescription (converts domain enum to GraphQL enum)
func (r *PrescriptionResolver) Status(ctx context.Context, obj *model.Prescription) (generated.PrescriptionStatus, error) {
	// Convert domain status to GraphQL enum
	switch obj.Status {
	case model.Draft:
		return generated.PrescriptionStatusDraft, nil
	case model.Active:
		return generated.PrescriptionStatusActive, nil
	case model.Paused:
		return generated.PrescriptionStatusPaused, nil
	case model.Completed:
		return generated.PrescriptionStatusCompleted, nil
	default:
		return generated.PrescriptionStatusDraft, nil
	}
}
