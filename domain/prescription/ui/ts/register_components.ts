/**
 * Prescription Domain Components Registry
 * Registers all prescription domain components
 */

import { registerComponent } from '@web/registry'
import { PrescriptionListComponent } from '@domain/prescription/ui/prescription_list/prescription_list.component'

// Register prescription domain components
export function registerPrescriptionComponents() {
  registerComponent('prescription.prescription-list', () => new PrescriptionListComponent())
}
