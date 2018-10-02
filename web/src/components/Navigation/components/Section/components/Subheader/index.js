import React from 'react'
import { Link } from 'react-router-dom'
import _ from 'lodash'
import './styles.scss'

export default function Subheader (props) {
  const {
    item: { title: subheader, path = '/', children }
  } = props
  return (
    <li className='navigation--left-item'>
      <div className='navigation-subheader'>{subheader}</div>
      {_.map(children, ({ title }) => (
        <div key={title} className='navigation-subheader-link'>
          <Link to={path}>{title}</Link>
        </div>
      ))}
    </li>
  )
}
