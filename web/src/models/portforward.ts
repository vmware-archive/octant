import { PortForwardModel, PortForwardSpec, TitleView, toTitle } from 'models'

export class JSONPortForward implements PortForwardModel {
  readonly type = 'portforward'
  readonly title: TitleView
  readonly text: string
  readonly id: string
  readonly action: string
  readonly status: string
  readonly ports: PortForwardSpec[]

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

    this.text = ct.config.text
    this.id = ct.config.id
    this.action = ct.config.action
    this.status = ct.config.status
    this.ports = ct.config.ports
  }
}
