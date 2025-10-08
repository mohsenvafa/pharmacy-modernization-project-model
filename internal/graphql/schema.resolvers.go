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

// Patient is the resolver for the patient field.
func (r *queryResolver) Patient(ctx context.Context, id string) (*model.Patient, error) {
	// Delegate to patient domain resolver
	return r.PatientResolver.Patient(ctx, id)
}

// Patients is the resolver for the patients field.
func (r *queryResolver) Patients(ctx context.Context, query *string, limit *int, offset *int) ([]model.Patient, error) {
	// Delegate to patient domain resolver
	return r.PatientResolver.Patients(ctx, query, limit, offset)
}

// Prescription is the resolver for the prescription field.
func (r *queryResolver) Prescription(ctx context.Context, id string) (*model1.Prescription, error) {
	// Delegate to prescription domain resolver
	return r.PrescriptionResolver.Prescription(ctx, id)
}

// Prescriptions is the resolver for the prescriptions field.
func (r *queryResolver) Prescriptions(ctx context.Context, status *string, limit *int, offset *int) ([]model1.Prescription, error) {
	// Delegate to prescription domain resolver
	return r.PrescriptionResolver.Prescriptions(ctx, status, limit, offset)
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
