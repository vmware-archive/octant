import { View } from 'models/View'

export interface LabelSelector {
  key: string
  value: string
  type: string
}

export interface ExpressionSelector {
  key: string
  operator: string
  values: string[]
  type: string
}

export interface SelectorsModel extends View {
  selectors: Array<LabelSelector | ExpressionSelector>
}

export class JSONSelectors implements SelectorsModel {
  readonly selectors: Array<LabelSelector | ExpressionSelector>
  readonly title: string
  readonly type = 'selectors'

  constructor(ct: ContentType) {
    this.title = ct.metadata.title

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
