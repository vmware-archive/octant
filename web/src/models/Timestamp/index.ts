import { View } from 'models/View'

export interface TimestampModel extends View {
  timestamp: number
}

export class JSONTimestamp implements TimestampModel {
  readonly isComparable = true

  readonly timestamp: number
  readonly title: string
  readonly type = 'timestamp'

  constructor(ct: ContentType) {
    this.title = ct.metadata.title
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
