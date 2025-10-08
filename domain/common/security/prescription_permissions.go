package security

// Prescription domain permissions
const (
	// Resource-based permissions
	PrescriptionPermissionRead = "prescription:read"
)

// Common permission sets for reuse in routes
var (
	// ReadAccess - user needs ANY of these permissions to read patient data
	PrescriptionReadAccess = []string{PrescriptionPermissionRead, "admin:all"}
)
