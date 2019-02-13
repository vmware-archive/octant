import { View } from 'models'
import React, { Component } from 'react'
import { Tab, TabList, TabPanel, Tabs } from 'react-tabs'

import Content from '../Content'

interface Props {
  views: View[]
  currentTab: number
  setTab(index: number): void
}

export default class Renderer extends Component<Props> {
  constructor(props: Props) {
    super(props)

    this.onSelect = this.onSelect.bind(this)
  }

  onSelect(index: number) {
    const { setTab } = this.props
    setTab(index)
  }

  render() {
    const { views } = this.props

    if (views.length > 1) {
      return this.renderViewsWithTabs(views)
    } else {
      return this.renderViewsWithoutTabs(views[0])
    }
  }

  private renderViewsWithoutTabs = (view: View) => {
    return (
      <div className='component--primary'>
        <Content view={view} />
      </div>
    )
  }

  private renderViewsWithTabs = (views: View[]) => {
    const { currentTab } = this.props

    const tabs = []
    const panels = []

    views.forEach((view, index) => {
      const contents = [view]
      const tabContents = contents.map((content, i) => (
        <div key={i} className='component--primary'>
          <Content view={content} />
        </div>
      ))

      if (!view.title) {
        throw new Error('view does not have a title')
      }

      if (view.title && view.title.length === 1 && view.title[0].type === 'text') {
        tabs.push(<Tab key={index}>{view.title[0].value}</Tab>)
      } else {
        throw new Error('invalid view title')
      }

      panels.push(<TabPanel key={index}>{tabContents}</TabPanel>)
    })

    return (
      <Tabs onSelect={this.onSelect} defaultIndex={currentTab}>
        <TabList key={0}>{tabs}</TabList>
        {panels}
      </Tabs>
    )
  }
}
