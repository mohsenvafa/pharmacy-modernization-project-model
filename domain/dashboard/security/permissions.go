package security

// Dashboard domain permissions
const (
	PermissionView      = "dashboard:view"
	PermissionAnalytics = "dashboard:analytics"
	PermissionReports   = "dashboard:reports"
)

// Common permission sets for reuse in routes
var (
	// ViewAccess - any authenticated user with dashboard permission can view
	ViewAccess = []string{PermissionView, "admin:all"}

	// AnalyticsAccess - needs dashboard view and analytics permissions
	AnalyticsAccess = []string{PermissionView, PermissionAnalytics}

	// ReportsAccess - needs dashboard view and reports permissions
	ReportsAccess = []string{PermissionView, PermissionReports}
)
