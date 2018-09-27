import React from 'react'
import { Link } from 'react-router-dom'

import './styles.scss'

export default function Item (props) {
  const { name, link = '/' } = props
  return (
    <h2 className='navigation--left-header'>
      <Link to={link}>{name}</Link>
    </h2>
  )
}
