package service

import (
	"context"

	commonmodel "pharmacy-modernization-project-model/domain/common/model"
	m "pharmacy-modernization-project-model/domain/prescription/contracts/model"
	repo "pharmacy-modernization-project-model/domain/prescription/repository"
	irisbilling "pharmacy-modernization-project-model/internal/integrations/iris_billing"
	irispharmacy "pharmacy-modernization-project-model/internal/integrations/iris_pharmacy"

	"go.uber.org/zap"
)

type PrescriptionService interface {
	List(ctx context.Context, status string, limit, offset int) ([]m.Prescription, error)
	GetByID(ctx context.Context, id string) (m.Prescription, error)
	CountByStatus(ctx context.Context, status string) (int, error)
	PatientPrescriptionListByPatientID(ctx context.Context, patientID string) ([]commonmodel.PatientPrescription, error)
}

type svc struct {
	repo     repo.PrescriptionRepository
	log      *zap.Logger
	pharmacy irispharmacy.PharmacyClient
	billing  irisbilling.BillingClient
}

func New(r repo.PrescriptionRepository, l *zap.Logger, pharmacy irispharmacy.PharmacyClient, billing irisbilling.BillingClient) PrescriptionService {
	return &svc{repo: r, log: l, pharmacy: pharmacy, billing: billing}
}

func (s *svc) List(ctx context.Context, status string, limit, offset int) ([]m.Prescription, error) {
	return s.repo.List(ctx, status, limit, offset)
}
func (s *svc) GetByID(ctx context.Context, id string) (m.Prescription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *svc) CountByStatus(ctx context.Context, status string) (int, error) {
	return s.repo.CountByStatus(ctx, status)
}

func (s *svc) PatientPrescriptionListByPatientID(ctx context.Context, patientID string) ([]commonmodel.PatientPrescription, error) {
	items, err := s.repo.ListByPatientID(ctx, patientID)
	if err != nil {
		if s.log != nil {
			s.log.Error("failed to list prescriptions by patient", zap.Error(err), zap.String("patient_id", patientID))
		}
		return nil, err
	}
	result := make([]commonmodel.PatientPrescription, 0, len(items))
	for _, item := range items {
		result = append(result, commonmodel.PatientPrescription{
			ID:        item.ID,
			Drug:      item.Drug,
			Dose:      item.Dose,
			Status:    string(item.Status),
			CreatedAt: item.CreatedAt,
		})
	}
	return result, nil
}
