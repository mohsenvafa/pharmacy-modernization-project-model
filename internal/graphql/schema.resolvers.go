package graphql

import (
	"context"
	"pharmacy-modernization-project-model/domain/patient/contracts/model"
	model1 "pharmacy-modernization-project-model/domain/prescription/contracts/model"
	"pharmacy-modernization-project-model/internal/graphql/generated"
)

// Empty is the resolver for the _empty field.
func (r *mutationResolver) Empty(ctx context.Context) (*string, error) {
	return nil, nil
}

// CreatePatient is the resolver for the createPatient field.
func (r *mutationResolver) CreatePatient(ctx context.Context, input generated.CreatePatientInput) (*model.Patient, error) {
	// Delegate to patient domain resolver
	return r.PatientResolver.CreatePatient(ctx, input)
}

// UpdatePatient is the resolver for the updatePatient field.
func (r *mutationResolver) UpdatePatient(ctx context.Context, id string, input generated.UpdatePatientInput) (*model.Patient, error) {
	// Delegate to patient domain resolver
	return r.PatientResolver.UpdatePatient(ctx, id, input)
}

// CreatePrescription is the resolver for the createPrescription field.
func (r *mutationResolver) CreatePrescription(ctx context.Context, input generated.CreatePrescriptionInput) (*model1.Prescription, error) {
	// Delegate to prescription domain resolver
	return r.PrescriptionResolver.CreatePrescription(ctx, input)
}

// UpdatePrescription is the resolver for the updatePrescription field.
func (r *mutationResolver) UpdatePrescription(ctx context.Context, id string, input generated.UpdatePrescriptionInput) (*model1.Prescription, error) {
	// Delegate to prescription domain resolver
	return r.PrescriptionResolver.UpdatePrescription(ctx, id, input)
}

// Addresses is the resolver for the addresses field.
func (r *patientResolver) Addresses(ctx context.Context, obj *model.Patient) ([]model.Address, error) {
	// Delegate to patient domain resolver
	return r.PatientResolver.Addresses(ctx, obj)
}

// Prescriptions is the resolver for the prescriptions field.
func (r *patientResolver) Prescriptions(ctx context.Context, obj *model.Patient) ([]model1.Prescription, error) {
	// Delegate to patient domain resolver
	return r.PatientResolver.Prescriptions(ctx, obj)
}

// Patient is the resolver for the patient field.
func (r *prescriptionResolver) Patient(ctx context.Context, obj *model1.Prescription) (*model.Patient, error) {
	// Delegate to prescription domain resolver
	return r.PrescriptionResolver.Patient(ctx, obj)
}

// Status is the resolver for the status field.
func (r *prescriptionResolver) Status(ctx context.Context, obj *model1.Prescription) (generated.PrescriptionStatus, error) {
	// Delegate to prescription domain resolver
	return r.PrescriptionResolver.Status(ctx, obj)
}

// Empty is the resolver for the _empty field.
func (r *queryResolver) Empty(ctx context.Context) (*string, error) {
	return nil, nil
}

// DashboardStats is the resolver for the dashboardStats field.
func (r *queryResolver) DashboardStats(ctx context.Context) (*generated.DashboardStats, error) {
	// Delegate to dashboard domain resolver
	return r.DashboardResolver.DashboardStats(ctx)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Patient returns generated.PatientResolver implementation.
func (r *Resolver) Patient() generated.PatientResolver { return &patientResolver{r} }

// Prescription returns generated.PrescriptionResolver implementation.
func (r *Resolver) Prescription() generated.PrescriptionResolver { return &prescriptionResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type patientResolver struct{ *Resolver }
type prescriptionResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
