import React from 'react'

import './styles.scss'

export default function Item (props) {
  const { params } = props
  const { key, value } = params
  return (
    <div className='summary--data summary--data-basic'>
      <div className='summary--data-key'>{key}</div>
      <div className='summary--data-basic'>{value}</div>
    </div>
  )
}
