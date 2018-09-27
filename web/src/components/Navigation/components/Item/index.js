import React from 'react'
import { Link } from 'react-router-dom'

import './styles.scss'

export default function Item (props) {
  const { name, link = '/' } = props
  return (
    <li className='navigation--left-item'>
      <Link className='link--gray' to={link}>
        {name}
      </Link>
    </li>
  )
}
