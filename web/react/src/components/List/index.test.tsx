import TextView from 'components/TextView'
import { mount, shallow } from 'enzyme'
import { ListModel } from 'models'
import React from 'react'

import List from '.'

describe('render list', () => {
  test('empty list', () => {
    const list: ListModel = {
      type: 'list',
      items: [],
    }

    const wrapper = shallow(<List view={list} />)

    expect(wrapper.find('[data-test="list"]').children()).toHaveLength(0)
  })

  test('items w/ view components', () => {
    const list: ListModel = {
      type: 'list',
      items: [{ type: 'text' }, { type: 'text' }],
    }

    const wrapper = shallow(<List view={list} />)

    expect(wrapper.children().length).toEqual(2)
  })
})
