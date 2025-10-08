package security

// Patient domain permissions
const (
	// Resource-based permissions
	PatientPermissionRead = "patient:read"
)

// Common permission sets for reuse in routes
var (
	// ReadAccess - user needs ANY of these permissions to read patient data
	PatientReadAccess = []string{PatientPermissionRead, "admin:all"}
)
