import React from 'react'
import ReactDOM from 'react-dom'
import { Router } from 'react-router-dom'
import createHistory from 'history/createBrowserHistory'
import './css/styles/styles.scss'
import App from './components/App'

const history = createHistory()

ReactDOM.render(
  <Router history={history}>
    <App />
  </Router>,
  document.getElementById('root')
)
