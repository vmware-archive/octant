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
      loading: false,
      error: false,
      navigation: [],
      currentNavLinkPath: [],
      namespaceOptions: [],
      title: '',
      namespaceOption: { label: 'default', value: 'default' }
    }
  }

  async componentDidMount () {
    const { location } = this.props
    const initialState = await getInitialState(location.pathname)
    this.setState(initialState)
  }

  onNamespaceChange = async (namespaceOption) => {
    this.setState({
      loading: true,
      error: false
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
        this.setState({ namespaceOption, loading: false, error: false })
      }
    } catch (e) {
      this.setState({ namespaceOption, loading: false, error: true })
    }
  }

  render () {
    const {
      loading,
      navigation,
      currentNavLinkPath,
      namespaceOptions,
      namespaceOption,
      title,
      error
    } = this.state
    const { location } = this.props

    const currentPath = location.pathname
    const currentNamespace = namespaceOption.value

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
                    path={currentPath}
                    namespace={currentNamespace}
                    loading={loading}
                    title={title}
                    error={error}
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
