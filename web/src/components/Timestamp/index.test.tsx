import { shallow } from 'enzyme'
import { JSONTimestamp } from 'models/View'
import React from 'react'

import Timestamp, { summarizeTimestamp } from '.'

const baseSeconds = 1548437212
const baseDate = new Date(0)
baseDate.setUTCSeconds(baseSeconds)

describe('render timestamp', () => {
  const ct = {
    config: {
      timestamp: 1548430000,
    },
    metadata: {
      type: 'timestamp',
      title: '',
    },
  }

  const view = new JSONTimestamp(ct)

  const timestamp = shallow(<Timestamp view={view} baseTime={baseDate} />)

  test('sets the anchor content', () => {
    expect(timestamp.text()).toBe('2h')
  })

  test('sets the full date', () => {
    expect(timestamp.prop('data-tip')).toBe(
      'Friday, January 25, 2019 3:26 PM UTC',
    )
  })
})

describe.each`
  timestamp                    | expected
  ${baseSeconds - 86400 * 365} | ${'365d'}
  ${baseSeconds - 60 * 122}    | ${'2h'}
  ${baseSeconds - 65}          | ${'1m'}
  ${baseSeconds - 30}          | ${'30s'}
`('$timestamp', ({ timestamp, expected }) => {
  test(`returns ${expected}`, () => {
    expect(summarizeTimestamp(timestamp, baseDate)).toBe(expected)
  })
})
