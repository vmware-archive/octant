import React from 'react'
import Text from './components/DataTypes/Text'
import JSON from './components/DataTypes/JSON'
import Link from './components/DataTypes/Link'
import Labels from './components/DataTypes/Labels'
import './styles.scss'

const dataTypeMap = {
  text: Text,
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
