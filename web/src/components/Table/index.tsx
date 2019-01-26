import './styles.scss'
import 'react-table/react-table.css'

import EmptyContent from 'components/EmptyContent'
import _ from 'lodash'
import { TableModel, viewFromContentType } from 'models/View'
import React from 'react'
import ReactTable from 'react-table'
import { renderView } from 'views'

interface Props {
  view: TableModel;
}

export default function Table({ view }: Props) {
  const { title, rows, columns, emptyContent } = view

  const tableColumns = _.map(columns, ({ name, accessor }: { name: string, accessor: string }, index: number) => ({
    Header: name,

    // TODO: re-enable when calculating column width
    // width: getColumnWidth(rows, accessor, name),
    id: `column_${index}`,
    Cell: (row) => {
      const entry = row.original
      const content = entry[accessor] as ContentType
      if (!_.isObject(content)) return content

      const cellView = viewFromContentType(content)
      return renderView(cellView)
    },
  }))

  const pageSize = rows && rows.length ? rows.length : null

  return (
    <div className='table--component'>
      <h2 className='table-component-title'>{title}</h2>
      {!rows || !rows.length ? (
        <EmptyContent emptyContent={emptyContent} />
      ) : (
        <ReactTable
          columns={tableColumns}
          data={rows}
          showPagination={false}
          pageSize={pageSize}
          defaultSorted={[
            {
              id: 'column_0',
            },
          ]}
        />
      )}
    </div>
  )
}

function getColumnWidth(rows: Array<{ [x: string]: ContentType; }>, accessor: string, headerText: string): number {
  const maxWidth = 600
  const magicSpacing = 10
  const cellLength = Math.max(
    ...rows.map((row) => (`${row[accessor]}` || '').length),
    headerText.length,
  )

  return Math.min(maxWidth, cellLength * magicSpacing)
}
