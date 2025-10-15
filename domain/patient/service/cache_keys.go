package service

import "fmt"

// CacheKeys provides centralized cache key management for the patient domain
type CacheKeys struct{}

// NewCacheKeys creates a new cache keys instance
func NewCacheKeys() *CacheKeys {
	return &CacheKeys{}
}

// PatientByID returns cache key for patient by ID
func (k *CacheKeys) PatientByID(id string) string {
	return fmt.Sprintf("patient:id:%s", id)
}

// PatientList returns cache key for patient list query
func (k *CacheKeys) PatientList(query string, limit, offset int) string {
	return fmt.Sprintf("patient:list:%s:%d:%d", query, limit, offset)
}

// PatientCount returns cache key for patient count
func (k *CacheKeys) PatientCount(query string) string {
	return fmt.Sprintf("patient:count:%s", query)
}

// AddressByID returns cache key for address by ID
func (k *CacheKeys) AddressByID(patientID, addressID string) string {
	return fmt.Sprintf("address:id:%s:%s", patientID, addressID)
}

// AddressesByPatientID returns cache key for addresses by patient ID
func (k *CacheKeys) AddressesByPatientID(patientID string) string {
	return fmt.Sprintf("address:patient:%s", patientID)
}
