import React from 'react'
import { Link } from 'react-router-dom'

import './styles.scss'

export default function Item (props) {
  const { title, link = '/' } = props
  return (
    <h2 className='navigation--left-header'>
      <Link to={link}>{title}</Link>
    </h2>
  )
}
