import { mount } from 'enzyme'
import { LinkModel } from 'models/View'
import React from 'react'
import { MemoryRouter } from 'react-router'

import WebLink from '.'

describe('render web link', () => {
  const view: LinkModel = {
    type: 'link',
    ref: 'ref',
    value: 'value',
    title: 'title',
  }

  const webLink = mount(
    <MemoryRouter>
      <WebLink view={view} />
    </MemoryRouter>,
  )

  test('creates an anchor', () => {
    expect(webLink.html()).toEqual('<a href="/ref">value</a>')
  })
})
