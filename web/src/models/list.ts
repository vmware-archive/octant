import { ListModel, TitleView, toTitle, View, viewFromContentType } from 'models'

export class JSONList implements ListModel {
  readonly type = 'list'
  readonly title: TitleView
  readonly items: View[] = []

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

    if (ct.config.items) {
      this.items = ct.config.items.map((item) => viewFromContentType(item))
    }
  }
}
