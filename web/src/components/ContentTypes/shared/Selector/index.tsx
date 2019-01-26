import './styles.scss'

import React from 'react'

export interface LabelSelector {
  metadata: {
    type: 'labelSelector';
  };
  config: {
    key: string;
    value: string;
  };
}

export interface ExpressionSelector {
  metadata: {
    type: 'expressionSelector';
  };
  config: {
    key: string;
    operator: string;
    values: string[];
  };
}

interface Props {
  config: {
    selectors: any[];
  };
}

export default function Selectors({ config }: Props) {
  if (!config.selectors) {
    return <div className='selectors' />
  }

  const selectors = config.selectors.map((selector, index) => {
    switch (selector.metadata.type) {
      case 'labelSelector':
        const labelSelector = selector as LabelSelector
        return (
          <div key={index} className='selectors--label'>
            {labelSelector.config.key}:{labelSelector.config.value}
          </div>
        )
      case 'expressionSelector':
        const expressionSelector = selector as ExpressionSelector
        return (
          <div key={index} className='selectors--expression'>
            {`${expressionSelector.config.key} ${expressionSelector.config.operator} `}
            <a data-tip={expressionSelector.config.values.join(',')}>[]</a>
          </div>
        )
      default:
        throw new Error(
          `unknown label selector ${JSON.stringify(selector.metadata)}`,
        )
    }
  })

  return <div className='selectors'>{selectors}</div>
}
