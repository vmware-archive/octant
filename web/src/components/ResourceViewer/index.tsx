import { ResourceViewerModel } from 'models'
import React, { Component } from 'react'

import Graph from './components/Graph'
import ViewSummary from './components/ViewerSummary'
import './index.scss'
import ResourceNode from './node'

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
      currentResource: this.props.view.selected,
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

    for (const [id, object] of Object.entries(objects)) {
      nodes[id] = new ResourceNode(id, object, this.state.currentResource === id).toDescriptor()
    }

    const edges = []
    if (adjacencyList) {
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
      }

    return (
      <div className='resourceViewer'>
        <Graph width='100%' height='100%' nodes={nodes} edges={edges} onNodeClick={this.setCurrentResource} />
        <ViewSummary node={currentObject} />
      </div>
    )
  }
}

export default ResourceViewer
