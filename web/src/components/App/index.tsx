import './styles.scss'

import { setNamespace } from 'api'
import Header from 'components/Header'
import _ from 'lodash'
import Overview from 'pages/Overview'
import React, { Component } from 'react'
import { Redirect, Route, RouteComponentProps, Switch, withRouter } from 'react-router-dom'

import Navigation from '../Navigation'
import getInitialState from './state/getInitialState'

interface AppState {
  isLoading: boolean;
  hasError: boolean;
  errorMessage: string;
  navigation: { sections: NavigationSectionType[] };
  currentNavLinkPath: NavigationSectionType[];
  namespaceOption: NamespaceOption;
  namespaceOptions: NamespaceOption[];
  title: string;
}

class App extends Component<RouteComponentProps, AppState> {
  private lastFetchedNamespace: string

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
    }
  }

  async componentDidMount() {
    const { location } = this.props
    const initialState = await getInitialState(location.pathname)
    this.setState(initialState as AppState)
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
    const { location } = this.props

    const currentPath = location.pathname

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
                path={rootNavigationPath}
                render={(props) => (
                  <Overview
                    {...props}
                    title={title}
                    path={currentPath}
                    namespace={currentNamespace}
                    isLoading={isLoading}
                    hasError={hasError}
                    errorMessage={errorMessage}
                    setIsLoading={(loading) =>
                      this.setState({ isLoading: loading })
                    }
                    setError={this.setError}
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
