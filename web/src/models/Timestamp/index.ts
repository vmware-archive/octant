import { View } from 'models/View'

export interface TimestampModel extends View {
  timestamp: number;
}

export class JSONTimestamp implements TimestampModel {
  readonly timestamp: number
  readonly title: string
  readonly type = 'timestamp'

  constructor(ct: ContentType) {
    this.title = ct.metadata.title
    this.timestamp = ct.config.timestamp
  }
}
