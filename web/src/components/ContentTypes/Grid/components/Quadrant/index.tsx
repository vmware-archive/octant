import './styles.scss'

import React from 'react'

interface QuadrantValue {
  value: number;
  label: string;
}

export interface IQuadrant {
  metadata: {
    type: 'quadrant';
    title: string;
  };
  config: {
    nw: QuadrantValue;
    ne: QuadrantValue;
    sw: QuadrantValue;
    se: QuadrantValue;
  };
}

interface Props {
  data: IQuadrant,
}

export default function Quadrant({ data }: Props) {
  const { metadata: { title }, config: { nw, ne, sw, se }} = data
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
