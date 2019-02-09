import { View } from 'models/View'

export interface LinkModel extends View {
  ref: string
  value: string
}

export class JSONLink implements LinkModel {
  readonly ref: string
  readonly value: string
  readonly title: string
  readonly type = 'link'

  constructor(private readonly ct: ContentType) {
    this.title = ct.metadata.title
    this.ref = ct.config.ref
    this.value = ct.config.value
  }
}
