import { ResourceViewerNode } from 'models'
import React from 'react'
import { renderView } from 'views'

import './styles.scss'

interface Props {
  node?: ResourceViewerNode
}

export default function ViewSummary(props: Props) {
  const { node } = props

  if (!node) {
    return <></>
  }

  const details = node.details || []

  const statusMessages = details.map((detail, index) => {
    return <li key={index}>{renderView(detail)}</li>
  })

  return (
    <div className='viewSummary'>
      <div className='viewSummary--title'>
        <span className={`status--${node.status}`} />
        {node.name}
      </div>
      <div className='viewSummary--content'>
        <ul>
          {statusMessages}
        </ul>
      </div>
    </div>
  )
}
