import { View } from 'models/View'
import React from 'react'
import { renderView } from 'views'

interface Props {
  view: View;
}

export default function({ view }: Props) {
  switch (view.type) {
    case 'table':
    case 'summary':
    // TODO: re-enable resource view
    // case 'resourceViewer':
    //   return <ResourceViewer data={content as IResourceViewer} />
    case 'grid':
    case 'list':
      return renderView(view)
    default:
      return <div>Can not render content type</div>
  }
}
