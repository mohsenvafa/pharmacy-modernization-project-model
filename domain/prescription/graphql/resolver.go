package graphql

import (
	"context"

	"go.uber.org/zap"

	patientmodel "pharmacy-modernization-project-model/domain/patient/contracts/model"
	patientservice "pharmacy-modernization-project-model/domain/patient/service"
	"pharmacy-modernization-project-model/domain/prescription/contracts/model"
	prescriptionservice "pharmacy-modernization-project-model/domain/prescription/service"
	"pharmacy-modernization-project-model/internal/graphql/generated"
	"pharmacy-modernization-project-model/internal/graphql/validation"
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
	// Validate ID parameter using bind validation
	idValidation := validation.PrescriptionQueryValidation{ID: id}
	_, validationErrors := validation.ValidateGraphQLInput(idValidation)
	if validationErrors != nil {
		r.Logger.Error("Prescription ID validation failed",
			zap.Any("validation_errors", validationErrors.Errors))
		return nil, validationErrors
	}

	prescription, err := r.PrescriptionService.GetByID(ctx, id)
	if err != nil {
		r.Logger.Error("Failed to fetch prescription",
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
	// Validate query parameters using bind validation
	queryValidation := validation.PrescriptionsQueryValidation{}
	if status != nil {
		queryValidation.Status = status
	}
	if limit != nil {
		queryValidation.Limit = limit
	}
	if offset != nil {
		queryValidation.Offset = offset
	}

	_, validationErrors := validation.ValidateGraphQLInput(queryValidation)
	if validationErrors != nil {
		r.Logger.Error("Prescriptions query validation failed",
			zap.Any("validation_errors", validationErrors.Errors))
		return nil, validationErrors
	}

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
			zap.Int("limit", lim),
			zap.Int("offset", off),
			zap.Error(err))
		return nil, err
	}

	return prescriptions, nil
}

// ============================================================================
// Mutation Resolvers
// ============================================================================

// CreatePrescription resolves the createPrescription mutation
func (r *PrescriptionResolver) CreatePrescription(ctx context.Context, input generated.CreatePrescriptionInput) (*model.Prescription, error) {
	// Validate input using bind validation
	validationInput := validation.ConvertCreatePrescriptionInput(input)
	_, validationErrors := validation.ValidateGraphQLInput(validationInput)
	if validationErrors != nil {
		r.Logger.Error("Prescription creation validation failed",
			zap.Any("validation_errors", validationErrors.Errors))
		return nil, validationErrors
	}

	// Convert GraphQL status to domain status
	var domainStatus model.Status
	switch input.Status {
	case generated.PrescriptionStatusDraft:
		domainStatus = model.Draft
	case generated.PrescriptionStatusActive:
		domainStatus = model.Active
	case generated.PrescriptionStatusPaused:
		domainStatus = model.Paused
	case generated.PrescriptionStatusCompleted:
		domainStatus = model.Completed
	default:
		domainStatus = model.Draft
	}

	// Convert GraphQL input to domain model
	prescription := model.Prescription{
		PatientID: input.PatientID,
		Drug:      input.Drug,
		Dose:      input.Dose,
		Status:    domainStatus,
	}

	// Create prescription
	createdPrescription, err := r.PrescriptionService.Create(ctx, prescription)
	if err != nil {
		r.Logger.Error("Failed to create prescription",
			zap.Error(err))
		return nil, err
	}

	return &createdPrescription, nil
}

// UpdatePrescription resolves the updatePrescription mutation
func (r *PrescriptionResolver) UpdatePrescription(ctx context.Context, id string, input generated.UpdatePrescriptionInput) (*model.Prescription, error) {
	// Validate ID parameter
	idValidation := validation.PrescriptionQueryValidation{ID: id}
	_, validationErrors := validation.ValidateGraphQLInput(idValidation)
	if validationErrors != nil {
		r.Logger.Error("Prescription ID validation failed",
			zap.Any("validation_errors", validationErrors.Errors))
		return nil, validationErrors
	}

	// Validate input using bind validation
	validationInput := validation.ConvertUpdatePrescriptionInput(input)
	_, validationErrors = validation.ValidateGraphQLInput(validationInput)
	if validationErrors != nil {
		r.Logger.Error("Prescription update validation failed",
			zap.Any("validation_errors", validationErrors.Errors))
		return nil, validationErrors
	}

	// Get existing prescription first
	existingPrescription, err := r.PrescriptionService.GetByID(ctx, id)
	if err != nil {
		r.Logger.Error("Failed to get existing prescription",
			zap.Error(err))
		return nil, err
	}

	// Update fields if provided
	if input.Drug != nil {
		existingPrescription.Drug = *input.Drug
	}
	if input.Dose != nil {
		existingPrescription.Dose = *input.Dose
	}
	if input.Status != nil {
		// Convert GraphQL status to domain status
		switch *input.Status {
		case generated.PrescriptionStatusDraft:
			existingPrescription.Status = model.Draft
		case generated.PrescriptionStatusActive:
			existingPrescription.Status = model.Active
		case generated.PrescriptionStatusPaused:
			existingPrescription.Status = model.Paused
		case generated.PrescriptionStatusCompleted:
			existingPrescription.Status = model.Completed
		}
	}

	// Update prescription
	err = r.PrescriptionService.Update(ctx, existingPrescription)
	if err != nil {
		r.Logger.Error("Failed to update prescription",
			zap.Error(err))
		return nil, err
	}

	return &existingPrescription, nil
}

// Patient resolves the patient field on Prescription
func (r *PrescriptionResolver) Patient(ctx context.Context, obj *model.Prescription) (*patientmodel.Patient, error) {
	patient, err := r.PatientService.GetByID(ctx, obj.PatientID)
	if err != nil {
		r.Logger.Error("Failed to fetch patient for prescription",
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
