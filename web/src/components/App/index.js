import React, { Component } from 'react'
import { Switch, Route } from 'react-router-dom'
import Header from '../Header'
import Navigation from '../Navigation'
import Home from '../../pages/Home'
import { getNavigation, getSummary, getTable } from '../../api'

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
          <Switch>
            <Route exact path='/' render={() => <Home summary={summary} table={table} />} />
          </Switch>
          <Navigation navigation={navigation} />
        </div>
      </div>
    )
  }
}

export default App
