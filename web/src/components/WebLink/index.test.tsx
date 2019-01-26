import { shallow } from 'enzyme'
import { LinkModel } from 'models/View'
import React from 'react'

import WebLink from '.'

describe('render web link', () => {
    const view: LinkModel = {
        type: 'link',
        ref: 'ref',
        value: 'value',
        title: 'title',
    }

    const webLink = shallow(<WebLink view={view}/>)

    test('creates an anchor', () => {
        expect(webLink.text()).toBe('value')
        expect(webLink.prop('href')).toBe('ref')
    })
})
