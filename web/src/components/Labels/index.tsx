import './styles.scss'

import { LabelsModel } from 'models/View'
import React from 'react'

interface Props {
  view: LabelsModel
}

export default function({ view }: Props) {
  return (
    <div className='content-labels'>
      {Object.entries(view.labels).map(([key, value], index) => (
        <div key={index} className='content-label'>
        {key}:{value}
      </div>
      ))}
    </div>
  )
}
