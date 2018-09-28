import React from 'react'
import Item from '../Item'
import Header from '../Header'

import './styles.scss'

export default function Section (props) {
  const { title, items = [], link = '/' } = props
  return (
    <div className='navigation--left-section'>
      <Header title={title} link={link} key={link} />
      <ul className='navigation--left-items'>
        {items.map(item => (
          <div className='navigation--left-item'>
            <Item title={item.title} link={item.link} key={item.key} />
          </div>
        ))}
      </ul>
    </div>
  )
}
