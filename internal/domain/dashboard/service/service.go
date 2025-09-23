package service

import (
	"context"

	"github.com/pharmacy-modernization-project-model/internal/domain/dashboard/providers"
)

type Summary struct {
	TotalPatients       int
	ActivePrescriptions int
}

type Service interface {
	Summary(ctx context.Context) (Summary, error)
}

type svc struct {
	patients      providers.PatientStatsProvider
	prescriptions providers.PrescriptionStatsProvider
}

func New(patients providers.PatientStatsProvider, prescriptions providers.PrescriptionStatsProvider) Service {
	return &svc{patients: patients, prescriptions: prescriptions}
}

func (s *svc) Summary(ctx context.Context) (Summary, error) {
	total, err := s.patients.Count(ctx, "")
	if err != nil {
		return Summary{}, err
	}

	active, err := s.prescriptions.CountByStatus(ctx, "Active")
	if err != nil {
		return Summary{}, err
	}

	return Summary{TotalPatients: total, ActivePrescriptions: active}, nil
}
