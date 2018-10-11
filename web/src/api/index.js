import queryString from 'query-string'
import mocks from './mock'

const { fetch } = window

function getAPIBase () {
  return window.API_BASE || process.env.API_BASE
}

async function buildRequest (params) {
  const apiBase = getAPIBase()

  const { endpoint } = params
  const headers = {
    Accept: 'application/json',
    'Content-Type': 'application/json'
  }

  if (apiBase) {
    const response = await fetch(`${apiBase}/${endpoint}`, {
      headers
    })
    return response.json()
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
