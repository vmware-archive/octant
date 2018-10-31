import React from 'react'
import Basic from './components/DataTypes/Basic'
import JSON from './components/DataTypes/JSON'
import Link from './components/DataTypes/Link'
import Labels from '../../../shared/Labels'
import './styles.scss'

const dataTypeMap = {
  text: Basic,
  json: JSON,
  link: Link,
  labels: Labels
}

export default function Section (props) {
  const { title, items } = props
  return (
    <div className='summary-component-section'>
      <div className='summary-component-title'>
        <h2>{title}</h2>
      </div>
      {items.map((item, index) => {
        const elem = dataTypeMap[item.type]
        if (!elem) return null
        return React.createElement(dataTypeMap[item.type], {
          key: index,
          params: item
        })
      })}
    </div>
  )
}
