import React from 'react'
import _ from 'lodash'
import './styles.scss'

interface Props {
  labels: ContentType[];
}

export default function({ labels }: Props) {
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
