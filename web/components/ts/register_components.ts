/**
 * Web Components Registry
 * Registers all web-level components
 */

import { registerComponent } from '@web/registry'
import { DataTableComponent } from '@components/elements/data_table.component'

// Register web components
export function registerWebComponents() {
  registerComponent('web.data-table', () => new DataTableComponent())
}
