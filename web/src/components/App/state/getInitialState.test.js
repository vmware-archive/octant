import mockNavigation from 'api/mock/_navigation'
import { getNavLinkPath } from './getInitialState'

describe('getNavLinkPath', () => {
  test('generates nav link path from a valid path', () => {
    const linkPath = getNavLinkPath(
      mockNavigation,
      '/content/overview/workloads/cron-jobs'
    )
    expect(linkPath.length).toBe(3)

    expect(linkPath[0].path).toBe('/content/overview')
    expect(linkPath[0].title).toBe('Overview')

    expect(linkPath[1].path).toBe('/content/overview/workloads')
    expect(linkPath[1].title).toBe('Workloads')

    expect(linkPath[2].path).toBe('/content/overview/workloads/cron-jobs')
    expect(linkPath[2].title).toBe('Cron Jobs')
  })

  test('returns undefined on an invalid path', () => {
    const linkPath = getNavLinkPath(
      mockNavigation,
      '/content/this/is/invalid/cron-jobs'
    )
    expect(linkPath).toBeUndefined()
  })

  test('returns undefined on empty navigation data', () => {
    const linkPath = getNavLinkPath([], '/content/this/is/invalid/cron-jobs')
    expect(linkPath).toBeUndefined()
  })
})
