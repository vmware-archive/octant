import { PortForwardState, PortModel, TitleView, toTitle } from 'models'

export class JSONPort implements PortModel {
  readonly type = 'port'
  readonly title: TitleView
  readonly namespace: string
  readonly apiVersion: string
  readonly kind: string
  readonly name: string
  readonly port: number
  readonly protocol: string
  readonly state: PortForwardState

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

    this.namespace = ct.config.namespace
    this.apiVersion = ct.config.apiVersion
    this.kind = ct.config.kind
    this.name = ct.config.name
    this.port = ct.config.port
    this.protocol = ct.config.protocol
    this.state = ct.config.state
  }
}
