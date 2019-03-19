import { GridModel, PanelModel, TitleView, toTitle, viewFromContentType } from 'models'

export class JSONGrid implements GridModel {
  readonly title: TitleView
  readonly type = 'grid'
  readonly panels: PanelModel[]

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

    this.panels = ct.config.panels.map((panelContentType) => viewFromContentType(panelContentType))
  }
}
