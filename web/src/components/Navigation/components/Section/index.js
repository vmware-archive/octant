import React from 'react'
import Subheader from './components/Subheader'
import Header from './components/Header'
import './styles.scss'

export default function Section (props) {
  const { title, items = [], link = '/' } = props
  return (
    <div className='navigation--left-section'>
      <Header title={title} link={link} key={link} />
      <ul className='navigation--left-items'>
        {items.map(item => (
          <div key={item.title} className='navigation--left-item'>
            <Subheader item={item} />
          </div>
        ))}
      </ul>
    </div>
  )
}
