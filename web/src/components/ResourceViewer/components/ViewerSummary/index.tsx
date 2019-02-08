import { ResourceViewerNode } from 'models/View'
import React from 'react'

import './styles.scss'

interface Props {
  node?: ResourceViewerNode
}

export default function ViewSummary(props: Props) {
  const { node } = props

  if (!node) {
    return <></>
  }

  return (
    <div className='viewSummary'>
      <div className='viewSummary--title'>
        <span className='status--ok' />
        {node.name}
      </div>
      <div className='viewSummary--content'>
        <div>multiple lines of content</div>
        <div>more content</div>
        <div>more content</div>
        <div>more content</div>
      </div>
    </div>
  )
}
