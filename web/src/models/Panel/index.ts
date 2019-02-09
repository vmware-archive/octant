import { View, viewFromContentType } from 'models/View'

export interface PanelPosition {
  x: number
  y: number
  w: number
  h: number
}

export interface PanelModel extends View {
  readonly position: PanelPosition
  readonly content: View
}

export class JSONPanel implements PanelModel {
  readonly type = 'panel'
  readonly title: string
  readonly position: PanelPosition
  readonly content: View

  constructor(ct: ContentType) {
    this.title = ct.metadata.title
    this.position = ct.config.position
    this.content = viewFromContentType(ct.config.content)
  }
}
