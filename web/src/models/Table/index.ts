import { View, viewFromContentType } from 'models/View'

export interface TableRow {
  [key: string]: View;
}
export interface TableModel extends View {
  readonly columns: Array<{ name: string; accessor: string }>;
  readonly rows: TableRow[];
  readonly emptyContent: string;
}

export class JSONTable implements TableModel {
  readonly type = 'table'
  readonly title: string
  readonly rows: TableRow[]
  readonly columns: Array<{ name: string; accessor: string }>
  readonly emptyContent: string

  constructor(ct: ContentType) {
    this.title = ct.metadata.title
    this.columns = ct.config.columns
    this.emptyContent = ct.config.emptyContent

    if (!ct.config.rows) {
      ct.config.rows = []
    }

    this.rows =
      ct.config.rows.map((row: { [key: string]: ContentType }) => {
        const viewRow: { [key: string]: View } = {}
        Object.entries(row).forEach(([column, vct]) => {
          viewRow[column] = viewFromContentType(vct)
        })

        return viewRow
      }) || []
  }
}
