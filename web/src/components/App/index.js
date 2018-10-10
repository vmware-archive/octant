import React, { Component } from 'react'
import { Switch, Route, withRouter } from 'react-router-dom'
import _ from 'lodash'
import Promise from 'promise'
import queryString from 'query-string'
import { getNavigation, getNamespaces } from 'api'
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
    const [navigation, namespacesPayload] = await Promise.all([
      getNavigation(),
      getNamespaces()
    ])

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

    this.setState({ navigation, namespaceOptions })
  }

  onNamespaceChange = (namespaceOption) => {
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
