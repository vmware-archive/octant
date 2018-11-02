import React from 'react'
import _ from 'lodash'
import Subheader from './components/Subheader'
import Header from './components/Header'
import './styles.scss'

export default function NavigationSection (props) {
  const { title, path, items } = props
  return (
    <div className='navigation--left-section'>
      <Header title={title} link={path} />
      <ul className='navigation--left-items'>
        {_.map(items, item => (
          <div key={item.title} className='navigation--left-item'>
            <Subheader item={item} />
          </div>
        ))}
      </ul>
    </div>
  )
}
