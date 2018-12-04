import { AdjacencyList, Edge, ResourceObjects } from './schema'

type Rows = Array<Set<string>>

export default class Grid {
  constructor(private adjacencyList: AdjacencyList, private objects: ResourceObjects) {}

  create(): Rows {
    const sorted = this.sort()
    const rows: Rows = []

    sorted.forEach((name: string) => {
      if (this.isSink(name)) {
        if (rows.length === 0) {
          rows.push(new Set<string>([name]))
        } else {
          rows[0].add(name)
        }
        return
      }

      const node = this.objects[name]
      this.adjacencyList[name].forEach((edge: Edge) => {
        const id = this.rowForNode(rows, edge.node)
        if (node.isNetwork) {
          if (id === 0) {
            // create new row in front of the current row
            rows.unshift(new Set<string>([name]))
          } else {
            // add to existing previous row
            rows[id - 1].add(name)
          }
        } else {
          if (rows.length - 1 === id) {
            // add a new row
            rows.push(new Set<string>([name]))
          } else {
            // add to existing row
            rows[id + 1].add(name)
          }
        }
      })
    })

    return rows
  }

  sort(): string[] {
    const visited: { [key: string]: boolean } = {}
    const sorted: string[] = []

    const visit = (id: string, ancestors: string[]) => {
      if (visited[id]) {
        return
      }

      ancestors.push(id)
      visited[id] = true

      if (this.adjacencyList[id]) {
        this.adjacencyList[id].forEach((edge: Edge) => {
          if (ancestors.indexOf(edge.node) >= 0) {
            throw new Error(`detected loop with ${edge.node}`)
          }

          visit(edge.node, ancestors.map((value) => value))
        })
      }

      sorted.unshift(id)
    }

    Object.keys(this.objects).forEach((id: string) => {
      visit(id, [])
    })

    return sorted.reverse()
  }

  isSink(name: string): boolean {
    const edgeKeys = new Set(Object.keys(this.adjacencyList))
    const objectKeys = new Set(Object.keys(this.objects))
    const sinks = new Set([...objectKeys].filter((key) => !edgeKeys.has(key)))
    return [...sinks].indexOf(name) >= 0
  }

  rowForNode = (rows: Rows, name: string): number => {
    const nodeRows = rows.map(
      (row: Set<string>, index: number) => ({
        row,
        index,
      }),
    )
    for (const { row, index } of nodeRows) {
      if ([...row].indexOf(name) >= 0) {
        return index
      }
    }

    throw new Error(`could not find node ${name}`)
  }
}
