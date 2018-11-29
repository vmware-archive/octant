import './styles.scss';
import 'react-tabs/style/react-tabs.css';

import { getAPIBase, getContentsUrl, POLL_WAIT } from 'api';
import cx from 'classnames';
import Loading from 'components/Icons/Loading';
import Title from 'components/Title';
import _ from 'lodash';
import React, { Component } from 'react';
import { Tab, TabList, TabPanel, Tabs } from 'react-tabs';

import Content from './components/Content';

interface ContentSection {
  contents: any;
  title: string;
}

interface ContentResponse {
  views: ContentSection[];
}

export interface OverviewProps {
  title: string;
  path: string;
  namespace: string;
  isLoading: boolean;
  hasError: boolean;

  setIsLoading(isLoading: boolean);
  setHasError(hasError: boolean);
}

interface OverviewState {
  data: any;
}

export default class Overview extends Component<OverviewProps, OverviewState> {
  private source: any

  constructor(props: OverviewProps) {
    super(props)
    this.state = { data: null }
  }

  componentDidMount() {
    const { path, namespace } = this.props
    this.setEventSourceStream(path, namespace)
  }

  componentDidUpdate({ path: previousPath, namespace: previousNamespace }) {
    const { path, namespace } = this.props
    if (path !== previousPath || namespace !== previousNamespace) {
      this.setEventSourceStream(path, namespace)
    }
  }

  setEventSourceStream(path: string, namespace: string) {
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

  renderViewsWithTabs = (views: any) => {
    const tabs = []
    const panels = []

    _.forEach(views, (view, index) => {
      const tabContents = view.contents.map((content, i) => (
        <div key={i} className='component--primary'>
          <Content content={content} />
        </div>
      ))

      tabs.push(<Tab key={index}>{view.title}</Tab>)
      panels.push(<TabPanel key={index}>{tabContents}</TabPanel>)
    })

    return (
      <Tabs>
        <TabList key={0}>{tabs}</TabList>
        {panels}
      </Tabs>
    )
  }

  renderViewsWithoutTabs = (contents: object[]) => {
    return contents.map((content, i) => (
      <div key={i} className='component--primary'>
        <Content content={content} />
      </div>
    ))
  }

  renderUnknownContent = (hasError: boolean) => {
    const classNames = cx({
      'content-text': true,
      'empty-content-text': hasError === false,
      'error-content-text': hasError === true,
    })

    return (
      <div className='component--primary'>
        <h3 className={classNames}>
          {hasError === true
            ? 'Oops, something is not right, try again.'
            : 'There is nothing to display here'}
        </h3>
      </div>
    )
  }

  render() {
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
      const { views } = data
      title = data.title

      if (views.length > 1) {
        // there are multiple views
        mainContent = this.renderViewsWithTabs(views)
      } else if (views.length === 1) {
        // only a single view
        mainContent = this.renderViewsWithoutTabs(views[0].contents)
      } else {
        mainContent = this.renderUnknownContent(true)
      }
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
