import EmptyContent from 'components/EmptyContent'
import { compareModel, instanceOfComparableView, TableModel, TableRow, TextModel, View } from 'models/View'
import React from 'react'
import ReactTable from 'react-table'
import 'react-table/react-table.css'
import { renderView } from 'views'

import './styles.scss'

interface Props {
  view: TableModel
}

export default function Table({ view }: Props) {
  const { title, rows, columns, emptyContent } = view

  const tableColumns = columns.map(({ name, accessor }, index) => {
    return {
      Header: name,
      id: `column_${index}`,
      accessor: (row) => row[name],
      Cell: (row) => {
        const entry = row.original
        const cellView = entry[accessor]
        return renderView(cellView)
      },
      maxWidth: getColumnWidth(rows, accessor),
      sortMethod: sortMethod(),
      sortable: isSortable(rows, accessor),
    }
  })

  const pageSize = rows && rows.length ? rows.length : null

  const noDataText = emptyContent || 'no data'
  if (rows.length > 0) {
    return (
      <div className='table--component'>
        <h2 className='table-component-title'>{title}</h2>
        <ReactTable
          noDataText={noDataText}
          columns={tableColumns}
          data={rows}
          showPagination={false}
          pageSize={pageSize}
          multiSort={false}
        />
      </div>
    )
  }

  return (
    <div className='table--component'>
      <h2 className='table-component-title'>{title}</h2>
      <EmptyContent emptyContent={emptyContent} />
    </div>
  )
}

export function getColumnWidth(
  rows: Array<{
    [key: string]: View
  }>,
  accessor: string
): number | undefined {
  const lens = rows
    .map((row) => {
      if (!row.hasOwnProperty(accessor)) {
        throw new Error(`table doesn't have a column named "${accessor}"`)
      }

      const view = row[accessor]
      switch (view.type) {
        case 'timestamp':
          return 60
        case 'text':
          return (view as TextModel).value.length * 45
        default:
          return 0
      }
    })
    .map((n) => Math.max(n, accessor.length * 45))

  const max = Math.max(...lens)
  return max === 0 ? undefined : max
}

export function isSortable(rows: TableRow[], accessor: string): boolean {
  return (rows.length > 0 && instanceOfComparableView(rows[0][accessor])) || false
}

export function sortMethod(): (a, b, desc) => number {
  return (a: View, b: View, desc: boolean) => {
    if (instanceOfComparableView(a)) {
      const n = compareModel(a, b)
      return desc ? n : n * -1
    }
    return 0
  }
}
