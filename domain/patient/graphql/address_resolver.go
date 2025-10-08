package graphql

import (
	"context"

	"go.uber.org/zap"

	"pharmacy-modernization-project-model/domain/patient/contracts/model"
	patientservice "pharmacy-modernization-project-model/domain/patient/service"
)

// AddressResolver handles address-specific GraphQL operations
// This demonstrates Phase 2: Splitting concerns into separate resolvers
type AddressResolver struct {
	AddressService patientservice.AddressService
	Logger         *zap.Logger
}

// NewAddressResolver creates a new address resolver
func NewAddressResolver(
	addressSvc patientservice.AddressService,
	logger *zap.Logger,
) *AddressResolver {
	return &AddressResolver{
		AddressService: addressSvc,
		Logger:         logger,
	}
}

// ============================================================================
// Field Resolvers
// ============================================================================

// Addresses resolves the addresses field on Patient
func (r *AddressResolver) Addresses(ctx context.Context, obj *model.Patient) ([]model.Address, error) {
	addresses, err := r.AddressService.GetByPatientID(ctx, obj.ID)
	if err != nil {
		r.Logger.Error("Failed to fetch addresses for patient",
			zap.String("patient_id", obj.ID),
			zap.Error(err))
		return []model.Address{}, nil // Return empty array on error to avoid null
	}
	return addresses, nil
}

// ValidateAddress could be a custom field resolver for address validation
// This demonstrates how sub-resources can have their own complex logic
func (r *AddressResolver) ValidateAddress(ctx context.Context, obj *model.Address) (bool, error) {
	// Example: Complex address validation logic
	// In a real app, this might call a validation service
	if obj.Zip == "" || obj.City == "" || obj.State == "" {
		return false, nil
	}
	return true, nil
}

// FormattedAddress demonstrates a computed field
func (r *AddressResolver) FormattedAddress(ctx context.Context, obj *model.Address) (string, error) {
	formatted := obj.Line1
	if obj.Line2 != "" {
		formatted += ", " + obj.Line2
	}
	formatted += ", " + obj.City + ", " + obj.State + " " + obj.Zip
	return formatted, nil
}
