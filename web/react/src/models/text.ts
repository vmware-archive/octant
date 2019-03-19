import { TextModel, TitleView, toTitle } from 'models'

export class JSONText implements TextModel {
  readonly isComparable = true

  readonly value: string
  readonly type = 'text'
  readonly title: TitleView

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

    this.value = ct.config.value
  }
}

export function compareTextModel(a: TextModel, b: TextModel): number {
  return a.value.localeCompare(b.value)
}
