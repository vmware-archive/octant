import { AnnotationsModel, TitleView, toTitle } from 'models'

export class JSONAnnotations implements AnnotationsModel {
  readonly title: TitleView
  readonly annotations: { [key: string]: string }
  readonly type = 'annotations'

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

    this.annotations = ct.config.annotations
  }
}
