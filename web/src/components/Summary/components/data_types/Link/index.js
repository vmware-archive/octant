import React from 'react'
import { Link } from 'react-router-dom'

import './styles.scss'

export default function Item (props) {
  const { params } = props
  const { key, link, value } = params
  return (
    <div className='summary--data summary--data-link'>
      <div className='summary--data-key'>{key}</div>
      <div className='summary--data-link'>
        <Link className='link--gray' to={link}>
          {value}
        </Link>
      </div>
    </div>
  )
}
