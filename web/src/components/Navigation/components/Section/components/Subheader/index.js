import React from 'react'
import { Link } from 'react-router-dom'
import './styles.scss'

export default function Subheader (props) {
  const {
    item: { title: subheader, path = '/', children = [] }
  } = props
  return (
    <li className='navigation--left-item'>
      <div className='navigation-subheader'>{subheader}</div>
      {children.map(({ title }) => (
        <div key={title} className='navigation-subheader-link'>
          <Link to={path}>{title}</Link>
        </div>
      ))}
    </li>
  )
}
