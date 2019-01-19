import React from 'react'
import './styles.scss'

interface Props {
  params: ContentType;
}

export default function Item(props: Props) {
  const { params } = props
  const {
    metadata: { title },
    config: { value },
  } = params
  return (
    <div className='summary--data summary--data-basic'>
      {title && <div className='summary--data-key'>{title}</div>}
      <div className='summary--data-value'>{value}</div>
    </div>
  )
}
