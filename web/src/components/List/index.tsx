import './styles.scss'

import { ListModel } from 'models/List'
import React from 'react'
import { renderView } from 'views'

interface Props {
  view: ListModel
}

export default function List(props: Props) {
  const { view } = props
  return (
    <div className='content-type-list'>
      {
        view.items().map((item, i) => {

          return (
            <div className='content-type-list-item' key={i} >
              {
                (() => {
                  switch (item.type) {
                    case 'quadrant':
                    case 'label':
                    case 'summary':
                    case 'table':
                        return renderView(item)
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
