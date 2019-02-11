import { ListModel } from 'models'
import React from 'react'
import { renderView } from 'views'

import './styles.scss'

interface Props {
  view: ListModel
}

export default function List(props: Props) {
  const { view } = props

  return (
    <div className='content-type-list' data-test='list'>
      {view.items.map((item, i) => {
        return (
          <div className='content-type-list-item' key={i}>
            {renderView(item)}
          </div>
        )
      })}
    </div>
  )
}
