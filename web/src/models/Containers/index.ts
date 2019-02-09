import { View } from 'models/View'

export interface ContainerDef {
  name: string
  image: string
}

export interface ContainersModel extends View {
  containerDefs: ContainerDef[]
}

export class JSONContainers implements ContainersModel {
  readonly containerDefs: ContainerDef[]
  readonly title: string
  readonly type = 'containers'

  constructor(ct: ContentType) {
    this.title = ct.metadata.title
    this.containerDefs = ct.config.containers
  }
}
