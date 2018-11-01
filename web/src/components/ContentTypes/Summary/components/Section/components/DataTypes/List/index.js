import React from 'react'
import ItemList from '../../ItemList'

export default function Item (props) {
  const { params } = props
  const {
    label,
    data: { items }
  } = params
  return (
    <div className='summary--data'>
      <div className='summary--data-key'>{label}</div>
      <div className='summary--data-basic'>
        <ItemList items={items} />
      </div>
    </div>
  )
}
