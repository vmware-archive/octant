import React from 'react'
import { NavLink } from 'react-router-dom'

import './styles.scss'

export default function Item (props) {
  const { title, link = '/' } = props
  return (
    <h2 className='navigation--left-header'>
      <NavLink exact to={link}>
        {title}
      </NavLink>
    </h2>
  )
}
