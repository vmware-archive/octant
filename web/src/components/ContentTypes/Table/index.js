import React from 'react'
import { Link } from 'react-router-dom'
import ReactTable from 'react-table'
import _ from 'lodash'
import moment from 'moment'
import EmptyContent from '../shared/EmptyContent'
import Labels from '../shared/Labels'
import './styles.scss'
import 'react-table/react-table.css'

export default function Table ({ data: { title, columns, rows } }) {
  const tableColumns = _.map(columns, ({ name, accessor }, index) => ({
    Header: name,
    accessor: (entry) => {
      if (!entry) return null
      const value = entry[accessor]
      if (!_.isObject(value)) return value
      switch (value.type) {
        case 'labels':
          return value.labels[0]
        case 'list':
          return value.list[0]
        case 'time':
          // currently a string, but should consider parsing into
          // a js date for sorting
          return value.time
        case 'link':
        case 'text':
        case 'string':
          return value.text
        default:
          return '-'
      }
    },
    id: `column_${index}`,
    Cell: (row) => {
      const entry = row.original
      const value = entry[accessor]
      if (!_.isObject(value)) return value
      switch (value.type) {
        case 'labels':
          return <Labels labels={value.labels} />
        case 'list':
          return value.list.join(', ')
        case 'link':
          return (
            <Link className='table--link' to={value.ref}>
              {value.text}
            </Link>
          )
        case 'time': {
          const t = moment(value.time)
          if (!t.isValid()) return value.time
          return t.toISOString()
        }
        default:
          return row.value
      }
    }
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
          data={rows}
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
