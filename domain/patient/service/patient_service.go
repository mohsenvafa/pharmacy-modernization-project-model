package service

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"

	m "pharmacy-modernization-project-model/domain/patient/contracts/model"
	"pharmacy-modernization-project-model/domain/patient/contracts/request"
	repo "pharmacy-modernization-project-model/domain/patient/repository"
	"pharmacy-modernization-project-model/internal/platform/cache"
)

type PatientService interface {
	List(ctx context.Context, req request.PatientListQueryRequest) ([]m.Patient, error)
	GetByID(ctx context.Context, id string) (m.Patient, error)
	Create(ctx context.Context, patient m.Patient) (m.Patient, error)
	Update(ctx context.Context, patient m.Patient) error
	Count(ctx context.Context, req request.PatientListQueryRequest) (int, error)
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

func (s *patientSvc) Create(ctx context.Context, patient m.Patient) (m.Patient, error) {
	s.log.Info("Creating patient")

	// Set creation tracking fields
	now := time.Now()
	patient.CreatedAt = now

	// Create patient in repository
	createdPatient, err := s.repo.Create(ctx, patient)
	if err != nil {
		s.log.Error("Failed to create patient",
			zap.Error(err))
		return m.Patient{}, err
	}

	s.log.Info("Patient created successfully")

	return createdPatient, nil
}
func (s *patientSvc) List(ctx context.Context, req request.PatientListQueryRequest) ([]m.Patient, error) {
	return s.repo.List(ctx, req)
}

func (s *patientSvc) GetByID(ctx context.Context, id string) (m.Patient, error) {
	cacheKey := s.cacheKeys.PatientByID(id)

	// Try cache first
	if s.cache != nil {
		if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
			var patient m.Patient
			if err := json.Unmarshal(cached, &patient); err == nil {
				s.log.Debug("Patient retrieved from cache")
				return patient, nil
			}
		}
	}

	s.log.Info("Getting patient from repository")

	patient, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("Failed to get patient",
			zap.Error(err))
		return m.Patient{}, err
	}

	// Cache the result
	if s.cache != nil {
		if data, err := json.Marshal(patient); err == nil {
			if err := s.cache.Set(ctx, cacheKey, data, 30*time.Minute); err != nil {
				s.log.Warn("Failed to cache patient", zap.Error(err))
			}
		}
	}

	s.log.Info("Patient retrieved successfully")
	return patient, nil
}

func (s *patientSvc) Update(ctx context.Context, patient m.Patient) error {
	s.log.Info("Updating patient")

	// Set edit tracking fields
	now := time.Now()
	patient.EditTime = &now

	// Update patient in repository
	_, err := s.repo.Update(ctx, patient.ID, patient)
	if err != nil {
		s.log.Error("Failed to update patient",
			zap.Error(err))
		return err
	}

	// Invalidate cache for this patient
	if s.cache != nil {
		cacheKey := s.cacheKeys.PatientByID(patient.ID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			s.log.Warn("Failed to invalidate patient cache",
				zap.Error(err))
		}
	}

	s.log.Info("Patient updated successfully")
	return nil
}

func (s *patientSvc) Count(ctx context.Context, req request.PatientListQueryRequest) (int, error) {
	cacheKey := s.cacheKeys.PatientCount(req.PatientName)

	// Try cache first
	if s.cache != nil {
		if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
			var count int
			if err := json.Unmarshal(cached, &count); err == nil {
				s.log.Debug("Patient count retrieved from cache")
				return count, nil
			}
		}
	}

	count, err := s.repo.Count(ctx, req)
	if err != nil {
		return 0, err
	}

	// Cache the result with shorter TTL for counts
	if s.cache != nil {
		if data, err := json.Marshal(count); err == nil {
			if err := s.cache.Set(ctx, cacheKey, data, 5*time.Minute); err != nil {
				s.log.Warn("Failed to cache patient count", zap.Error(err))
			}
		}
	}

	return count, nil
}
