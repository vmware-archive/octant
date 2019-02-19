import Annotations from 'components/Annotations'
import Containers from 'components/Containers'
import FlexLayout from 'components/FlexLayout'
import Grid from 'components/Grid'
import Labels from 'components/Labels'
import List from 'components/List'
import Logs from 'components/Logs'
import Quadrant from 'components/Quadrant'
import ResourceViewer from 'components/ResourceViewer'
import Selectors from 'components/Selector'
import Summary from 'components/Summary'
import Table from 'components/Table'
import TextView from 'components/TextView'
import Timestamp from 'components/Timestamp'
import WebLink from 'components/WebLink'
import YAML from 'components/YAML'
import {
  AnnotationsModel,
  ContainersModel,
  FlexLayoutModel,
  GridModel,
  LabelsModel,
  LinkModel,
  ListModel,
  LogsModel,
  QuadrantModel,
  ResourceViewerModel,
  SelectorsModel,
  SummaryModel,
  TableModel,
  TextModel,
  TimestampModel,
  View,
  YAMLViewerModel,
} from 'models'
import React from 'react'

export function renderView(view: View, extraProps?: any): JSX.Element {
  switch (view.type) {
    case 'annotations':
      return <Annotations view={view as AnnotationsModel} {...extraProps} />
    case 'containers':
      return <Containers view={view as ContainersModel} {...extraProps} />
    case 'flexlayout':
      return <FlexLayout view={view as FlexLayoutModel} {...extraProps} />
    case 'grid':
      return <Grid view={view as GridModel} {...extraProps} />
    case 'labels':
      return <Labels view={view as LabelsModel} {...extraProps} />
    case 'link':
      return <WebLink view={view as LinkModel} {...extraProps} />
    case 'list':
      return <List view={view as ListModel} {...extraProps} />
    case 'logs':
      return <Logs view={view as LogsModel} {...extraProps} />
    case 'quadrant':
      return <Quadrant view={view as QuadrantModel} {...extraProps} />
    case 'resourceViewer':
      return <ResourceViewer view={view as ResourceViewerModel} {...extraProps} />
    case 'selectors':
      return <Selectors view={view as SelectorsModel} {...extraProps} />
    case 'summary':
      return <Summary view={view as SummaryModel} {...extraProps} />
    case 'table':
      return <Table view={view as TableModel} {...extraProps} />
    case 'text':
      return <TextView view={view as TextModel} {...extraProps} />
    case 'timestamp':
      return <Timestamp view={view as TimestampModel} {...extraProps} />
    case 'yaml':
      return <YAML view={view as YAMLViewerModel} {...extraProps} />
    default:
      throw new Error(`unable to render view of type ${view.type}`)
  }
}
