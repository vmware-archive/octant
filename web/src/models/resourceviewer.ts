import { Edge, ResourceViewerModel, ResourceViewerNode, TitleView, toTitle } from 'models'

export class JSONResourceViewer implements ResourceViewerModel {
  readonly edges: { [key: string]: Edge[] }
  readonly nodes: { [key: string]: ResourceViewerNode } = {}
  readonly type = 'resourceViewer'
  readonly title: TitleView

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

    this.edges = ct.config.edges

    Object.keys(ct.config.nodes).forEach((name) => {
      const ctNode = ct.config.nodes[name]
      const details = ctNode.details || []
      const node: ResourceViewerNode = {
        name: ctNode.name,
        apiVersion: ctNode.apiVersion,
        kind: ctNode.kind,
        status: ctNode.status,
        details: toTitle(details),
      }

      this.nodes[name] = node
    })
  }
}
