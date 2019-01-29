import './styles.scss'

import { SummaryItem } from 'models/SummaryItem'
import { SummaryModel } from 'models/View'
import React from 'react'
import { renderView } from 'views'

interface Props {
  view: SummaryModel;
}

export default function Summary({ view }: Props) {
  return (
    <div className='summary-component'>
      <h2 className='summary-component-title'>{view.title}</h2>
      <div className='summary-component-section'>
        {view.items.map((section, index) => summaryContent(section, index))}
      </div>
    </div>
  )
}

function summaryContent(item: SummaryItem, key: number): JSX.Element {
  let content: JSX.Element

  switch (item.content.type) {
    case 'text':
      content = renderView(item.content)
      break
    default:
      throw new Error(`unsupported summary content type '${item.content.type}'`)
  }

  return (
    <div key={key} className='summary--data summary--data-basic'>
      <div className='summary--data-key'>{item.header} </div>
      <div className='summary--data-value'>{content} </div>
    </div>
  )
}
