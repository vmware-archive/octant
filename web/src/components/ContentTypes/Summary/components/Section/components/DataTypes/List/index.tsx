import React from 'react'
import ItemList from '../../ItemList'

interface Props {
  params: ListContentType;
}

export default function Item(props: Props) {
  const { params } = props
  const {
    label,
    data: { items },
  } = params
  return (
    <div className='summary--data'>
      <div className='summary--data-key'>{label}</div>
      <div className='summary--data-list'>
        <ItemList items={items} />
      </div>
    </div>
  )
}
