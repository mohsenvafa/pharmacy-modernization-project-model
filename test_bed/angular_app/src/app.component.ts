import { Component, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core'

@Component({
  selector: 'app-root',
  standalone: true,
  template: `
    <main class="container">
      <header class="page-header">
        <img
          src="https://angular.io/assets/images/logos/angular/angular.svg"
          alt="Angular logo"
          width="48"
          height="48"
        />
        <h1>Angular Test Bed – Prescription Info Component</h1>
      </header>
      <form class="controls" (submit)="applyPrescription($event)">
        <label>
          Prescription ID
          <input
            type="text"
            [value]="pendingPrescriptionId"
            (input)="onPrescriptionInput($event)"
            placeholder="Enter prescription ID"
          />
        </label>
        <button type="submit">Load Prescription</button>
      </form>
      <section class="card">
        <prescription-info
          [attr.prescription-id]="currentPrescriptionId"
          env="local"
          auth-token="demo-token"
        ></prescription-info>
      </section>

      <form class="controls" (submit)="applyPatient($event)">
        <label>
          Patient ID
          <input
            type="text"
            [value]="pendingPatientId"
            (input)="onPatientInput($event)"
            placeholder="Enter patient ID"
          />
        </label>
        <button type="submit">Load Patient</button>
      </form>
      <section class="card">
        <patient-info
          [attr.patient-id]="currentPatientId"
          env="local"
          auth-token="demo-token"
        ></patient-info>
      </section>
    </main>
  `,
  styles: [
    `
      :host {
        font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
        padding: 0 2rem 2rem;
        display: block;
      }
      .container {
        max-width: 720px;
        margin: 0 auto;
        display: grid;
        gap: 1.5rem;
        padding: 0 0 2rem;
      }
      .page-header {
        position: sticky;
        top: 0;
        z-index: 10;
        display: flex;
        align-items: center;
        gap: 1rem;
        padding: 1rem 0;
        background: #ffffff;
        border-bottom: 1px solid rgba(0, 0, 0, 0.05);
      }
      .page-header img {
        display: block;
      }
      .controls {
        display: flex;
        align-items: flex-end;
        gap: 1rem;
      }
      .controls label {
        display: grid;
        gap: 0.25rem;
        font-size: 0.85rem;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        color: #4b5563;
      }
      .controls input {
        padding: 0.5rem;
        border-radius: 0.5rem;
        border: 1px solid rgba(0, 0, 0, 0.1);
        font-size: 1rem;
        min-width: 220px;
      }
      .controls button {
        padding: 0.5rem 1.25rem;
        border-radius: 0.5rem;
        background: #2563eb;
        color: white;
        border: none;
        font-size: 0.95rem;
        font-weight: 600;
        cursor: pointer;
      }
      .controls button:hover {
        background: #1d4ed8;
      }
      .card {
        padding: 1rem;
        border: 1px solid rgba(0, 0, 0, 0.1);
        border-radius: 0.75rem;
        background: #f9fafb;
      }
    `
  ],
  schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class AppComponent {
  readonly defaultPrescriptionId = 'RX002'
  pendingPrescriptionId = this.defaultPrescriptionId
  currentPrescriptionId = this.defaultPrescriptionId

  readonly defaultPatientId = 'P001'
  pendingPatientId = this.defaultPatientId
  currentPatientId = this.defaultPatientId

  onPrescriptionInput(event: Event) {
    const target = event.target as HTMLInputElement | null
    this.pendingPrescriptionId = target?.value ?? ''
  }

  applyPrescription(event: Event) {
    event.preventDefault()
    const nextId = this.pendingPrescriptionId.trim() || this.defaultPrescriptionId
    this.currentPrescriptionId = nextId
  }

  onPatientInput(event: Event) {
    const target = event.target as HTMLInputElement | null
    this.pendingPatientId = target?.value ?? ''
  }

  applyPatient(event: Event) {
    event.preventDefault()
    const nextId = this.pendingPatientId.trim() || this.defaultPatientId
    this.currentPatientId = nextId
  }
}

