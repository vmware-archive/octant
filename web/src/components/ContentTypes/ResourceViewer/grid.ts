import _ from 'lodash'
import React from 'react'

import { AdjacencyList, Edge, ResourceObjects } from './schema'

type Rows = Array<Set<string>>

export default class Grid {
  constructor(
    private adjacencyList: AdjacencyList,
    private objects: ResourceObjects,
  ) {}

  isSink = (name: string): boolean => {
    const edgeKeys = new Set(Object.keys(this.adjacencyList))
    const objectKeys = new Set(Object.keys(this.objects))
    const sinks = new Set([...objectKeys].filter((key) => !edgeKeys.has(key)))
    return [...sinks].indexOf(name) >= 0
  }

  create(): Rows {
    const sorted = this.sort()
    const rows: Rows = []

    sorted.forEach((name: string) => {
      if (this.isSink(name)) {
        // handle the sink and continue
        if (rows.length === 0) {
          rows.push(new Set<string>([name]))
        } else {
          rows[0].add(name)
        }
        return
      }

      const curRow = this.rowForNode(rows, name)
      if (curRow > 0) {
        // append to rows
        if (rows[curRow]) {
          // use existing row
          rows[curRow].add(name)
        } else {
          // create a new row
          rows[curRow] = new Set<string>([name])
        }
      } else {
        // prepend to rows
        rows.unshift(new Set<string>([name]))
      }
    })

    return rows
  }

  sort(): string[] {
    const visited: { [key: string]: boolean } = {}
    const sorted: string[] = []

    const visit = (id: string, ancestors: { [key: string]: boolean }) => {
      if (visited[id]) {
        return
      }

      ancestors[id] = true
      visited[id] = true

      if (this.adjacencyList[id]) {
        this.adjacencyList[id].forEach((edge: Edge) => {
          if (ancestors[edge.node]) {
            throw new Error(`detected loop with ${edge.node}`)
          }

          visit(edge.node, ancestors)
        })
      }

      sorted.unshift(id)
    }
    Object.keys(this.objects).forEach((id: string) => {
      visit(id, {})
    })

    return sorted.reverse()
  }

  rowForNode(rows: Rows, name: string): number {
    const edges = _.map(this.adjacencyList[name], (edge) => edge.node)

    let high = 0
    _.forEach(rows, (s, i) => {
      _.forEach([...s], (edgeName) => {
        if (edges.indexOf(edgeName) >= 0) {
          high = i
        }
      })
    })

    const object = this.objects[name]
    if (object.isNetwork) {
      return high - 1
    } else {
      return high + 1
    }
  }
}
