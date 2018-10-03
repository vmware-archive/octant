import React, { Component } from 'react'
import { Switch, Route } from 'react-router-dom'
import { getNavigation } from 'api'
import Home from 'pages/Home'
import Header from '../Header'
import Navigation from '../Navigation'

import './styles.scss'

class App extends Component {
  constructor (props) {
    super(props)
    this.state = {
      navigation: []
    }
  }

  async componentDidMount () {
    const navigation = await getNavigation()
    this.setState({ navigation })
  }

  render () {
    const { navigation } = this.state
    return (
      <div className='app'>
        <Header />
        <div className='app-page'>
          <div className='app-nav'>
            <Navigation navigation={navigation} />
          </div>
          <div className='app-main'>
            <Switch>
              <Route path='/' component={Home} />
            </Switch>
          </div>
        </div>
      </div>
    )
  }
}

export default App
