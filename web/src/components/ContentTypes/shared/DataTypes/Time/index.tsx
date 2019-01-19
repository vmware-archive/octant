import React from 'react'
import moment from 'moment'
import './styles.scss'

interface Props {
  params: ContentType;
}

export default function Time(props: Props) {
  const {
    metadata: { title },
    config: { timestamp },
  } = props.params
  let text = timestamp
  const t = moment(timestamp)
  if (t.isValid()) {
    text = `${t.fromNow()} - ${t.toString()}`
  }
  return (
    <div className='summary--data summary-data-time'>
      {title && <div className='summary--data-key'>{title}</div>}
      <div className='summary--data-value'>{text}</div>
    </div>
  )
}
