import React from 'react'
import _ from 'lodash'
import ItemList from '../shared/ItemList'
import './styles.scss'

export interface ISummary {
  metadata: {
    type: 'summary';
    title: string;
  },
  config: {
    empty_content: string;
    sections: Array<{
      header: string;
      content: ContentType;
    }>;
  };
}

interface Props {
  data: ISummary,
}

export default function Summary({ data }: Props) {
  const { metadata: { title }, config: { sections }  } = data
  const items = _.map(sections, 'content')
  return (
    <div className='summary-component'>
      <h2 className='summary-component-title'>{title}</h2>
      <div className='summary-component-section'>
        <ItemList items={items} />
      </div>
    </div>
  )
}
