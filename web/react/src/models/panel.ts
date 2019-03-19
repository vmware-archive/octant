import { PanelModel, PanelPosition, TitleView, toTitle, View, viewFromContentType } from 'models'

export class JSONPanel implements PanelModel {
  readonly type = 'panel'
  readonly position: PanelPosition
  readonly content: View
  readonly title: TitleView

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

    this.position = ct.config.position
    this.content = viewFromContentType(ct.config.content)
  }
}
