import { ViewTitle } from 'components/ViewTitle'
import { SummaryItem, SummaryModel } from 'models'
import React from 'react'
import { renderView } from 'views'

import './styles.scss'

interface Props {
  view: SummaryModel
}

export default function Summary({ view }: Props) {
  return (
    <div className='summary-component'>
      <h2 className='summary-component-title'>
        <ViewTitle parts={view.title} />
      </h2>
      <div className='summary-component-section'>
        <table>
          <tbody>{view.items.map((section, index) => summaryContent(section, index))}</tbody>
        </table>
      </div>
    </div>
  )
}

function summaryContent(item: SummaryItem, key: number): JSX.Element {
  let content: JSX.Element

  switch (item.content.type) {
    case 'annotations':
    case 'labels':
    case 'link':
    case 'list':
    case 'ports':
    case 'selectors':
    case 'table':
    case 'text':
    case 'timestamp':
      content = renderView(item.content, { noHeader: true, noBorder: true })
      break
    default:
      throw new Error(`unsupported summary content type '${item.content.type}'`)
  }

  return (
    <tr key={key}>
      <th>{item.header}</th>
      <td>{content}</td>
    </tr>
  )
}
