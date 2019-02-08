import _ from 'lodash'
import React from 'react'
import { NavLink } from 'react-router-dom'

import './styles.scss'

interface Props {
  linkPath: NavigationSectionType[]
  childLinks: NavigationSectionType[]
  onNavChange: (NavigationSectionType) => void
}

export default function NavigationSubheader(props: Props) {
  const { linkPath: parentLinkPath, childLinks, onNavChange } = props
  const { title, path } = _.last(parentLinkPath)
  return (
    <li className='navigation--left-item'>
      <div className='navigation-subheader'>
        <NavLink exact to={path} onClick={() => onNavChange(parentLinkPath)}>
          {title}
        </NavLink>
      </div>
      {_.map(childLinks, (link) => {
        const { title: childTitle, path: childPath } = link
        return (
          <div key={childPath} className='navigation-subheader-link'>
            <NavLink to={childPath} onClick={() => onNavChange([...parentLinkPath, link])}>
              {childTitle}
            </NavLink>
          </div>
        )
      })}
    </li>
  )
}
