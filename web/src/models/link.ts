import { LinkModel, TitleView, toTitle } from 'models'

export class JSONLink implements LinkModel {
  readonly ref: string
  readonly value: string
  readonly type = 'link'
  readonly title: TitleView

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

    this.ref = ct.config.ref
    this.value = ct.config.value
  }
}
