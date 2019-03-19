import _ from 'lodash'
import React from 'react'
import { NavLink } from 'react-router-dom'

import './styles.scss'

interface Props {
  linkPath: NavigationSectionType[]
  onNavChange: (NavigationSectionType) => void
}

export default function NavigationHeader(props: Props) {
  const { linkPath, onNavChange } = props
  const { title, path } = _.last(linkPath)
  return (
    <h2 className='navigation--left-header'>
      <NavLink exact to={path} onClick={() => onNavChange(linkPath)}>
        {title}
      </NavLink>
    </h2>
  )
}
