import React from 'react'
import ReactTable from 'react-table'
import _ from 'lodash'
import './styles.scss'
import 'react-table/react-table.css'

export default function Table ({ data: { name: tableTitle, columns, rows } }) {
  const tableColumns = _.map(columns, ({ name, accessor }) => ({ Header: name, accessor }))
  return (
    <div className='table--component'>
      <ReactTable
        columns={tableColumns}
        data={rows}
        showPagination={false}
        pageSize={rows.length}
      />
    </div>
  )
}
