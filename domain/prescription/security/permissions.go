package security

import (
	commonsecurity "pharmacy-modernization-project-model/domain/common/security"
)

// Prescription domain permissions
const (
	PermissionRead     = commonsecurity.PrescriptionPermissionRead
	PermissionWrite    = "prescription:write"
	PermissionApprove  = "prescription:approve"
	PermissionDispense = "prescription:dispense"
	PermissionCancel   = "prescription:cancel"
)

// Common permission sets for reuse in routes
var (
	// ReadAccess - healthcare roles and admins can read prescriptions
	ReadAccess = commonsecurity.PrescriptionReadAccess

	// WriteAccess - only doctors or admins can create/edit prescriptions
	WriteAccess = []string{PermissionWrite, "doctor:role", "admin:all"}

	// ApproveAccess - needs ALL of these permissions to approve prescriptions
	ApproveAccess = []string{PermissionWrite, PermissionApprove}

	// DispenseAccess - only pharmacists or admins can dispense
	DispenseAccess = []string{PermissionDispense, "pharmacist:role", "admin:all"}

	// CancelAccess - needs both write and cancel permissions
	CancelAccess = []string{PermissionWrite, PermissionCancel}
)
