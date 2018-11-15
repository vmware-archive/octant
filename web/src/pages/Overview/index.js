import React, { Component } from 'react'
import cx from 'classnames'
import Loading from 'components/Icons/Loading'
import Title from 'components/Title'
import { getAPIBase, getContentsUrl, POLL_WAIT } from 'api'
import Content from './components/Content'
import './styles.scss'

export default class Overview extends Component {
  constructor (props) {
    super(props)
    this.state = { contents: null }
  }

  async componentDidUpdate ({
    path: previousPath,
    namespace: previousNamespace
  }) {
    const { path, namespace } = this.props
    if (path !== previousPath || namespace !== previousNamespace) {
      // clear state and this.source on change
      if (this.source) {
        this.source.close()
        this.source = null
      }

      this.props.toggleIsLoading(true)
      this.setState({ contents: null })

      const url = getContentsUrl(path, namespace, POLL_WAIT)

      this.source = new window.EventSource(`${getAPIBase()}/${url}`)

      this.source.addEventListener('message', (e) => {
        const data = JSON.parse(e.data)
        this.setState({ contents: data.contents })
        this.props.toggleIsLoading(false)
      })

      // if EventSource error clear close
      this.source.addEventListener('error', () => {
        this.setState({ contents: null })
        this.props.toggleIsLoading(false)
        this.props.toggleHasError(true)

        this.source.close()
        this.source = null
      })
    }
  }

  render () {
    const { title, isLoading, hasError } = this.props
    const { contents } = this.state
    let mainContent
    if (isLoading) {
      mainContent = (
        <div className='loading-parent'>
          <Loading />
        </div>
      )
    } else if (contents && contents.length) {
      mainContent = contents.map((content, i) => (
        <div key={i} className='component--primary'>
          <Content content={content} />
        </div>
      ))
    } else {
      const cnames = cx({
        'content-text': true,
        'empty-content-text': hasError === false,
        'error-content-text': hasError === true
      })
      mainContent = (
        <div className='component--primary'>
          <h3 className={cnames}>
            {hasError === true
              ? "Oops, something's not right, try again."
              : 'There is nothing to display here'}
          </h3>
        </div>
      )
    }

    return (
      <div className='overview'>
        <Title title={title} />
        <div className='main'>{mainContent}</div>
      </div>
    )
  }
}
