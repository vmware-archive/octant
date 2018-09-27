import React from 'react'
import { Switch, Route } from 'react-router-dom'
import Home from 'pages/Home'
import Header from '../Header'
import Navigation from '../Navigation'
import './styles.scss'

function App () {
  return (
    <div className='app'>
      <Header />
      <div className='app-page'>
        <div className='app-nav'>
          <Navigation />
        </div>
        <div className='app-main'>
          <Switch>
            <Route exact path='/' component={Home} />
          </Switch>
        </div>
      </div>
    </div>
  )
}

export default App
