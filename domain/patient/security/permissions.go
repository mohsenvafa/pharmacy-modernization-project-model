package security

// Patient domain permissions
const (
	// Resource-based permissions
	PermissionRead   = "patient:read"
	PermissionWrite  = "patient:write"
	PermissionDelete = "patient:delete"
	PermissionExport = "patient:export"
)

// Common permission sets for reuse in routes
var (
	// ReadAccess - user needs ANY of these permissions to read patient data
	ReadAccess = []string{PermissionRead, "admin:all"}

	// WriteAccess - user needs ANY of these permissions to create/update patients
	WriteAccess = []string{PermissionWrite, "admin:all"}

	// ExportAccess - user needs ALL of these permissions to export patient data
	ExportAccess = []string{PermissionRead, PermissionExport}

	// DeleteAccess - user needs ALL of these permissions to delete patients
	DeleteAccess = []string{PermissionWrite, PermissionDelete}
)
