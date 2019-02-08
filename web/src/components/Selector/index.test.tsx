import { shallow } from 'enzyme'
import React from 'react'

import Selectors from '.'

describe('render label selector', () => {
  const view = {
    type: 'selectors',
    selectors: [
      { key: 'key', value: 'value', type: 'labelSelector' },
      { key: 'key', operator: 'In', values: ['value'], type: 'expressionSelector' },
    ],
  }

  const selectors = shallow(<Selectors view={view} />)

  test('create two selectors', () => {
    expect(selectors.children().length).toBe(2)
  })

  test('first selector is a label selector', () => {
    const selector = selectors.childAt(0)
    expect(selector.hasClass('selectors--label')).toBeTruthy()
    expect(selector.text()).toBe('key:value')
  })

  test('second selector is a expression selector', () => {
    const selector = selectors.childAt(1)
    expect(selector.hasClass('selectors--expression')).toBeTruthy()
    expect(selector.text()).toBe('key In []')
  })
})
