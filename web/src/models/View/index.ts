import { JSONContainers } from 'models/Containers'
import { JSONGrid } from 'models/Grid'
import { JSONLabels } from 'models/Labels'
import { JSONLink } from 'models/Link'
import { JSONList } from 'models/List'
import { JSONPanel } from 'models/Panel'
import { JSONQuadrant } from 'models/Quadrant'
import { JSONSelectors } from 'models/Selectors'
import { JSONSummary } from 'models/Summary'
import { JSONTable } from 'models/Table'
import { JSONText } from 'models/Text'
import { JSONTimestamp } from 'models/Timestamp'

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

export interface View {
  readonly type: string;
  readonly title: string;
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
