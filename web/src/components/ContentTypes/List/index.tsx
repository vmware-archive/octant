import _ from 'lodash'
import React from 'react'

import Label, { ILabel } from '../Grid/components/Label'
import Quadrant, { IQuadrant } from '../Grid/components/Quadrant'
import Summary, { ISummary } from '../Summary'
import Table, { ITable } from '../Table'

export interface IList {
  metadata: {
    type: 'list';
    title: string;
  };
  config: {
    items: ContentType[];
  };
}

interface Props {
  data: IList
}

export default function List(props: Props) {
  const { data: { config: { items } } } = props
  return (
    <div className='content-type-list'>
      {
        _.map(items, (item, i) => {
          const { metadata: { type } } = item
          return (
            <div className='content-type-list-item' key={i} >
              {
                (() => {
                  switch (type) {
                    case 'quadrant':
                      return <Quadrant data={item as IQuadrant} />
                    case 'label':
                      return <Label data={item as ILabel} />
                    case 'summary':
                      return <Summary data={item as ISummary} />
                    case 'table':
                      return <Table data={item as ITable} />
                    default:
                    return <div/>
                  }
                })()
              }
            </div>
          )
        })
      }
    </div>
  )
}
