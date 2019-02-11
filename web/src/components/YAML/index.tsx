import { YAMLViewerModel } from 'models'
import React, { Component } from 'react'
import SyntaxHighlighter from 'react-syntax-highlighter/dist'
import { atomOneDark } from 'react-syntax-highlighter/dist/styles/hljs'

import './styles.scss'

interface Props {
  view: YAMLViewerModel
}

export default class extends Component<Props> {
  constructor(props: Props) {
    super(props)
  }

  render() {
    const { view } = this.props
    return (
      <div className='yamlViewer'>
        <SyntaxHighlighter language='yaml' style={atomOneDark}>
          {view.data}
        </SyntaxHighlighter>
      </div>
    )
  }
}
