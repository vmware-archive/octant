import './index.scss'

import * as React from 'react'

import Graph from './graph'
import Grid from './grid'
import QuickView from './quickview'
import { Schema } from './schema'

interface Props {
  schema: Schema;
}

interface State {
  currentResource: string;
}

class ResourceViewer extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props)

    this.state = {
      currentResource: props.schema.selected,
    }

    this.setCurrentResource = this.setCurrentResource.bind(this)
  }

  setCurrentResource(name: string) {
    this.setState({ currentResource: name })
  }

  render() {
    const adjacencyList = this.props.schema.adjacencyList
    const objects = this.props.schema.objects

    const grid = new Grid(adjacencyList, objects)
    const rows = grid.create()

    const currentObject = objects[this.state.currentResource]

    return (
      <div className='resourceViewer'>
        <Graph
          rows={rows}
          schema={this.props.schema}
          setCurrentResource={this.setCurrentResource}
        />
        {currentObject ? <QuickView object={currentObject} /> : null }
      </div>
    )
  }
}

export default ResourceViewer
