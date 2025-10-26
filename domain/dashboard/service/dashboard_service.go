package service

import (
	"context"

	model "pharmacy-modernization-project-model/domain/dashboard/contracts/model"
	"pharmacy-modernization-project-model/domain/dashboard/providers"
	"pharmacy-modernization-project-model/domain/patient/contracts/request"
)

type IDashboardService interface {
	Summary(ctx context.Context) (model.DashboardSummary, error)
}

type dashboardService struct {
	patients      providers.PatientStatsProvider
	prescriptions providers.PrescriptionStatsProvider
}

func New(patients providers.PatientStatsProvider, prescriptions providers.PrescriptionStatsProvider) IDashboardService {
	return &dashboardService{patients: patients, prescriptions: prescriptions}
}

func (s *dashboardService) Summary(ctx context.Context) (model.DashboardSummary, error) {
	req := request.PatientListQueryRequest{
		Limit:  0, // Get all patients for count
		Offset: 0,
	}
	total, err := s.patients.Count(ctx, req)
	if err != nil {
		return model.DashboardSummary{}, err
	}

	active, err := s.prescriptions.CountByStatus(ctx, "Active")
	if err != nil {
		return model.DashboardSummary{}, err
	}

	return model.DashboardSummary{TotalPatients: total, ActivePrescriptions: active}, nil
}
