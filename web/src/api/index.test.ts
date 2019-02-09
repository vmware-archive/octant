import { getContentsUrl } from './index'

describe('getContentsUrl', () => {
  test('invalid paths return null', () => {
    expect(getContentsUrl(undefined)).toBe(null)
    expect(getContentsUrl(null)).toBe(null)
    expect(getContentsUrl('')).toBe(null)
    expect(getContentsUrl('/')).toBe(null)
  })

  test('valid path without params', () => {
    const fakePath = '/path'
    expect(getContentsUrl(fakePath)).toBe('api/v1/path')
  })

  test('adds query params', () => {
    const fakePath = '/path'
    expect(getContentsUrl(fakePath, { poll: 10 })).toBe('api/v1/path?poll=10')
    expect(getContentsUrl(fakePath, { filter: ['app:nginx', 'deployment:dev'] })).toBe(
      'api/v1/path?filter=app%3Anginx&filter=deployment%3Adev'
    )
  })
})
