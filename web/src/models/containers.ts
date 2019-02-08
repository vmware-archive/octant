import { ContainerDef, ContainersModel, TitleView, toTitle } from 'models'

export class JSONContainers implements ContainersModel {
  readonly title: TitleView
  readonly containerDefs: ContainerDef[]
  readonly type = 'containers'

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

    this.containerDefs = ct.config.containers
  }
}
