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
    if (_.isObject(value)) {
      if (_.includes(['array', 'list', 'labels'], value.type)) {
        const arr = _.find([value.array, value.list, value.labels])
        if (arr) return arr.join(', ')
        return '-'
      }
      if (value.text) return value.text
      return '-'
    }
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
