import queryString from 'query-string'
import mocks from './mock'

const { fetch } = window

function getAPIBase () {
  return window.API_BASE || process.env.API_BASE
}

async function buildRequest (params) {
  const apiBase = getAPIBase()

  const { endpoint, method, data } = params
  if (!endpoint) throw new Error('endpoint not specified in buildRequest')
  const headers = {
    Accept: 'application/json',
    'Content-Type': 'application/json'
  }

  if (apiBase) {
    const fetchOptions = {
      headers,
      method: method || 'GET'
    }
    if (data) fetchOptions.body = JSON.stringify(data)
    try {
      const response = await fetch(`${apiBase}/${endpoint}`, fetchOptions)
      if (response.status !== 204) {
        return await response.json()
      }
    } catch (e) {
      console.error('Failed fetch response: ', e) // eslint-disable-line no-console
      // Note(marlon): should consider throwing again here so that
      // specific messaging can be handled in the UI
    }
  }

  return Promise.resolve(mocks[endpoint])
}

export function getNavigation () {
  const params = {
    endpoint: 'api/v1/navigation'
  }
  return buildRequest(params)
}

export function getNamespaces () {
  const params = {
    endpoint: 'api/v1/namespaces'
  }
  return buildRequest(params)
}

export function getContents (path, namespace) {
  if (!path || path === '/') return null
  let query = ''
  if (namespace) query = `?${queryString.stringify({ namespace })}`
  const params = {
    endpoint: `api/v1${path}${query}`
  }
  return buildRequest(params)
}

export function setNamespace (namespace) {
  return buildRequest({
    endpoint: 'api/v1/namespace',
    method: 'POST',
    data: { namespace }
  })
}

export function getNamespace () {
  return buildRequest({ endpoint: 'api/v1/namespace' })
}
