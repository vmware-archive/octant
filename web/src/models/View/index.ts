import { JSONContainers } from 'models/Containers'
import { JSONGrid } from 'models/Grid'
import { JSONLabels } from 'models/Labels'
import { JSONLink } from 'models/Link'
import { JSONList } from 'models/List'
import { JSONPanel } from 'models/Panel'
import { JSONQuadrant } from 'models/Quadrant'
import { JSONResourceViewer } from 'models/ResourceViewer'
import { JSONSelectors } from 'models/Selectors'
import { JSONSummary } from 'models/Summary'
import { JSONTable } from 'models/Table'
import { compareTextModel, JSONText, TextModel } from 'models/Text'
import { compareTimestampModel, JSONTimestamp, TimestampModel } from 'models/Timestamp'

export * from 'models/List'
export * from 'models/Summary'
export * from 'models/Table'
export * from 'models/Quadrant'
export * from 'models/Labels'
export * from 'models/Grid'
export * from 'models/Panel'
export * from 'models/Text'
export * from 'models/Link'
export * from 'models/Timestamp'
export * from 'models/Containers'
export * from 'models/Selectors'
export * from 'models/ResourceViewer'

export interface View {
  readonly type: string
  readonly title: string

  readonly isComparable?: boolean
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
    case 'grid':
      return new JSONGrid(ct)
    case 'labels':
      return new JSONLabels(ct)
    case 'link':
      return new JSONLink(ct)
    case 'list':
      return new JSONList(ct)
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
    default:
      throw new Error(
        `can't handle content response view '${ct.metadata.type}'`,
      )
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
