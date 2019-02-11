import { mount } from 'enzyme'
import { TitleView } from 'models'
import React from 'react'

import { ViewTitle } from '.'

describe('render view title', () => {
  test('with a single part', () => {
    const title: TitleView = [{ type: 'text', value: 'title' }]

    const view = mount(<ViewTitle parts={title} />)

    expect(view.text()).toEqual('title')
  })

  test('with multiple parts', () => {
    const title: TitleView = [{ type: 'text', value: 'part1' }, { type: 'text', value: 'part2' }]

    const view = mount(<ViewTitle parts={title} />)

    expect(view.children().length).toEqual(3)
    expect(view.childAt(0).text()).toEqual('part1')
    expect(view.childAt(1).text()).toEqual('›')
    expect(view.childAt(2).text()).toEqual('part2')
  })
})
