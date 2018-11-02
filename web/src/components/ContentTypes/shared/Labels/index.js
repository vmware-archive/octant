import React from 'react'
import _ from 'lodash'
import './styles.scss'

export default function ({ labels }) {
  return (
    <div className='content-labels'>
      {_.map(labels, ({ data }) => (
        <div key={data.value} className='content-label'>
          {data.value}
        </div>
      ))}
    </div>
  )
}
