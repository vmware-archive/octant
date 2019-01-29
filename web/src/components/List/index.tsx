import React from 'react'
import _ from 'lodash'
import { ListModel } from 'models/List'
import { renderView } from 'views'
import './styles.scss'

interface Props {
  view: ListModel
}

export default function List(props: Props) {
  const { view } = props
  return (
    <div className='content-type-list' data-test='list'>
      {
        view.items().map((item, i) => {
          return (
            <div className='content-type-list-item' key={i} >
              {
                (() => {
                  if (_.includes(['quadrant', 'label', 'summary', 'table'], item.type)) {
                    return renderView(item)
                  }
                  return <div />
                })()
              }
            </div>
          )
        })
      }
    </div>
  )
}
