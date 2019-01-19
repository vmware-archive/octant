import React from 'react'
import ItemList from 'components/ContentTypes/shared/ItemList'
import './styles.scss'

export interface ILabel {
  metadata: {
    type: 'label';
    title: string;
  },
  config: {
    contents: [LabelsContentType];
  }
}

interface Props {
  data: ILabel;
}

export default function Label(props: Props) {
  const { metadata: { title }, config: { contents } } = props.data
  return (
    <div className='grid-label'>
      <div className='grid-label-title'>
        {title}
      </div>
      <div className='grid-label-sub'>
        <ItemList items={contents} />
      </div>
    </div>
  )
}
