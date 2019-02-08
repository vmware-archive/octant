import queryString from 'query-string'

declare global {
  interface Window {
    API_BASE: string
    EventSource(uri: string): void
  }
}

const { fetch } = window

export function getAPIBase() {
  return window.API_BASE || process.env.API_BASE
}

export const POLL_WAIT: number = 5

interface BuildRequestParams {
  endpoint: string
  method?: string
  data?: object
}

async function buildRequest(params: BuildRequestParams): Promise<any> {
  const apiBase = getAPIBase()

  const { endpoint, method, data } = params
  if (!endpoint) throw new Error('endpoint not specified in buildRequest')
  const headers = {
    Accept: 'application/json',
    'Content-Type': 'application/json',
  }

  if (apiBase) {
    const fetchOptions: {
      headers: HeadersInit
      method: string
      body?: string
    } = {
      headers,
      method: method || 'GET',
    }
    if (data) fetchOptions.body = JSON.stringify(data)
    try {
      const response = await fetch(`${apiBase}/${endpoint}`, fetchOptions)
      if (response.status !== 204) {
        return await response.json()
      }
    } catch (e) {
      console.error('Failed fetch response: ', e) // eslint-disable-line no-console
      // Note(marlon): add logging for network errors here
      throw e
    }
  }

  return Promise.resolve(null)
}

export function getNavigation() {
  const params = {
    endpoint: 'api/v1/navigation',
  }
  return buildRequest(params)
}

export function getNamespaces() {
  const params = {
    endpoint: 'api/v1/namespaces',
  }
  return buildRequest(params)
}

export type ContentsUrlParams = Partial<{
  namespace: string
  poll: number
  filter: string[]
}>

export function getContentsUrl(path: string, params?: ContentsUrlParams): string | null {
  if (!path || path === '/') return null
  if (params) path += `?${queryString.stringify(params)}`
  return `api/v1${path}`
}

export function setNamespace(namespace: string) {
  return buildRequest({
    endpoint: 'api/v1/namespace',
    method: 'POST',
    data: { namespace },
  })
}

export function getNamespace() {
  return buildRequest({ endpoint: 'api/v1/namespace' })
}
