# Adding TypeScript Components

This guide explains how to add TypeScript components to the rxintake_scaffold project following the established patterns.

## Overview

Each templ component can have an associated TypeScript file that handles client-side interactions. The TypeScript components are organized by domain and follow a consistent naming convention.

## File Structure

```
domain/
├── [domain_name]/
│   └── ui/
│       ├── [component_name]/
│       │   ├── [component_name].component.templ
│       │   ├── [component_name].component.ts          # ← TypeScript component
│       │   └── [component_name].component.handler.go
│       └── ts/
│           └── register_components.ts                  # ← Component registry
```

## Step-by-Step Guide

### 1. Create the TypeScript Component File

Create a new TypeScript file in the same directory as your templ component:

**File**: `domain/[domain]/ui/[component]/[component].component.ts`

```typescript
/**
 * [Component Name] Component
 * Handles interactions for the [component description]
 */
export class [ComponentName]Component {
  constructor() {
    this.init()
  }

  private init(): void {
    console.log('[Component Name] component initialized')
    this.setupEventListeners()
  }

  private setupEventListeners(): void {
    // Add your event listeners here
    document.addEventListener('htmx:afterRequest', (event: any) => {
      if (event.detail.xhr.responseURL.includes('[domain]') ||
          event.target?.closest('[data-component="[domain].[component]"]')) {
        this.onComponentLoaded()
      }
    })
  }

  private onComponentLoaded(): void {
    console.log('[Component name] loaded')
    // Add initialization logic here
  }
}
```

### 2. Add data-component Attribute to Template

Update your templ file to include the `data-component` attribute:

**File**: `domain/[domain]/ui/[component]/[component].component.templ`

```go
templ [componentName](pageParam [ComponentName]Param) {
	<div class="..." data-component="[domain].[component]">
		<!-- Your template content -->
	</div>
}
```

**Example**:
```go
templ patientDetail(pageParam PatientDetailPageParam) {
	<div class="flex flex-col space-y-6" data-component="patient.patient-detail">
		<!-- Patient detail content -->
	</div>
}
```

### 3. Register the Component

#### Option A: Create/Update Domain Registry

**File**: `domain/[domain]/ui/ts/register_components.ts`

```typescript
/**
 * [Domain] Domain Components Registry
 * Registers all [domain] domain components
 */

import { registerComponent } from '@web/registry'
import { [ComponentName]Component } from '@domain/[domain]/ui/[component]/[component].component'

// Register [domain] domain components
export function register[Domain]Components() {
  registerComponent('[domain].[component]', () => new [ComponentName]Component())
}
```

#### Option B: Update Existing Domain Registry

If the domain registry already exists, add your component:

```typescript
import { [ComponentName]Component } from '@domain/[domain]/ui/[component]/[component].component'

export function register[Domain]Components() {
  // ... existing components
  registerComponent('[domain].[component]', () => new [ComponentName]Component())
}
```

### 4. Register in Main Application

Update the main TypeScript entry point:

**File**: `web/ts/main.ts`

```typescript
import { register[Domain]Components } from '@domain/[domain]/ui/ts/register_components'

export function registerAllComponents() {
  // ... existing registrations
  register[Domain]Components()
}
```

### 5. Build and Test

```bash
cd web && npm run build
```

## Naming Conventions

### File Names
- **Template**: `[component].component.templ`
- **TypeScript**: `[component].component.ts`
- **Handler**: `[component].component.handler.go`
- **Registry**: `register_components.ts`

### Component Names
- **Class**: `[ComponentName]Component` (PascalCase)
- **Registry Key**: `[domain].[component]` (lowercase with dots)
- **Data Attribute**: `data-component="[domain].[component]"`

### Examples
- **Domain**: `patient`
- **Component**: `patient-detail`
- **Class**: `PatientDetailComponent`
- **Registry Key**: `patient.patient-detail`
- **Data Attribute**: `data-component="patient.patient-detail"`

## Component Features

### Standard Features
- **Console Logging**: Initialization and load events
- **HTMX Integration**: Automatic re-initialization after HTMX requests
- **Event Handling**: Click handlers for component actions
- **Visual Feedback**: Loading states and hover effects

### Common Patterns

