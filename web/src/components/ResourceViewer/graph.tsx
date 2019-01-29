import './graph.scss'

import React from 'react'
import isEqual from 'react-fast-compare'
import dagreD3 from 'dagre-d3'
import * as d3 from 'd3'

interface Props {
  nodes: any;
  edges: any[];
  height?: string;
  width?: string;
  shapeRenders?: any;
  onNodeClick(name: string): void;
}

class Graph extends React.Component<Props> {
  private nodeTree = React.createRef<SVGSVGElement>()
  private nodeTreeGroup = React.createRef<SVGGElement>()

  shouldComponentUpdate(nextProps: Props, _): boolean {
    return (
      !isEqual(this.props.nodes, nextProps.nodes) ||
      !isEqual(this.props.edges, nextProps.edges)
    )
  }

  componentDidMount() {
    this.renderDag()
  }

  componentDidUpdate() {
    this.renderDag()
  }

  renderDag() {
    const g = new dagreD3.graphlib.Graph().setGraph({})

    for (const [id, node] of Object.entries(this.props.nodes)) {
      g.setNode(id, node)
    }

    for (const edge of this.props.edges) {
      g.setEdge(edge[0], edge[1], edge[2])
    }

    const svg = d3.select(this.nodeTree.current)
    const inner = d3.select(this.nodeTreeGroup.current)

    const render = new dagreD3.render()

    // @ts-ignore
    render(inner, g)

    const { height: gHeight, width: gWidth } = g.graph()
    const { height, width } = this.nodeTree.current.getBBox()
    const transX = width - gWidth + 40
    const transY = height - gHeight + 40
    svg.attr('height', height + 80)
    svg.attr('width', width + 80)
    // @ts-ignore
    inner.attr('transform', d3.zoomIdentity.translate(transX, transY))

    if (this.props.onNodeClick) {
      svg.selectAll('g.node').on('click', (id) => {
        this.props.onNodeClick(id as string)
      })
    }
  }

  render() {
    return (
      <svg
        className='dagre-d3'
        ref={this.nodeTree}
        width={this.props.height}
        height={this.props.width}
      >
        <g ref={this.nodeTreeGroup} />
      </svg>
    )
  }
}

export default Graph
