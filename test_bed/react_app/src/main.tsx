import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'

import '@rx/micro-prescription-info'
import '@rx/micro-patient-info'

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
)