#### Event Listeners
```typescript
private setupEventListeners(): void {
  // HTMX integration
  document.addEventListener('htmx:afterRequest', (event: any) => {
    if (event.detail.xhr.responseURL.includes('patient') ||
        event.target?.closest('[data-component="patient.patient-detail"]')) {
      this.onComponentLoaded()
    }
  })

  // Click handlers
  document.addEventListener('click', (event: Event) => {
    const target = event.target as HTMLElement
    
    if (target.matches('.action-button')) {
      this.handleAction(target)
    }
  })
}
```

#### Action Handling
```typescript
private handleAction(button: HTMLElement): void {
  const action = button.getAttribute('data-action')
  const id = button.getAttribute('data-id')
  
  // Visual feedback
  button.classList.add('loading')
  
  // Handle action
  switch (action) {
    case 'edit':
      this.editItem(id)
      break
    case 'delete':
      this.deleteItem(id)
      break
  }
  
  // Remove loading state
  setTimeout(() => {
    button.classList.remove('loading')
  }, 500)
}
```

#### Status Indicators
```typescript
private updateStatusIndicator(element: HTMLElement, status: string | null): void {
  if (!status) return
  
  // Remove existing classes
  element.classList.remove('badge-success', 'badge-warning', 'badge-error')
  
  // Add appropriate class
  switch (status.toLowerCase()) {
    case 'active':
      element.classList.add('badge-success')
      break
    case 'pending':
      element.classList.add('badge-warning')
      break
    case 'inactive':
      element.classList.add('badge-error')
      break
  }
}
```

## Path Aliases

The project uses TypeScript path aliases for clean imports:

```typescript
// In tsconfig.json
"paths": {
  "@web/*": ["ts/*"],
  "@components/*": ["components/*"],
  "@domain/*": ["../domain/*"]
}

// Usage in components
import { registerComponent } from '@web/registry'
import { MyComponent } from '@domain/patient/ui/component/component.component'
```

## Testing

### Console Logs
After building, check that your component logs appear:
```bash
curl -s http://localhost:8080/assets/js/dist/main.js | grep -o "Your Component component initialized"
```

### Component Registration
Verify the component is registered:
```bash
curl -s http://localhost:8080/your-page | grep -o 'data-component="your.domain.component"'
```

## Examples

### Complete Example: Patient Detail Component

**File**: `domain/patient/ui/patient_detail/patient_detail.component.ts`
```typescript
export class PatientDetailComponent {
  constructor() {
    this.init()
  }

  private init(): void {
    console.log('Patient Detail component initialized')
  }
}
```

**File**: `domain/patient/ui/ts/register_components.ts`
```typescript
import { registerComponent } from '@web/registry'
import { PatientDetailComponent } from '@domain/patient/ui/patient_detail/patient_detail.component'

export function registerPatientComponents() {
  registerComponent('patient.patient-detail', () => new PatientDetailComponent())
}
```

**Template**: `domain/patient/ui/patient_detail/patient_detail_page.component.templ`
```go
templ patientDetail(pageParam PatientDetailPageParam) {
	<div class="flex flex-col space-y-6" data-component="patient.patient-detail">
		<!-- Content -->
	</div>
}
```

## Troubleshooting

### Component Not Initializing
1. Check that `data-component` attribute is present in template
2. Verify component is registered in domain registry
3. Ensure domain registry is called in `main.ts`
4. Check browser console for errors

### Console Logs Not Appearing
1. Verify TypeScript compiled successfully
2. Check that component class is instantiated
3. Ensure `init()` method is called in constructor

### HTMX Integration Not Working
1. Check event listener for correct URL pattern
2. Verify `data-component` attribute matches registry key
3. Ensure component re-initializes after HTMX requests

## Best Practices

1. **Keep Components Simple**: Focus on UI interactions, not business logic
2. **Use Consistent Naming**: Follow the established patterns
3. **Add Console Logs**: Help with debugging and development
4. **Handle Errors Gracefully**: Don't let component errors break the page
5. **Use TypeScript**: Leverage type safety for better development experience
6. **Document Complex Logic**: Add comments for non-obvious functionality

## Related Files

- `web/ts/registry.ts` - Component registry system
- `web/ts/main.ts` - Main application entry point
- `web/tsconfig.json` - TypeScript configuration
- `web/esbuild.config.js` - Build configuration
- `web/package.json` - Build scripts and dependencies
