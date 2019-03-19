import { TableModel, TableRow, TitleView, toTitle, View, viewFromContentType } from 'models'

export class JSONTable implements TableModel {
  readonly type = 'table'
  readonly title: TitleView
  readonly rows: TableRow[]
  readonly columns: Array<{ name: string; accessor: string }>
  readonly emptyContent: string

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

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
