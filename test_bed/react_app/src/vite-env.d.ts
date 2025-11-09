/// <reference types="vite/client" />

declare namespace JSX {
  interface IntrinsicElements {
    'prescription-info': {
      'prescription-id'?: string
      env?: 'local' | 'dev' | 'stg' | 'prod'
      'auth-token'?: string
      'base-path'?: string
      [key: string]: any
    }
    'patient-info': {
      'patient-id'?: string
      env?: 'local' | 'dev' | 'stg' | 'prod'
      'auth-token'?: string
      'base-path'?: string
      [key: string]: any
    }
  }
}

