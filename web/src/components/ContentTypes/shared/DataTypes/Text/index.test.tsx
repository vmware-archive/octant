import {shallow} from 'enzyme'
import React from 'react'
import Text from './index'

describe('render text', () => {
  const params = {
    metadata: {
      type: 'text',
      title: 'Revision History Limit',
    },
    config: {
      value: '10',
    },
  }

  const component = shallow(<Text params={params}/>)

  test('renders title and value', () => {
    const componentText = component.text()
    expect(componentText).toEqual(expect.stringContaining('Revision History Limit'))
    expect(componentText).toEqual(expect.stringContaining('10'))
  })
})
