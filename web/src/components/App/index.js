import React, { Component } from 'react'
import { Switch, Route, withRouter } from 'react-router-dom'
import Promise from 'promise'
import _ from 'lodash'
import {
  getNavigation, getNamespaces, getNamespace, setNamespace
} from 'api'
import Home from 'pages/Home'
import Header from '../Header'
import Navigation from '../Navigation'
import './styles.scss'

class App extends Component {
  constructor (props) {
    super(props)
    this.state = {
      navigation: [],
      namespaceOptions: [],
      namespaceOption: { label: 'default', value: 'default' }
    }
  }

  async componentDidMount () {
    // Note(marlon): this logic for this should not live in <App />. it
    // might be better handled in a <Namespace /> container component or
    // in an HOC
    const [navigation, namespacesPayload, namespacePayload] = await Promise.all(
      [getNavigation(), getNamespaces(), getNamespace()]
    )

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
      if (option) namespaceOption = option
    }

    this.setState({ navigation, namespaceOption, namespaceOptions })
  }

  onNamespaceChange = async (namespaceOption) => {
    const { value } = namespaceOption
    await setNamespace(value)
    this.setState({ namespaceOption })
  }

  render () {
    const { navigation, namespaceOptions, namespaceOption } = this.state
    return (
      <div className='app'>
        <Header />
        <div className='app-page'>
          <div className='app-nav'>
            <Navigation
              navigation={navigation}
              namespaceOptions={namespaceOptions}
              namespaceValue={namespaceOption}
              onNamespaceChange={this.onNamespaceChange}
            />
          </div>
          <div className='app-main'>
            <Switch>
              <Route
                path='/'
                component={props => (
                  <Home {...props} namespace={namespaceOption.value} />
                )}
              />
            </Switch>
          </div>
        </div>
      </div>
    )
  }
}

export default withRouter(App)
