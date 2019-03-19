import { TimestampModel, TitleView, toTitle } from 'models'

export class JSONTimestamp implements TimestampModel {
  readonly isComparable = true

  readonly timestamp: number
  readonly type = 'timestamp'
  readonly title: TitleView

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

    this.timestamp = ct.config.timestamp
  }
}

export function compareTimestampModel(a: TimestampModel, b: TimestampModel): number {
  if (a.timestamp < b.timestamp) {
    return -1
  } else if (a.timestamp > b.timestamp) {
    return 1
  } else {
    return 0
  }
}
