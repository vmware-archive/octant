import { buildRequest, BuildRequestParams, getAPIBase } from 'api'
import { pointRadial } from 'd3'
import { PortModel, PortsModel } from 'models'
import React from 'react'

import './styles.scss'

interface Props {
  view: PortsModel
}

export default class Ports extends React.Component<Props> {
  async add(port: PortModel) {
    const data = {
      apiVersion: port.apiVersion,
      kind: port.kind,
      name: port.name,
      namespace: port.namespace,
      port: port.port,
    }

    const params: BuildRequestParams = {
      endpoint: 'api/v1/content/overview/port-forwards',
      method: 'POST',
      data,
    }

    await buildRequest(params)

  }

  async remove(port: PortModel) {
    const params: BuildRequestParams = {
      endpoint:  `api/v1/content/overview/port-forwards/${port.state.id}`,
      method: 'DELETE',
    }

    await buildRequest(params)
  }

portLink(port: PortModel) {
    if (port.state.isForwarded) {
      const link = `localhost:${port.state.port}`
      return <a onClick={() => window.open(`http://${link}`, '_blank')}>
      {port.port}/{port.protocol} forwarded to {link}</a>
    }

    return `${port.port}/${port.protocol}`
  }

render() {
    const view = this.props.view

    const ports = view.ports.map((port: PortModel, index: number) => {
      const state = port.state.isForwardable ? (
        port.state.isForwarded ? (
          <a onClick={() => this.remove(port)}>remove</a>
        ) : (
          <a onClick={() => this.add(port)}>port forward</a>
        )
      ) : (
        <></>
      )

      return (
        <div key={index} className='port'>
          {this.portLink(port)} {state}
        </div>
      )
    })

    return <div className='ports'>{ports}</div>
  }
}
