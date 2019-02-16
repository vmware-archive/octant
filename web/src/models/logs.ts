import { LogsModel, TitleView, toTitle } from 'models'

export class JSONLogs implements LogsModel {
  readonly type = 'logs'
  readonly title: TitleView
  readonly namespace: string
  readonly name: string
  readonly containers: string[]
  readonly accessor: string

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

    this.accessor = ct.metadata.accessor
    this.namespace = ct.config.namespace
    this.name = ct.config.name
    this.containers = ct.config.containers
  }
}
