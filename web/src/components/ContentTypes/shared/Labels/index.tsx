import './styles.scss'

import React from 'react'

interface Props {
  labels: {[key: string]: string};
}

export default function({ labels }: Props) {
  return (
    <div className='content-labels'>
      {Object.entries(labels).map(([key, value], index) => (
        <div key={index} className='content-label'>
        {key}:{value}
      </div>
      ))}
    </div>
  )
}
