import './connections.scss'

import React from 'react'

import { Rect } from './resource'
import { Edge } from './schema'

const linePadding = 4
const markerHeight = 10

export interface Props {
  offsets: { [key: string]: Rect };
  adjacencyList: { [key: string]: Edge[] };
}

class Connections extends React.Component<Props, any> {
  render() {
    const connections: JSX.Element[] = []

    for (const key in this.props.adjacencyList) {
      if (this.props.offsets[key]) {
        const parentRect = this.props.offsets[key]
        if (parentRect) {
          const list = this.props.adjacencyList[key]

          connections.push(
            ...list.map((edge: Edge, id: number) => {
              const childRect = this.props.offsets[edge.node]
              let start = { x: parentRect.left, y: parentRect.top }
              let end = { x: childRect.left, y: childRect.top }
              const edgeID = `${key}-${id}`

              const isUp = start.y > end.y
              if (start.y > end.y) {
                const temp = start
                start = end
                end = temp
              }

              start.x += childRect.width / 2
              start.y += childRect.height
              end.x += childRect.width / 2

              let markerEnd = 'url(#down-arrow)'
              let markerStart = ''
              if (isUp) {
                markerEnd = ''
                markerStart = 'url(#up-arrow)'
                start.y += linePadding
                end.y -= linePadding
              } else {
                end.y -= markerHeight
              }

              return (
                <line
                  key={edgeID}
                  className='connection'
                  x1={start.x}
                  y1={start.y}
                  x2={end.x}
                  y2={end.y}
                  markerEnd={markerEnd}
                  markerStart={markerStart}
                />
              )
            }),
          )
        }
      }
    }

    return (
      <svg width='0' height='0' className='svgContainer'>
        <defs>
          <marker
            id='down-arrow'
            markerWidth='10'
            markerHeight='10'
            refX='0'
            refY='3'
            orient='auto'
            markerUnits='strokeWidth'
          >
            <path d='M0,0 L0,6 L9,3 z' fill='rgb(242,88,45)' />
          </marker>
          <marker
            id='up-arrow'
            markerWidth='10'
            markerHeight='10'
            refX='0'
            refY='3'
            orient='auto'
            markerUnits='strokeWidth'
          >
            <path d='M0,3 L9,6 L9,0 z' fill='rgb(242,88,45)' />
          </marker>
        </defs>
        {connections}
      </svg>
    )
  }
}

export default Connections
