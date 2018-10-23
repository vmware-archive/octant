import React from 'react'
import { Link } from 'react-router-dom'

import './styles.scss'

export default function Item (props) {
  const { params } = props
  const {
    label,
    data: { value, ref }
  } = params
  return (
    <div className='summary--data summary--data-link'>
      <div className='summary--data-key'>{label}</div>
      <div className='summary--data-link'>
        <Link className='link--gray' to={ref}>
          {value || ref}
        </Link>
      </div>
    </div>
  )
}
