const DEFAULT_BASE_URL_BY_ENV = {
  local: 'http://localhost:8080',
  dev: 'https://dev.pharmacy.com',
  stg: 'https://stg.pharmacy.com',
  prod: 'https://pharmacy.com'
} as const

export type EnvironmentKey = keyof typeof DEFAULT_BASE_URL_BY_ENV

export function resolveBaseUrl(
  env: EnvironmentKey,
  override?: string | null
): string {
  if (override) {
    return normalizeBaseUrl(override)
  }

  return DEFAULT_BASE_URL_BY_ENV[env] ?? DEFAULT_BASE_URL_BY_ENV.local
}

export function normalizeBaseUrl(value: string): string {
  try {
    const url = new URL(value)
    url.hash = ''
    url.search = ''
    return `${url.origin}${trimTrailingSlash(url.pathname)}`
  } catch (error) {
    return trimTrailingSlash(value)
  }
}

export function buildServiceUrl(baseUrl: string, ...segments: string[]): URL {
  const normalizedBase = ensureTrailingSlash(normalizeBaseUrl(baseUrl))
  const normalizedPath = segments
    .map(segment => segment.replace(/^\/+|\/+$/g, ''))
    .filter(Boolean)
    .join('/')

  return new URL(normalizedPath, normalizedBase)
}

function ensureTrailingSlash(value: string): string {
  return value.endsWith('/') ? value : `${value}/`
}

function trimTrailingSlash(value: string): string {
  if (!value) return value
  return value.replace(/\/+$/, '')
}

