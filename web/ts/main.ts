/**
 * Main TypeScript entry point for rxintake_scaffold
 * Registers all components and initializes the application
 */

import { initializeComponents } from '@web/registry'
import { registerWebComponents } from '@components/ts/register_components'
import { registerPatientComponents } from '@domain/patient/ui/ts/register_components'
import { registerDashboardComponents } from '@domain/dashboard/ui/ts/register_components'

// Register all components
export function registerAllComponents() {
  // Register web-level components
  registerWebComponents()
  
  // Register domain-specific components
  registerPatientComponents()
  registerDashboardComponents()
  
  console.log('All components registered')
}

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
  console.log('Initializing rxintake_scaffold application')
  
  // Register all components
  registerAllComponents()
  
  // Initialize components
  initializeComponents()
  
  console.log('Application initialized successfully')
})

// Re-initialize components after HTMX requests
document.addEventListener('htmx:afterRequest', () => {
  console.log('HTMX request completed, re-initializing components')
  initializeComponents()
})
