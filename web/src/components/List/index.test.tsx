import React from 'react'
import { shallow } from 'enzyme'
import List from './index'
import { JSONList } from 'models/List'
import Quadrant from 'components/Quadrant'
import Summary from 'components/Summary'
import Table from 'components/Table'

describe('render list', () => {

  test('empty list', () => {
    const listJSON: ContentType = {
      metadata: {
        type: 'list',
        title: 'List',
      },
      config: {
        items: [
        ],
      },
    }

    const listModel = new JSONList(listJSON)

    const wrapper = shallow(<List view={listModel} />)

    expect(wrapper.find('[data-test="list"]').children()).toHaveLength(0)
  })

  test('items w/ view components', () => {
    const listJSON: ContentType = {
      metadata: {
        type: 'list',
        title: 'List',
      },
      config: {
        items: [
          {
            metadata: {
              type: 'quadrant',
              title: 'Status',
            },
            config: {
              nw: {
                value: 1,
                label: 'Total',
              },
              ne: {
                value: 1,
                label: 'Updated',
              },
              sw: {
                value: 1,
                label: 'Available',
              },
              se: {
                value: 1,
                label: 'Unavailable',
              },
            },
          },
          {
            metadata: {
              type: 'summary',
              title: 'Additional properties',
            },
            config: {
              sections: [
                {
                  header: 'Image',
                  content: {
                    metadata: {
                      type: 'text',
                      title: 'Image',
                    },
                    config: {
                      value: 'nginx:1.15',
                    },
                  },
                },
                {
                  header: 'Port',
                  content: {
                    metadata: {
                      type: 'text',
                      title: 'Port',
                    },
                    config: {
                      value: '80/TCP',
                    },
                  },
                },
                {
                  header: 'Host Port',
                  content: {
                    metadata: {
                      type: 'text',
                      title: 'Host Port',
                    },
                    config: {
                      value: '0/TCP',
                    },
                  },
                },
                {
                  header: 'Environment',
                  content: {
                    metadata: {
                      type: 'text',
                      title: 'Environment',
                    },
                    config: {
                      value: 'none',
                    },
                  },
                },
                {
                  header: 'Mounts',
                  content: {
                    metadata: {
                      type: 'text',
                      title: 'Mounts',
                    },
                    config: {
                      value: '/usr/share/nginx/html=www(rw)',
                    },
                  },
                },
              ],
            },
          },
          {
            metadata: {
              type: 'table',
              title: 'Conditions',
            },
            config: {
              columns: [
                {
                  name: 'Name',
                  accessor: 'Name',
                },
                {
                  name: 'Labels',
                  accessor: 'Labels',
                },
                {
                  name: 'Desired',
                  accessor: 'Desired',
                },
                {
                  name: 'Current',
                  accessor: 'Current',
                },
                {
                  name: 'Ready',
                  accessor: 'Ready',
                },
                {
                  name: 'Age',
                  accessor: 'Age',
                },
                {
                  name: 'Containers',
                  accessor: 'Containers',
                },
                {
                  name: 'Images',
                  accessor: 'Images',
                },
                {
                  name: 'Selector',
                  accessor: 'Selector',
                },
              ],
              empty_content: 'Namespace overview does not contain any events for this Deployment',
            },
          },
        ],
      },
    }

    const listModel = new JSONList(listJSON)

    const wrapper = shallow(<List view={listModel} />)

    expect(wrapper.find(Quadrant)).toHaveLength(1)
    expect(wrapper.find(Summary)).toHaveLength(1)
    expect(wrapper.find(Table)).toHaveLength(1)
  })
})
