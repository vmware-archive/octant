import React, { Component } from 'react'
import {
  Switch, Route, withRouter, Redirect
} from 'react-router-dom'
import _ from 'lodash'
import { setNamespace } from 'api'
import Overview from 'pages/Overview'
import Header from '../Header'
import Navigation from '../Navigation'
import fetchContents from './state/fetchContents'
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
      contents: [],
      title: '',
      namespaceOption: { label: 'default', value: 'default' }
    }
  }

  async componentDidMount () {
    const { location } = this.props
    const initialState = await getInitialState(location.pathname)
    this.setState(initialState)
  }

  async componentDidUpdate ({ location: { pathname: lastPath } }) {
    const {
      location: { pathname: thisPath }
    } = this.props

    if (thisPath && lastPath !== thisPath) {
      await this.setContents()
    }
  }

  // Note(marlon): this is an overview concept, not a dev dash concept.
  // This logic should move to the overview component child.
  setContents = async (namespace) => {
    this.setState({
      contents: [],
      title: '',
      loading: true,
      error: false
    })
    if (!namespace) {
      const { namespaceOption } = this.state
      namespace = namespaceOption.value
    }
    const { location } = this.props
    const state = await fetchContents(location.pathname, namespace)
    this.setState(state)
  }

  onNamespaceChange = async (namespaceOption) => {
    this.setState({
      namespaceOption,
      loading: true,
      contents: [],
      error: false
    })
    const { value } = namespaceOption
    try {
      await setNamespace(value)
      // Note(marlon): this is needed because user might switch namespaces
      // before the previous namespace request and we want to make sure
      // we render the correct contents
      const {
        namespaceOption: _namespaceOption,
        currentNavLinkPath
      } = this.state
      if (value === _namespaceOption.value) {
        const currentLink = _.last(currentNavLinkPath)
        this.props.history.push(currentLink.path)
        await this.setContents(value)
      }
    } catch (e) {
      this.setState({ loading: false, error: true })
    }
  }

  render () {
    const {
      loading,
      contents,
      navigation,
      currentNavLinkPath,
      namespaceOptions,
      namespaceOption,
      title,
      error
    } = this.state
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
                path='/content/overview'
                render={props => (
                  <Overview
                    {...props}
                    contents={contents}
                    loading={loading}
                    title={title}
                    error={error}
                  />
                )}
              />
              <Redirect exact from='/' to='/content/overview' />
            </Switch>
          </div>
        </div>
      </div>
    )
  }
}

export default withRouter(App)
