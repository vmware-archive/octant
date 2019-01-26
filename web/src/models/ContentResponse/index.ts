import { View, viewFromContentType } from 'models/View'

interface ContentResponse {
  content: {
    title: string;
    viewComponents: ContentType[];
  };
}

export function Parse(data: any): JSONContentResponse {
  const cr: ContentResponse = JSON.parse(data)

  return new JSONContentResponse(cr)
}

export default class JSONContentResponse {
  readonly title: string
  readonly views: View[]

  constructor(private readonly cr: ContentResponse) {
    this.title = cr.content.title

    this.views = cr.content.viewComponents.map((view) => {
      return viewFromContentType(view)
    })
  }
}
