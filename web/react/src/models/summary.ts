import _ from 'lodash'
import { SummaryItem, SummaryModel, TitleView, toTitle, viewFromContentType } from 'models'

export class JSONSummary implements SummaryModel {
  readonly type = 'summary'
  readonly title: TitleView
  readonly items: SummaryItem[]

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

    this.items = _.map(ct.config.sections, (section) => {
      return {
        header: section.header,
        content: viewFromContentType(section.content),
      }
    })
  }
}
