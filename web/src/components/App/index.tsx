import { ContentsUrlParams, getAPIBase, getContentsUrl, POLL_WAIT, setNamespace } from 'api'
import Header from 'components/Header'
import ResourceFiltersContext from 'contexts/resource-filters'
import _ from 'lodash'
import JSONContentResponse, { Parse } from 'models/contentresponse'
import Overview from 'pages/Overview'
import React, { Component } from 'react'
import { Redirect, Route, RouteComponentProps, Switch, withRouter } from 'react-router-dom'
import ReactTooltip from 'react-tooltip'

import Navigation from '../Navigation'

import getInitialState from './state/getInitialState'
import './styles.scss'

interface AppState {
  isLoading: boolean
  hasError: boolean
  errorMessage: string
  navigation: { sections: NavigationSectionType[] }
  currentNavLinkPath: NavigationSectionType[]
  namespaceOption: NamespaceOption
  namespaceOptions: NamespaceOption[]
  title: string
  contentResponse: JSONContentResponse
  resourceFilters: string[]
}

class App extends Component<RouteComponentProps, AppState> {
  private lastFetchedNamespace: string

  private source: any

  constructor(props) {
    super(props)
    this.state = {
      title: '',
      isLoading: true, // to do the initial data fetch
      hasError: false,
      errorMessage: '',
      navigation: null,
      currentNavLinkPath: [],
      namespaceOption: { label: 'default', value: 'default' },
      namespaceOptions: [],
      contentResponse: null,
      resourceFilters: [],
    }
  }

  async componentDidMount() {
    const { location: initialLocation } = this.props
    const initialState = await getInitialState(initialLocation.pathname)
    this.setState(initialState as AppState)
    this.setEventSourceStream(this.props.location.pathname, this.state.namespaceOption.value)
  }

  componentDidUpdate({ location: previousLocation }, { namespaceOption: previousNamespace }) {
    const { location } = this.props
    const { namespaceOption } = this.state

    const namespace = namespaceOption.value
    const prevNamespace = previousNamespace.value

    if (location.pathname !== previousLocation.pathname || namespace !== prevNamespace) {
      this.setEventSourceStream(location.pathname, namespace)
    }

    // this is required to make tool tips show.
    ReactTooltip.rebuild()
  }

  componentWillUnmount(): void {
    if (this.source) {
      this.source.close()
      this.source = null
    }
  }

  setEventSourceStream(path: string, namespace: string) {
    // clear state and this.source on change
    if (this.source) {
      this.source.close()
      this.source = null
    }

    if (!path || !namespace) return

    const params: ContentsUrlParams = {
      poll: POLL_WAIT,
    }

    const { resourceFilters } = this.state
    if (resourceFilters && resourceFilters.length) params.filter = resourceFilters

    const url = getContentsUrl(path, namespace, params)

    this.source = new window.EventSource(`${getAPIBase()}/${url}`)
    this.setState({ isLoading: true })

    this.source.addEventListener('message', (e) => {
      const contentResponse = Parse(e.data)

      this.setState({
        contentResponse,
        isLoading: false,
      })
    })

    this.source.addEventListener('navigation', (e) => {
      const data = JSON.parse(e.data)
      this.setState({ navigation: data })
    })

    this.source.addEventListener('namespaces', (e) => {
      const data = JSON.parse(e.data)
      const updated = data.namespaces.map((ns) => ({
        label: ns,
        value: ns,
      }))

      // TODO if current namespace is not in list, redirect to the
      // the first item in the list.
      this.setState({ namespaceOptions: updated })
    })

    this.source.addEventListener('error', () => {
      this.setState({ isLoading: false })
      this.setError(true, 'Looks like the back end source has gone away. Retrying...')
    })
  }

  onNamespaceChange = async (namespaceOption) => {
    this.setState({
      isLoading: true,
      hasError: false,
    })

    const { value } = namespaceOption
    const { history } = this.props
    try {
      const { currentNavLinkPath } = this.state
      const path = this.removeObjectNameFromPath(history.location.pathname)
      const fixedPath = this.injectNamespace(path, value)
      history.push(fixedPath)
      this.setState({ namespaceOption, isLoading: false, hasError: false })
    } catch (e) {
      this.setState({ namespaceOption, isLoading: false, hasError: true })
    }
  }

