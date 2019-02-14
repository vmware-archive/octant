import { JSONContainers } from './containers'
import { JSONFlexLayout } from './flexlayout'
import { JSONGrid } from './grid'
import { JSONLabels } from './labels'
import { JSONLink } from './link'
import { JSONList } from './list'
import { JSONLogs } from './logs'
import { JSONPanel } from './panel'
import { JSONQuadrant } from './quadrant'
import { JSONResourceViewer } from './resourceviewer'
import { JSONSelectors } from './selectors'
import { JSONSummary } from './summary'
import { JSONTable } from './table'
import { compareTextModel, JSONText } from './text'
import { compareTimestampModel, JSONTimestamp } from './timestamp'
import { JSONYAMLViewer } from './yaml'

export type TitleView = Array<TextModel | LinkModel>

export interface View {
  readonly type: string
  readonly title?: TitleView

  readonly isComparable?: boolean
}

export interface ContainerDef {
  name: string
  image: string
}

export interface ContainersModel extends View {
  containerDefs: ContainerDef[]
}

export interface GridModel extends View {
  panels: PanelModel[]
}

export interface FlexLayoutItem {
  width: number
  view: View
}

export type FlexLayoutSection = FlexLayoutItem[]

export interface FlexLayoutModel extends View {
  sections: FlexLayoutSection[]
}

export interface LabelsModel extends View {
  labels: { [key: string]: string }
}

export interface LinkModel extends View {
  ref: string
  value: string
}

export interface PanelPosition {
  x: number
  y: number
  w: number
  h: number
}

export interface PanelModel extends View {
  readonly position: PanelPosition
  readonly content: View
}

export interface QuadrantModel extends View {
  nw: QuadrantSector
  ne: QuadrantSector
  sw: QuadrantSector
  se: QuadrantSector
}

export interface QuadrantSector {
  value: string
  label: string
}

export interface ListModel extends View {
  readonly items: View[]
}

export interface LogsModel extends View {
  readonly namespace: string
  readonly name: string
  readonly containers: string[]
}

export interface ResourceViewerNode {
  readonly name: string
  readonly apiVersion: string
  readonly kind: string
  readonly status: string
  readonly details: TitleView
}

export interface Edge {
  readonly node: string
  readonly type: string
}

export interface ResourceViewerModel extends View {
  readonly edges: { [key: string]: Edge[] }
  readonly nodes: { [key: string]: ResourceViewerNode }
  readonly type: 'resourceViewer'
}

export interface LabelSelector {
  key: string
  value: string
  type: string
}

export interface ExpressionSelector {
  key: string
  operator: string
  values: string[]
  type: string
}

export interface SelectorsModel extends View {
  selectors: Array<LabelSelector | ExpressionSelector>
}

export interface SummaryItem {
  readonly header: string
  readonly content: View
}

export interface SummaryModel extends View {
  items: SummaryItem[]
}

export interface TableRow {
  [key: string]: View
}
export interface TableModel extends View {
  readonly columns: Array<{ name: string; accessor: string }>
  readonly rows: TableRow[]
  readonly emptyContent: string
}

export interface TextModel extends View {
  value: string
}

export interface TimestampModel extends View {
  timestamp: number
}

export interface YAMLViewerModel extends View {
  data: string
}

export function toTitle(parts?: ContentType[]): TitleView | undefined {
  return parts.map((part) => {
    const view = viewFromContentType(part)
    switch (view.type) {
      case 'text':
        return view as TextModel
      case 'link':
        return view as LinkModel
      default:
        throw new Error(`invalid title type ${view.type}`)
    }
  })
}

export function instanceOfComparableView(object: any): object is View {
  return object.isComparable
}

export function viewFromContentType(ct: ContentType): View {
  if (!ct) {
    return null
  }

  switch (ct.metadata.type) {
    case 'containers':
      return new JSONContainers(ct)
    case 'flexlayout':
      return new JSONFlexLayout(ct)
    case 'grid':
      return new JSONGrid(ct)
    case 'labels':
      return new JSONLabels(ct)
    case 'link':
      return new JSONLink(ct)
    case 'list':
      return new JSONList(ct)
    case 'logs':
      return new JSONLogs(ct)
    case 'panel':
      return new JSONPanel(ct)
    case 'quadrant':
      return new JSONQuadrant(ct)
    case 'resourceViewer':
      return new JSONResourceViewer(ct)
    case 'selectors':
      return new JSONSelectors(ct)
    case 'summary':
      return new JSONSummary(ct)
    case 'table':
      return new JSONTable(ct)
    case 'text':
      return new JSONText(ct)
    case 'timestamp':
      return new JSONTimestamp(ct)
    case 'yaml':
      return new JSONYAMLViewer(ct)
    default:
      throw new Error(`can't handle content response view '${ct.metadata.type}'`)
  }
}

export function compareModel(a: View, b: View): number {
  if (a.type !== b.type) {
    throw new Error(`unable to compare ${a.type} to ${b.type}`)
  }

  if (!a.isComparable && !b.isComparable) {
    throw new Error(`views are not comparable`)
  }

  switch (a.type) {
    case 'text':
      return compareTextModel(a as TextModel, b as TextModel)
    case 'timestamp':
      return compareTimestampModel(a as TimestampModel, b as TimestampModel)
    default:
      return 0
  }
}
