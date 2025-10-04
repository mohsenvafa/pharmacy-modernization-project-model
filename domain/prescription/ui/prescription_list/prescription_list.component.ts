/**
 * Prescription List Component
 * Handles interactions for the prescription list page
 */
export class PrescriptionListComponent {
  constructor() {
    this.init()
  }

  private init(): void {
    console.log('Prescription List component initialized')
    this.setupEventListeners()
  }

  private setupEventListeners(): void {
    // Listen for prescription list-specific events
    document.addEventListener('htmx:afterRequest', (event: any) => {
      if (event.detail.xhr.responseURL.includes('prescription') ||
          event.target?.closest('[data-component="prescription.prescription-list"]')) {
        this.onPrescriptionListLoaded()
      }
    })

    // Listen for prescription actions
    document.addEventListener('click', (event: Event) => {
      const target = event.target as HTMLElement

      if (target.matches('.prescription-action')) {
        this.handlePrescriptionAction(target)
      } else if (target.matches('.prescription-view')) {
        this.handleView(target)
      } else if (target.matches('.prescription-edit')) {
        this.handleEdit(target)
      } else if (target.matches('.prescription-delete')) {
        this.handleDelete(target)
      } else if (target.matches('.prescription-filter')) {
        this.handleFilter(target)
      }
    })
  }

  private onPrescriptionListLoaded(): void {
    console.log('Prescription list loaded')
    this.initializePrescriptionListInteractions()
  }

  private initializePrescriptionListInteractions(): void {
    // Add hover effects to prescription items
    document.querySelectorAll('[data-component="prescription.prescription-list"] .prescription-item').forEach(item => {
      item.addEventListener('mouseenter', () => {
        item.classList.add('hover:bg-base-200')
      })

      item.addEventListener('mouseleave', () => {
        item.classList.remove('hover:bg-base-200')
      })
    })

    // Initialize prescription status indicators
    this.initializeStatusIndicators()

    // Initialize filter controls
    this.initializeFilterControls()
  }

  private initializeStatusIndicators(): void {
    const statusElements = document.querySelectorAll('[data-component="prescription.prescription-list"] .prescription-status')
    statusElements.forEach(element => {
      const status = element.getAttribute('data-status')
      this.updateStatusIndicator(element as HTMLElement, status)
    })
  }

  private initializeFilterControls(): void {
    // Add filter dropdown handlers
    document.querySelectorAll('[data-component="prescription.prescription-list"] .filter-dropdown').forEach(dropdown => {
      dropdown.addEventListener('change', (event) => {
        const target = event.target as HTMLSelectElement
        this.handleFilterChange(target.value)
      })
    })

    // Add search input handlers
    document.querySelectorAll('[data-component="prescription.prescription-list"] .search-input').forEach(input => {
      input.addEventListener('input', (event) => {
        const target = event.target as HTMLInputElement
        this.handleSearch(target.value)
      })
    })
  }

  private handlePrescriptionAction(button: HTMLElement): void {
    const action = button.getAttribute('data-action')
    const prescriptionId = button.getAttribute('data-prescription-id')
    console.log(`Prescription action: ${action} for ID: ${prescriptionId}`)

    // Add visual feedback
    button.classList.add('loading')

    // Handle different actions
    switch (action) {
      case 'view':
        this.viewPrescription(prescriptionId)
        break
      case 'edit':
        this.editPrescription(prescriptionId)
        break
      case 'delete':
        this.deletePrescription(prescriptionId)
        break
      case 'refill':
        this.refillPrescription(prescriptionId)
        break
      case 'cancel':
        this.cancelPrescription(prescriptionId)
        break
      default:
        console.log(`Unknown prescription action: ${action}`)
    }

    // Remove loading state after a short delay
    setTimeout(() => {
      button.classList.remove('loading')
    }, 500)
  }