  // Injects a namespace into an optionally namespaced overview content url
  injectNamespace = (url: string, namespace: string) => {
    // Insert a /namespace/... segment if it was missing
    const addNamespace = /(\/api\/v1)?\/content\/overview(\/namespace\/[^/]+)?(.*)/
    const withNamespace = url.replace(addNamespace, '$1/content/overview/namespace/...$3')

    // Now replace the actual namespace with the new namespace
    const re = /(\/api\/v1)?\/content\/overview\/namespace\/[^/]+(.*)/
    const final = withNamespace.replace(re, '$1/content/overview/namespace/' + namespace + '$2')

    return final
  }

  // Returns whether an overview content url is namespaced
  hasNamespace = (url: string) => {
    const re = /(\/api\/v1)?\/content\/overview\/namespace\/[^/]+(.*)/
    return re.test(url)
  }

  removeObjectNameFromPath = (url: string) => {
    const parts = url.split('/')

    // e.g. /content/overview/namespace/default/workloads/pods
    if (this.hasNamespace(url)) {
      const ret = _.slice(parts, 0, 7).join('/')
      return ret
    }
    // e.g. /content/overview/workloads/pods
    return _.slice(parts, 0, 5).join('/')
  }

  refreshEventStream = () => {
    const { location } = this.props
    const { namespaceOption } = this.state
    this.setEventSourceStream(location.pathname, namespaceOption.value)
  }

  onTagsChange = (newFilterTags) => {
    this.setState({ resourceFilters: newFilterTags }, this.refreshEventStream)
  }

  onLabelClick = (key: string, value: string) => {
    const tag = `${key}:${value}`
    const { resourceFilters } = this.state
    this.setState({ resourceFilters: [...resourceFilters, tag] }, this.refreshEventStream)
  }

  setError = (hasError: boolean, errorMessage?: string): void => {
    errorMessage = errorMessage || 'Oops, something is not right, try again.'
    this.setState({ hasError, errorMessage })
  }

  render() {
    const {
      isLoading,
      hasError,
      errorMessage,
      navigation,
      currentNavLinkPath,
      namespaceOptions,
      namespaceOption,
      title,
      resourceFilters,
    } = this.state

    let currentNamespace = null
    if (namespaceOption) {
      currentNamespace = namespaceOption.value
    }

    let navSections = null
    let rootNavigationPath = `/content/overview/namespace/${currentNamespace}/`
    if (navigation && navigation.sections) {
      navSections = navigation.sections
      rootNavigationPath = navigation.sections[0].path
    }

    return (
      <div className='app'>
        <Header
          namespaceOptions={namespaceOptions}
          namespace={currentNamespace}
          namespaceValue={namespaceOption}
          onNamespaceChange={this.onNamespaceChange}
          resourceFilters={resourceFilters}
          onResourceFiltersChange={this.onTagsChange}
        />
        <ResourceFiltersContext.Provider value={{ onLabelClick: this.onLabelClick }}>
          <div className='app-page'>
            <div className='app-nav'>
              <Navigation
                navSections={navSections}
                currentNavLinkPath={currentNavLinkPath}
                onNavChange={(linkPath) => this.setState({ currentNavLinkPath: linkPath })}
              />
            </div>
            <div className='app-main'>
              <Switch>
                <Redirect exact from='/' to={rootNavigationPath} />
                <Route
                  render={(props) => (
                    <Overview
                      {...props}
                      title={title}
                      isLoading={isLoading}
                      hasError={hasError}
                      errorMessage={errorMessage}
                      setError={this.setError}
                      data={this.state.contentResponse}
                    />
                  )}
                />
              </Switch>
            </div>
            <ReactTooltip html />
          </div>
        </ResourceFiltersContext.Provider>
      </div>
    )
  }
}

export default withRouter(App)
