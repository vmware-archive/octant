import React from 'react'
import './styles.scss'

export default function Item (props) {
  const { params } = props
  const {
    label,
    data: { value }
  } = params
  return (
    <div className='summary--data summary--data-basic'>
      {label && <div className='summary--data-key'>{label}</div>}
      <div className='summary--data-value'>{value}</div>
    </div>
  )
}
