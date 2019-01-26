import { PanelModel, View, viewFromContentType } from 'models/View'

export interface GridModel extends View {
  panels: PanelModel[];
}

export class JSONGrid implements GridModel {
  readonly type = 'grid'
  readonly title: string
  readonly panels: PanelModel[]

  constructor(ct: ContentType) {
    this.title = ct.metadata.title
    this.panels = ct.config.panels.map((panelContentType) =>
      viewFromContentType(panelContentType),
    )
  }
}
