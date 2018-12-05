import React from 'react'
import _ from 'lodash'
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
  time: Time,
}

interface Props {
  items: ContentType[]
}

export default function({ items }: Props) {
  return (
    <React.Fragment>
      {_.map(items, (item: ContentType, index: number) => {
        const elem = dataTypeMap[item.type]
        if (!elem) return null
        return React.createElement(dataTypeMap[item.type], {
          key: index,
          params: item,
        })
      })}
    </React.Fragment>
  )
}
