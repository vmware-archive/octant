import React from 'react'
import Text from '../DataTypes/Text'
import JSON from '../DataTypes/JSON'
import Link from '../DataTypes/Link'
import List from '../DataTypes/List'
import Labels from '../DataTypes/Labels'
import Time from '../DataTypes/Time'

const dataTypeMap = {
  text: Text,
  json: JSON,
  link: Link,
  labels: Labels,
  list: List,
  time: Time
}

export default function ({ items }) {
  return items.map((item, index) => {
    const elem = dataTypeMap[item.type]
    if (!elem) return null
    return React.createElement(dataTypeMap[item.type], {
      key: index,
      params: item
    })
  })
}
