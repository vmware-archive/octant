import React, { Component } from 'react'
import {
  Switch, Route, withRouter, Redirect
} from 'react-router-dom'
import _ from 'lodash'
import { setNamespace } from 'api'
import Overview from 'pages/Overview'
import Header from '../Header'
import Navigation from '../Navigation'
import getInitialState from './state/getInitialState'
import './styles.scss'

class App extends Component {
  constructor (props) {
    super(props)
    this.state = {
      isLoading: true, // to do the initial data fetch
      hasError: false,
      navigation: [],
      currentNavLinkPath: [],
      namespaceOptions: [],
      title: '',
      namespaceOption: null
    }
  }

  async componentDidMount () {
    const { location } = this.props
    const initialState = await getInitialState(location.pathname)
    this.setState(initialState)
  }

  onNamespaceChange = async (namespaceOption) => {
    this.setState({
      isLoading: true,
      hasError: false
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

  render () {
    const {
      isLoading,
      hasError,
      navigation,
      currentNavLinkPath,
      namespaceOptions,
      namespaceOption,
      title
    } = this.state
    const { location } = this.props

    const currentPath = location.pathname

    let currentNamespace = null
    if (namespaceOption) {
      currentNamespace = namespaceOption.value
    }

    let rootNavigationPath = '/content/overview/'
    if (navigation && navigation.sections) {
      rootNavigationPath = navigation.sections[0].path
    }

    return (
      <div className='app'>
        <Header />
        <div className='app-page'>
          <div className='app-nav'>
            <Navigation
              navSections={navigation.sections}
              currentNavLinkPath={currentNavLinkPath}
              onNavChange={linkPath => this.setState({ currentNavLinkPath: linkPath })
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
                render={props => (
                  <Overview
                    {...props}
                    title={title}
                    path={currentPath}
                    namespace={currentNamespace}
                    isLoading={isLoading}
                    hasError={hasError}
                    setIsLoading={loading => this.setState({ isLoading: loading })
                    }
                    setHasError={error => this.setState({ hasError: error })}
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
