import { View } from 'models/View'

export interface QuadrantModel extends View {
  nw: QuadrantSector;
  ne: QuadrantSector;
  sw: QuadrantSector;
  se: QuadrantSector;
}

export interface QuadrantSector {
  value: string;
  label: string;
}

export class JSONQuadrant implements QuadrantModel {
  readonly type = 'quadrant'
  readonly title: string
  readonly nw: QuadrantSector
  readonly ne: QuadrantSector
  readonly sw: QuadrantSector
  readonly se: QuadrantSector

  constructor(ct: ContentType) {
    this.nw = ct.config.nw
    this.ne = ct.config.ne
    this.sw = ct.config.sw
    this.se = ct.config.se
    this.title = ct.metadata.title
  }
}
