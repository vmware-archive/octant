import './styles.scss'

import { getAPIBase, getContentsUrl, POLL_WAIT, setNamespace } from 'api'
import Header from 'components/Header'
import _ from 'lodash'
import Overview from 'pages/Overview'
import React, { Component } from 'react'
import { Redirect, Route, RouteComponentProps, Switch, withRouter } from 'react-router-dom'

import Navigation from '../Navigation'
import getInitialState from './state/getInitialState'

interface ContentResponse {
  content: {
    title: string;
    viewComponents: ContentType[];
  };
}

interface AppState {
  isLoading: boolean;
  hasError: boolean;
  errorMessage: string;
  navigation: { sections: NavigationSectionType[] };
  currentNavLinkPath: NavigationSectionType[];
  namespaceOption: NamespaceOption;
  namespaceOptions: NamespaceOption[];
  title: string;

  overviewData: ContentResponse;
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
      overviewData: null,
    }
  }

  async componentDidMount() {
    let namespace = 'default'

    const { location } = this.props
    const initialState = await getInitialState(location.pathname)
    this.setState(initialState as AppState)

    if (this.state.namespaceOption) {
      namespace = this.state.namespaceOption.value
    }

    this.setEventSourceStream(location.pathname, namespace)
  }

  componentDidUpdate(
    { location: previousLocation },
    { namespaceOption: previousNamespace },
  ) {
    const { location } = this.props
    const { namespaceOption } = this.state

    const namespace = namespaceOption ? namespaceOption.value : 'default'
    const prevNamespace = previousNamespace ? previousNamespace.value : ''

    if (
      location.pathname !== previousLocation.pathname ||
      namespace !== prevNamespace
    ) {
      this.setEventSourceStream(location.pathname, namespace)
    }
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

    // this.setState({ isLoading: true, overviewData: null })

    const url = getContentsUrl(path, namespace, POLL_WAIT)

    this.source = new window.EventSource(`${getAPIBase()}/${url}`)

    this.source.addEventListener('message', (e) => {
      const data: ContentResponse = JSON.parse(e.data)
      this.setState({ overviewData: data, isLoading: false })
    })

    this.source.addEventListener('navigation', (e) => {
      const data = JSON.parse(e.data)
      this.setState({ navigation: data })
    })

    this.source.addEventListener('error', () => {
      this.setState({ isLoading: false })
      this.setError(
        true,
        'Looks like the back end source has gone away. Retrying...',
      )
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
        />
        <div className='app-page'>
          <div className='app-nav'>
            <Navigation
              navSections={navSections}
              currentNavLinkPath={currentNavLinkPath}
              onNavChange={(linkPath) =>
                this.setState({ currentNavLinkPath: linkPath })
              }
              namespaceOptions={namespaceOptions}
              namespaceValue={namespaceOption}
              onNamespaceChange={this.onNamespaceChange}
            />
          </div>
          <div className='app-main'>
            <Switch>
              <Route
                render={(props) => (
                  <Overview
                    {...props}
                    title={title}
                    isLoading={isLoading}
                    hasError={hasError}
                    errorMessage={errorMessage}
                    setError={this.setError}
                    data={this.state.overviewData}
                  />
                )}
              />
              <Redirect exact from='/' to={rootNavigationPath} />
            </Switch>
          </div>
        </div>
      </div>
    )
  }
}

export default withRouter(App)
