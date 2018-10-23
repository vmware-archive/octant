import React from 'react'
import CONFIG from './_config'

import './styles.scss'

export default function Section (props) {
  const { title, items } = props
  return (
    <div className='summary--component-section'>
      <div className='summary--component-title'>
        <h2>{title}</h2>
      </div>
      {items.map((item, index) => React.cloneElement(CONFIG.dataMap[item.type], {
        key: index,
        params: item
      }))}
    </div>
  )
}
