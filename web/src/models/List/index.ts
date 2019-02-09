import { View, viewFromContentType } from 'models/View'

export interface ListModel extends View {
  items(): View[]
}

export class JSONList implements ListModel {
  readonly type = 'list'
  readonly title: string

  constructor(private readonly ct: ContentType) {
    this.title = ct.metadata.title
  }

  items(): View[] {
    return this.ct.config.items.map((ct) => viewFromContentType(ct))
  }
}
