import Navigation from './_navigation'
import Summary from './_summary'
import Table from './_table'

const { fetch } = window

async function buildRequest (params) {
  const { endpoint } = params
  const headers = {
    Accept: 'application/json',
    'Content-Type': 'application/json'
  }
  const response = await fetch(endpoint, {
    headers
  })
  const json = await response.json()
  return json
}

export function getNavigation () {
  const params = {
    endpoint: 'http://127.0.0.1:52181/api/v1/navigation'
  }
  buildRequest(params)
  return Navigation
}

export function getSummary () {
  const params = {
    endpoint: 'http://127.0.0.1:52181/api/v1/content'
  }
  buildRequest(params)
  return Summary
}

export function getTable () {
  const params = {
    endpoint: 'http://127.0.0.1:52181/api/v1/content'
  }
  buildRequest(params)
  return Table
}
