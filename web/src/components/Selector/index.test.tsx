import { shallow } from 'enzyme'
import { JSONSelectors } from 'models/View'
import React from 'react'

import Selectors from '.'

describe('render label selector', () => {
  const view = new JSONSelectors({
    config: {
      selectors: [
        {
          metadata: {
            type: 'labelSelector',
            title: '',
          },
          config: {
            key: 'key',
            value: 'value',
          },
        },
        {
          metadata: {
            type: 'expressionSelector',
            title: '',
          },
          config: {
            key: 'key',
            operator: 'In',
            values: ['value'],
          },
        },
      ],
    },
    metadata: {
      type: 'selectors',
      title: 'selectors',
    },
  })

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
