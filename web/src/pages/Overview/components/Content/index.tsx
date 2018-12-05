import ResourceViewer from 'components/ContentTypes/ResourceViewer'
import { Schema } from 'components/ContentTypes/ResourceViewer/schema'
import Summary from 'components/ContentTypes/Summary'
import Table from 'components/ContentTypes/Table'
import React from 'react'

interface Props {
  content: Content | Schema
}

export default function({ content }: Props) {
  const { type } = content
  switch (type) {
    case 'table': {
      return <Table data={content as ContentTable} />
    }
    case 'summary': {
      return <Summary data={content as ContentSummary} />
    }
    case 'resourceviewer': {
      return <ResourceViewer schema={content as Schema} />
    }
    default: {
      return <div>Can not render content type</div>
    }
  }
}
