import Containers from 'components/Containers'
import Grid from 'components/Grid'
import Labels from 'components/Labels'
import List from 'components/List'
import Quadrant from 'components/Quadrant'
import ResourceViewer from 'components/ResourceViewer'
import Selectors from 'components/Selector'
import Summary from 'components/Summary'
import Table from 'components/Table'
import TextView from 'components/TextView'
import Timestamp from 'components/Timestamp'
import WebLink from 'components/WebLink'
import {
  ContainersModel,
  GridModel,
  LabelsModel,
  LinkModel,
  ListModel,
  QuadrantModel,
  ResourceViewerModel,
  SelectorsModel,
  SummaryModel,
  TableModel,
  TextModel,
  TimestampModel,
  View,
} from 'models'
import React from 'react'

export function renderView(view: View): JSX.Element {
  switch (view.type) {
    case 'grid':
      return <Grid view={view as GridModel} />
    case 'containers':
      return <Containers view={view as ContainersModel} />
    case 'labels':
      return <Labels view={view as LabelsModel} />
    case 'link':
      return <WebLink view={view as LinkModel} />
    case 'list':
      return <List view={view as ListModel} />
    case 'quadrant':
      return <Quadrant view={view as QuadrantModel} />
    case 'resourceViewer':
      return <ResourceViewer view={view as ResourceViewerModel} />
    case 'selectors':
      return <Selectors view={view as SelectorsModel} />
    case 'summary':
      return <Summary view={view as SummaryModel} />
    case 'table':
      return <Table view={view as TableModel} />
    case 'text':
      return <TextView view={view as TextModel} />
    case 'timestamp':
      return <Timestamp view={view as TimestampModel} />
    default:
      throw new Error(`unable to render view of type ${view.type}`)
  }
}
