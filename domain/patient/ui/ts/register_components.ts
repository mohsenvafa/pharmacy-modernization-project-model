/**
 * Patient Domain Components Registry
 * Registers all patient domain components
 */

import { registerComponent } from '@web/registry'
import { PatientPrescriptionsComponent } from '@domain/patient/ui/components/patient_prescriptions/patient_prescriptions.component'
import { AddressListComponent } from '@domain/patient/ui/components/address_list/address_list.component'
import { PatientDetailComponent } from '@domain/patient/ui/patient_detail/patient_detail.component'

// Register patient domain components
export function registerPatientComponents() {
  registerComponent('patient.patient-prescriptions', () => new PatientPrescriptionsComponent())
  registerComponent('patient.address-list', () => new AddressListComponent())
  registerComponent('patient.patient-detail', () => new PatientDetailComponent())
}
