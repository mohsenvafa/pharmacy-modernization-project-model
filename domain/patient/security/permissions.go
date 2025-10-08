package security

import (
	commonsecurity "pharmacy-modernization-project-model/domain/common/security"
)

// Patient domain permissions
const (
	// Resource-based permissions
	PermissionRead   = commonsecurity.PatientPermissionRead
	PermissionWrite  = "patient:write"
	PermissionDelete = "patient:delete"
	PermissionExport = "patient:export"
)

// Common permission sets for reuse in routes
var (
	// ReadAccess - user needs ANY of these permissions to read patient data
	ReadAccess = commonsecurity.PatientReadAccess

	// WriteAccess - user needs ANY of these permissions to create/update patients
	WriteAccess = []string{PermissionWrite, "admin:all"}

	// ExportAccess - user needs ALL of these permissions to export patient data
	ExportAccess = []string{commonsecurity.PatientPermissionRead, PermissionExport}

	// DeleteAccess - user needs ALL of these permissions to delete patients
	DeleteAccess = []string{PermissionWrite, PermissionDelete}
)
