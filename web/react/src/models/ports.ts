import { PortModel, PortsModel, TitleView, toTitle, viewFromContentType } from 'models'

export class JSONPorts implements PortsModel {
  readonly type = 'ports'
  readonly title: TitleView
  readonly ports: PortModel[]

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

    this.ports = ct.config.ports.map((portContentType: ContentType) => {
      return viewFromContentType(portContentType)
    })
  }
}
