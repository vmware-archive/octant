import { View } from 'models/View'

export interface TextModel extends View {
  value: string;
}

export class JSONText implements TextModel {
  readonly value: string
  readonly title: string
  readonly type = 'text'

  constructor(ct: ContentType) {
    this.title = ct.metadata.title
    this.value = ct.config.value
  }
}
