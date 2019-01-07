import React, { Component } from 'react'
import QuickView from './quickview'
import { IResourceViewer } from './schema'
import Graph from './graph'
import ResourceNode from './node'
import './index.scss'
import './resource'

interface Props {
  data: IResourceViewer;
}

interface State {
  currentResource: string;
}

class ResourceViewer extends Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = {
      currentResource: props.data.config.selected,
    }
  }

  setCurrentResource = (name: string) => {
    this.setState({ currentResource: name })
  }

  render() {
    const { data: { config } } = this.props
    const adjacencyList = config.adjacencyList
    const objects = config.objects
    const currentObject = objects[this.state.currentResource]
    const nodes = {}
    for (const [id, object] of Object.entries(config.objects)) {
      nodes[id] = new ResourceNode(
        object,
        this.state.currentResource === id,
      ).toDescriptor()
    }

    const edges = []
    for (const [node, nodeEdges] of Object.entries(adjacencyList)) {
      edges.push(...nodeEdges.map((e) => [node, e.node, { arrowhead: 'vee' }]))
    }

    return (
      <div className='resourceViewer'>
        <Graph
          width='100%'
          height='100%'
          nodes={nodes}
          edges={edges}
          onNodeClick={this.setCurrentResource}
        />
        {currentObject ? <QuickView object={currentObject} /> : null}
      </div>
    )
  }
}

export default ResourceViewer
