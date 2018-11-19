import React from 'react'
import moment from 'moment'
import './styles.scss'

export default function Time (props) {
  const { params } = props
  const {
    label,
    data: { value }
  } = params
  let text = value
  const t = moment(value)
  if (t.isValid()) {
    text = `${t.fromNow()} - ${t.toString()}`
  }
  return (
    <div className='summary--data summary-data-time'>
      {label && <div className='summary--data-key'>{label}</div>}
      <div className='summary--data-value'>{text}</div>
    </div>
  )
}
