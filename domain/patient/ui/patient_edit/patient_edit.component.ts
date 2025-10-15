/**
 * Patient Edit Component
 * Handles patient edit form interactions and validation
 */

export class PatientEditComponent {
  private form: HTMLFormElement | null = null

  constructor() {
    this.init()
  }

  private init(): void {
    console.log('Patient Edit component initialized')
    this.setupFormValidation()
    this.setupPhoneFormatting()
    this.setupFormSubmission()
  }

  /**
   * Setup client-side form validation
   */
  private setupFormValidation(): void {
    this.form = document.querySelector('form[method="POST"]')
    if (!this.form) return

    // Add real-time validation on input change
    const inputs = this.form.querySelectorAll('input[required]')
    inputs.forEach(input => {
      input.addEventListener('blur', () => this.validateField(input as HTMLInputElement))
      input.addEventListener('input', () => this.clearFieldError(input as HTMLInputElement))
    })
  }

  /**
   * Setup phone number formatting
   */
  private setupPhoneFormatting(): void {
    const phoneInput = document.getElementById('phone') as HTMLInputElement
    if (!phoneInput) return

    phoneInput.addEventListener('input', (e) => {
      let value = (e.target as HTMLInputElement).value.replace(/\D/g, '')
      
      if (value.length >= 6) {
        value = `(${value.slice(0, 3)}) ${value.slice(3, 6)}-${value.slice(6, 10)}`
      } else if (value.length >= 3) {
        value = `(${value.slice(0, 3)}) ${value.slice(3)}`
      } else if (value.length > 0) {
        value = `(${value}`
      }
      
      phoneInput.value = value
    })
  }

  /**
   * Setup form submission with client-side validation
   */
  private setupFormSubmission(): void {
    if (!this.form) return

    this.form.addEventListener('submit', (e) => {
      e.preventDefault()
      
      if (this.validateForm()) {
        // If validation passes, submit the form
        this.form?.submit()
      }
    })
  }

  /**
   * Validate a single field
   */
  private validateField(input: HTMLInputElement): boolean {
    const fieldName = input.name
    const value = input.value.trim()
    let isValid = true
    let errorMessage = ''

    // Clear previous error styling
    this.clearFieldError(input)

    // Validate based on field type
    switch (fieldName) {
      case 'name':
        if (value.length < 2) {
          errorMessage = 'Name must be at least 2 characters'
          isValid = false
        } else if (value.length > 100) {
          errorMessage = 'Name must be less than 100 characters'
          isValid = false
        }
        break

      case 'phone':
        const phoneDigits = value.replace(/\D/g, '')
        if (phoneDigits.length !== 10) {
          errorMessage = 'Please enter a valid 10-digit phone number'
          isValid = false
        }
        break

      case 'dob':
        const dob = new Date(value)
        const today = new Date()
        if (dob > today) {
          errorMessage = 'Date of birth cannot be in the future'
          isValid = false
        } else if (dob.getFullYear() < 1900) {
          errorMessage = 'Please enter a valid date of birth'
          isValid = false
        }
        break

      case 'state':
        if (value.length > 50) {
          errorMessage = 'State must be less than 50 characters'
          isValid = false
        }
        break
    }

    // Show error if validation failed
    if (!isValid) {
      this.showFieldError(input, errorMessage)
    }

    return isValid
  }

  /**
   * Validate the entire form
   */
  private validateForm(): boolean {
    if (!this.form) return false

    const inputs = this.form.querySelectorAll('input[required]')
    let isFormValid = true

    inputs.forEach(input => {
      const isValid = this.validateField(input as HTMLInputElement)
      if (!isValid) {
        isFormValid = false
      }
    })

    return isFormValid
  }

  /**
   * Show error message for a field
   */
  private showFieldError(input: HTMLInputElement, message: string): void {
    input.classList.add('input-error')
    
    // Find or create error label
    let errorLabel = input.parentElement?.querySelector('.field-error')
    if (!errorLabel) {
      errorLabel = document.createElement('label')
      errorLabel.className = 'label field-error'
      errorLabel.innerHTML = `<span class="label-text-alt text-error">${message}</span>`
      input.parentElement?.appendChild(errorLabel)
    } else {
      errorLabel.innerHTML = `<span class="label-text-alt text-error">${message}</span>`
    }
  }

  /**
   * Clear error message for a field
   */
  private clearFieldError(input: HTMLInputElement): void {
    input.classList.remove('input-error')
    
    const errorLabel = input.parentElement?.querySelector('.field-error')
    if (errorLabel) {
      errorLabel.remove()
    }
  }

  /**
   * Show loading state on form submission
   */
  public showLoading(): void {
    if (!this.form) return

    const submitButton = this.form.querySelector('button[type="submit"]') as HTMLButtonElement
    if (submitButton) {
      submitButton.disabled = true
      submitButton.innerHTML = `
        <span class="loading loading-spinner loading-sm"></span>
        Saving...
      `
    }
  }

  /**
   * Hide loading state
   */
  public hideLoading(): void {
    if (!this.form) return

    const submitButton = this.form.querySelector('button[type="submit"]') as HTMLButtonElement
    if (submitButton) {
      submitButton.disabled = false
      submitButton.innerHTML = `
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        Save Changes
      `
    }
  }
}

// Auto-initialize when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
  new PatientEditComponent()
})
