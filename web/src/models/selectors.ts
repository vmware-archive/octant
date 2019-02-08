import { ExpressionSelector, LabelSelector, SelectorsModel, TitleView, toTitle } from 'models'

export class JSONSelectors implements SelectorsModel {
  readonly selectors: Array<LabelSelector | ExpressionSelector>
  readonly type = 'selectors'
  readonly title: TitleView

  constructor(ct: ContentType) {
    if (ct.metadata.title) {
      this.title = toTitle(ct.metadata.title)
    }

    if (ct.config.selectors) {
      this.selectors = ct.config.selectors.map((selector) => {
        switch (selector.metadata.type) {
          case 'labelSelector':
            return {
              key: selector.config.key,
              value: selector.config.value,
              type: 'labelSelector',
            }
          case 'expressionSelector':
            return {
              key: selector.config.key,
              operator: selector.config.operator,
              values: selector.config.values,
              type: 'expressionSelector',
            }
          default:
            return new Error(`unknown selector ${selector.metadata.type}`)
        }
      })
    }
  }
}
