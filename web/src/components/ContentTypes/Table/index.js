import React from 'react'
import { Link } from 'react-router-dom'
import ReactTable from 'react-table'
import _ from 'lodash'
import EmptyContent from '../shared/EmptyContent'
import Labels from '../shared/Labels'
import './styles.scss'
import 'react-table/react-table.css'

export default function Table ({ data: { title, columns, rows } }) {
  // Note(marlon):this lives here while the API keeps changing
  // Ideally a lot of this lives in a component or several
  const tableColumns = _.map(columns, ({ name, accessor }, index) => ({
    Header: name,
    accessor,
    id: `column_${index}`,
    Cell: (row) => {
      if (row && row.value) {
        const data = row.value
        if (data.type === 'labels') {
          return <Labels labels={data.labels} />
        }
      }
      return row.value
    }
  }))

  const tableRows = _.map(rows, row => _.mapValues(row, (value) => {
    if (_.isObject(value)) {
      if (value.type === 'labels' && value.labels) {
        return value
      }
      if (value.type === 'link' && value.ref) {
        return (
          <Link className='table--link' to={value.ref}>
            {value.text}
          </Link>
        )
      }
      if (_.includes(['array', 'list'], value.type)) {
        const arr = _.find([value.array, value.list])
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
      {!rows || !rows.length ? (
        <EmptyContent title={title} />
      ) : (
        <ReactTable
          columns={tableColumns}
          data={tableRows}
          showPagination={false}
          pageSize={pageSize}
          defaultSorted={[
            {
              id: 'column_0'
            }
          ]}
        />
      )}
    </div>
  )
}
