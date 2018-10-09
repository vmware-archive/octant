import React from 'react'
import ReactTable from 'react-table'
import _ from 'lodash'
import './styles.scss'
import 'react-table/react-table.css'

export default function Table ({ data: { title, columns, rows } }) {
  const tableColumns = _.map(columns, ({ name, accessor }) => ({
    Header: name,
    accessor
  }))
  const textRows = _.map(rows, row => _.mapValues(row, (value) => {
    if (_.isObject(value)) return value.text
    return value
  }))
  const pageSize = rows && rows.length ? rows.length : null
  return (
    <div className='table--component'>
      <h2 className='table-component-title'>{title}</h2>
      <ReactTable
        columns={tableColumns}
        data={textRows}
        showPagination={false}
        pageSize={pageSize}
      />
    </div>
  )
}
