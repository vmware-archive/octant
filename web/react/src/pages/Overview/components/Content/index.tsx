import { ListModel, View } from 'models'
import React from 'react'
import { renderView } from 'views'

interface Props {
  view: View
}

export default class Content extends React.Component<Props> {
  private ref: React.RefObject<HTMLDivElement>
  private scrolled = false

  constructor(props: Props) {
    super(props)
    this.ref = React.createRef()
  }

  componentDidUpdate() {
    if (!this.scrolled) {
      window.scrollTo(0, 0)
      this.scrolled = true
    }
  }

  render() {
    const { view } = this.props

    let rendered: JSX.Element

    const supportedTypes = new Set(['table', 'summary', 'resourceViewer', 'grid', 'list', 'flexlayout', 'yaml', 'logs'])
    if (supportedTypes.has(view.type)) {
      if (view.type === 'list') {
        const list = view as ListModel
        if (list.items.length > 1) {
          rendered = <div ref={this.ref}>{renderView(view, { isOverview: true })}</div>
        } else {
          rendered = <div ref={this.ref}>{renderView(view)}</div>
        }
      } else {
        rendered = <div ref={this.ref}>{renderView(view)}</div>
      }

      return rendered
    }

    return <div>Can not render content type {view.type}</div>
  }
}
