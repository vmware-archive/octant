import ResourceViewer from 'components/ContentTypes/ResourceViewer'
import Summary from 'components/ContentTypes/Summary'
import Table from 'components/ContentTypes/Table'
import React from 'react'

export default function ({ content }) {
  const { type } = content
  switch (type) {
    case 'table': {
      return <Table data={content} />
    }
    case 'summary': {
      return <Summary data={content} />
    }
    case 'resourceviewer': {
      return <ResourceViewer schema={content} />
    }
    default: {
      return <div>Can not render content type</div>
    }
  }
}
