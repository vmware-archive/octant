import { FlexLayoutItem, FlexLayoutModel, FlexLayoutSection, TitleView, toTitle, viewFromContentType } from 'models'

export class JSONFlexLayout implements FlexLayoutModel {
  readonly title: TitleView
  readonly type = 'flexlayout'
  readonly sections: FlexLayoutSection[]

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

    this.sections = ct.config.sections.map((section) => {
      return section.map((item) => {
        return {
          width: item.width,
          view: viewFromContentType(item.view),
        } as FlexLayoutItem
      }) as FlexLayoutSection
    })
  }
}
