import { LabelsModel, TitleView, toTitle } from 'models'

export class JSONLabels implements LabelsModel {
  readonly title: TitleView
  readonly labels: { [key: string]: string }
  readonly type = 'labels'

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

    this.labels = ct.config.labels
  }
}
