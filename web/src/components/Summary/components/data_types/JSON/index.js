import React from 'react'

import './styles.scss'

export default function Item (props) {
  const { params } = props
  const { key, data } = params
  return (
    <div className='summary--data summary--data-json'>
      <div className='summary--data-key'>{key}</div>
      <div className='summary--data-json-value'>{JSON.stringify(data, null, 2, 50)}</div>
    </div>
  )
}
