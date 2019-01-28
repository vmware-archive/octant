import {mount} from 'enzyme'
import React from 'react'
import { MemoryRouter, Link as RouterLink } from 'react-router-dom'
import Link from './index'

describe('summary\'s Link ', () => {
  const params = {
    config: {
      ref: '/content/overview/workloads/replica-sets/nginx-deployment-7778c58546',
      value: 'nginx-deployment-7778c58546',
    },
    metadata: {
      type: 'link',
      title: 'nginx deployment',
    },
  }

  const component = mount(<MemoryRouter><Link params={params}/></MemoryRouter>)
  const linkComponent = component.find(Link)

  test('renders title', () => {
    const titleText = linkComponent.find('[data-test="title"]').text()

    expect(titleText).toEqual(expect.stringContaining('nginx deployment'))
  })

  test('renders link', () => {
    const routerLinkComponent = linkComponent.find(RouterLink).find('a')

    expect(routerLinkComponent).toHaveLength(1)
    expect(routerLinkComponent.text()).toBe('nginx-deployment-7778c58546')
    expect(routerLinkComponent.prop('href'))
      .toBe('/content/overview/workloads/replica-sets/nginx-deployment-7778c58546')
  })
})
