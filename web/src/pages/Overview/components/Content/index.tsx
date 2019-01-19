import ResourceViewer from 'components/ContentTypes/ResourceViewer'
import { IResourceViewer } from 'components/ContentTypes/ResourceViewer/schema'
import Summary, { ISummary } from 'components/ContentTypes/Summary'
import Table, { ITable } from 'components/ContentTypes/Table'
import Grid, { IGrid } from 'components/ContentTypes/Grid'
import List, { IList } from 'components/ContentTypes/List'
import React from 'react'

interface Props {
  content: ContentType
}

export default function({ content }: Props) {
  const { metadata: { type } } = content
  switch (type) {
    case 'table':
      return <Table data={content as ITable} />
    case 'summary':
      return <Summary data={content as ISummary} />
    case 'resourceViewer':
      return <ResourceViewer data={content as IResourceViewer} />
    case 'grid':
      return <Grid data={content as IGrid} />
    case 'list':
      return <List data={content as IList} />
    default:
      return <div>Can not render content type</div>
  }
}
