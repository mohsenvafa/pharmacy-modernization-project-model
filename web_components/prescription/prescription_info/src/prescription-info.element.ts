import { buildServiceUrl, resolveBaseUrl } from '../../../shared/base-url'

const COMPONENT_TAG = 'prescription-info'
const MICRO_UI_PATH_SEGMENTS = ['micro-ui', 'prescriptions']

type PrescriptionInfoAttributes = {
  prescriptionId: string
  env: 'local' | 'dev' | 'stg' | 'prod'
  authToken: string
  basePath?: string | null
}

const observedAttributes = ['prescription-id', 'env', 'auth-token', 'base-path'] as const

function normalizeEnv(
  value: string | null
): PrescriptionInfoAttributes['env'] | null {
  if (!value) return null
  const normalized = value.toLowerCase()

  if (
    normalized === 'local' ||
    normalized === 'dev' ||
    normalized === 'stg' ||
    normalized === 'prod'
  ) {
    return normalized
  }

  return null
}

export class PrescriptionInfoElement extends HTMLElement {
  static get observedAttributes() {
    return observedAttributes
  }

  #shadow: ShadowRoot
  #isConnected = false

  constructor() {
    super()
    this.#shadow = this.attachShadow({ mode: 'open' })
  }

  connectedCallback() {
    this.#isConnected = true
    void this.#load()
  }

  disconnectedCallback() {
    this.#isConnected = false
  }

  attributeChangedCallback() {
    if (this.#isConnected) {
      void this.#load()
    }
  }

  get attributesData(): PrescriptionInfoAttributes | null {
    const prescriptionId = this.getAttribute('prescription-id')
    const env = normalizeEnv(this.getAttribute('env'))
    const authToken = this.getAttribute('auth-token')
    const basePathAttrRaw = this.getAttribute('base-path')
    const basePathAttr =
      basePathAttrRaw && basePathAttrRaw.trim().length > 0
        ? basePathAttrRaw.trim()
        : null

    if (!prescriptionId || !env || !authToken) {
      return null
    }

    return {
      prescriptionId,
      env,
      authToken,
      basePath: basePathAttr
    }
  }

  async #load() {
    const attrs = this.attributesData

    if (!attrs) {
      this.#renderError(
        'Missing required attributes. Please provide prescription-id, env, and auth-token.'
      )
      return
    }

    this.dispatchEvent(new CustomEvent('prescription-info:loading', { bubbles: true }))

    try {
      const endpoint = this.#buildEndpoint(attrs)
      const response = await fetch(endpoint.toString(), {
        headers: {
          Accept: 'text/html'
        },
        credentials: 'omit'
      })

      if (!response.ok) {
        throw new Error(`Request failed with status ${response.status}`)
      }

      const content = await response.text()
      this.#renderContent(content)
      this.dispatchEvent(new CustomEvent('prescription-info:loaded', { bubbles: true }))
    } catch (error) {
      console.error('Failed to load prescription info', error)
      this.#renderError('Unable to load prescription details.')
      this.dispatchEvent(
        new CustomEvent('prescription-info:error', {
          bubbles: true,
          detail: { error }
        })
      )
    }
  }

  #buildEndpoint(attrs: PrescriptionInfoAttributes): URL {
    const baseUrl = resolveBaseUrl(attrs.env, attrs.basePath)
    const endpoint = buildServiceUrl(
      baseUrl,
      ...MICRO_UI_PATH_SEGMENTS,
      encodeURIComponent(attrs.prescriptionId)
    )
    endpoint.searchParams.set('env', attrs.env)
    endpoint.searchParams.set('auth_token', attrs.authToken)
    return endpoint
  }

  #renderContent(content: string) {
    this.#shadow.innerHTML = `
      <style>
        :host {
          display: block;
          font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
        }

        .rx-prescription-info {
          border: 1px solid rgba(0, 0, 0, 0.1);
          border-radius: 0.5rem;
          padding: 1rem;
          background: white;
          color: #1f2937;
          box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
        }

        .rx-prescription-info__header {
          display: flex;
          justify-content: space-between;
          align-items: baseline;
          gap: 0.75rem;
          margin-bottom: 1rem;
        }

        .rx-prescription-info__header h2 {
          font-size: 1.1rem;
          font-weight: 600;
          margin: 0;
        }

        .rx-prescription-info__env {
          font-size: 0.75rem;
          text-transform: uppercase;
          letter-spacing: 0.08em;
          color: #4b5563;
        }

        .rx-prescription-info__details {
          display: grid;
          gap: 0.75rem;
        }

        .rx-prescription-info__details div {
          display: grid;
          gap: 0.25rem;
        }

        .rx-prescription-info__details dt {
          font-size: 0.75rem;
          letter-spacing: 0.05em;
          text-transform: uppercase;
          color: #6b7280;
        }

        .rx-prescription-info__details dd {
          margin: 0;
          font-size: 0.95rem;
        }

        .rx-prescription-info--error {
          border-color: rgba(220, 38, 38, 0.2);
          background: rgba(254, 242, 242, 0.8);
          color: #991b1b;
        }
      </style>
      ${content}
    `
  }

  #renderError(message: string) {
    this.#shadow.innerHTML = `
      <style>
        :host {
          display: block;
          font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
        }
        .rx-prescription-info {
          border: 1px solid rgba(220, 38, 38, 0.2);
          border-radius: 0.5rem;
          padding: 1rem;
          background: rgba(254, 242, 242, 0.8);
          color: #991b1b;
        }
      </style>
      <section class="rx-prescription-info rx-prescription-info--error">
        <strong>Prescription Info Error</strong>
        <p>${message}</p>
      </section>
    `
  }
}

export function registerPrescriptionInfoComponent(tagName = COMPONENT_TAG) {
  if (!customElements.get(tagName)) {
    customElements.define(tagName, PrescriptionInfoElement)
  }
}

registerPrescriptionInfoComponent()

