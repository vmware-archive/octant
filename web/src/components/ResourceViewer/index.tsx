import { ResourceViewerModel } from 'models/View'
import React, { Component } from 'react'

import { Tooltip } from './components/Tooltip'
import ViewSummary from './components/ViewerSummary'
import Graph from './graph'
import './index.scss'
import ResourceNode from './node'
import './resource'

interface Props {
  view: ResourceViewerModel
}

interface State {
  currentResource: string
}

class ResourceViewer extends Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = {
      currentResource: '',
    }
  }

  setCurrentResource = (name: string) => {
    this.setState({ currentResource: name })
  }

  render() {
    const adjacencyList = this.props.view.edges
    const objects = this.props.view.nodes
    const currentObject = objects[this.state.currentResource]
    const nodes = {}

    const tooltips: JSX.Element[] = []

    for (const [id, object] of Object.entries(objects)) {
      nodes[id] = new ResourceNode(id, object, this.state.currentResource === id).toDescriptor()

      tooltips.push(
        <Tooltip key={id} id={id}>
          <ViewSummary node={currentObject} />
        </Tooltip>
      )
    }

    const edges = []
    for (const [node, nodeEdges] of Object.entries(adjacencyList)) {
      edges.push(
        ...nodeEdges.map((e) => [
          node,
          e.node,
          {
            arrowhead: 'undirected',
            arrowheadStyle: 'fill: rgba(173, 187, 196, 0.3)',
          },
        ])
      )
    }

    return (
      <div className='resourceViewer'>
        <Graph width='100%' height='100%' nodes={nodes} edges={edges} onNodeClick={this.setCurrentResource} />
        {tooltips}
      </div>
    )
  }
}

export default ResourceViewer
