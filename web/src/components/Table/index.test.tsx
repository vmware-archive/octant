import { mount } from 'enzyme'
import { TableRow } from 'models/View'
import React from 'react'

import Table, { getColumnWidth, isSortable, sortMethod } from '.'

describe('render table', () => {
  describe('creates a table', () => {
    const view = {
      type: 'table',
      title: 'my table',
      columns: [
        { name: 'one', accessor: 'one' },
        { name: 'two', accessor: 'two' },
      ],
      rows: [
        {
          one: { type: 'text', value: 'a', title: '', isComparable: true },
          two: { type: 'labels', labels: {a: 'a'}, title: ''},
        },
        {
          one: { type: 'text', value: 'c', title: '', isComparable: true },
          two: { type: 'labels', labels: {b: 'b'}, title: ''  },
        },
      ],
      emptyContent: 'is empty',
    }

    const table = mount(<Table view={view} />)

    test('it draws a title', () => {
      expect(table.find('.table-component-title').text()).toEqual('my table')
    })

    test('it has two rows', () => {
      expect(table.find('.rt-tbody').children().length).toEqual(2)
    })

    test('it has two columns', () => {
      expect(table.find('.rt-thead .rt-tr').children().length).toEqual(2)
    })
  })

  describe('getColumnWidth', () => {
    const ts = {
      type: 'timestamp',
      timestamp: 5,
      title: '',
      isComparable: true,
    }

    const text = {
      type: 'text',
      value: 'a',
      title: '',
      isComparable: true,
    }

    const other = {
      type: 'over',
      title: '',
    }

    test('invalid column accessor', () => {
      const rows: TableRow[] = [{ ts }]
      expect(() => getColumnWidth(rows, 'invalid')).toThrow()
    })

    test('timestamp with short column header', () => {
      const rows: TableRow[] = [{ t: ts }]
      expect(getColumnWidth(rows, 't')).toEqual(60)
    })

    test('timestamp with long column header', () => {
      const rows: TableRow[] = [{ 'long text': ts }]
      expect(getColumnWidth(rows, 'long text')).toEqual(45 * 9)
    })

    test('text column', () => {
      const rows: TableRow[] = [{ text }]
      expect(getColumnWidth(rows, 'text')).toEqual(45 * 4)
    })

    test('other column', () => {
      const rows: TableRow[] = [{ other }]
      expect(getColumnWidth(rows, 'other')).toEqual(45 * 5)
    })
  })

  describe('isSortable', () => {
    const textView = {
      type: 'text',
      title: '',
      isComparable: true,
    }

    const labelsView = {
      type: 'labels',
      title: '',
    }

    const rows: TableRow[] = [{ text: textView, labels: labelsView }]

    test('views that are comparable are sortable', () => {
      expect(isSortable(rows, 'text')).toBe(true)
    })

    test('views that are not comparable are not sortable', () => {
      expect(isSortable(rows, 'labels')).toBe(false)
    })
  })

  describe('sortMethod', () => {
    const view1 = {
      type: 'text',
      value: 'a',
      title: '',
      isComparable: true,
    }
    const view2 = {
      type: 'text',
      value: 'b',
      title: '',
      isComparable: true,
    }

    test('sort in descending order', () => {
      const res = sortMethod()(view1, view2, true)
      expect(res).toEqual(-1)
    })

    test('sort in ascending order', () => {
      const res = sortMethod()(view1, view2, false)
      expect(res).toEqual(1)
    })
  })
})
