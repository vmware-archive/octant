import React from 'react'
import ReactDOM from 'react-dom'
import { HashRouter } from 'react-router-dom'
import createHistory from 'history/createBrowserHistory'
import './css/styles/styles.scss'
import App from './components/App'

const history = createHistory()

ReactDOM.render(
  <HashRouter history={history}>
    <App />
  </HashRouter>,
  document.getElementById('root')
)
