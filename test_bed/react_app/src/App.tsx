import { useState } from 'react'

export default function App() {
  const [env, setEnv] = useState<'local' | 'dev' | 'stg' | 'prod'>('local')
  const [pendingId, setPendingId] = useState('RX002')
  const [currentId, setCurrentId] = useState('RX002')
  const [pendingPatientId, setPendingPatientId] = useState('P001')
  const [currentPatientId, setCurrentPatientId] = useState('P001')

  return (
    <main style={styles.container}>
      <header style={styles.header}>
        <img
          style={styles.logo}
          src="https://upload.wikimedia.org/wikipedia/commons/a/a7/React-icon.svg"
          alt="React logo"
          width={48}
          height={48}
        />
        <h1 style={styles.title}>React Test Bed – Prescription Info Component</h1>
      </header>

      <section style={styles.controls}>
        <label style={styles.label}>
          Environment
          <select
            value={env}
            onChange={event => setEnv(event.target.value as typeof env)}
            style={styles.select}
          >
            <option value="local">local</option>
            <option value="dev">dev</option>
            <option value="stg">stg</option>
            <option value="prod">prod</option>
          </select>
        </label>
      </section>

      <form
        style={styles.form}
        onSubmit={event => {
          event.preventDefault()
          setCurrentId(pendingId.trim() || 'RX002')
        }}
      >
        <label style={styles.label}>
          Prescription ID
          <input
            style={styles.input}
            value={pendingId}
            onChange={event => setPendingId(event.target.value)}
            placeholder="Enter prescription ID"
          />
        </label>
        <button style={styles.button} type="submit">
          Load Prescription
        </button>
      </form>

      <section style={styles.card}>
        <prescription-info
          prescription-id={currentId}
          env={env}
          auth-token="demo-token"
        ></prescription-info>
      </section>

      <form
        style={styles.form}
        onSubmit={event => {
          event.preventDefault()
          setCurrentPatientId(pendingPatientId.trim() || 'P001')
        }}
      >
        <label style={styles.label}>
          Patient ID
          <input
            style={styles.input}
            value={pendingPatientId}
            onChange={event => setPendingPatientId(event.target.value)}
            placeholder="Enter patient ID"
          />
        </label>
        <button style={styles.button} type="submit">
          Load Patient
        </button>
      </form>

      <section style={styles.card}>
        <patient-info
          patient-id={currentPatientId}
          env={env}
          auth-token="demo-token"
        ></patient-info>
      </section>
    </main>
  )
}

const styles: Record<string, React.CSSProperties> = {
  container: {
    fontFamily: `system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif`,
    padding: '0 2rem 2rem',
    maxWidth: '720px',
    margin: '0 auto',
    display: 'grid',
    gap: '1.5rem'
  },
  header: {
    position: 'sticky',
    top: 0,
    zIndex: 10,
    display: 'flex',
    alignItems: 'center',
    gap: '1rem',
    padding: '1rem 0',
    background: '#ffffff',
    borderBottom: '1px solid rgba(0, 0, 0, 0.05)'
  },
  logo: {
    display: 'block'
  },
  title: {
    fontSize: '1.75rem',
    fontWeight: 600,
    margin: 0
  },
  controls: {
    display: 'flex',
    gap: '1rem',
    alignItems: 'center'
  },
  form: {
    display: 'flex',
    gap: '1rem',
    alignItems: 'flex-end'
  },
  label: {
    display: 'grid',
    gap: '0.25rem',
    fontSize: '0.85rem',
    textTransform: 'uppercase',
    letterSpacing: '0.05em',
    color: '#4b5563'
  },
  select: {
    padding: '0.5rem',
    borderRadius: '0.5rem',
    border: '1px solid rgba(0,0,0,0.1)',
    fontSize: '1rem'
  },
  input: {
    padding: '0.5rem',
    borderRadius: '0.5rem',
    border: '1px solid rgba(0,0,0,0.1)',
    fontSize: '1rem',
    minWidth: '220px'
  },
  button: {
    padding: '0.5rem 1.25rem',
    borderRadius: '0.5rem',
    background: '#2563eb',
    color: 'white',
    border: 'none',
    fontSize: '0.95rem',
    fontWeight: 600,
    cursor: 'pointer'
  },
  card: {
    padding: '1rem',
    borderRadius: '0.75rem',
    border: '1px solid rgba(0,0,0,0.1)',
    background: '#f9fafb'
  }
}

