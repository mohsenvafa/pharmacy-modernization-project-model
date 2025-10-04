/**
 * Dashboard Domain Components Registry
 * Registers all dashboard domain components
 */

import { registerComponent } from '@web/registry'
import { DashboardPageComponent } from '@domain/dashboard/ui/dashboard_page/dashboard_page.component'

// Register dashboard domain components
export function registerDashboardComponents() {
  registerComponent('dashboard.dashboard-page', () => new DashboardPageComponent())
}
