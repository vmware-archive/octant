import React, { Component } from 'react'
import {
  Switch, Route, withRouter, Redirect
} from 'react-router-dom'
import Promise from 'promise'
import _ from 'lodash'
import {
  getNavigation,
  getNamespaces,
  getNamespace,
  setNamespace,
  getContents
} from 'api'
import Overview from 'pages/Overview'
import Header from '../Header'
import Navigation from '../Navigation'
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
    // Note(marlon): this logic for this should not live in <App />. it
    // might be better handled in a <Namespace /> container component or
    // in an HOC
    let navigation,
      namespacesPayload,
      namespacePayload
    try {
      [navigation, namespacesPayload, namespacePayload] = await Promise.all([
        getNavigation(),
        getNamespaces(),
        getNamespace()
      ])
    } catch (e) {
      this.setState({ loading: false, error: true })
    }

    let namespaceOptions = []
    if (
      namespacesPayload &&
      namespacesPayload.namespaces &&
      namespacesPayload.namespaces.length
    ) {
      namespaceOptions = namespacesPayload.namespaces.map(ns => ({
        label: ns,
        value: ns
      }))
    }

    let { namespaceOption } = this.state
    if (namespacePayload && namespaceOptions.length) {
      const option = _.find(namespaceOptions, {
        value: namespacePayload.namespace
      })
      if (option) {
        namespaceOption = option
        await this.fetchContents(namespaceOption.value)
      }
    }

    this.setState({
      navigation,
      namespaceOption,
      namespaceOptions
    })
  }

  async componentDidUpdate ({ location: { pathname: lastPath } }) {
    const {
      location: { pathname: thisPath }
    } = this.props

    if (thisPath && lastPath !== thisPath) {
      await this.fetchContents()
    }
  }

  // Note(marlon): this is an overview concept, not a dev dash concept.
  // This logic should move to the overview component child.
  fetchContents = async (namespace) => {
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
    const {
      location: { pathname }
    } = this.props
    try {
      const payload = await getContents(pathname, namespace)
      if (payload) {
        this.setState({
          contents: payload.contents,
          title: payload.title,
          loading: false,
          error: false
        })
      }
    } catch (e) {
      this.setState({ loading: false, error: true })
    }
  }

  onNamespaceChange = async (namespaceOption) => {
    this.setState({
      namespaceOption,
      loading: true,
      contents: [],
      error: false
    })
    const { value } = namespaceOption
    await setNamespace(value)
    // Note(marlon): this is needed because user might switch namespaces
    // before the previous namespace request and we want to make sure
    // we render the correct contents
    const { namespaceOption: _namespaceOption } = this.state
    if (value === _namespaceOption.value) {
      await this.fetchContents(value)
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
