import * as d3 from 'd3'
import dagreD3 from 'dagre-d3'
import React from 'react'
import isEqual from 'react-fast-compare'

import './graph.scss'

interface Props {
  nodes: any
  edges: any[]
  height?: string
  width?: string
  shapeRenders?: any
  onNodeClick(name: string): void
}

class Graph extends React.Component<Props> {
  private nodeTree = React.createRef<SVGSVGElement>()
  private nodeTreeGroup = React.createRef<SVGGElement>()

  shouldComponentUpdate(nextProps: Props, _): boolean {
    return !isEqual(this.props.nodes, nextProps.nodes) || !isEqual(this.props.edges, nextProps.edges)
  }

  componentDidMount() {
    this.renderDag()
  }

  componentDidUpdate() {
    this.renderDag()
  }

  renderDag() {
    const g = new dagreD3.graphlib.Graph().setGraph({})

    for (const [id, n] of Object.entries(this.props.nodes)) {
      g.setNode(id, n)
    }

    g.nodes().forEach((v) => {
      const node = g.node(v)
      node.rx = node.ry = 4
    })

    for (const edge of this.props.edges) {
      g.setEdge(edge[0], edge[1], edge[2])
    }

    const svg = d3.select(this.nodeTree.current)
    const inner = d3.select(this.nodeTreeGroup.current)

    const render = new dagreD3.render()

    // swallow type error that can happen if edges are in transition.
    try {
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
        svg
          .selectAll('g.node')
          .attr('data-for', (v) => g.node(v).id)
          .attr('data-tip', (v) => g.node(v).description)
          .attr('data-event', 'click')
          .on('click', (id: string) => {
            this.props.onNodeClick(id)
          })
      }
    } catch (e) {
      if (!(e instanceof TypeError)) {
        throw e
      }
    }
  }

  render() {
    return (
      <svg className='dagre-d3' ref={this.nodeTree} width={this.props.height} height={this.props.width}>
        <g ref={this.nodeTreeGroup} />
      </svg>
    )
  }
}

export default Graph
