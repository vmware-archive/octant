import { getAPIBase } from 'api'
import { LogsModel } from 'models'
import React, { Component } from 'react'

import { LogEntry, LogResponse } from './log'
import Message from './message'
import './styles.scss'

const fetchInterval = 3000

interface Props {
  view: LogsModel
}

interface State {
  logs: LogEntry[]
  selectedContainer: string
}

export default class Logs extends Component<Props, State> {
  private timer: number

  private scrollRef = React.createRef<HTMLDivElement>()

  constructor(props: Props) {
    super(props)

    this.state = {
      logs: [],
      selectedContainer: props.view.containers[0],
    }
  }

  componentWillUnmount() {
    this.stopTimer()
  }

  componentDidMount() {
    this.fetchLogs()
  }

  componentDidUpdate() {
    const { logs } = this.state
    if (logs && logs.length > 0) {
      this.scrollRef.current.scrollIntoView()
    }
  }

  handleChange(event) {
    this.setState({ selectedContainer: event.target.value })
  }

  fetchLogs() {
    const { view } = this.props

    const urlParts = [
      getAPIBase(),
      'api/v1/content/overview',
      `namespace/${view.namespace}`,
      'logs',
      `pod/${view.name}`,
      `container/${this.state.selectedContainer}`,
    ]

    const url = urlParts.join('/')

    const component = this

    fetch(url)
      .then((res) => res.json())
      .then((logResponse: LogResponse) => {
        this.setState({ logs: logResponse.entries })
        component.startTimer()
      })
  }

  tick() {
    this.fetchLogs()
  }

  startTimer() {
    clearInterval(this.timer)

    // there is a nodejs setInterval that returns a NodeJS.Timeout instead
    // of a number. Editors can get confused, so explicitly ensure this
    // is a number.
    this.timer = (setTimeout(this.tick.bind(this), fetchInterval) as unknown) as number
  }

  stopTimer() {
    clearInterval(this.timer)
  }

  render() {
    const { view } = this.props

    const options = view.containers.map((container, index) => {
      return <option key={index}>{container}</option>
    })

    let messages: JSX.Element[] = []
    if (this.state.logs) {
      messages = this.state.logs.map((log, index) => <Message key={index} log={log} />)
    }

    return (
      <div className='logs'>
        <div>
          <select defaultValue={this.state.selectedContainer} onChange={this.handleChange}>
            {options}
          </select>
        </div>
        <div className='logs--messages'>
          {messages}
          <div ref={this.scrollRef} />
        </div>
      </div>
    )
  }
}
