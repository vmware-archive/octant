import React from 'react'
import './styles.scss'

interface Props {
  params: ContentType;
}

export default function Item(props: Props) {
  const { metadata: { title }, config } = props.params
  return (
    <div className='summary--data summary--data-json'>
      <div className='summary--data-key'>{title}</div>
      <div className='summary--data-json-value'>
        {JSON.stringify(config, null, 2)}
      </div>
    </div>
  )
}
