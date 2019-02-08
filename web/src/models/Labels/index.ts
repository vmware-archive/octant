import { View } from 'models/View'

export interface LabelsModel extends View {
  labels: { [key: string]: string }
}

export class JSONLabels implements LabelsModel {
  readonly labels: { [key: string]: string }
  readonly title: string
  readonly type = 'labels'

  constructor(private readonly ct: ContentType) {
    this.title = ct.metadata.title
    this.labels = ct.config.labels
  }
}
