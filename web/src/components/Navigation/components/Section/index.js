import React from 'react'
import Item from '../Item'
import Header from '../Header'

import './styles.scss'

export default function Section (props) {
  const { name, items = [], link = '/' } = props
  return (
    <div className='navigation--left-section'>
      <Header name={name} link={link} key={link} />
      <ul className='navigation--left-items'>
        {items.map(item => (
          <Item name={item.name} link={item.link} key={item.key} />
        ))}
      </ul>
    </div>
  )
}
