import { YAMLViewerModel } from 'models'
import React, { Component } from 'react'
import yaml from 'react-syntax-highlighter/dist/cjs/languages/hljs/yaml'
import SyntaxHighlighter from 'react-syntax-highlighter/dist/cjs/light'
import atomOneDark from 'react-syntax-highlighter/dist/cjs/styles/hljs/atom-one-dark'

import './styles.scss'

interface Props {
  view: YAMLViewerModel
}

SyntaxHighlighter.registerLanguage('yaml', yaml)

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
