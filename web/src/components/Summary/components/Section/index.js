import React from 'react'
import CONFIG from './_config'

import './styles.scss'

export default function Section (props) {
  const { title, data } = props
  return (
    <div className='summary--component-section'>
      <div className='summary--component-title'>
        <h2>{title}</h2>
      </div>
      {data.map(item => React.cloneElement(CONFIG.dataMap[item.type], { key: item.key, params: item }))}
    </div>
  )
}
