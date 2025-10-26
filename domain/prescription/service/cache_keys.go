package service

import (
	"fmt"
	"pharmacy-modernization-project-model/internal/platform/cache"
)

// CacheKeys provides centralized cache key management for the prescription domain
type CacheKeys struct{}

// NewCacheKeys creates a new cache keys instance
func NewCacheKeys() *CacheKeys {
	return &CacheKeys{}
}

// PrescriptionByID returns cache key for prescription by ID
func (k *CacheKeys) PrescriptionByID(id string) string {
	if !cache.ValidateID(id) {
		return "prescription:id:invalid"
	}
	sanitizedID := cache.SanitizeKey(id)
	return fmt.Sprintf("prescription:id:%s", sanitizedID)
}

// PrescriptionList returns cache key for prescription list query
func (k *CacheKeys) PrescriptionList(status string, limit, offset int) string {
	sanitizedStatus := cache.SanitizeKey(status)
	return fmt.Sprintf("prescription:list:%s:%d:%d", sanitizedStatus, limit, offset)
}

// PrescriptionCountByStatus returns cache key for prescription count by status
func (k *CacheKeys) PrescriptionCountByStatus(status string) string {
	sanitizedStatus := cache.SanitizeKey(status)
	return fmt.Sprintf("prescription:count:status:%s", sanitizedStatus)
}

// PrescriptionsByPatientID returns cache key for prescriptions by patient ID
func (k *CacheKeys) PrescriptionsByPatientID(patientID string) string {
	if !cache.ValidateID(patientID) {
		return "prescription:patient:invalid"
	}
	sanitizedPatientID := cache.SanitizeKey(patientID)
	return fmt.Sprintf("prescription:patient:%s", sanitizedPatientID)
}
