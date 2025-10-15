package service

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"

	commonmodel "pharmacy-modernization-project-model/domain/common/model"
	m "pharmacy-modernization-project-model/domain/prescription/contracts/model"
	repo "pharmacy-modernization-project-model/domain/prescription/repository"
	irisbilling "pharmacy-modernization-project-model/internal/integrations/iris_billing"
	irispharmacy "pharmacy-modernization-project-model/internal/integrations/iris_pharmacy"
	"pharmacy-modernization-project-model/internal/platform/cache"
)

type PrescriptionService interface {
	List(ctx context.Context, status string, limit, offset int) ([]m.Prescription, error)
	GetByID(ctx context.Context, id string) (m.Prescription, error)
	CountByStatus(ctx context.Context, status string) (int, error)
	PatientPrescriptionListByPatientID(ctx context.Context, patientID string) ([]commonmodel.PatientPrescription, error)
}

type svc struct {
	repo      repo.PrescriptionRepository
	cache     cache.Cache
	cacheKeys *CacheKeys
	log       *zap.Logger
	pharmacy  irispharmacy.PharmacyClient
	billing   irisbilling.BillingClient
}

func New(r repo.PrescriptionRepository, c cache.Cache, l *zap.Logger, pharmacy irispharmacy.PharmacyClient, billing irisbilling.BillingClient) PrescriptionService {
	return &svc{
		repo:      r,
		cache:     c,
		cacheKeys: NewCacheKeys(),
		log:       l,
		pharmacy:  pharmacy,
		billing:   billing,
	}
}

func (s *svc) List(ctx context.Context, status string, limit, offset int) ([]m.Prescription, error) {
	return s.repo.List(ctx, status, limit, offset)
}
func (s *svc) GetByID(ctx context.Context, id string) (m.Prescription, error) {
	cacheKey := s.cacheKeys.PrescriptionByID(id)

	// Try cache first
	if s.cache != nil {
		if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
			var prescription m.Prescription
			if err := json.Unmarshal(cached, &prescription); err == nil {
				if s.log != nil {
					s.log.Debug("Prescription retrieved from cache", zap.String("prescription_id", id))
				}
				return prescription, nil
			}
		}
	}

	prescription, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return m.Prescription{}, err
	}

	// Cache the result
	if s.cache != nil {
		if data, err := json.Marshal(prescription); err == nil {
			if err := s.cache.Set(ctx, cacheKey, data, 15*time.Minute); err != nil && s.log != nil {
				s.log.Warn("Failed to cache prescription", zap.String("prescription_id", id), zap.Error(err))
			}
		}
	}

	return prescription, nil
}

func (s *svc) CountByStatus(ctx context.Context, status string) (int, error) {
	cacheKey := s.cacheKeys.PrescriptionCountByStatus(status)

	// Try cache first
	if s.cache != nil {
		if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
			var count int
			if err := json.Unmarshal(cached, &count); err == nil {
				if s.log != nil {
					s.log.Debug("Prescription count retrieved from cache", zap.String("status", status))
				}
				return count, nil
			}
		}
	}

	count, err := s.repo.CountByStatus(ctx, status)
	if err != nil {
		return 0, err
	}

	// Cache the result with shorter TTL for counts
	if s.cache != nil {
		if data, err := json.Marshal(count); err == nil {
			if err := s.cache.Set(ctx, cacheKey, data, 5*time.Minute); err != nil && s.log != nil {
				s.log.Warn("Failed to cache prescription count", zap.String("status", status), zap.Error(err))
			}
		}
	}

	return count, nil
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
