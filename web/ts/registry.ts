/**
 * Component Registry for rxintake_scaffold
 * Simple component registry for lazy loading and initialization
 */

// Component registry
const componentRegistry = new Map<string, () => void>()

// Register a component
export function registerComponent(name: string, componentClass: () => void) {
  componentRegistry.set(name, componentClass)
}

// Initialize components based on DOM elements
export function initializeComponents() {
  // Find all elements with data-component attribute
  const elements = document.querySelectorAll('[data-component]')
  
  elements.forEach(element => {
    const componentName = element.getAttribute('data-component')
    if (componentName && componentRegistry.has(componentName)) {
      // Only initialize if not already initialized
      if (!element.hasAttribute('data-component-initialized')) {
        const componentClass = componentRegistry.get(componentName)
        if (componentClass) {
          componentClass()
          element.setAttribute('data-component-initialized', 'true')
        }
      }
    }
  })
}

// Auto-initialize on DOM changes (for HTMX)
document.addEventListener('htmx:afterRequest', initializeComponents)
document.addEventListener('DOMContentLoaded', initializeComponents)
