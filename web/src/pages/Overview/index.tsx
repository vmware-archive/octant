import cx from 'classnames'
import Loading from 'components/Icons/Loading'
import Title from 'components/Title'
import JSONContentResponse from 'models/ContentResponse'
import React, { Component } from 'react'

import Renderer from './renderer'
import './styles.scss'

interface Props {
  title: string
  isLoading: boolean
  hasError: boolean
  errorMessage: string

  data: JSONContentResponse

  setError(hasError: boolean, errorMessage?: string): void
}

export default class Overview extends Component<Props> {
  constructor(props: Props) {
    super(props)
  }

  renderUnknownContent = (hasError: boolean) => {
    const classNames = cx({
      'content-text': true,
      'empty-content-text': hasError === false,
      'error-content-text': hasError === true,
    })

    const { errorMessage } = this.props

    return (
      <div className='component--primary'>
        <h3 className={classNames}>{hasError === true ? errorMessage : 'There is nothing to display here'}</h3>
      </div>
    )
  }

  render() {
    const { isLoading, hasError, data } = this.props
    let title = ''
    let mainContent = <div />
    if (isLoading) {
      mainContent = (
        <div className='loading-parent'>
          <Loading />
        </div>
      )
    } else if (data) {
      title = data.title
      const renderer = new Renderer(data.views)
      mainContent = renderer.content()
    } else {
      // no views or an error
      mainContent = this.renderUnknownContent(hasError)
    }

    return (
      <div className='overview'>
        <Title title={title} />
        <div className='main'>{mainContent}</div>
      </div>
    )
  }
}
