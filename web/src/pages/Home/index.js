import React, { Component } from 'react'
import { getContents } from 'api'
import ContentSwitcher from './components/ContentSwitcher'
import './styles.scss'

class Home extends Component {
  state = {
    contents: []
  }

  async componentDidMount () {
    await this.fetchContents()
  }

  async componentDidUpdate ({ location: { pathname: lastPath } }) {
    const {
      location: { pathname: thisPath }
    } = this.props

    if (thisPath && lastPath !== thisPath) {
      await this.fetchContents()
    }
  }

  static getDerivedStateFromProps ({ location: { pathname } }) {
    if (!pathname || pathname === '/') {
      return { contents: [] }
    }
    return null
  }

  fetchContents = async () => {
    const {
      location: { pathname }
    } = this.props
    const payload = await getContents(pathname)
    if (payload) {
      this.setState({ contents: payload.contents || [] })
    }
  }

  render () {
    const { contents = [] } = this.state
    return (
      <div className='home'>
        <div className='main'>
          {contents.map((content, i) => (
            <div key={i} className='component--primary'>
              <ContentSwitcher content={content} />
            </div>
          ))}
        </div>
      </div>
    )
  }
}
export default Home
