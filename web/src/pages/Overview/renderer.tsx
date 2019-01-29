import { View } from 'models/View'
import React from 'react'
import { Tab, TabList, TabPanel, Tabs } from 'react-tabs'

import Content from './components/Content'

export default class Renderer {
  constructor(private readonly views: View[]) {}

  content(): JSX.Element {
    if (this.views.length > 1) {
      return this.renderViewsWithTabs(this.views)
    } else {
      return this.renderViewsWithoutTabs(this.views[0])
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
    const tabs = []
    const panels = []

    views.forEach((view, index) => {
      const contents = [view]
      const tabContents = contents.map((content, i) => (
        <div key={i} className='component--primary'>
          <Content view={content} />
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
}
