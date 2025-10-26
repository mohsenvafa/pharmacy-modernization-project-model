package service

import (
	"fmt"
	"pharmacy-modernization-project-model/internal/platform/cache"
)

// CacheKeys provides centralized cache key management for the patient domain
type CacheKeys struct{}

// NewCacheKeys creates a new cache keys instance
func NewCacheKeys() *CacheKeys {
	return &CacheKeys{}
}

// PatientByID returns cache key for patient by ID
func (k *CacheKeys) PatientByID(id string) string {
	if !cache.ValidateID(id) {
		return "patient:id:invalid"
	}
	sanitizedID := cache.SanitizeKey(id)
	return fmt.Sprintf("patient:id:%s", sanitizedID)
}

// PatientList returns cache key for patient list query
func (k *CacheKeys) PatientList(query string, limit, offset int) string {
	sanitizedQuery := cache.SanitizeKey(query)
	return fmt.Sprintf("patient:list:%s:%d:%d", sanitizedQuery, limit, offset)
}

// PatientCount returns cache key for patient count
func (k *CacheKeys) PatientCount(query string) string {
	sanitizedQuery := cache.SanitizeKey(query)
	return fmt.Sprintf("patient:count:%s", sanitizedQuery)
}

// AddressByID returns cache key for address by ID
func (k *CacheKeys) AddressByID(patientID, addressID string) string {
	if !cache.ValidateID(patientID) || !cache.ValidateID(addressID) {
		return "address:id:invalid"
	}
	sanitizedPatientID := cache.SanitizeKey(patientID)
	sanitizedAddressID := cache.SanitizeKey(addressID)
	return fmt.Sprintf("address:id:%s:%s", sanitizedPatientID, sanitizedAddressID)
}

// AddressesByPatientID returns cache key for addresses by patient ID
func (k *CacheKeys) AddressesByPatientID(patientID string) string {
	if !cache.ValidateID(patientID) {
		return "address:patient:invalid"
	}
	sanitizedPatientID := cache.SanitizeKey(patientID)
	return fmt.Sprintf("address:patient:%s", sanitizedPatientID)
}
