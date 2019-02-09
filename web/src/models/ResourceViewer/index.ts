import { View } from 'models/View'

export interface ResourceViewerNode {
  name: string
  apiVersion: string
  kind: string
  status: string
}

interface Edge {
  node: string
  type: string
}

export interface ResourceViewerModel extends View {
  readonly edges: { [key: string]: Edge[] }
  readonly nodes: { [key: string]: ResourceViewerNode }
  readonly title: string
  readonly type: 'resourceViewer'
}

export class JSONResourceViewer implements ResourceViewerModel {
  readonly edges: { [key: string]: Edge[] }
  readonly nodes: { [key: string]: ResourceViewerNode }
  readonly title: string
  readonly type = 'resourceViewer'

  constructor(ct: ContentType) {
    this.title = ct.metadata.title
    this.edges = ct.config.edges
    this.nodes = ct.config.nodes
  }
}
