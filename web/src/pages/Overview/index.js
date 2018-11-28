import './styles.scss'
import { getAPIBase, getContentsUrl, POLL_WAIT } from 'api'
import cx from 'classnames'
import Loading from 'components/Icons/Loading'
import Title from 'components/Title'
import React, { Component } from 'react'
import Content from './components/Content'

export default class Overview extends Component {
  constructor (props) {
    super(props)
    this.state = { data: null }
  }

  componentDidMount () {
    const { path, namespace } = this.props
    this.setEventSourceStream(path, namespace)
  }

  componentDidUpdate ({ path: previousPath, namespace: previousNamespace }) {
    const { path, namespace } = this.props
    if (path !== previousPath || namespace !== previousNamespace) {
      this.setEventSourceStream(path, namespace)
    }
  }

  setEventSourceStream (path, namespace) {
    // clear state and this.source on change
    if (this.source) {
      this.source.close()
      this.source = null
    }

    if (!path || !namespace) return

    this.props.setIsLoading(true)
    this.setState({ data: null })

    const url = getContentsUrl(path, namespace, POLL_WAIT)

    this.source = new window.EventSource(`${getAPIBase()}/${url}`)

    this.source.addEventListener('message', (e) => {
      const data = JSON.parse(e.data)
      this.setState({ data })
      this.props.setIsLoading(false)
    })

    // if EventSource error clear close
    this.source.addEventListener('error', () => {
      this.setState({ data: null })
      this.props.setIsLoading(false)
      this.props.setHasError(true)

      this.source.close()
      this.source = null
    })
  }

  render () {
    const { isLoading, hasError } = this.props
    const { data } = this.state
    let title
    let mainContent
    if (isLoading) {
      mainContent = (
        <div className='loading-parent'>
          <Loading />
        </div>
      )
    } else if (data) {
      const { views, default_view: defaultView } = data
      const view = views[defaultView]
      const { contents } = view
      title = view.title // eslint-disable-line prefer-destructuring
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
