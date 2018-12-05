import React from 'react'
import ItemList from './components/ItemList'
import './styles.scss'

export default function Section(props: ContentSection) {
  const { title, items } = props
  return (
    <div className='summary-component-section'>
      <div className='summary-component-title'>
        <h2>{title}</h2>
      </div>
      <ItemList items={items} />
    </div>
  )
}
