import { TitleView, toTitle, YAMLViewerModel } from 'models'

export class JSONYAMLViewer implements YAMLViewerModel {
  readonly data: string
  readonly type = 'yaml'
  readonly title: TitleView
  readonly accessor: string

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

    this.accessor = ct.metadata.accessor
    this.data = ct.config.data
  }
}
