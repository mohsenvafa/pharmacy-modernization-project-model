package service

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"

	m "pharmacy-modernization-project-model/domain/patient/contracts/model"
	repo "pharmacy-modernization-project-model/domain/patient/repository"
	"pharmacy-modernization-project-model/internal/platform/cache"
)

type PatientService interface {
	List(ctx context.Context, query string, limit, offset int) ([]m.Patient, error)
	GetByID(ctx context.Context, id string) (m.Patient, error)
	Count(ctx context.Context, query string) (int, error)
}

type patientSvc struct {
	repo      repo.PatientRepository
	cache     cache.Cache
	cacheKeys *CacheKeys
	log       *zap.Logger
}

func New(r repo.PatientRepository, c cache.Cache, l *zap.Logger) PatientService {
	return &patientSvc{
		repo:      r,
		cache:     c,
		cacheKeys: NewCacheKeys(),
		log:       l,
	}
}

func (s *patientSvc) List(ctx context.Context, query string, limit, offset int) ([]m.Patient, error) {
	return s.repo.List(ctx, query, limit, offset)
}
func (s *patientSvc) GetByID(ctx context.Context, id string) (m.Patient, error) {
	cacheKey := s.cacheKeys.PatientByID(id)

	// Try cache first
	if s.cache != nil {
		if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
			var patient m.Patient
			if err := json.Unmarshal(cached, &patient); err == nil {
				s.log.Debug("Patient retrieved from cache", zap.String("patient_id", id))
				return patient, nil
			}
		}
	}

	s.log.Info("Getting patient from repository", zap.String("patient_id", id))

	patient, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("Failed to get patient",
			zap.String("patient_id", id),
			zap.Error(err))
		return m.Patient{}, err
	}

	// Cache the result
	if s.cache != nil {
		if data, err := json.Marshal(patient); err == nil {
			if err := s.cache.Set(ctx, cacheKey, data, 30*time.Minute); err != nil {
				s.log.Warn("Failed to cache patient", zap.String("patient_id", id), zap.Error(err))
			}
		}
	}

	s.log.Info("Patient retrieved successfully", zap.String("patient_id", id))
	return patient, nil
}

func (s *patientSvc) Count(ctx context.Context, query string) (int, error) {
	cacheKey := s.cacheKeys.PatientCount(query)

	// Try cache first
	if s.cache != nil {
		if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
			var count int
			if err := json.Unmarshal(cached, &count); err == nil {
				s.log.Debug("Patient count retrieved from cache", zap.String("query", query))
				return count, nil
			}
		}
	}

	count, err := s.repo.Count(ctx, query)
	if err != nil {
		return 0, err
	}

	// Cache the result with shorter TTL for counts
	if s.cache != nil {
		if data, err := json.Marshal(count); err == nil {
			if err := s.cache.Set(ctx, cacheKey, data, 5*time.Minute); err != nil {
				s.log.Warn("Failed to cache patient count", zap.String("query", query), zap.Error(err))
			}
		}
	}

	return count, nil
}
