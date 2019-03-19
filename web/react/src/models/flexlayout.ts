import { FlexLayoutItem, FlexLayoutModel, FlexLayoutSection, TitleView, toTitle, viewFromContentType } from 'models'

interface ContentSection {
  map: (arg0: (item: SectionItem) => FlexLayoutItem) => FlexLayoutItem[]
}
interface SectionItem {
  width: number
  view: ContentType
}

export class JSONFlexLayout implements FlexLayoutModel {
  readonly title: TitleView
  readonly type = 'flexlayout'
  readonly sections: FlexLayoutSection[]
  readonly accessor: string

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

    this.accessor = ct.metadata.accessor

    this.sections = ct.config.sections.map(
      (section: ContentSection) =>
        section.map((item: SectionItem) => {
          return {
            width: item.width,
            view: viewFromContentType(item.view),
          }
        }) as FlexLayoutSection
    )
  }
}