  private handleView(button: HTMLElement): void {
    const prescriptionId = button.getAttribute('data-prescription-id')
    console.log(`Viewing prescription: ${prescriptionId}`)

    // Add visual feedback
    button.classList.add('loading')

    this.viewPrescription(prescriptionId)

    // Remove loading state after a short delay
    setTimeout(() => {
      button.classList.remove('loading')
    }, 500)
  }

  private handleEdit(button: HTMLElement): void {
    const prescriptionId = button.getAttribute('data-prescription-id')
    console.log(`Editing prescription: ${prescriptionId}`)

    // Add visual feedback
    button.classList.add('loading')

    this.editPrescription(prescriptionId)

    // Remove loading state after a short delay
    setTimeout(() => {
      button.classList.remove('loading')
    }, 500)
  }

  private handleDelete(button: HTMLElement): void {
    const prescriptionId = button.getAttribute('data-prescription-id')
    console.log(`Deleting prescription: ${prescriptionId}`)

    // Add visual feedback
    button.classList.add('loading')

    // Show confirmation dialog
    if (confirm('Are you sure you want to delete this prescription? This action cannot be undone.')) {
      this.deletePrescription(prescriptionId)
    }

    // Remove loading state after a short delay
    setTimeout(() => {
      button.classList.remove('loading')
    }, 500)
  }

  private handleFilter(button: HTMLElement): void {
    const filterType = button.getAttribute('data-filter-type')
    const filterValue = button.getAttribute('data-filter-value')
    console.log(`Filtering prescriptions by ${filterType}: ${filterValue}`)

    // Add visual feedback
    button.classList.add('loading')

    this.applyFilter(filterType, filterValue)

    // Remove loading state after a short delay
    setTimeout(() => {
      button.classList.remove('loading')
    }, 500)
  }

  private updateStatusIndicator(element: HTMLElement, status: string | null): void {
    if (!status) return

    // Remove existing status classes
    element.classList.remove('badge-success', 'badge-warning', 'badge-error', 'badge-info')

    // Add appropriate status class
    switch (status.toLowerCase()) {
      case 'active':
      case 'filled':
      case 'completed':
        element.classList.add('badge-success')
        break
      case 'pending':
      case 'processing':
        element.classList.add('badge-warning')
        break
      case 'cancelled':
      case 'expired':
        element.classList.add('badge-error')
        break
      case 'new':
      case 'draft':
        element.classList.add('badge-info')
        break
      default:
        element.classList.add('badge-ghost')
    }
  }

  private viewPrescription(prescriptionId: string | null): void {
    if (prescriptionId) {
      console.log(`Viewing prescription: ${prescriptionId}`)
      // Add view functionality here (e.g., open modal, navigate to detail page)
    }
  }

  private editPrescription(prescriptionId: string | null): void {
    if (prescriptionId) {
      console.log(`Editing prescription: ${prescriptionId}`)
      // Add edit functionality here
    }
  }

  private deletePrescription(prescriptionId: string | null): void {
    if (prescriptionId) {
      console.log(`Deleting prescription: ${prescriptionId}`)
      // Add delete functionality here
    }
  }

  private refillPrescription(prescriptionId: string | null): void {
    if (prescriptionId) {
      console.log(`Refilling prescription: ${prescriptionId}`)
      // Add refill functionality here
    }
  }

  private cancelPrescription(prescriptionId: string | null): void {
    if (prescriptionId) {
      console.log(`Cancelling prescription: ${prescriptionId}`)
      // Add cancel functionality here
    }
  }

  private handleFilterChange(filterValue: string): void {
    console.log(`Filter changed to: ${filterValue}`)
    // Add filter functionality here
  }

  private handleSearch(searchTerm: string): void {
    console.log(`Searching for: ${searchTerm}`)
    // Add search functionality here
  }

  private applyFilter(filterType: string | null, filterValue: string | null): void {
    if (filterType && filterValue) {
      console.log(`Applying filter: ${filterType} = ${filterValue}`)
      // Add filter application logic here
    }
  }
}
