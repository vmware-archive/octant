import './index.scss'
import './resource'

import * as React from 'react'

import QuickView from './quickview'
import { Schema } from './schema'
import Graph from './graph'
import ResourceNode from './node'
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

    const currentObject = objects[this.state.currentResource]

    const nodes = {}
    for (const [id, object] of Object.entries(this.props.schema.objects)) {
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
