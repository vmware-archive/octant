import { TitleView, toTitle, View, viewFromContentType } from 'models'

interface ContentResponse {
  content: {
    title: ContentType[]
    viewComponents: ContentType[]
  }
}

export function Parse(data: any): JSONContentResponse {
  const cr: ContentResponse = JSON.parse(data)

  return new JSONContentResponse(cr)
}

export default class JSONContentResponse {
  readonly title: TitleView
  readonly views: View[]

  constructor(cr: ContentResponse) {
    if (cr.content.title) {
      this.title = toTitle(cr.content.title)
    }

    this.views = cr.content.viewComponents.map((view) => {
      return viewFromContentType(view)
    })
  }
}
