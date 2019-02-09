import { ContentsUrlParams, getAPIBase, getContentsUrl, POLL_WAIT, setNamespace } from 'api'
import Header from 'components/Header'
import ResourceFiltersContext from 'contexts/resource-filters'
import _ from 'lodash'
import JSONContentResponse, { Parse } from 'models/ContentResponse'
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
      namespaceOption: null,
      namespaceOptions: [],
      contentResponse: null,
      resourceFilters: [],
    }
  }

  async componentDidMount() {
    let namespace = 'default'

    const { location: initialLocation } = this.props
    const initialState = await getInitialState(initialLocation.pathname)
    this.setState(initialState as AppState)

    if (this.state.namespaceOption) {
      namespace = this.state.namespaceOption.value
    }

    this.setEventSourceStream(this.props.location.pathname, namespace)
  }

  componentDidUpdate({ location: previousLocation }, { namespaceOption: previousNamespace }) {
    const { location } = this.props
    const { namespaceOption } = this.state

    const namespace = namespaceOption ? namespaceOption.value : 'default'
    const prevNamespace = previousNamespace ? previousNamespace.value : ''

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
      namespace,
      poll: POLL_WAIT,
    }

    const { resourceFilters } = this.state
    if (resourceFilters && resourceFilters.length) params.filter = resourceFilters

    const url = getContentsUrl(path, params)

    this.source = new window.EventSource(`${getAPIBase()}/${url}`)

    this.source.addEventListener('message', (e) => {
      const cr2 = Parse(e.data)

      this.setState({
        contentResponse: cr2,
        isLoading: false,
      })
    })

    this.source.addEventListener('navigation', (e) => {
      const data = JSON.parse(e.data)
      this.setState({ navigation: data })
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
      this.lastFetchedNamespace = value

      await setNamespace(value)

      if (this.lastFetchedNamespace === value) {
        const { currentNavLinkPath } = this.state
        const { path } = _.last(currentNavLinkPath)
        history.push(path)
        this.setState({ namespaceOption, isLoading: false, hasError: false })
      }
    } catch (e) {
      this.setState({ namespaceOption, isLoading: false, hasError: true })
    }
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
    let rootNavigationPath = '/content/overview/'
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
