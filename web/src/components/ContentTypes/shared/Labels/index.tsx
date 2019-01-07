import React from 'react'
import _ from 'lodash'
import './styles.scss'

interface Props {
  labels: ContentType[];
}

export default function({ labels }: Props) {
  return (
    <div className='content-labels'>
      {_.map(labels, (value, key) => (
        <div key={key} className='content-label'>
          {value}
        </div>
      ))}
    </div>
  )
}
