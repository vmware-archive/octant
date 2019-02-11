import { QuadrantModel, QuadrantSector, TitleView, toTitle } from 'models'

export class JSONQuadrant implements QuadrantModel {
  readonly type = 'quadrant'
  readonly title: TitleView
  readonly nw: QuadrantSector
  readonly ne: QuadrantSector
  readonly sw: QuadrantSector
  readonly se: QuadrantSector

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

    this.nw = ct.config.nw
    this.ne = ct.config.ne
    this.sw = ct.config.sw
    this.se = ct.config.se
  }
}
