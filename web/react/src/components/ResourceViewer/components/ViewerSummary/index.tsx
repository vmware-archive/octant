import { ResourceViewerNode } from 'models'
import React, { Component } from 'react'
import { renderView } from 'views'

import './styles.scss'

interface Props {
  node?: ResourceViewerNode
}

class ViewSummary extends Component<Props> {
  render() {
    const { node } = this.props

    if (!node) {
      return <></>
    }

    const details = node.details || []

    const statusMessages = details.map((detail, index) => {
      return <li key={index}>{renderView(detail)}</li>
    })

    let title:
      | string
      | number
      | boolean
      | {}
      | JSX.Element
      | React.ReactElement<any>
      | React.ReactNodeArray
      | React.ReactPortal
    if (node.path) {
      title = renderView(node.path)
    } else {
      title = <span>{node.name}</span>
    }

    return (
      <div className='viewSummary'>
        <div className='viewSummary--title'>
          <span className={`status--${node.status}`} />
          {title}
        </div>
        <div className='viewSummary--content'>
          <ul>{statusMessages}</ul>
        </div>
      </div>
    )
  }
}

export default ViewSummary
