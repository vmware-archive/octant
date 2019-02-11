import { View } from 'models'
import React from 'react'
import { renderView } from 'views'

interface Props {
  view: View
}

export default function({ view }: Props) {
  const supportedTypes = new Set(['table', 'summary', 'resourceViewer', 'grid', 'list', 'flexlayout'])
  if (supportedTypes.has(view.type)) {
    return renderView(view)
  }

  return <div>Can not render content type {view.type}</div>
}
