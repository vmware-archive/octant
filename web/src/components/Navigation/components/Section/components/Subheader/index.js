import React from 'react'
import { Link } from 'react-router-dom'
import './styles.scss'

export default function Subheader (props) {
  const {
    item: { title, path, children = [] }
  } = props
  return (
    <li className='navigation--left-item'>
      <div className='navigation-subheader'>
        <Link to={path}>{title}</Link>
      </div>
      {children.map(({ title: childTitle, path: childPath }) => (
        <div key={childPath} className='navigation-subheader-link'>
          <Link to={childPath}>{childTitle}</Link>
        </div>
      ))}
    </li>
  )
}
