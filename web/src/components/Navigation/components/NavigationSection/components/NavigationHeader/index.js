import React from 'react'
import _ from 'lodash'
import { NavLink } from 'react-router-dom'
import './styles.scss'

export default function NavigationHeader (props) {
  const { currentLinkPath, linkPath, onNavChange } = props
  const { title, path } = _.last(linkPath)
  return (
    <h2 className='navigation--left-header'>
      <NavLink exact to={path} onClick={() => onNavChange(linkPath)}>
        {title}
      </NavLink>
    </h2>
  )
}
