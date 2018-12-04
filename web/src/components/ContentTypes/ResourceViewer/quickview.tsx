import './quickview.scss'

import * as React from 'react'

import { ResourceObject } from './schema'

interface Props {
  name: string;
  object: ResourceObject;
}

class QuickView extends React.Component<Props, any> {
  render() {
    return <div className='quickView'>quick view for {this.props.name}</div>
  }
}

export default QuickView
