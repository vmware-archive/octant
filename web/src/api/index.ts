import queryString from 'query-string'
import mocks from './mock'

declare global {
  interface Window { API_BASE: string }
}

const { fetch } = window

export function getAPIBase() {
  return window.API_BASE || process.env.API_BASE
}

export const POLL_WAIT = 5

interface BuildRequestParams {
  endpoint: string;
  method?: string;
  data?: object;
}

async function buildRequest(params: BuildRequestParams) {
  const apiBase = getAPIBase()

  const { endpoint, method, data } = params
  if (!endpoint) throw new Error('endpoint not specified in buildRequest')
  const headers = {
    Accept: 'application/json',
    'Content-Type': 'application/json',
  }

  if (apiBase) {
    const fetchOptions: {
      headers: HeadersInit;
      method: string;
      body?: string;
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

  return Promise.resolve(mocks[endpoint])
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

export function getContentsUrl(path: string, namespace: string, poll?: string) {
  if (!path || path === '/') return null
  let query = ''
  if (namespace) query = `?${queryString.stringify({ namespace })}`
  // if poll is set poll the API
  if (poll) query += `&poll=${poll}`

  return `api/v1${path}${query}`
}

export function getContents(path: string, namespace: string) {
  const endpoint = getContentsUrl(path, namespace)
  return buildRequest({ endpoint })
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
