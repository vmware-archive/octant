import { View } from 'models/View'

export interface TableModel extends View {
  readonly columns: string[];
  readonly rows: { [key: string]: View };
  readonly emptyContent: string;
}

export class JSONTable implements TableModel {
  readonly type = 'table'
  readonly title: string
  readonly rows: { [key: string]: View }
  readonly columns: string[]
  readonly emptyContent: string

  constructor(ct: ContentType) {
    this.title = ct.metadata.title
    this.rows = ct.config.rows
    this.columns = ct.config.columns
    this.emptyContent = ct.config.emptyContent
  }
}
