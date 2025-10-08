package security

// Dashboard domain permissions
const (
	// Resource-based permissions
	DashboardPermissionView = "dashboard:view"
	PermissionAdminAll      = "admin:all"
)

// Common permission sets for reuse in routes
var (
	// DashboardAccess - user needs ANY of these permissions to view dashboard
	DashboardAccess = []string{DashboardPermissionView, PermissionAdminAll}
)
