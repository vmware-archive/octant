import './styles.scss'
import 'react-tabs/style/react-tabs.css'

import cx from 'classnames'
import Loading from 'components/Icons/Loading'
import Title from 'components/Title'
import _ from 'lodash'
import React, { Component } from 'react'
import { Tab, TabList, TabPanel, Tabs } from 'react-tabs'

import Content from './components/Content'

interface Props {
  title: string;
  isLoading: boolean;
  hasError: boolean;
  errorMessage: string;

  data: {
    content: {
      title: string;
      viewComponents: ContentType[];
    };
  };

  setError(hasError: boolean, errorMessage?: string): void;
}

export default class Overview extends Component<Props> {
  constructor(props: Props) {
    super(props)
  }

  renderViewsWithTabs = (views: ContentType[]) => {
    const tabs = []
    const panels = []

    _.forEach(views, (view, index) => {
      const contents = [view]
      const tabContents = _.map(contents, (content, i) => (
        <div key={i} className='component--primary'>
          <Content content={content} />
        </div>
      ))

      tabs.push(<Tab key={index}>{view.metadata.title}</Tab>)
      panels.push(<TabPanel key={index}>{tabContents}</TabPanel>)
    })

    return (
      <Tabs>
        <TabList key={0}>{tabs}</TabList>
        {panels}
      </Tabs>
    )
  }

  renderViewsWithoutTabs = (view: ContentType) => {
    return (
      <div className='component--primary'>
        <Content content={view} />
      </div>
    )
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
        <h3 className={classNames}>
          {hasError === true
            ? errorMessage
            : 'There is nothing to display here'}
        </h3>
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
    } else if (data && data.content.viewComponents) {
      const views = data.content.viewComponents
      title = data.content.title

      if (views.length > 1) {
        // there are multiple views
        mainContent = this.renderViewsWithTabs(views)
      } else if (views.length === 1) {
        // only a single view
        mainContent = this.renderViewsWithoutTabs(views[0])
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
