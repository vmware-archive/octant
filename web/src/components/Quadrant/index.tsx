import { QuadrantModel } from 'models/View'
import React from 'react'

import './styles.scss'

interface QuadrantValue {
  value: number
  label: string
}

interface Props {
  view: QuadrantModel
}

export default function Quadrant({ view }: Props) {
  const { nw, ne, sw, se, title } = view
  return (
    <div className='quadrant'>
      <div className='quadrant-header'>{title}</div>
      <div className='quadrant-body'>
        <div className='quadrant-ne'>
          <div className='quadrant-value'>{ne.value}</div>
          <div className='quadrant-label'>{ne.label}</div>
        </div>
        <div className='quadrant-nw'>
          <div className='quadrant-value'>{nw.value}</div>
          <div className='quadrant-label'>{nw.label}</div>
        </div>
        <div className='quadrant-se'>
          <div className='quadrant-value'>{se.value}</div>
          <div className='quadrant-label'>{se.label}</div>
        </div>
        <div className='quadrant-sw'>
          <div className='quadrant-value'>{sw.value}</div>
          <div className='quadrant-label'>{sw.label}</div>
        </div>
      </div>
    </div>
  )
}
