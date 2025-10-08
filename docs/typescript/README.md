# TypeScript Documentation

Documentation for TypeScript components and frontend development.

## 📚 Documentation Files

- **[Adding TypeScript Components](./ADDING_TYPESCRIPT_COMPONENTS.md)** - Guide for creating and integrating TypeScript components

## 🎯 Overview

TypeScript is used selectively for complex client-side interactions that require:
- State management
- Complex user interactions
- Real-time updates
- Client-side validation
- Rich UI components

## 🏗️ Architecture

### Hybrid Approach
The application uses a hybrid frontend approach:
- **HTMX** - For simple, server-driven interactions
- **TypeScript** - For complex, client-driven components
- **Templ** - For server-side rendering

### When to Use TypeScript
Use TypeScript components when:
- Complex client-side state management needed
- Real-time updates required
- Heavy client-side computation
- Interactive visualizations
- Rich text editors or complex forms

### When to Use HTMX
Use HTMX when:
- Simple CRUD operations
- Form submissions
- Page navigation
- Server-driven updates
- Most typical web interactions

## 📁 File Structure

```
web/
├── assets/
│   ├── js/
│   │   ├── src/
│   │   │   ├── components/     # TypeScript components
│   │   │   │   └── example.ts
│   │   │   └── main.ts         # Entry point
│   │   └── dist/               # Compiled JavaScript
│   │       └── main.js
│   └── vendor/                 # Third-party libraries
│       └── htmx.min.js
└── components/                 # Templ components
```

## 🚀 Quick Start

### 1. Install Dependencies
```bash
npm install
```

### 2. Create a Component
```typescript
// web/assets/js/src/components/myComponent.ts
export class MyComponent {
    constructor(private element: HTMLElement) {
        this.init();
    }

    private init() {
        this.element.addEventListener('click', this.handleClick.bind(this));
    }

    private handleClick(event: Event) {
        console.log('Clicked!', event);
    }
}
```

### 3. Register in main.ts
```typescript
import { MyComponent } from './components/myComponent';

document.addEventListener('DOMContentLoaded', () => {
    const elements = document.querySelectorAll('[data-component="my-component"]');
    elements.forEach(el => new MyComponent(el as HTMLElement));
});
```

### 4. Use in Templ
```go
templ MyPage() {
    <div data-component="my-component">
        <button>Click Me</button>
    </div>
}
```

### 5. Build
```bash
npm run build
```

## 🔧 Development Workflow

### Watch Mode
```bash
npm run watch
```

### Build for Production
```bash
npm run build:prod
```

### Type Checking
```bash
npm run type-check
```

## 📦 Component Patterns

### Data Attributes Pattern
Use `data-*` attributes for component initialization:

```html
<div 
    data-component="search" 
    data-api-url="/api/patients"
    data-debounce="300">
</div>
```

```typescript
const apiUrl = element.dataset.apiUrl;
const debounce = parseInt(element.dataset.debounce || '0');
```

### Event-Driven Components
Components communicate through custom events:

```typescript
// Dispatch
this.element.dispatchEvent(new CustomEvent('data-loaded', {
    detail: { data: result },
    bubbles: true
}));

// Listen
document.addEventListener('data-loaded', (e) => {
    console.log(e.detail.data);
});
```

### Lifecycle Hooks
```typescript
export class Component {
    constructor(element: HTMLElement) {
        this.init();
    }

    private init() {
        // Initialize component
    }

    public destroy() {
        // Cleanup
    }
}
```

## 🎨 Integration with HTMX

TypeScript components can work alongside HTMX:

```html
<div data-component="live-search">
    <input 
        type="text" 
        hx-post="/search"
        hx-trigger="input changed delay:500ms"
        hx-target="#results">
    <div id="results"></div>
</div>
```

## 🧪 Testing

### Unit Tests
```typescript
import { MyComponent } from './myComponent';

describe('MyComponent', () => {
    it('should initialize', () => {
        const element = document.createElement('div');
        const component = new MyComponent(element);
        expect(component).toBeDefined();
    });
});
```

## 📖 Related Documentation

- [Architecture Overview](../architecture/ARCHITECTURE.md)
- [Adding TypeScript Components - Detailed Guide](./ADDING_TYPESCRIPT_COMPONENTS.md)

## 🛠️ Tools & Libraries

- **TypeScript** - Type-safe JavaScript
- **esbuild** - Fast bundler
- **HTMX** - HTML-first approach
- **npm** - Package management

## 💡 Best Practices

1. **Keep Components Small** - Single responsibility
2. **Use Type Safety** - Leverage TypeScript's type system
3. **Avoid Overuse** - Use HTMX for simple interactions
4. **Clean Up** - Implement destroy methods
5. **Document Data Attributes** - Clear API for integration
6. **Progressive Enhancement** - Components should enhance, not replace

## 🔍 Common Use Cases

- **Autocomplete/Search** - Real-time filtering
- **Charts/Graphs** - Data visualization
- **Rich Text Editors** - Complex input
- **Drag & Drop** - Interactive UI
- **Real-time Updates** - WebSocket integration
- **Complex Forms** - Multi-step, dynamic validation

