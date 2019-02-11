import { View } from 'models'
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
      <Tabs>
        <TabList key={0}>{tabs}</TabList>
        {panels}
      </Tabs>
    )
  }
}
