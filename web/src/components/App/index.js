import React, { Component } from 'react'
import { Switch, Route } from 'react-router-dom'
import { getNavigation, getSummary, getTable } from 'api'
import Home from 'pages/Home'
import Header from '../Header'
import Navigation from '../Navigation'

import './styles.scss'

class App extends Component {
  constructor (props) {
    super(props)
    this.state = {
      navigation: [],
      summary: [],
      table: []
    }
  }

  async componentDidMount () {
    const navigation = getNavigation()
    const summary = getSummary()
    const table = getTable()

    this.setState({ navigation, summary, table })
  }

  render () {
    const { navigation, summary, table } = this.state
    return (
      <div className='app'>
        <Header />
        <div className='app-page'>
          <div className='app-nav'>
            <Navigation navigation={navigation} />
          </div>
          <div className='app-main'>
            <Switch>
              <Route exact path='/' render={() => <Home summary={summary} table={table} />} />
            </Switch>
          </div>
        </div>
      </div>
    )
  }
}

export default App
