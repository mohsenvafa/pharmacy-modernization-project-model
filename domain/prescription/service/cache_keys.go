package service

import "fmt"

// CacheKeys provides centralized cache key management for the prescription domain
type CacheKeys struct{}

// NewCacheKeys creates a new cache keys instance
func NewCacheKeys() *CacheKeys {
	return &CacheKeys{}
}

// PrescriptionByID returns cache key for prescription by ID
func (k *CacheKeys) PrescriptionByID(id string) string {
	return fmt.Sprintf("prescription:id:%s", id)
}

// PrescriptionList returns cache key for prescription list query
func (k *CacheKeys) PrescriptionList(status string, limit, offset int) string {
	return fmt.Sprintf("prescription:list:%s:%d:%d", status, limit, offset)
}

// PrescriptionCountByStatus returns cache key for prescription count by status
func (k *CacheKeys) PrescriptionCountByStatus(status string) string {
	return fmt.Sprintf("prescription:count:status:%s", status)
}

// PrescriptionsByPatientID returns cache key for prescriptions by patient ID
func (k *CacheKeys) PrescriptionsByPatientID(patientID string) string {
	return fmt.Sprintf("prescription:patient:%s", patientID)
}
