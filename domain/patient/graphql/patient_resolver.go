package graphql

import (
	"context"

	"go.uber.org/zap"

	"pharmacy-modernization-project-model/domain/patient/contracts/model"
	"pharmacy-modernization-project-model/domain/patient/contracts/request"
	patientservice "pharmacy-modernization-project-model/domain/patient/service"
	model1 "pharmacy-modernization-project-model/domain/prescription/contracts/model"
	prescriptionservice "pharmacy-modernization-project-model/domain/prescription/service"
	"pharmacy-modernization-project-model/internal/platform/errors"
)

// PatientResolver handles Patient domain GraphQL operations
// Phase 2: We've separated address resolution into AddressResolver
// This keeps patient-specific logic focused and manageable
type PatientResolver struct {
	PatientService      patientservice.PatientService
	PrescriptionService prescriptionservice.PrescriptionService
	AddressResolver     *AddressResolver // Delegates address operations
	Logger              *zap.Logger
}

// NewPatientResolver creates a new patient resolver
func NewPatientResolver(
	patientSvc patientservice.PatientService,
	addressSvc patientservice.AddressService,
	prescriptionSvc prescriptionservice.PrescriptionService,
	logger *zap.Logger,
) *PatientResolver {
	return &PatientResolver{
		PatientService:      patientSvc,
		PrescriptionService: prescriptionSvc,
		AddressResolver:     NewAddressResolver(addressSvc, logger),
		Logger:              logger,
	}
}

// ============================================================================
// Query Resolvers
// ============================================================================

// Patient resolves the patient query
func (r *PatientResolver) Patient(ctx context.Context, id string) (*model.Patient, error) {
	patient, err := r.PatientService.GetByID(ctx, id)
	if err != nil {
		r.Logger.Error("Failed to fetch patient",
			zap.String("patient_id", id),
			zap.Error(err))

		if errors.IsNotFoundError(err) {
			return nil, nil // Return null for not found
		}
		return nil, err
	}
	return &patient, nil
}

// Patients resolves the patients query
func (r *PatientResolver) Patients(ctx context.Context, query *string, limit *int, offset *int) ([]model.Patient, error) {
	req := request.PatientListQueryRequest{
		Limit:  50, // default limit
		Offset: 0,
	}

	if query != nil {
		req.PatientName = *query
	}

	if limit != nil {
		req.Limit = *limit
	}

	if offset != nil {
		req.Offset = *offset
	}

	patients, err := r.PatientService.List(ctx, req)
	if err != nil {
		r.Logger.Error("Failed to list patients",
			zap.String("patientName", req.PatientName),
			zap.Int("limit", req.Limit),
			zap.Int("offset", req.Offset),
			zap.Error(err))
		return nil, err
	}

	return patients, nil
}

// ============================================================================
// Field Resolvers
// ============================================================================

// Addresses resolves the addresses field on Patient
// Phase 2: Delegates to AddressResolver for better separation of concerns
func (r *PatientResolver) Addresses(ctx context.Context, obj *model.Patient) ([]model.Address, error) {
	return r.AddressResolver.Addresses(ctx, obj)
}

// Prescriptions resolves the prescriptions field on Patient
func (r *PatientResolver) Prescriptions(ctx context.Context, obj *model.Patient) ([]model1.Prescription, error) {
	prescriptions, err := r.PrescriptionService.List(ctx, "", 100, 0)
	if err != nil {
		r.Logger.Error("Failed to fetch prescriptions for patient",
			zap.String("patient_id", obj.ID),
			zap.Error(err))
		return []model1.Prescription{}, nil
	}

	// Filter by patient ID
	var patientPrescriptions []model1.Prescription
	for _, p := range prescriptions {
		if p.PatientID == obj.ID {
			patientPrescriptions = append(patientPrescriptions, p)
		}
	}
	return patientPrescriptions, nil
}
