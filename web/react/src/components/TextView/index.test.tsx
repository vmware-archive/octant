import { shallow } from 'enzyme'
import React from 'react'

import TextView from '.'

describe('render text', () => {
  const view = {
    type: 'text',
    value: '10',
    title: [{ value: 'Revision History Limit', type: 'text' }],
  }

  const component = shallow(<TextView view={view} />)

  test('renders title and value', () => {
    const componentText = component.text()
    expect(componentText).toEqual(expect.stringContaining('10'))
  })
})
