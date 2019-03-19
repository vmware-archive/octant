import { View } from 'models'
import queryString from 'query-string'
import React, { Component } from 'react'
import { RouteComponentProps, withRouter } from 'react-router'
import { Tab, TabList, TabPanel, Tabs } from 'react-tabs'

import Content from '../Content'

const tabIndicator = 'view'

interface Props extends RouteComponentProps {
  views: View[]
}

interface State {
  currentTab: number
}

class Renderer extends Component<Props, State> {
  constructor(props: Props) {
    super(props)

    this.onSelect = this.onSelect.bind(this)

    this.state = {
      currentTab: 0,
    }
  }

  componentDidMount() {
    this.checkTab()
  }

  componentDidUpdate() {
    const { views } = this.props

    const currentTabName = this.tabName()

    const currentTabID = views.map((view) => view.accessor).indexOf(currentTabName)

    if (this.state.currentTab !== currentTabID) {
      this.setState({ currentTab: currentTabID })
    }
  }

  tabName(): string {
    const values = queryString.parse(this.props.location.search)

    let currentTabName: string

    const keys = Object.keys(values)
    if (keys.indexOf(tabIndicator) > -1) {
      if (typeof values[tabIndicator] === 'string') {
        currentTabName = values[tabIndicator] as string
      }
    }

    return currentTabName
  }

  checkTab() {
    let currentTabName = this.tabName()

    let currentTabID: number

    const { views } = this.props
    if (currentTabName) {
      currentTabID = views.map((view) => view.accessor).indexOf(currentTabName)
    } else {
      currentTabID = 0
      currentTabName = views[0].accessor
    }

    this.setState({ currentTab: currentTabID })
    this.setTab(currentTabName)
  }

  setTab(accessor: string) {
    const { history } = this.props

    const newSearch = `?${queryString.stringify({ [tabIndicator]: accessor })}`

    if (history.location.search !== newSearch) {
      history.push({
        search: newSearch,
      })
    }
  }

  onSelect(index: number) {
    const { views } = this.props
    if (views[index]) {
      this.setState({ currentTab: index })
      this.setTab(views[index].accessor)
    }
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
    let { currentTab } = this.state

    if (currentTab < 0) {
      currentTab = 0
    }

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
      <Tabs onSelect={this.onSelect} selectedIndex={currentTab}>
        <TabList key={0}>{tabs}</TabList>
        {panels}
      </Tabs>
    )
  }
}

export default withRouter(Renderer)
