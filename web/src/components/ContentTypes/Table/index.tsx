import './styles.scss'
import 'react-table/react-table.css'

import _ from 'lodash'
import React from 'react'
import { Link } from 'react-router-dom'
import ReactTable from 'react-table'

import EmptyContent from '../shared/EmptyContent'
import Labels from '../shared/Labels'
import Timestamp from '../shared/Timestamp'

export interface ITable {
  metadata: {
    type: 'table';
    title: string;
  },
  config: {
    empty_content: string;
    columns: Array<{
      name: string;
      accessor: string;
    }>;
    rows: Array<{
      [x: string]: ContentType;
    }>;
  };
}

interface Props {
  data: ITable;
}

export default function Table({ data }: Props) {
  const { metadata: { title }, config: { rows, columns, empty_content: emptyContent } } = data

  const tableColumns = _.map(columns, ({ name, accessor }: { name: string, accessor: string }, index: number) => ({
    Header: name,
    accessor: (entry) => {
      if (!entry) return null
      const content = entry[accessor] as ContentType
      if (!_.isObject(content)) return content
      const { metadata: { type }, config } = content
      switch (type) {
        case 'labels':
          return _.entries(config.labels)[0]
        case 'list':
          return config.list[0]
        case 'link':
        case 'text':
        case 'string':
          return config.value
        default:
          return '-'
      }
    },
    id: `column_${index}`,
    Cell: (row) => {
      const entry = row.original
      const content = entry[accessor] as ContentType
      if (!_.isObject(content)) return content
      const { metadata: { type }, config } = content
      switch (type) {
        case 'labels':
          return <Labels labels={config.labels} />
        case 'link':
          return (
            <Link className='table--link' to={config.ref}>
              {config.value}
            </Link>
          )
        case 'timestamp': {
          return(
            <Timestamp config={config} />
          )
        }
        default:
          return row.value
      }
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
