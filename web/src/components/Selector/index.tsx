import './styles.scss'

import { ExpressionSelector, LabelSelector, SelectorsModel } from 'models/View'
import React from 'react'

interface Props {
  view: SelectorsModel
}

export default function Selectors({ view }: Props) {
  if (!view.selectors) {
    return <div className='selectors' />
  }

  const selectors = view.selectors.map((selector, index) => {
    switch (selector.type) {
      case 'labelSelector':
        const labelSelector = selector as LabelSelector
        return (
          <div key={index} className='selectors--label'>
            {labelSelector.key}:{labelSelector.value}
          </div>
        )
      case 'expressionSelector':
        const expressionSelector = selector as ExpressionSelector
        return (
          <div key={index} className='selectors--expression'>
            {`${expressionSelector.key} ${expressionSelector.operator} `}
            <a data-tip={expressionSelector.values.join(',')}>[]</a>
          </div>
        )
      default:
        throw new Error(
          `unknown label selector ${JSON.stringify(selector)}`,
        )
    }
  })

  return <div className='selectors'>{selectors}</div>
}
