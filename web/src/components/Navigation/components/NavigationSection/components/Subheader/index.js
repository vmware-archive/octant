import React from 'react'
import { NavLink } from 'react-router-dom'
import './styles.scss'

export default function Subheader (props) {
  const {
    item: { title, path, children = [] }
  } = props
  return (
    <li className='navigation--left-item'>
      <div className='navigation-subheader'>
        <NavLink exact to={path}>
          {title}
        </NavLink>
      </div>
      {children.map(({ title: childTitle, path: childPath }) => (
        <div key={childPath} className='navigation-subheader-link'>
          <NavLink to={childPath}>{childTitle}</NavLink>
        </div>
      ))}
    </li>
  )
}
