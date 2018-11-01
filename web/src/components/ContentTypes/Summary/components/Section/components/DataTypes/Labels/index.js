import React from 'react'
import Labels from 'components/ContentTypes/shared/Labels'

export default function Item (props) {
  const { params } = props
  const {
    label,
    data: { items }
  } = params
  return (
    <div className='summary--data'>
      <div className='summary--data-key'>{label}</div>
      <div className='summary--data-labels'>
        <Labels labels={items} />
      </div>
    </div>
  )
}
